package main

import (
	"context"
	"log"
	"time"

	"github.com/emicklei/hazana"
)

// Perform one sample call
// go run zombie.go -t 1
//
// go run zombie.go -v -rps 10

type zombieAttack struct{}

func (z zombieAttack) Setup(c hazana.Config) error {
	return nil
}

func (z zombieAttack) Do(ctx context.Context) hazana.DoResult {
	time.Sleep(100 * time.Millisecond)
	return hazana.DoResult{RequestLabel: "sample"}
}

func (z zombieAttack) Teardown() error {
	return nil
}

func (z zombieAttack) Clone() hazana.Attack {
	return z
}

// optional you can implement BeforeRunAware
func (z zombieAttack) BeforeRun(c hazana.Config) error {
	log.Println("before run")
	return nil
}

// optional you can implement AfterRunAware
func (z zombieAttack) AfterRun(r *hazana.RunReport) error {
	log.Println("after run")
	r.Failed = true // he survived
	return nil
}

func main() {
	r := hazana.Run(new(zombieAttack), hazana.ConfigFromFlags())
	hazana.PrintReport(r)
	hazana.PrintSummary(r)
}
