package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextTag(t *testing.T) {
	tr := NewTagResolver("v1.0.0", "d2iq")

	testCases := []struct {
		nextTag      string
		existingTags []string
	}{
		{
			existingTags: []string{"v1.0.0-d2iq.3"},
			nextTag:      "v1.0.0-d2iq.4",
		},
		{
			existingTags: []string{},
			nextTag:      "v1.0.0-d2iq.0",
		},
		{
			existingTags: []string{"foo", "bar"},
			nextTag:      "v1.0.0-d2iq.0",
		},
		{
			existingTags: []string{"v1.0.0-d2iq.20"},
			nextTag:      "v1.0.0-d2iq.21",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.nextTag, func(t *testing.T) {
			assert.Equal(
				t,
				tc.nextTag,
				tr.Next(tc.existingTags),
			)
		})
	}
}

func TestParseBaseTag(t *testing.T) {
	testCases := []struct {
		tag      string
		expected string
	}{
		{"v1.0.0-d2iq.0", "v1.0.0"},
		{"v1.0.0", ""},
		{"v1", ""},
		{"v1-d2iq.0", "v1"},
		{"v1.0.0-d2iq.10-d2iq.5", "v1.0.0-d2iq.10"},
	}

	for _, tc := range testCases {
		t.Run(tc.tag, func(t *testing.T) {
			assert.Equal(t, tc.expected, ParseBaseTag(tc.tag))
		})
	}
}
