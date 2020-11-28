package azureusageprom

import (
	"strings"
	"testing"
)

func TestParseResourceProvider(t *testing.T) {
	cases := []struct {
		s           string
		expected    ResourceProvider
		expectError bool
	}{
		{
			s:           "",
			expectError: true,
		},
		{
			s:        string(ResourceProviderNetwork),
			expected: ResourceProviderNetwork,
		},
		{
			s:        strings.ToLower(string(ResourceProviderNetwork)),
			expected: ResourceProviderNetwork,
		},
		{
			s:        string(ResourceProviderCompute),
			expected: ResourceProviderCompute,
		},
	}

	for _, c := range cases {
		v, err := ParseResourceProvider(c.s)
		if c.expectError {
			if err == nil {
				t.Errorf("ParseResourceProvider(%s): expect error, got nil", c.s)
			}
		} else {
			if err != nil {
				t.Errorf("ParseResourceProvider(%s): expect error, got: %s", c.s, err)
			}
			if v != c.expected {
				t.Errorf("ParseResourceProvider(%s) != %s, got: %s", c.s, c.expected, v)
			}
		}
	}
}
