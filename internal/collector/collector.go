package collector

import (
	"context"
	"sync"
	"time"

	azureusageprom "github.com/b4fun/azure-usage-prom"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusCollector implements the prometheus collector.
type PrometheusCollector struct {
	usageLister azureusageprom.AzureUsageLister
	queryTarget azureusageprom.QueryTarget

	descriptorsInit         *sync.Once
	descriptorsLimit        map[string]*prometheus.Desc
	descriptorsCurrentValue map[string]*prometheus.Desc
}

// New creates a prometheus collector.
func New(
	usageLister azureusageprom.AzureUsageLister,
	queryTarget azureusageprom.QueryTarget,
) prometheus.Collector {
	return &PrometheusCollector{
		usageLister: usageLister,
		queryTarget: queryTarget,

		descriptorsInit:         new(sync.Once),
		descriptorsLimit:        map[string]*prometheus.Desc{},
		descriptorsCurrentValue: map[string]*prometheus.Desc{},
	}
}

var _ prometheus.Collector = (*PrometheusCollector)(nil)

func (p *PrometheusCollector) listResourceProviderUsage() (azureusageprom.AzureUsageList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	usageList, err := p.usageLister.ListByResourceProvider(ctx, p.queryTarget)
	if err != nil {
		return azureusageprom.AzureUsageList{}, nil
	}

	p.descriptorsInit.Do(func() {
		p.initDescriptors(usageList)
	})

	return usageList, nil
}

func (p *PrometheusCollector) initDescriptors(usageList azureusageprom.AzureUsageList) {
	for _, usage := range usageList.Value {
		limit, currentValue := azureusageprom.AzureUsageToPromDescs(p.queryTarget, usage)
		p.descriptorsLimit[usage.ID] = limit
		p.descriptorsCurrentValue[usage.ID] = currentValue
	}
}

// Collect - implements prometheus.Collector.
func (p *PrometheusCollector) Collect(ch chan<- prometheus.Metric) {
	usageList, err := p.listResourceProviderUsage()
	if err != nil {
		glog.Errorf("listResourceProviderUsage: %s", err)
		return
	}

	for _, usage := range usageList.Value {
		descLimit, ok := p.descriptorsLimit[usage.ID]
		if !ok {
			continue
		}
		descCurrentValue, ok := p.descriptorsCurrentValue[usage.ID]
		if !ok {
			continue
		}

		labels := azureusageprom.AzureUsageToMetricLabels(p.queryTarget, usage)
		ch <- prometheus.MustNewConstMetric(
			descLimit, prometheus.GaugeValue, float64(usage.Limit),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			descCurrentValue, prometheus.GaugeValue, float64(usage.CurrentValue),
			labels...,
		)
	}
}

// Describe - implements prometheus.Describe.
func (p *PrometheusCollector) Describe(ch chan<- *prometheus.Desc) {
	if _, err := p.listResourceProviderUsage(); err != nil {
		glog.Errorf("listResourceProviderUsage: %s", err)
		return
	}

	for _, desc := range p.descriptorsLimit {
		ch <- desc
	}
	for _, desc := range p.descriptorsCurrentValue {
		ch <- desc
	}
}
