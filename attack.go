package hazana

import (
	"errors"
	"time"
)

// Attack must be implemented by a service client.
type Attack interface {
	// Setup should establish the connection to the service
	// It may want to access the config of the runner.
	Setup(c Config) error
	// Do performs one request and is executed in one fixed goroutine.
	Do() DoResult
	// Teardown should close the connection of the service
	Teardown() error
	// Clone should return a new fresh Attack
	Clone() Attack
}

var errAttackDoTimedOut = errors.New("Attack Do() timedout")

// attack calls attacker.Do upon each received next token, forever
// attack aborts the loop on a quit receive
// attack sends a result on the results channel after each call.
func attack(attacker Attack, next, quit <-chan bool, results chan<- result, timeout time.Duration) {
	for {
		select {
		case <-next:
			begin := time.Now()
			done := make(chan DoResult)
			go func() {
				done <- attacker.Do()
			}()
			var dor DoResult
			// either get the result from the attacker or from the timeout
			select {
			case <-time.After(timeout):
				dor = DoResult{Error: errAttackDoTimedOut}
			case dor = <-done:
			}
			end := time.Now()
			results <- result{
				doResult: dor,
				begin:    begin,
				end:      end,
				elapsed:  end.Sub(begin),
			}
		case <-quit:
			return
		}
	}
}
