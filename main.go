package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/DazWilkin/gohome"
	"github.com/DazWilkin/home-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// GitCommit is the git commit value and is expected to be set during build
	GitCommit string
	// GoVersion is the Golang runtime version
	GoVersion = runtime.Version()
	// OSVersion is the OS version (uname --kernel-release) and is expected to be set during build
	OSVersion string
	// StartTime is the start time of the exporter represented as a UNIX epoch
	StartTime = time.Now().Unix()
)
var (
	device      = flag.String("device", "localhost:8008", "The endpoint of the Google Home device")
	endpoint    = flag.String("endpoint", ":9999", "The endpoint of the HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)

const (
	rootContent string = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Google Home Exporter</title>
  <style>body { font-family: sans-serif; }</style>
</head>
<body>
  <h2>Google Home Exporter</h2>
  <ul>
    <li><a href="{{ .MetricsPath }}">metrics</a></li>
	<li><a href="/healthz">healthz</a></li>
  </ul>
</body>
</html>
`
)

var (
	rootTemplate = template.Must(template.New("root").Parse(rootContent))
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Println(err)
	}
}
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if err := rootTemplate.Execute(w, struct {
		MetricsPath string
	}{
		MetricsPath: *metricsPath,
	}); err != nil {
		log.Println(err)
	}
}
func main() {
	flag.Parse()

	if GitCommit == "" {
		log.Println("[main] GitCommit value unchanged: expected to be set during build")
	}
	if OSVersion == "" {
		log.Println("[main] OSVersion value unchanged: expected to be set during build")
	}

	registry := prometheus.NewRegistry()

	client := gohome.NewClient(&http.Client{}, *device)

	log.Printf("[main] Exporting metrics for Google Home device (%s)", *device)
	registry.MustRegister(collector.NewExporterCollector(OSVersion, GoVersion, GitCommit, StartTime))
	registry.MustRegister(collector.NewDeviceInfoCollector(client))

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRoot))
	mux.Handle("/healthz", http.HandlerFunc(handleHealthz))
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Printf("[main] Server starting (%s)", *endpoint)
	log.Printf("[main] metrics served on: %s", *metricsPath)
	log.Fatal(http.ListenAndServe(*endpoint, mux))
}
