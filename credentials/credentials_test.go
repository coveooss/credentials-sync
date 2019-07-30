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
			name: "Bad filter (not string or list), ignored",
			creds: &Base{
				TargetTags: targetTagsMatcher{
					DontMatch: map[string]interface{}{
						"MyTag": 123,
					},
				},
			},
			targetTags: map[string]string{
				"MyTag": "123",
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

func TestBaseValidateCredentials(t *testing.T) {
	t.Parallel()

	credWithoutType := &Base{
		ID: "test",
	}
	assert.EqualError(t, credWithoutType.BaseValidate(), "Credentials (test) has no type. This is probably a bug in the software")

	credWithoutID := &Base{
		CredType:    "test",
		Description: "test2",
	}
	assert.EqualError(t, credWithoutID.BaseValidate(), "Credentials ( -> Type: test, Description: test2) has no defined ID")
}

func TestParseCredentials(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		credMaps []map[string]interface{}
		result   []Credentials
		wantErr  bool
	}{
		{
			name: "Invalid type",
			credMaps: []map[string]interface{}{
				{
					"id":          "stuff",
					"type":        "bad",
					"description": "test-desc",
					"secret":      "my secret",
				},
			},
			result:  nil,
			wantErr: true,
		},
		{
			name: "Invalid type (not a string)",
			credMaps: []map[string]interface{}{
				{
					"id":          "stuff",
					"type":        1234,
					"description": "test-desc",
					"secret":      "my secret",
				},
			},
			result:  nil,
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gottenCreds, err := ParseCredentials(tt.credMaps)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.result, gottenCreds)
		})
	}
}
