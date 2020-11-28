package main

import (
	"testing"

	azureusageprom "github.com/b4fun/azure-usage-prom"
)

func TestParseQueryTargets(t *testing.T) {
	cases := []struct {
		s           string
		expectError bool
		expected    []azureusageprom.QueryTarget
	}{
		{
			s:        "",
			expected: []azureusageprom.QueryTarget{},
		},
		{
			s:           "foo|bar",
			expectError: true,
		},
		{
			s:           "foo|bar|baz",
			expectError: true,
		},
		{
			s: "microsoft.compute|test-sub-id|japaneast",
			expected: []azureusageprom.QueryTarget{
				{
					ResourceProvider: azureusageprom.ResourceProviderCompute,
					SubscriptionID:   "test-sub-id",
					Location:         "japaneast",
				},
			},
		},
		{
			s:           "microsoft.compute|test-sub-id|japaneast,foo,bar,baz",
			expectError: true,
		},
		{
			s: "microsoft.compute|test-sub-id|japaneast,Microsoft.Network|test-sub-id-2|westus",
			expected: []azureusageprom.QueryTarget{
				{
					ResourceProvider: azureusageprom.ResourceProviderCompute,
					SubscriptionID:   "test-sub-id",
					Location:         "japaneast",
				},
				{
					ResourceProvider: azureusageprom.ResourceProviderNetwork,
					SubscriptionID:   "test-sub-id-2",
					Location:         "westus",
				},
			},
		},
	}

	for _, c := range cases {
		v, err := parseQueryTargets(c.s)
		if c.expectError {
			if err == nil {
				t.Errorf("parseQueryTargets(%s) expected error, got nil", c.s)
			}
		} else {
			if err != nil {
				t.Errorf("parseQueryTargets(%s) unexpected error: %s", c.s, err)
			}
			if len(v) != len(c.expected) {
				t.Errorf("parseQueryTargets(%s) expected: %v, got: %v", c.s, c.expected, v)
			}
		}
	}
}
