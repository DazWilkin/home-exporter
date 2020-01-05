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

	Up          *prometheus.Desc
	BuildInfo   *prometheus.Desc
	NoiseLevel  *prometheus.Desc
	SignalLevel *prometheus.Desc
	Uptime      *prometheus.Desc
}

// NewDeviceInfoCollector returns a new DeviceInfoCollector.
func NewDeviceInfoCollector(client *gohome.Client) *DeviceInfoCollector {
	subsystem := "device_info"
	labelNames := []string{
		"name",
		"build_version",
		"version",
	}
	return &DeviceInfoCollector{
		client: client,
		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "up"),
			"If 1 the Home device is accessible, 0 otherwise",
			labelNames,
			nil,
		),
		BuildInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "build_info"),
			"A metric with a constant '1' value labeled by X of the device",
			[]string{"build_version", "cast_build_revision", "release_track"},
			nil,
		),
		NoiseLevel: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "noise_level"),
			"Noise Level dB",
			labelNames,
			nil,
		),
		SignalLevel: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "signal_level"),
			"Signal Level dB",
			labelNames,
			nil,
		),
		Uptime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "uptime"),
			"Uptime in seconds",
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
	ch <- prometheus.MustNewConstMetric(
		c.BuildInfo,
		prometheus.CounterValue,
		1.0,
		info.BuildVersion, info.CastBuildRevision, info.ReleaseTrack)
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
