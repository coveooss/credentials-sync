package sync

import (
	"fmt"

	"github.com/coveooss/credentials-sync/credentials"
	"github.com/coveooss/credentials-sync/targets"
	log "github.com/sirupsen/logrus"
)

// Configuration represents the parsed configuration file given to the application
type Configuration struct {
	CredentialsToDelete []string                          `mapstructure:"credentials_to_delete"`
	Sources             *credentials.SourcesConfiguration `mapstructure:"sources"`
	StopOnError         bool                              `mapstructure:"stop_on_error"`
	TargetParallelism   int                               `mapstructure:"target_parallelism"`
	Targets             *targets.Configuration            `mapstructure:"targets"`
}

// NewConfiguration creates a new configuration with default values
func NewConfiguration() *Configuration {
	return &Configuration{
		StopOnError:       false,
		TargetParallelism: 4,
	}
}

// Sync syncs credentials from the configured sources to the configured targets
func (config *Configuration) Sync() {
	// Start reading credentials
	creds, err := config.Sources.Credentials()
	if err != nil {
		log.Fatalf("Caught an error while fetching credentials: %v", err)
	}

	// Initialize targets
	validTargets := []targets.Target{}
	allTargets := config.Targets.AllTargets()
	initChannel := make(chan targets.Target)
	for _, target := range allTargets {
		go initTarget(target, creds, initChannel, config.StopOnError)
	}
	for i := 0; i < len(allTargets); i++ {
		initTarget := <-initChannel
		if initTarget != nil {
			validTargets = append(validTargets, initTarget)
		}
	}

	syncChannel := make(chan bool, config.TargetParallelism)
	for _, target := range validTargets {
		syncChannel <- true
		go syncCredentials(target, creds, syncChannel, config.StopOnError)
	}

	for i := 0; i < cap(syncChannel); i++ {
		syncChannel <- true
	}

}

func initTarget(target targets.Target, creds []credentials.Credentials, channel chan targets.Target, panicOnError bool) {
	err := target.Initialize(creds)
	if err == nil {
		log.Infof("Connected to %s", target.ToString())
		channel <- target
	} else {
		message := fmt.Sprintf("Target `%s` has failed initialization: %v", target.GetName(), err)
		if panicOnError {
			log.Fatal(message)
		}
		log.Warning(message)
		channel <- nil
	}
}

func syncCredentials(target targets.Target, credentialsList []credentials.Credentials, channel chan bool, panicOnError bool) {
	defer func() { <-channel }()

	filteredCredentials := []credentials.Credentials{}
	for _, cred := range credentialsList {
		if cred.ShouldSync(target.GetName(), target.GetTags()) {
			filteredCredentials = append(filteredCredentials, cred)
		}
	}

	if err := targets.UpdateListOfCredentials(target, filteredCredentials); err != nil {
		message := fmt.Sprintf("Failed to send credentials to %s: %v", target.GetName(), err)
		if panicOnError {
			log.Fatal(message)
		}
		log.Error(message)
	}

	log.Infof("Finished sync to %s", target.GetName())
}
