package models

import "time"

// JSONReportResponse provides the model for a report response object
type JSONReportResponse struct {
	ID        string `json:"id"`
	Latencies struct {
		Total int `json:"total"`
		Mean  int `json:"mean"`
		Max   int `json:"max"`
		P50th int `json:"50th"`
		P95th int `json:"95th"`
		P99th int `json:"99th"`
	} `json:"latencies"`
	BytesIn struct {
		Total int     `json:"total"`
		Mean  float64 `json:"mean"`
	} `json:"bytes_in"`
	BytesOut struct {
		Total int     `json:"total"`
		Mean  float64 `json:"mean"`
	} `json:"bytes_out"`
	Earliest    string         `json:"earliest"`
	Latest      string         `json:"latest"`
	End         string         `json:"end"`
	Duration    int            `json:"duration"`
	Wait        int            `json:"wait"`
	Requests    int            `json:"requests"`
	Rate        float64        `json:"rate"`
	Success     float64        `json:"success"`
	StatusCodes map[string]int `json:"status_codes"`
	Errors      []string       `json:"errors"`
}

type SeqResult struct {
	Attack    string        `json:"attack"`
	Seq       uint64        `json:"seq"`
	Code      uint16        `json:"code"`
	Timestamp time.Time     `json:"timestamp"`
	Latency   time.Duration `json:"latency"`
	BytesOut  uint64        `json:"bytes_out"`
	BytesIn   uint64        `json:"bytes_in"`
	Error     string        `json:"error"`
	Body      []byte        `json:"body"`
}
