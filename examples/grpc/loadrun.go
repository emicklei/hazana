package main

import (
	"log"

	"context"

	"github.com/emicklei/hazana"
	grpc "google.golang.org/grpc"
)

func main() {
	hazana.PrintSummary(hazana.Run(new(clockAttack), hazana.ConfigFromFlags()))
}

type clockAttack struct {
	conn   *grpc.ClientConn
	client ClockServiceClient
}

func (c *clockAttack) Setup(hc hazana.Config) error {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial failed:", err)
	}
	c.conn = conn
	c.client = NewClockServiceClient(conn)
	return nil
}

func (c *clockAttack) Do(ctx context.Context) hazana.DoResult {
	_, err := c.client.GetTime(ctx, new(GetTimeRequest))
	return hazana.DoResult{RequestLabel: "now", Error: err}
}

func (c *clockAttack) Teardown() error {
	return c.conn.Close()
}

func (c *clockAttack) Clone() hazana.Attack {
	// do not reuse client or connection
	return new(clockAttack)
}
