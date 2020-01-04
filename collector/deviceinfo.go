package collector

import (
	"log"
	"strconv"

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
	labelNames := []string{
		"name",
		"build_version",
		"version",
	}
	return &DeviceInfoCollector{
		client: client,
		NoiseLevel: prometheus.NewDesc(
			prometheus.BuildFQName("home", "device_info", "noise_level"),
			"TBD",
			labelNames,
			nil,
		),
		SignalLevel: prometheus.NewDesc(
			prometheus.BuildFQName("home", "device_info", "signal_level"),
			"TBD",
			labelNames,
			nil,
		),
		Uptime: prometheus.NewDesc(
			prometheus.BuildFQName("home", "device_info", "uptime"),
			"TBD",
			labelNames,
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
	labelValues := []string{
		info.Name,
		info.BuildVersion,
		strconv.Itoa(int(info.Version)),
	}
	ch <- prometheus.MustNewConstMetric(
		c.NoiseLevel,
		prometheus.GaugeValue,
		info.NoiseLevel,
		labelValues...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.SignalLevel,
		prometheus.GaugeValue,
		info.SignalLevel,
		labelValues...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.Uptime,
		prometheus.GaugeValue,
		info.Uptime,
		labelValues...,
	)

}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *DeviceInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.NoiseLevel
	ch <- c.SignalLevel
	ch <- c.Uptime
}
