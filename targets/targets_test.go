package targets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidateConfiguration(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		targets  []*JenkinsTarget
		expected bool
	}{
		{
			name:     "no targets",
			targets:  []*JenkinsTarget{},
			expected: true,
		},
		{
			name: "valid target",
			targets: []*JenkinsTarget{
				&JenkinsTarget{
					Base: Base{Name: "test"},
					URL:  "https://test.com",
				},
			},
			expected: true,
		},
		{
			name: "no name",
			targets: []*JenkinsTarget{
				&JenkinsTarget{
					Base: Base{Name: ""},
					URL:  "https://test.com",
				},
			},
			expected: false,
		},
		{
			name: "two actions for unsynced credentials",
			targets: []*JenkinsTarget{
				&JenkinsTarget{
					Base: Base{Name: "test", DeleteUnsynced: true, TagUnsynced: true},
					URL:  "https://test.com",
				},
			},
			expected: false,
		},
		{
			name: "bad url (validated by Jenkins)",
			targets: []*JenkinsTarget{
				&JenkinsTarget{
					Base: Base{Name: "test"},
					URL:  "bad",
				},
			},
			expected: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			config := &Configuration{JenkinsTargets: tt.targets}
			assert.Equal(t, tt.expected, config.ValidateConfiguration())
			for i, gottenItem := range config.AllTargets() {
				assert.Equal(t, tt.targets[i], gottenItem)
			}
		})
	}
}
