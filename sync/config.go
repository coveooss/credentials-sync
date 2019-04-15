package sync

import (
	"github.com/coveo/credentials-sync/credentials"
	"github.com/coveo/credentials-sync/targets"
	log "github.com/sirupsen/logrus"
)

type Configuration struct {
	Sources     *credentials.SourcesConfiguration
	StopOnError bool `mapstructure:"stop_on_error"`
	Targets     *targets.Configuration
}

func NewConfiguration() *Configuration {
	return &Configuration{
		StopOnError: false,
	}
}

func (config *Configuration) Sync() {
	validTargets := []targets.Target{}
	// Initialize targets
	allTargets := config.Targets.AllTargets()
	initChannel := make(chan targets.Target)
	for _, target := range allTargets {
		go initTarget(target, initChannel)
	}
	for i := 0; i < len(allTargets); i++ {
		initTarget := <-initChannel
		if initTarget != nil {
			validTargets = append(validTargets, initTarget)
		} else if config.StopOnError {
			return
		}
	}
}

func initTarget(target targets.Target, channel chan targets.Target) {
	err := target.Initialize([]credentials.Credentials{})
	if err == nil {
		channel <- target
	} else {
		log.Errorf("Target `%s` has failed initialization. Ignoring it.", target.GetName())
		channel <- nil
	}
}
