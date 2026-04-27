package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	UUIDRegExp = regexp.MustCompile(`[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)
)

const (
	ServiceName = "ServiceName"
	URL         = "Url"
	Method      = "Method"
	StatusCode  = "StatusCode"
	RequestID   = "RequestID"
)

type writer struct {
	http.ResponseWriter
	statusCode int
}

func NewWriter(w http.ResponseWriter) *writer {
	return &writer{w, http.StatusOK}
}

func (w *writer) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

type MetricsMiddleware struct {
	metric          *prometheus.GaugeVec
	durations       *prometheus.HistogramVec
	errors          *prometheus.CounterVec
	durationNew     *prometheus.SummaryVec
	name            string
	cpuUsage        prometheus.Gauge
	memoryUsage     prometheus.Gauge
	diskUsage       *prometheus.GaugeVec
	diskReadBytes   prometheus.Gauge
	diskWriteBytes  prometheus.Gauge
	collectorTicker *time.Ticker
}

func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{}
}

func (m *MetricsMiddleware) Register(name string) {
	m.name = name

	labels := []string{ServiceName, URL, Method, StatusCode}

	m.metric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name + "_requests_total",
			Help: fmt.Sprintf("Total requests for service %s", name),
		},
		labels,
	)

	m.durations = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name + "_duration_seconds",
			Help:    "Request duration distribution.",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms -> ~16s
		},
		labels,
	)

	m.errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_errors_total",
			Help: "Counter of errors.",
		},
		labels,
	)

	m.durationNew = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       name + "_duration_summary_seconds",
			Help:       "Summary of request durations.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		labels,
	)

	m.cpuUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name + "_cpu_usage_percent",
			Help: "Current CPU usage in percent",
		},
	)

	m.memoryUsage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name + "_memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
	)

	m.diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name + "_disk_usage_percent",
			Help: "Disk usage in percent by mount point",
		},
		[]string{"mount"},
	)

	m.diskReadBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name + "_disk_read_bytes_total",
			Help: "Total bytes read from disk",
		},
	)

	m.diskWriteBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name + "_disk_write_bytes_total",
			Help: "Total bytes written to disk",
		},
	)

	prometheus.MustRegister(m.metric)
	prometheus.MustRegister(m.durations)
	prometheus.MustRegister(m.errors)
	prometheus.MustRegister(m.durationNew)
	prometheus.MustRegister(m.cpuUsage)
	prometheus.MustRegister(m.memoryUsage)
	prometheus.MustRegister(m.diskUsage)
	prometheus.MustRegister(m.diskReadBytes)
	prometheus.MustRegister(m.diskWriteBytes)

	m.collectorTicker = time.NewTicker(10 * time.Second)
	go m.collectSystemMetrics()
}

func (m *MetricsMiddleware) collectSystemMetrics() {
	for range m.collectorTicker.C {
		if cpuPercent, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercent) > 0 {
			m.cpuUsage.Set(cpuPercent[0])
		}

		if memInfo, err := mem.VirtualMemory(); err == nil {
			m.memoryUsage.Set(float64(memInfo.Used))
		}

		if partitions, err := disk.Partitions(false); err == nil {
			for _, partition := range partitions {
				if usage, err := disk.Usage(partition.Mountpoint); err == nil {
					m.diskUsage.WithLabelValues(partition.Mountpoint).Set(usage.UsedPercent)
				}
			}
		}

		if ioCounters, err := disk.IOCounters(); err == nil {
			for _, counter := range ioCounters {
				m.diskReadBytes.Set(float64(counter.ReadBytes))
				m.diskWriteBytes.Set(float64(counter.WriteBytes))
			}
		}
	}
}

func (m *MetricsMiddleware) Close() {
	if m.collectorTicker != nil {
		m.collectorTicker.Stop()
	}
}
func getRoutePattern(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if route == nil {
		return r.URL.Path
	}

	pathTemplate, err := route.GetPathTemplate()
	if err != nil {
		return r.URL.Path
	}

	return pathTemplate
}

func (m *MetricsMiddleware) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := NewWriter(w)

		next.ServeHTTP(wrapper, r)

		tm := time.Since(start)

		routePattern := getRoutePattern(r)
		fmt.Println(routePattern)

		labels := prometheus.Labels{
			ServiceName: m.name,
			URL:         routePattern,
			Method:      r.Method,
			StatusCode:  fmt.Sprintf("%d", wrapper.statusCode),
		}

		m.metric.With(labels).Inc()
		m.durations.With(labels).Observe(tm.Seconds())
		m.durationNew.With(labels).Observe(tm.Seconds())

		if wrapper.statusCode >= 500 {
			m.errors.With(labels).Inc()
		}
	})
}
