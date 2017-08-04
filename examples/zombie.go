package main

import (
	"time"

	"github.com/emicklei/hazana"
)

// go run zombie.go -rps 10

type zombieAttack struct{}

func (z zombieAttack) Setup(c hazana.Config) error {
	return nil
}

func (z zombieAttack) Do() hazana.DoResult {
	time.Sleep(100 * time.Millisecond)
	return hazana.DoResult{}
}

func (z zombieAttack) TearDown() error {
	return nil
}

func (z zombieAttack) Clone() hazana.Attack {
	return z
}

func main() {
	hazana.Run(zombieAttack{}, hazana.ConfigFromFlags())
}
