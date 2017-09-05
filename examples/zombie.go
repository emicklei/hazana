package main

import (
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

func (z zombieAttack) Do() hazana.DoResult {
	time.Sleep(100 * time.Millisecond)
	return hazana.DoResult{RequestLabel: "sample"}
}

func (z zombieAttack) Teardown() error {
	return nil
}

func (z zombieAttack) Clone() hazana.Attack {
	return z
}

func main() {
	hazana.PrintReport(hazana.Run(zombieAttack{}, hazana.ConfigFromFlags()))
}
