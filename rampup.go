package hazana

import (
	"math"
	"time"

	"go.uber.org/ratelimit"
)

const defaultRampupStrategy = "exp2 keep=1 factor=2.0"

type rampupStrategy interface {
	execute(r *runner)
}

type linearIncreasingGoroutinesAndRequestsPerSecondStrategy struct{}

func (s linearIncreasingGoroutinesAndRequestsPerSecondStrategy) execute(r *runner) {
	r.spawnAttacker() // start at least one
	for i := 1; i <= r.config.RampupTimeSec; i++ {
		spawnAttackersToSize(r, i*r.config.MaxAttackers/r.config.RampupTimeSec)
		takeDuringRampupSeconds(r, i, 1)
		select {
		case <-r.abort:
			goto end
		default:
		}
	}
end: // aborted
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

// takeDuringRampupSecond puts all attackers to work during one second with a reduced RPS.
func takeDuringRampupSeconds(r *runner, second int, durationSeconds int) (int, *Metrics) {
	// collect metrics for each second
	rampMetrics := new(Metrics)
	// rampup can only proceed when at least one attacker is waiting for rps tokens
	if len(r.attackers) == 0 {
		Printf("no attackers available to start rampup or full attack")
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
	secondsAhead := time.Now().Add(time.Duration(time.Duration(durationSeconds) * time.Second))
	// put the attackers to work
	for time.Now().Before(secondsAhead) {
		limiter.Take()
		r.next <- true
	}
	limiter.Take() // to compensate for the first Take of the new limiter
	rampMetrics.updateLatencies()

	if r.config.Verbose {
		Printf("rate [%4f -> %v], mean response [%v], requests [%d], attackers [%d], success [%d %%]\n",
			rampMetrics.Rate, rps, rampMetrics.meanLogEntry(), rampMetrics.Requests, len(r.attackers), rampMetrics.successLogEntry())
	}
	return rps, rampMetrics
}

// exp2 keep=1 max-factor=2.0
type spawnAsWeNeedStrategy struct {
	parameters strategyParameters
}

func (s spawnAsWeNeedStrategy) execute(r *runner) {
	keep := s.parameters.intParam("keep", 1)
	maxFactor := s.parameters.floatParam("max-factor", 2.0)

	r.spawnAttacker() // start at least one
	for i := 1; i <= r.config.RampupTimeSec; i++ {
		select {
		case <-r.abort:
			goto end
		default:
		}
		targetRate, lastMetrics := takeDuringRampupSeconds(r, i, keep)
		currentRate := lastMetrics.Rate
		if currentRate < float64(targetRate) {
			currentFactor := float64(targetRate) / currentRate
			if currentFactor > maxFactor {
				currentFactor = maxFactor
			}
			spawnAttackersToSize(r, int(math.Ceil(float64(len(r.attackers))*currentFactor)))
		}
	}
end: // aborted
}
