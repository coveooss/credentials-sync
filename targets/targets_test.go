package targets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidateConfiguration(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		targets     []*JenkinsTarget
		expectError bool
	}{
		{
			name:        "no targets",
			targets:     []*JenkinsTarget{},
			expectError: false,
		},
		{
			name: "valid target",
			targets: []*JenkinsTarget{
				{
					Base: Base{Name: "test"},
					URL:  "https://test.com",
				},
			},
			expectError: false,
		},
		{
			name: "no name",
			targets: []*JenkinsTarget{
				{
					Base: Base{Name: ""},
					URL:  "https://test.com",
				},
			},
			expectError: true,
		},
		{
			name: "two actions for unsynced credentials",
			targets: []*JenkinsTarget{
				{
					Base: Base{Name: "test", DeleteUnsynced: true, TagUnsynced: true},
					URL:  "https://test.com",
				},
			},
			expectError: true,
		},
		{
			name: "bad url (validated by Jenkins)",
			targets: []*JenkinsTarget{
				{
					Base: Base{Name: "test"},
					URL:  "bad",
				},
			},
			expectError: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config := &Configuration{JenkinsTargets: tt.targets}
			assert.Equal(t, tt.expectError, config.ValidateConfiguration() != nil)
			for i, gottenItem := range config.AllTargets() {
				assert.Equal(t, tt.targets[i], gottenItem)
			}
		})
	}
}
