package hazana

import (
	"log"
	"math"
	"time"

	"go.uber.org/ratelimit"
)

const defaultRampupStrategy = "exp2"

type rampupStrategy interface {
	execute(r *runner)
}

type linearIncreasingGoroutinesAndRequestsPerSecondStrategy struct{}

func (s linearIncreasingGoroutinesAndRequestsPerSecondStrategy) execute(r *runner) {
	r.spawnAttacker() // start at least one
	for i := 1; i <= r.config.RampupTimeSec; i++ {
		spawnAttackersToSize(r, i*r.config.MaxAttackers/r.config.RampupTimeSec)
		takeDuringOneRampupSecond(r, i)
	}
}

func spawnAttackersToSize(r *runner, count int) {
	routines := count
	if count > r.config.MaxAttackers {
		routines = r.config.MaxAttackers
	}
	// spawn extra goroutines
	for s := len(r.attackers); s < routines; s++ {
		r.spawnAttacker()
	}
}

// takeDuringOneRampupSecond puts all attackers to work during one second with a reduced RPS.
func takeDuringOneRampupSecond(r *runner, second int) (int, *Metrics) {
	// collect metrics for each second
	rampMetrics := new(Metrics)
	// rampup can only proceed when at least one attacker is waiting for rps tokens
	if len(r.attackers) == 0 {
		log.Println("no attackers available to start rampup or full attack")
		return 0, rampMetrics
	}
	// change pipeline function to collect local metrics
	r.resultsPipeline = func(rs result) result {
		rampMetrics.add(rs)
		return rs
	}
	// for each second start a new reduced rate limiter
	rps := second * r.config.RPS / r.config.RampupTimeSec
	if rps == 0 { // minimal 1
		rps = 1
	}
	limiter := ratelimit.New(rps)
	oneSecondAhead := time.Now().Add(time.Duration(1 * time.Second))
	// put the attackers to work
	for time.Now().Before(oneSecondAhead) {
		limiter.Take()
		r.next <- true
	}
	limiter.Take() // to compensate for the first Take of the new limiter
	rampMetrics.updateLatencies()

	if r.config.Verbose {
		log.Printf("rate [%4f -> %v], mean response [%v], # requests [%d], # attackers [%d], %% success [%d]\n",
			rampMetrics.Rate, rps, rampMetrics.meanLogEntry(), rampMetrics.Requests, len(r.attackers), rampMetrics.successLogEntry())
	}
	return rps, rampMetrics
}

type spawnAsWeNeedStrategy struct{}

func (s spawnAsWeNeedStrategy) execute(r *runner) {
	r.spawnAttacker() // start at least one
	for i := 1; i <= r.config.RampupTimeSec; i++ {
		targetRate, lastMetrics := takeDuringOneRampupSecond(r, i)
		currentRate := lastMetrics.Rate
		if currentRate < float64(targetRate) {
			factor := float64(targetRate) / currentRate
			if factor > 2.0 {
				factor = 2.0
			}
			spawnAttackersToSize(r, int(math.Ceil(float64(len(r.attackers))*factor)))
		}
	}
}
