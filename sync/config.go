package sync

import (
	"fmt"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/logger"
	"github.com/coveooss/credentials-sync/targets"
	"github.com/hashicorp/go-multierror"
)

// Configuration represents the parsed configuration file given to the application
type Configuration struct {
	CredentialsToDelete []string                     `mapstructure:"credentials_to_delete"`
	Sources             credentials.SourceCollection `mapstructure:"-"`
	StopOnError         bool                         `mapstructure:"stop_on_error"`
	TargetParallelism   int                          `mapstructure:"target_parallelism"`
	Targets             targets.TargetCollection     `mapstructure:"-"`
}

// NewConfiguration creates a new configuration with default values
func NewConfiguration() *Configuration {
	return &Configuration{
		StopOnError:       false,
		TargetParallelism: 4,
	}
}

// SetSources sets the source configuration on synchronization configuration
func (config *Configuration) SetSources(sources credentials.SourceCollection) {
	config.Sources = sources
}

// SetTargets sets the target configuration on synchronization configuration
func (config *Configuration) SetTargets(targets targets.TargetCollection) {
	config.Targets = targets
}

// Sync syncs credentials from the configured sources to the configured targets
func (config *Configuration) Sync() error {
	// Start reading credentials
	creds, err := config.Sources.Credentials()
	if err != nil {
		return fmt.Errorf("Caught an error while fetching credentials: %v", err)
	}

	// Initialize targets
	validTargets := []targets.Target{}
	allTargets := config.Targets.AllTargets()
	initChannel := make(chan interface{})
	for _, target := range allTargets {
		go config.initTarget(target, creds, initChannel)
	}
	// We will use this to accumulate errors that happen if config.StopOnError is set to false
	// the multierror.Error implements error so we use the interface to type the accumulator
	var errorAccumulator error
	for i := 0; i < len(allTargets); i++ {
		initTarget := <-initChannel
		if err, ok := initTarget.(error); ok {
			if config.StopOnError {
				return err
			}
			errorAccumulator = multierror.Append(errorAccumulator, err)
			logger.Log.Error(err)
		} else {
			validTargets = append(validTargets, initTarget.(targets.Target))
		}
	}

	// Sync credentials with as many targets as the config allows
	parallelismChannel := make(chan bool, config.TargetParallelism)
	errorChannel := make(chan error)
	for _, target := range validTargets {
		parallelismChannel <- true
		go config.syncCredentials(target, creds, parallelismChannel, errorChannel)

		// Check for errors. Errors are only passed back if StopOnError is true so this should always return
		err := <-errorChannel
		if err != nil {
			if config.StopOnError {
				return err
			}
			errorAccumulator = multierror.Append(errorAccumulator, err)
		}
	}

	// Ensure that the sync method is completely done for all targets
	for i := 0; i < cap(parallelismChannel); i++ {
		parallelismChannel <- true
	}

	// This is either a nil, or a collection of past errors which we want to bubble up
	return errorAccumulator
}

func (config *Configuration) initTarget(target targets.Target, creds []credentials.Credentials, channel chan interface{}) {
	var channelValue interface{}

	defer func() {
		channel <- channelValue
	}()

	err := target.Initialize(creds)
	if err == nil {
		logger.Log.Infof("Connected to %s", target.ToString())
		channelValue = target
	} else {
		channelValue = fmt.Errorf("Target `%s` has failed initialization: %v", target.GetName(), err)
	}
}

func (config *Configuration) syncCredentials(target targets.Target, credentialsList []credentials.Credentials, parallelismChannel chan bool, errorChannel chan error) {
	// We will use this to accumulate errors that happen if config.StopOnError is set to false
	// the multierror.Error implements error so we use the interface to type the accumulator
	var errorAccumulator error
	defer func() {
		errorChannel <- errorAccumulator
		<-parallelismChannel
	}()

	filteredCredentials := []credentials.Credentials{}
	for _, cred := range credentialsList {
		if cred.ShouldSync(target.GetName(), target.GetTags()) {
			filteredCredentials = append(filteredCredentials, cred)
		}
	}

	if err := config.UpdateListOfCredentials(target, filteredCredentials); err != nil {
		errorAccumulator = multierror.Append(errorAccumulator, err)
		if config.StopOnError {
			return
		}
	}
	if err := config.DeleteListOfCredentials(target); err != nil {
		errorAccumulator = multierror.Append(errorAccumulator, err)
		if config.StopOnError {
			return
		}
	}
	logger.Log.Infof("Finished sync to %s", target.GetName())
}
