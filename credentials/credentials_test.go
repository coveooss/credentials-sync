package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSyncCredentials(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		creds      *Base
		targetName string
		targetTags map[string]string
		expected   bool
	}{
		{
			name:     "No tags",
			creds:    &Base{},
			expected: true,
		},
		{
			name:     "Cred that shouldn't be synced",
			creds:    &Base{NoSync: true},
			expected: false,
		},
		{
			name:       "Non matching target name",
			targetName: "Target",
			creds:      &Base{TargetName: "Other target"},
			expected:   false,
		},
		{
			name:       "Matching target name",
			targetName: "Target",
			creds:      &Base{TargetName: "Target"},
			expected:   true,
		},
		{
			name:  "No filter",
			creds: &Base{},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: true,
		},
		{
			name: "Match",
			creds: &Base{TargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyFirstTag": "MyValue",
					"MyTag":      "MyValue",
				},
			}},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: true,
		},
		{
			name: "Match but not target name",
			creds: &Base{
				TargetName: "Target",
				TargetTags: targetTagsMatcher{
					DoMatch: map[string]interface{}{
						"MyFirstTag": "MyValue",
						"MyTag":      "MyValue",
					},
				}},
			targetName: "Other",
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "Match List Item",
			creds: &Base{TargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyTag": []string{"FirstValue", "MyValue"},
				},
			}},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: true,
		},
		{
			name: "List Without Matches",
			creds: &Base{TargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyTag": []string{"FirstValue", "SecondValue"},
				},
			}},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "String Doesnt Match",
			creds: &Base{TargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyTag": "AValue",
				},
			}},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "String that shouldn't match",
			creds: &Base{TargetTags: targetTagsMatcher{
				DontMatch: map[string]interface{}{
					"MyTag": "MyValue",
				},
			}},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "String in list that shouldn't match",
			creds: &Base{TargetTags: targetTagsMatcher{
				DontMatch: map[string]interface{}{
					"MyTag": []string{"Test", "MyValue"},
				},
			}},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "Match and exclude",
			creds: &Base{TargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyTag": "Value",
				},
				DontMatch: map[string]interface{}{
					"MyOtherTag": []string{"Test", "MyValue"},
				},
			}},
			targetTags: map[string]string{
				"MyTag":      "Value",
				"MyOtherTag": "MyValue",
			},
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.creds.ShouldSync(tt.targetName, tt.targetTags))
		})
	}
}

func TestGetDescriptionOrID(t *testing.T) {
	t.Parallel()

	cred := &Base{
		ID:       "test",
		CredType: "aType",
	}

	assert.Equal(t, "test", cred.GetDescriptionOrID())
	assert.Equal(t, "test -> Type: aType", cred.BaseToString())

	cred.Description = "other"
	assert.Equal(t, "other", cred.GetDescriptionOrID())
	assert.Equal(t, "test -> Type: aType, Description: other", cred.BaseToString())
}
