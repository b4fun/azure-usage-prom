package collector

import (
	"context"
	"errors"
	"testing"

	azureusageprom "github.com/b4fun/azure-usage-prom"
	"github.com/prometheus/client_golang/prometheus"
)

type mockAzureUsageLister struct {
	listByResourceProvider func(
		ctx context.Context,
		queryTarget azureusageprom.QueryTarget,
	) (azureusageprom.AzureUsageList, error)
}

func (m mockAzureUsageLister) ListByResourceProvider(
	ctx context.Context,
	queryTarget azureusageprom.QueryTarget,
) (azureusageprom.AzureUsageList, error) {
	return m.listByResourceProvider(ctx, queryTarget)
}

func TestCollector_AsPrometheusCollector(t *testing.T) {
	queryTarget := azureusageprom.QueryTarget{
		ResourceProvider: azureusageprom.ResourceProviderCompute,
		SubscriptionID:   "test-subscription-id",
		Location:         "test-location",
	}

	{
		t.Logf("returning empty ussage...")
		usageLister := &mockAzureUsageLister{}
		usageLister.listByResourceProvider = func(
			ctx context.Context,
			queryTarget azureusageprom.QueryTarget,
		) (azureusageprom.AzureUsageList, error) {
			return azureusageprom.AzureUsageList{}, nil
		}

		descChan := make(chan *prometheus.Desc, 100)
		metricChan := make(chan prometheus.Metric, 100)
		defer close(descChan)
		defer close(metricChan)

		collector := New(usageLister, queryTarget)
		collector.Describe(descChan)
		collector.Collect(metricChan)
		// no input
	}

	{
		t.Logf("returning 1 ussage...")
		usageLister := &mockAzureUsageLister{}
		usageLister.listByResourceProvider = func(
			ctx context.Context,
			queryTarget azureusageprom.QueryTarget,
		) (azureusageprom.AzureUsageList, error) {
			return azureusageprom.AzureUsageList{
				Value: []azureusageprom.AzureUsage{
					{
						ID:           "test-id",
						CurrentValue: 1,
						Limit:        100,
						Name: azureusageprom.AzureUsageName{
							Value:          "test-name",
							LocalizedValue: "test-name",
						},
						Unit: azureusageprom.AzureUsageUnit("Count"),
					},
				},
			}, nil
		}

		descChan := make(chan *prometheus.Desc, 100)
		metricChan := make(chan prometheus.Metric, 100)
		defer close(descChan)
		defer close(metricChan)

		collector := New(usageLister, queryTarget)
		collector.Describe(descChan)
		collector.Collect(metricChan)
		for i := 0; i < 2; i++ {
			// NOTE: should be able to read 2 descriptors
			<-descChan
		}
		for i := 0; i < 2; i++ {
			// NOTE: should be able to read 2 metrics
			<-metricChan
		}
	}

	{
		t.Logf("returning error...")
		usageLister := &mockAzureUsageLister{}
		usageLister.listByResourceProvider = func(
			ctx context.Context,
			queryTarget azureusageprom.QueryTarget,
		) (azureusageprom.AzureUsageList, error) {
			return azureusageprom.AzureUsageList{}, errors.New("error")
		}

		descChan := make(chan *prometheus.Desc, 100)
		metricChan := make(chan prometheus.Metric, 100)
		defer close(descChan)
		defer close(metricChan)

		collector := New(usageLister, queryTarget)
		collector.Describe(descChan)
		collector.Collect(metricChan)
		// no input
	}
}
