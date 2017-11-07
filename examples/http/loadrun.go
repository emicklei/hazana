package main

import (
	"net/http"

	"github.com/emicklei/hazana"
)

func main() {
	hazana.PrintSummary(hazana.Run(new(siteAttack), hazana.ConfigFromFlags()))
}

type siteAttack struct {
	client *http.Client
}

func (a *siteAttack) Setup(hc hazana.Config) error {
	a.client = new(http.Client)
	return nil
}

func (a *siteAttack) Do() hazana.DoResult {
	_, err := a.client.Get("http://ubanita.org")
	if err != nil {
		return hazana.DoResult{Error: err}
	}
	return hazana.DoResult{RequestLabel: "ubanita.org"}
}

func (a *siteAttack) Teardown() error { return nil }

func (a *siteAttack) Clone() hazana.Attack {
	// do not reuse client or connection
	return new(siteAttack)
}
