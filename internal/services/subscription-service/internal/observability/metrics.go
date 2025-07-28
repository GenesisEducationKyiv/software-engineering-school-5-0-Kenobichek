package observability

import (
    "net/http"
    "os"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    MessageProcessedTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "subscription_messages_processed_total",
            Help: "Total number of Kafka messages processed by the subscription service.",
        },
        []string{"topic", "status"},
    )

    MessageProcessingDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "subscription_message_processing_seconds",
            Help:    "Processing duration distribution for subscription messages.",
            Buckets: prometheus.DefBuckets,
        },
        []string{"topic"},
    )

    WeatherJobRuns = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "subscription_weather_job_runs_total",
            Help: "Number of executions of the weather update job.",
        },
    )
)

func init() {
    // Register metrics with default registry.
    prometheus.MustRegister(MessageProcessedTotal, MessageProcessingDuration, WeatherJobRuns)
}

// StartMetricsServer launches an HTTP server exposing /metrics. Returns immediately.
func StartMetricsServer() {
    addr := os.Getenv("METRICS_ADDR")
    if addr == "" {
        addr = ":2112"
    }

    mux := http.NewServeMux()
    mux.Handle("/metrics", promhttp.Handler())

    go func() {
        Logger.Infof("metrics server listening on %s", addr)
        if err := http.ListenAndServe(addr, mux); err != nil {
            Logger.Errorf("metrics server exited: %v", err)
        }
    }()
}

// ObserveDuration is a helper to measure time since start.
func ObserveDuration(hist *prometheus.HistogramVec, start time.Time, labels prometheus.Labels) {
    if hist == nil {
        return
    }
    d := time.Since(start).Seconds()
    hist.With(labels).Observe(d)
}
