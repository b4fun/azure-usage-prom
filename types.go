package azureusageprom

import (
	"context"
	"fmt"
	"strings"
)

// AzureUsageName - name of the usage.
type AzureUsageName struct {
	Value          string `json:"value"`
	LocalizedValue string `json:"localizedValue"`
}

// AzureUsageUnit - the usage unit enum.
type AzureUsageUnit string

// AzureUsage - Azure usage data.
// ref: https://docs.microsoft.com/en-us/rest/api/virtualnetwork/usages/list#usage
type AzureUsage struct {
	ID           string         `json:"id"`
	CurrentValue int            `json:"currentValue"`
	Limit        int            `json:"limit"`
	Name         AzureUsageName `json:"name"`
	Unit         AzureUsageUnit `json:"unit"`
}

// AzureUsageList - Azure usage data list.
type AzureUsageList struct {
	Value []AzureUsage `json:"value"`
}

// ResourceProvider - Azure resource provider enum.
type ResourceProvider string

// ParseResourceProvider parses supported resource provider.
func ParseResourceProvider(s string) (ResourceProvider, error) {
	for _, rp := range []ResourceProvider{
		ResourceProviderCompute,
		ResourceProviderNetwork,
	} {
		if strings.EqualFold(s, string(rp)) {
			return rp, nil
		}
	}

	return ResourceProvider(""), fmt.Errorf("unknown resource provider: %s", s)
}

const (
	// ResourceProviderNetwork - network RP
	// https://docs.microsoft.com/en-us/rest/api/virtualnetwork/usages/list
	ResourceProviderNetwork ResourceProvider = "Microsoft.Network"
	// ResourceProviderCompute - compute RP
	// https://docs.microsoft.com/en-us/rest/api/compute/usage/list
	ResourceProviderCompute ResourceProvider = "Microsoft.Compute"
)

// APIVersion returns the API version to use for a resource provider.
func (rp ResourceProvider) APIVersion() string {
	switch rp {
	case ResourceProviderNetwork:
		return "2020-05-01"
	case ResourceProviderCompute:
		return "2020-06-01"
	default:
		return ""
	}
}

// ToPromNamespace converts resource provider to prometheus namespace.
func (rp ResourceProvider) ToPromNamespace() string {
	// Microsoft.Network -> microsoft_network
	return strings.ToLower(strings.ReplaceAll(string(rp), ".", "_"))
}

// QueryTarget defines the usage query target.
type QueryTarget struct {
	ResourceProvider ResourceProvider
	SubscriptionID   string
	Location         string
}

// AzureUsageLister lists Azure usage metrics.
type AzureUsageLister interface {
	// ListByResourceProvider lists usage by resource provider.
	ListByResourceProvider(
		ctx context.Context,
		target QueryTarget,
	) (AzureUsageList, error)
}
