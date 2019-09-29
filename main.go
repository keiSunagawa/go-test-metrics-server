package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr = flag.String("listen-address", ":9999", "The address to listen on for HTTP requests.")
)

const (
	namespace = "test-metrics"
)

type weatherData struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity float64 `json:"humidity"`
	}
}

interface topCommand struct {
    GetMemoryUsagePer() int
    GetCPUUsagePer() int
}

struct topMock {}
func (t *topMock) GetMemoryUsagePer(){
    50
}

func (t *topMock) GetCPUUsagePer(){
    70
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
			Name:      "cpu_usage",
			Help:      "cpu usage percent form top command",
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
		c.temp.Desc(),
		c.memoryUsageG,
		c.topCommandInterface.GetMemoryUsagePer(),
	)
	ch <- prometheus.MustNewConstMetric(
		c.pressure.Desc(),
		c.cpuUsageG,
		c.topCommandInterface.GetUPUUsagePer(),
	)}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	c := newOpenWeatherMapCollector(topMock {})
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
