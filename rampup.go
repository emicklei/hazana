package hazana

import (
	"log"
	"time"

	"go.uber.org/ratelimit"
)

type rampupStrategy interface {
	execute(r *runner)
}

type linearIncreasingGoroutinesAndRequestsPerSecondStrategy struct{}

func (s linearIncreasingGoroutinesAndRequestsPerSecondStrategy) execute(r *runner) {
	var rampMetrics *Metrics
	r.spawnAttacker() // start at least one
	for i := 1; i <= r.config.RampupTimeSec; i++ {
		routines := i * r.config.MaxAttackers / r.config.RampupTimeSec
		// spawn extra goroutines
		for s := len(r.attackers); s < routines; s++ {
			r.spawnAttacker()
		}
		// collect metrics for each second
		rampMetrics = new(Metrics)
		// change pipeline function to collect local metrics
		r.resultsPipeline = func(rs result) result {
			rampMetrics.add(rs)
			return rs
		}
		// for each second start a new reduced rate limiter
		rps := i * r.config.RPS / r.config.RampupTimeSec
		if rps == 0 { // minimal 1
			rps = 1
		}
		limiter := ratelimit.New(rps) // per second
		oneSecond := time.Now().Add(time.Duration(1 * time.Second))
		for time.Now().Before(oneSecond) {
			limiter.Take()
			r.next <- true
		}
		limiter.Take() // to compensate for the first Take of the new limiter
		rampMetrics.updateLatencies()

		if r.config.Verbose {
			log.Printf("current rate [%v], target rate [%v], attackers [%v], mean response time [%v]\n", rampMetrics.Rate, rps, len(r.attackers), time.Duration(rampMetrics.Latencies.Mean))
		}
	}
}

type spawnAsWeNeedStrategy struct{}

// TODO
func (s spawnAsWeNeedStrategy) execute(r *runner) {}
