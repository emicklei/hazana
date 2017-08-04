package hazana

import "time"

type result struct {
	begin, end time.Time
	elapsed    time.Duration
	doResult   DoResult
}

// DoResult is the return value of a Do call on an Attack.
type DoResult struct {
	// Index in the list of request samples, 0 if not applicable
	RequestIndex int
	// The error that happened when sending the request or receiving the response.
	Error error
	// The HTTP status code.
	StatusCode int
	// Number of bytes transferred when sending the request.
	BytesIn int64
	// Number of bytes transferred when receiving the response.
	BytesOut int64
}
