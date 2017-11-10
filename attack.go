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
	// Do performs one request and is executed in a separate goroutine.
	// The context is used to cancel the request on timeout.
	Do(ctx context.Context) DoResult
	// Teardown can be used to close the connection to the service
	Teardown() error
	// Clone should return a fresh new Attack
	// Make sure the new Attack has values for shared struct fields initialized at Setup.
	Clone() Attack
}

var errAttackDoTimedOut = errors.New("Attack Do(ctx) timedout")

// attack calls attacker.Do upon each received next token, forever
// attack aborts the loop on a quit receive
// attack sends a result on the results channel after each call.
func attack(attacker Attack, next, quit <-chan bool, results chan<- result, timeout time.Duration) {
	for {
		select {
		case <-next:
			begin := time.Now()
			done := make(chan DoResult)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			go func() {
				done <- attacker.Do(ctx)
			}()
			var dor DoResult
			// either get the result from the attacker or from the timeout
			select {
			case <-ctx.Done():
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
