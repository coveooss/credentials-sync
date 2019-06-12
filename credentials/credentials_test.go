package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSyncCredentials(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		wantedTargetTags targetTagsMatcher
		targetTags       map[string]string
		expected         bool
	}{
		{
			name: "No tags",
			wantedTargetTags: targetTagsMatcher{
				DoMatch:   map[string]interface{}{},
				DontMatch: map[string]interface{}{},
			},
			targetTags: map[string]string{},
			expected:   true,
		},
		{
			name: "No filter",
			wantedTargetTags: targetTagsMatcher{
				DoMatch:   map[string]interface{}{},
				DontMatch: map[string]interface{}{},
			},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: true,
		},
		{
			name: "Match",
			wantedTargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyFirstTag": "MyValue",
					"MyTag":      "MyValue",
				},
				DontMatch: map[string]interface{}{},
			},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: true,
		},
		{
			name: "Match List Item",
			wantedTargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyTag": []string{"FirstValue", "MyValue"},
				},
				DontMatch: map[string]interface{}{},
			},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: true,
		},
		{
			name: "List Without Matches",
			wantedTargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyTag": []string{"FirstValue", "SecondValue"},
				},
				DontMatch: map[string]interface{}{},
			},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "String Doesnt Match",
			wantedTargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyTag": "AValue",
				},
				DontMatch: map[string]interface{}{},
			},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "String that shouldn't match",
			wantedTargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{},
				DontMatch: map[string]interface{}{
					"MyTag": "MyValue",
				},
			},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "String in list that shouldn't match",
			wantedTargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{},
				DontMatch: map[string]interface{}{
					"MyTag": []string{"Test", "MyValue"},
				},
			},
			targetTags: map[string]string{
				"MyTag": "MyValue",
			},
			expected: false,
		},
		{
			name: "Match and exclude",
			wantedTargetTags: targetTagsMatcher{
				DoMatch: map[string]interface{}{
					"MyTag": "Value",
				},
				DontMatch: map[string]interface{}{
					"MyOtherTag": []string{"Test", "MyValue"},
				},
			},
			targetTags: map[string]string{
				"MyTag":      "Value",
				"MyOtherTag": "MyValue",
			},
			expected: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			credentials := &Base{TargetTags: tt.wantedTargetTags}
			assert.Equal(t, tt.expected, credentials.ShouldSync(tt.targetTags))
		})
	}
}
