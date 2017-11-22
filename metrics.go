package hazana

import (
	"strconv"
	"time"

	"github.com/streadway/quantile"
)

// this file is a modified version from https://github.com/tsenart/vegeta/blob/master/lib/metrics.go

type (
	// Metrics holds metrics computed out of a slice of Results which are used
	// in some of the Reporters
	Metrics struct {
		// Latencies holds computed request latency metrics.
		Latencies LatencyMetrics `json:"latencies"`
		// First is the earliest timestamp in a Result set.
		Earliest time.Time `json:"earliest"`
		// Latest is the latest timestamp in a Result set.
		Latest time.Time `json:"latest"`
		// End is the latest timestamp in a Result set plus its latency.
		End time.Time `json:"end"`
		// Duration is the duration of the attack.
		Duration time.Duration `json:"duration"`
		// Wait is the extra time waiting for responses from targets.
		Wait time.Duration `json:"wait"`
		// Requests is the total number of requests executed.
		Requests uint64 `json:"requests"`
		// Rate is the rate of requests per second.
		Rate float64 `json:"rate"`
		// Success is the percentage of non-error responses.
		Success float64 `json:"success"`
		// StatusCodes is a histogram of the responses' status codes.
		StatusCodes map[string]int `json:"status_codes"`
		// Errors is a set of unique errors returned by the targets during the attack.
		Errors []string `json:"errors"`

		errors    map[string]struct{}
		success   uint64
		latencies *quantile.Estimator
	}

	// LatencyMetrics holds computed request latency metrics.
	LatencyMetrics struct {
		// Total is the total latency sum of all requests in an attack.
		Total time.Duration `json:"total"`
		// Mean is the mean request latency.
		Mean time.Duration `json:"mean"`
		// P50 is the 50th percentile request latency.
		P50 time.Duration `json:"50th"`
		// P95 is the 95th percentile request latency.
		P95 time.Duration `json:"95th"`
		// P99 is the 99th percentile request latency.
		P99 time.Duration `json:"99th"`
		// Max is the maximum observed request latency.
		Max time.Duration `json:"max"`
	}
)

func (m Metrics) successLogEntry() int {
	s := int(m.Success * 100.0)
	if s < 0 {
		return 0
	}
	return s
}

func (m Metrics) meanLogEntry() time.Duration {
	lm := m.Latencies.Mean
	if lm < 0 {
		return time.Duration(0)
	}
	return time.Duration(lm)
}

func (m *Metrics) add(r result) {
	m.init()

	m.Requests++
	// StatusCode is optional
	if r.doResult.StatusCode > 0 {
		m.StatusCodes[strconv.Itoa(r.doResult.StatusCode)]++
	}
	m.Latencies.Total += r.elapsed

	m.latencies.Add(float64(r.elapsed))

	if m.Earliest.IsZero() || m.Earliest.After(r.begin) {
		m.Earliest = r.begin
	}

	if r.begin.After(m.Latest) {
		m.Latest = r.begin
	}

	if end := r.end; end.After(m.End) {
		m.End = end
	}

	if r.elapsed > m.Latencies.Max {
		m.Latencies.Max = r.elapsed
	}

	if r.doResult.Error != nil {
		if _, ok := m.errors[r.doResult.Error.Error()]; !ok {
			m.errors[r.doResult.Error.Error()] = struct{}{}
			m.Errors = append(m.Errors, r.doResult.Error.Error())
		}
	} else {
		if r.doResult.StatusCode == 0 || (r.doResult.StatusCode >= 200 && r.doResult.StatusCode < 400) {
			m.success++
		}
	}
}

// updateLatencies computes derived summary metrics which don't need to be run on every add call.
func (m *Metrics) updateLatencies() {
	m.init()
	fRequests := float64(m.Requests)
	m.Duration = m.Latest.Sub(m.Earliest)
	if secs := m.Duration.Seconds(); secs > 0 {
		m.Rate = fRequests / secs
	}
	m.Wait = m.End.Sub(m.Latest)
	m.Success = float64(m.success) / fRequests
	m.Latencies.Mean = time.Duration(float64(m.Latencies.Total) / fRequests)
	m.Latencies.P50 = time.Duration(m.latencies.Get(0.50))
	m.Latencies.P95 = time.Duration(m.latencies.Get(0.95))
	m.Latencies.P99 = time.Duration(m.latencies.Get(0.99))
}

func (m *Metrics) init() {
	if m.latencies == nil {
		m.StatusCodes = map[string]int{}
		m.errors = map[string]struct{}{}
		m.latencies = quantile.New(
			quantile.Known(0.50, 0.01),
			quantile.Known(0.95, 0.001),
			quantile.Known(0.99, 0.0005),
		)
	}
}
