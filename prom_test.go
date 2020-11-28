package azureusageprom

import "testing"

func TestAzureUsageToMetricLabels(t *testing.T) {
	target := QueryTarget{}
	usage := AzureUsage{
		Name: AzureUsageName{
			Value: "foobar",
		},
	}

	labels := AzureUsageToMetricLabels(target, usage)
	if len(labels) != 1 {
		t.Errorf("invalid metric labels: %v", labels)
	}
	if v := labels[0]; v != usage.Name.Value {
		t.Errorf("invalid metric label: %s, expected: %s", v, usage.Name.Value)
	}
}

func TestAzureUsageToPromDescs(t *testing.T) {
	target := QueryTarget{
		ResourceProvider: ResourceProviderCompute,
		SubscriptionID:   "test-subscription-id",
		Location:         "test-location",
	}
	usage := AzureUsage{
		Name: AzureUsageName{
			Value: "foobar",
		},
	}

	limitDesc, currentValueDesc := AzureUsageToPromDescs(target, usage)
	if limitDesc == nil || currentValueDesc == nil {
		t.Errorf("nil descriptor returned")
	}
}
