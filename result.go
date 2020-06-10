package hazana

import (
	"strings"
	"time"
)

type result struct {
	begin, end time.Time
	elapsed    time.Duration
	doResult   DoResult
}

// DoResult is the return value of a Do call on an Attack.
type DoResult struct {
	// Label identifying the request that was send which is only used for reporting the metrics.
	// Use the RequesLabelSeparator to have multiple labels.
	RequestLabel string
	// The error that happened when sending the request or receiving the response.
	Error error
	// The HTTP status code.
	StatusCode int
	// Number of bytes transferred when sending the request.
	BytesIn int64
	// Number of bytes transferred when receiving the response.
	BytesOut int64
}

// RequesLabelSeparator is what it says
var RequesLabelSeparator = ","

// WithLabels return a copy with multiple Request labels set.
func (d DoResult) WithLabels(labels ...string) DoResult {
	d.RequestLabel = strings.Join(labels, RequesLabelSeparator)
	return d
}
