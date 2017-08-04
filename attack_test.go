package hazana

import (
	"testing"
	"time"
)

func TestAttack(t *testing.T) {
	attacker := new(attackMock)
	dur := 10 * time.Millisecond
	attacker.sleep = dur
	next := make(chan bool)
	quit := make(chan bool)
	results := make(chan result)

	go attack(attacker, next, quit, results)

	next <- true
	r := <-results
	quit <- true
	if got, want := r.doResult.Error, error(nil); got != want {
		t.Fatalf("got %v want %v", got, want)
	}
	if got, want := r.elapsed, dur; got < want {
		t.Fatalf("got %v want >= %v", got, want)
	}
}