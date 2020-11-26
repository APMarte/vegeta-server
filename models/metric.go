package models

import (
	"github.com/prometheus/client_golang/prometheus"
)

type AttackBaseInfo struct {
	ID string `json:"id,omitempty"`
	// Params captures the attack parameters
	Params AttackParams `json:"params,omitempty"`
}

type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
}

var ReqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total",
	Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqRt = &Metric{
	ID:          "reqRt",
	Name:        "requests_rate",
	Description: "Number request rate sustained during the attack period",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqDur = &Metric{
	ID:          "reqDur",
	Name:        "request_duration_total",
	Description: "Time taken in the attack.",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqAttck = &Metric{
	ID:          "reqAttck",
	Name:        "request_duration_attack",
	Description: "Time taken issuing all requests",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqWait = &Metric{
	ID:          "reqWait",
	Name:        "request_duration_wait",
	Description: "Time taken issuing all requests",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqLatMean = &Metric{
	ID:          "reqLatMean",
	Name:        "request_latencies_mean",
	Description: "Average of the latencies of all requests in an attack.",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqLat50th = &Metric{
	ID:          "reqLat50th",
	Name:        "request_latencies_50thpercentile",
	Description: "50th percentile of all requests in an attack.",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqLat95th = &Metric{
	ID:          "reqLat95th",
	Name:        "request_latencies_mean_95thpercentile",
	Description: "95th percentile of all requests in an attack.",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqLat99th = &Metric{
	ID:          "reqLat99th",
	Name:        "request_latencies_mean_99thpercentile",
	Description: "99th percentile of all requests in an attack.",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqLatMax = &Metric{
	ID:          "reqLatMax",
	Name:        "request_latencies_max",
	Description: "Maximum latency of all requests in an attack.",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ResSuccessRatio = &Metric{
	ID:          "resSuccessRatio",
	Name:        "response_success_ratio",
	Description: "The percentage of requests whose responses didn't error",
	Type:        "gauge_vec",
	Args:        []string{"id", "rate", "duration"},
}

var ReqStsCode = &Metric{
	ID:   "reqStsCode",
	Name: "request_status_code",
	Type: "gauge_vec",
	Args: []string{"id", "rate", "duration", "code"},
}

var Histogram = &Metric{
	ID:   "Histogram",
	Name: "request_duration_histogram",
	Type: "histogram_vec",
	Args: []string{"id"},
}

var StandardMetrics = []*Metric{
	ReqCnt,
	ReqDur,
	ReqRt,
	ReqAttck,
	ReqWait,
	ReqLatMean,
	ReqLat50th,
	ReqLat95th,
	ReqLat99th,
	ReqLatMax,
	ResSuccessRatio,
	ReqStsCode,
	Histogram,
}

// NewMetric associates prometheus.Collector based on Metric.Type
func NewMetric(m *Metric, subsystem string) prometheus.Collector {
	var metric prometheus.Collector
	switch m.Type {
	case "counter_vec":
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "counter":
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "gauge_vec":
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "gauge":
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "histogram_vec":
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
				Buckets:   []float64{0, 20, 50, 100, 500, 1000},
			},
			m.Args,
		)
	case "histogram":
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "summary_vec":
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "summary":
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	}
	return metric
}
