package targets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJenkinsToString(t *testing.T) {
	jenkins := &JenkinsTarget{
		Base: Base{Name: "targetName", Tags: map[string]string{"my_tag": "tag_value"}},
		URL:  "test.com",
	}

	assert.Equal(t, "targetName [Tags: my_tag=tag_value] (Jenkins) - test.com", jenkins.ToString())
}
