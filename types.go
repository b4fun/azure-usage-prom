package azureusageprom

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
