package main

import (
	"time"

	"github.com/emicklei/hazana"
)

// Perform one sample call
// go run zombie.go -t
//
// go run zombie.go -rps 10

type zombieAttack struct{}

func (z zombieAttack) Setup(c hazana.Config) error {
	return nil
}

func (z zombieAttack) Do() hazana.DoResult {
	time.Sleep(100 * time.Millisecond)
	return hazana.DoResult{}
}

func (z zombieAttack) Teardown() error {
	return nil
}

func (z zombieAttack) Clone() hazana.Attack {
	return z
}

func main() {
	hazana.Run(zombieAttack{}, hazana.ConfigFromFlags())
}
