package sync

import (
	"github.com/coveo/credentials-sync/credentials"
	"github.com/coveo/credentials-sync/targets"
)

type Configuration struct {
	Sources *credentials.SourcesConfiguration
	Targets *targets.Configuration
}
