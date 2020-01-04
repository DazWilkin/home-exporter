package collector

import (
	"log"

	"github.com/DazWilkin/gohome"
	"github.com/prometheus/client_golang/prometheus"
)

// DeviceInfoCollector collects metrics, mostly runtime, about this exporter in general.
type DeviceInfoCollector struct {
	client *gohome.Client

	NoiseLevel  *prometheus.Desc
	SignalLevel *prometheus.Desc
	Uptime      *prometheus.Desc
}

// NewDeviceInfoCollector returns a new DeviceInfoCollector.
func NewDeviceInfoCollector(client *gohome.Client) *DeviceInfoCollector {
	return &DeviceInfoCollector{
		client: client,
		NoiseLevel: prometheus.NewDesc(
			prometheus.BuildFQName("home", "device_info", "noise_level"),
			"TBD",
			nil,
			nil,
		),
		SignalLevel: prometheus.NewDesc(
			prometheus.BuildFQName("home", "device_info", "signal_level"),
			"TBD",
			nil,
			nil,
		),
		Uptime: prometheus.NewDesc(
			prometheus.BuildFQName("home", "device_info", "uptime"),
			"TBD",
			nil,
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *DeviceInfoCollector) Collect(ch chan<- prometheus.Metric) {
	info, err := c.client.DeviceInfo()
	if err != nil {
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.NoiseLevel,
		prometheus.GaugeValue,
		info.NoiseLevel,
		[]string{}...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.SignalLevel,
		prometheus.GaugeValue,
		info.SignalLevel,
		[]string{}...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.Uptime,
		prometheus.GaugeValue,
		info.Uptime,
		[]string{}...,
	)

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *DeviceInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.NoiseLevel
	ch <- c.SignalLevel
	ch <- c.Uptime
}
