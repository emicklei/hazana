package hazana

import "time"

type result struct {
	begin, end time.Time
	request    int // index in list from requests
	elapsed    time.Duration
	err        error
}
