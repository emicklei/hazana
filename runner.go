package hazana

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"go.uber.org/ratelimit"
)

// BeforeRunner can be implemented by an Attacker
// and its method is called before a test or run.
type BeforeRunner interface {
	BeforeRun(c Config) error
}

// AfterRunner can be implemented by an Attacker
// and its method is called after a test or run.
// The report is passed to compute the Failed field and/or store values in Output.
type AfterRunner interface {
	AfterRun(r *RunReport) error
}

type runner struct {
	config            Config
	attackers         []Attack
	next, quit, abort chan bool
	results           chan result
	prototype         Attack
	metrics           map[string]*Metrics
	resultsPipeline   func(r result) result
}

// Run starts attacking a service using an Attack implementation and a configuration.
// Return a report with statistics per sample and the configuration used.
func Run(a Attack, c Config) *RunReport {
	if c.Verbose {
		Printf("*** Hazana load runner ready to attack ***\n")
		Printf("%v", c)
		Printf("[%d] available logical CPUs\n", runtime.NumCPU())
	}
	r := new(runner)
	r.config = c
	r.prototype = a

	// validate the configuration
	if msg := c.Validate(); len(msg) > 0 {
		for _, each := range msg {
			fmt.Println("[hazana] - a configuration error was found", each)
		}
		fmt.Println()
		flag.Usage()
		os.Exit(0)
	}
	r.init()

	// is the attacker interested in the run lifecycle?
	if lifecycler, ok := a.(BeforeRunner); ok {
		if err := lifecycler.BeforeRun(c); err != nil {
			Printf("BeforeRun failed:%v\n", err)
		}
	}

	// do a test if the flag says so
	if *oSample > 0 {
		report := r.test(*oSample)
		if lifecycler, ok := a.(AfterRunner); ok {
			if err := lifecycler.AfterRun(report); err != nil {
				Printf("AfterRun failed:%v\n", err)
			}
		}
		return report
	}
	report := r.run()
	if lifecycler, ok := a.(AfterRunner); ok {
		if err := lifecycler.AfterRun(report); err != nil {
			Printf("AfterRun failed:%v\n", err)
		}
	}
	return report
}

func (r *runner) listenForAbort() {
	r.abort = make(chan bool)
	// abort will stop test,rampup or full attack
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		sig := <-ch
		Printf("caught signal %v. aborting run...\n", sig)
		r.abort <- true
		Printf("aborted. clean up...\n")
	}()
}

func (r *runner) init() {
	r.next = make(chan bool)
	r.quit = make(chan bool)
	r.results = make(chan result)
	r.attackers = []Attack{}
	r.metrics = map[string]*Metrics{}
	r.resultsPipeline = r.addResult
	r.listenForAbort()
}

func (r *runner) spawnAttacker() {
	if r.config.Verbose {
		Printf("setup and spawn new attacker [%d]\n", len(r.attackers)+1)
	}
	attacker := r.prototype.Clone()
	if err := attacker.Setup(r.config); err != nil {
		Printf("attacker [%d] setup failed with [%v]\n", len(r.attackers)+1, err)
		return
	}
	r.attackers = append(r.attackers, attacker)
	go attack(attacker, r.next, r.quit, r.results, r.config.timeout())
}

// addResult is called from a dedicated goroutine.
func (r *runner) addResult(s result) result {
	m, ok := r.metrics[s.doResult.RequestLabel]
	if !ok {
		m = new(Metrics)
		r.metrics[s.doResult.RequestLabel] = m
	}
	m.add(s)
	return s
}

// test uses the Attack to perform {count} calls and report its result
// it is intended for development of an Attack implementation.
func (r *runner) test(count int) *RunReport {
	probe := r.prototype.Clone()
	if err := probe.Setup(r.config); err != nil {
		Printf("test attack setup failed [%v]", err)
		return &RunReport{Configuration: r.config, Output: map[string]interface{}{}}
	}
	defer probe.Teardown()
	for s := count; s > 0; s-- {
		now := time.Now()
		doResult := probe.Do(context.Background())
		end := time.Now()
		r.addResult(result{
			doResult: doResult,
			begin:    now,
			end:      end,
			elapsed:  end.Sub(now),
		})
		errorString := "no error"
		if doResult.Error != nil {
			errorString = fmt.Sprintf("error [%v:%T]", doResult.Error, doResult.Error)
		}
		Printf("test attack call [%s] took [%v] with status [%v] and %s\n", doResult.RequestLabel, end.Sub(now), doResult.StatusCode, errorString)
		select {
		case <-r.abort:
			goto end
		default:
		}
	}
end:
	return r.reportMetrics()
}

// run offers the complete flow of a load test.
func (r *runner) run() *RunReport {
	go r.collectResults()

	r.rampUp()
	r.fullAttack()

	r.quitAttackers()
	r.tearDownAttackers()
	return r.reportMetrics()
}

func (r *runner) fullAttack() {
	// attack can only proceed when at least one attacker is waiting for rps tokens
	if len(r.attackers) == 0 {
		// rampup probably has failed too
		return
	}
	if r.config.Verbose {
		Printf("begin full attack of [%d] remaining seconds\n", r.config.AttackTimeSec-r.config.RampupTimeSec)
	}
	fullAttackStartedAt = time.Now()
	limiter := ratelimit.New(r.config.RPS) // per second
	doneDeadline := time.Now().Add(time.Duration(r.config.AttackTimeSec-r.config.RampupTimeSec) * time.Second)
	for time.Now().Before(doneDeadline) {
		limiter.Take()
		select {
		case <-r.abort:
			goto end
		default:
		}
	}
end:
	if r.config.Verbose {
		Printf("end full attack")
	}
}

func (r *runner) rampUp() {
	strategy := strategyParameters{line: r.config.rampupStrategy()}
	if r.config.Verbose {
		Printf("begin rampup of [%d] seconds to RPS [%d] within attack of [%d] seconds\n", r.config.RampupTimeSec, r.config.RPS, r.config.AttackTimeSec)
	}
	if strategy.is("linear") {
		linearIncreasingGoroutinesAndRequestsPerSecondStrategy{}.execute(r)
	}
	if strategy.is("exp2") {
		spawnAsWeNeedStrategy{parameters: strategy}.execute(r)
	}
	// restore pipeline function incase it was changed by the rampup strategy
	r.resultsPipeline = r.addResult
	if r.config.Verbose {
		Printf("end rampup ending up with [%d] attackers\n", len(r.attackers))
	}
}

func (r *runner) quitAttackers() {
	if r.config.Verbose {
		Printf("stopping attackers [%d]\n", len(r.attackers))
	}
	for range r.attackers {
		r.quit <- true
	}
}

func (r *runner) tearDownAttackers() {
	if r.config.Verbose {
		Printf("tearing down attackers [%d]\n", len(r.attackers))
	}
	for i, each := range r.attackers {
		if err := each.Teardown(); err != nil {
			Printf("failed to teardown attacker [%d]:%v\n", i, err)
		}
	}
}

func (r *runner) reportMetrics() *RunReport {
	for _, each := range r.metrics {
		each.updateLatencies()
	}
	return &RunReport{
		StartedAt:     fullAttackStartedAt,
		FinishedAt:    time.Now(),
		Configuration: r.config,
		Metrics:       r.metrics,
		Failed:        false, // must be overwritten by program
		Output:        map[string]interface{}{},
	}
}

func (r *runner) collectResults() {
	for {
		r.resultsPipeline(<-r.results)
	}
}
