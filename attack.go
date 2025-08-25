package hazana

import (
	"context"
	"errors"
	"time"
)

// Attack must be implemented by a service client.
type Attack interface {
	// Setup should establish the connection to the service
	// It may want to access the config of the runner.
	Setup(c Config) error
	// Do performs one request.
	// The context is used to cancel the request on timeout.
	// This method must honor the context cancellation.
	Do(ctx context.Context) DoResult
	// Teardown can be used to close the connection to the service
	Teardown() error
	// Clone should return a fresh new Attack
	// Make sure the new Attack has values for shared struct fields initialized at Setup.
	Clone() Attack
}

var errAttackDoTimedOut = errors.New("hazana.Attack Do(ctx) timed out")

// attack calls attacker.Do upon each received next token, forever
// attack aborts the loop on a quit receive
// attack sends a result on the results channel after each call.
func attack(attacker Attack, next, quit <-chan bool, results chan<- result, timeout time.Duration) {
	for {
		select {
		case <-next:
			begin := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)

			// call Do and block
			dor := attacker.Do(ctx)

			// call cancel to avoid context leak
			cancel()

			if dor.Error == context.DeadlineExceeded {
				dor.Error = errAttackDoTimedOut
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
