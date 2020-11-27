package azureusageprom

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// PrometheusSubsystemUsage - sub system of the metrics.
	PrometheusSubsystemUsage = "usage"

	// PrometheusLabelUsageName - usage name label
	PrometheusLabelUsageName = "usage_name"
	// PrometheusLabelSubscriptionID - subscription id label
	PrometheusLabelSubscriptionID = "subscription_id"
	// PrometheusLabelLocation - location label
	PrometheusLabelLocation = "location"

	// PrometheusMetricLimit - usage limit metric
	PrometheusMetricLimit = "limit"
	// PrometheusMetricCurrentValue - usage limit value
	PrometheusMetricCurrentValue = "current_value"
)

// AzureUsageToMetricLabels converts usage to prometheus metric labels.
func AzureUsageToMetricLabels(
	target QueryTarget,
	usage AzureUsage,
) []string {
	return []string{usage.Name.Value}
}

// AzureUsageToPromDescs converts usage to prometheus descriptors.
func AzureUsageToPromDescs(
	target QueryTarget,
	usage AzureUsage,
) (*prometheus.Desc, *prometheus.Desc) {
	limitDesc := prometheus.NewDesc(
		prometheus.BuildFQName(
			target.ResourceProvider.ToPromNamespace(),
			PrometheusSubsystemUsage,
			PrometheusMetricLimit,
		),
		"Current limit for a resource usage.",
		[]string{PrometheusLabelUsageName},
		prometheus.Labels{
			PrometheusLabelSubscriptionID: target.SubscriptionID,
			PrometheusLabelLocation:       target.Location,
		},
	)
	currentValueDesc := prometheus.NewDesc(
		prometheus.BuildFQName(
			target.ResourceProvider.ToPromNamespace(),
			PrometheusSubsystemUsage,
			PrometheusMetricCurrentValue,
		),
		"Current value for a resource usage.",
		[]string{PrometheusLabelUsageName},
		prometheus.Labels{
			PrometheusLabelSubscriptionID: target.SubscriptionID,
			PrometheusLabelLocation:       target.Location,
		},
	)

	return limitDesc, currentValueDesc
}
