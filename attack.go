package hazana

import (
	"time"
)

// Attack must be implemented by a service client.
type Attack interface {
	// Setup should establish the connection to the service
	// It may want to access the config of the runner.
	Setup(c Config) error
	// Do performs one request and is executed in one fixed goroutine.
	// Returns the index of request used and error received.
	Do() (requestIndex int, err error)
	// Teardown should close the connection of the service
	TearDown() error
	// Clone should return a new fresh Attack
	Clone() Attack
}

type result struct {
	begin, end time.Time
	request    int // index in list from requests
	elapsed    time.Duration
	err        error
}

// attack calls attacker.Do upon each received next token, forever
// attack aborts the loop on a quit receive
// attack sends a result on the results channel after each call.
func attack(attacker Attack, next, quit <-chan bool, results chan<- result) {
	for {
		select {
		case <-next:
			begin := time.Now()
			ri, err := attacker.Do()
			end := time.Now()
			results <- result{
				request: ri,
				begin:   begin,
				end:     end,
				elapsed: end.Sub(begin),
				err:     err,
			}
		case <-quit:
			return
		}
	}
}
