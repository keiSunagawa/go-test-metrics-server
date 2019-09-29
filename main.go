package main

import (
	"flag"
	"log"
	"net/http"
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":9999", "The address to listen on for HTTP requests.")
)

const (
	namespace = "test_metrics"
)

type topCommand interface {
    GetMemoryUsagePer() float64
    GetCPUUsagePer() float64
}


type topMock struct {}

func (t topMock) GetMemoryUsagePer() float64 {
    return float64(rand.Intn(100))
}

func (t topMock) GetCPUUsagePer() float64 {
    return float64(rand.Intn(100))
}

type topCommandCollector struct {
	topCommandInterface topCommand
	memoryUsageG prometheus.Gauge
	cpuUsageG prometheus.Gauge
}

func newTopCommandCollector(cmd topCommand) *topCommandCollector {
	return &topCommandCollector{
		topCommandInterface: cmd,
		memoryUsageG: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_usage",
			Help:      "memory usage percent form top command",
		}),
		cpuUsageG: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cpu_usage",
			Help:      "cpu usage percent form top command",
		}),
	}
}

func (c *topCommandCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.memoryUsageG .Desc()
	ch <- c.cpuUsageG.Desc()
}

func (c *topCommandCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.memoryUsageG.Desc(),
		prometheus.GaugeValue,
		c.topCommandInterface.GetMemoryUsagePer(),
	)
	ch <- prometheus.MustNewConstMetric(
		c.cpuUsageG.Desc(),
		prometheus.GaugeValue,
		c.topCommandInterface.GetCPUUsagePer(),
	)}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	c := newTopCommandCollector(topMock {})
	registry := prometheus.NewRegistry()
	registry.MustRegister(c)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	flag.Parse()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metricsHandler(w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
