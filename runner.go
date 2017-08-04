package hazana

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type runner struct {
	config     Config
	attackers  []Attack
	next, quit chan bool
	results    chan result
	prototype  Attack
	metrics    map[int]*Metrics
}

// Run starts attacking a service using an Attack implementation and a configuration.
func Run(a Attack, c Config) {
	if msg := c.Validate(); len(msg) > 0 {
		for _, each := range msg {
			fmt.Println("[config error]", each)
		}
		fmt.Println()
		flag.Usage()
		os.Exit(0)
	}
	r := new(runner)
	r.config = c
	r.prototype = a
	r.init()
	r.run()
}

func (r *runner) init() {
	r.next = make(chan bool)
	r.quit = make(chan bool)
	r.results = make(chan result)
	r.attackers = []Attack{}
	r.metrics = map[int]*Metrics{}
}

func (r *runner) spawnAttacker() {
	if r.config.Verbose {
		log.Printf("setup and spawn new attacker [%d]\n", len(r.attackers)+1)
	}
	attacker := r.prototype.Clone()
	if err := attacker.Setup(r.config); err != nil {
		log.Printf("attacker [%d] setup failed with [%v]\n", len(r.attackers)+1, err)
		return
	}
	r.attackers = append(r.attackers, attacker)
	go attack(attacker, r.next, r.quit, r.results)
}

// addResult is called from a dedicated goroutine.
func (r *runner) addResult(s result) {
	m, ok := r.metrics[s.request]
	if !ok {
		m = new(Metrics)
		r.metrics[s.request] = m
	}
	m.add(s)
}

func (r *runner) run() {
	// TODO which request takes longest?
	// launch a probe
	r.spawnAttacker()
	r.next <- true
	result := <-r.results
	if r.config.Verbose {
		log.Printf("probe response time [%v]\n", result.elapsed)
		if result.err != nil {
			log.Fatal("probe failed ", result.err)
		}
	}

	// based on initial response, we spawn attackers
	// they wait on receive from the next channel
	probeMs := (result.elapsed.Nanoseconds() / 1e06) + 1
	attackerCount := int(probeMs * int64(r.config.RPS) / 1000)
	if attackerCount > r.config.MaxAttackers {
		attackerCount = r.config.MaxAttackers
	}
	for i := 0; i < attackerCount-1; i++ { // minus 1 because the probe is still active
		r.spawnAttacker()
	}
	slot := 1000 / *oRPS // time between requests to be send in milliseconds
	slotTicker := time.Tick(time.Duration(slot) * time.Millisecond)

	rampupDeadline := time.Now().Add(time.Duration(r.config.RampupTimeSec) * time.Second)
	delayMs := int(probeMs) - slot
	requests := 0
	for time.Now().Before(rampupDeadline) {
		<-slotTicker
		if delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
			delayMs-- // TODO compute the delta, use 1 for now
		}
		r.next <- true
		<-r.results
		requests++
	}
	if r.config.Verbose {
		log.Printf("sent [%d] requests during rampup of [%v] (average %v rps)", requests, time.Duration(r.config.RampupTimeSec)*time.Second, float64(requests)/float64((time.Duration(*oRampupTime)*time.Second).Seconds()))
	}

	go r.collectResults()

	doneDeadline := time.Now().Add(time.Duration(r.config.AttackTimeSec-r.config.RampupTimeSec) * time.Second)
	for time.Now().Before(doneDeadline) {
		<-slotTicker
		r.next <- true
	}

	r.quitAttackers()
	r.tearDownAttackers()
	r.reportMetrics()
}

func (r *runner) quitAttackers() {
	for i := range r.attackers {
		if r.config.Verbose {
			log.Printf("stopping attacker [%d]\n", i+1)
		}
		r.quit <- true
	}
}

func (r *runner) tearDownAttackers() {
	for i, each := range r.attackers {
		if r.config.Verbose {
			log.Printf("tearing down attacker [%d]\n", i+1)
		}
		_ = each.TearDown()
	}
}

func (r *runner) reportMetrics() {
	for _, metrics := range r.metrics {
		metrics.updateLatencies()
		json.NewEncoder(os.Stdout).Encode(metrics)
	}
}

func (r *runner) collectResults() {
	for {
		result := <-r.results
		if result.err != nil && r.config.Verbose {
			log.Println("WARN ", result.err)
		}
		r.addResult(result)
	}
}
