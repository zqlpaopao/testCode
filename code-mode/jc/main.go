package main

import (
	"context"
	"fmt"
)

type StatusCmd struct {
	val string
}

func (s *StatusCmd) Name() string {
	return "StatusCmd-impl-Cmder"
}

type Cmdable interface {
	Ping(ctx context.Context) *StatusCmd
}

type UniversalClient interface {
	Cmdable
	Close() error
}

type Cmder interface {
	Name() string
}

type cmdable func(ctx context.Context, cmd Cmder) error

func (c cmdable) Ping(ctx context.Context) *StatusCmd {
	fmt.Println("cmdable-Ping")
	s := &StatusCmd{val: "ping"}
	c(ctx, s)
	return s
}

//////

type ClusterClient struct {
	Name string
	cmdable
}

func (c *ClusterClient) Process(ctx context.Context, cmd Cmder) (err error) {
	fmt.Println("ClusterClient-Process")
	return nil
}

func main() {
	cc := &ClusterClient{
		Name: "cluster",
	}
	cc.cmdable = cc.Process

	cc.Ping(context.Background())
}
