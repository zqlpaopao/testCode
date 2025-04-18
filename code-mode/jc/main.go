package main

import (
	"context"
	"fmt"
)

type Cmder interface {
	Build() error
	Check() error
	Do() error
	GetTag() string
	Error() error
	SerError(err error)
}

type StatusCmd struct {
	tag string
	err error
}

func (s *StatusCmd) GetTag() string {
	return s.tag
}

func (s *StatusCmd) Build() error {
	fmt.Println("StatusCmd-impl-Cmder-Build")
	return nil
}

func (s *StatusCmd) Check() error {
	//return errors.New("StatusCmd-impl-Chec？k")
	fmt.Println("StatusCmd-impl-Check")
	return nil
}
func (s *StatusCmd) Do() error {
	fmt.Println("StatusCmd-impl-Do")
	return nil
}

func (s *StatusCmd) Error() error {
	return fmt.Errorf("%v : %v", s.GetTag(), s.err)
}

func (s *StatusCmd) SerError(err error) {
	s.err = err
}

type CmdAble interface {
	Run(ctx context.Context, cmder Cmder) Cmder
}

type cmdAble func(ctx context.Context, cmd Cmder)

func (c cmdAble) Run(ctx context.Context, able Cmder) Cmder {
	var err error
	if err = able.Build(); err != nil {
		goto END
	}
	if err = able.Check(); err != nil {
		goto END
	}
	c(ctx, able)
END:
	able.SerError(err)
	return able
}

type Proxy struct {
	cmdAble
	Cmder
}

func (c *Proxy) Process(ctx context.Context, cmd Cmder) {
	cmd.SerError(cmd.Do())
	fmt.Println("ClusterClient-Process")
}

func NewProxy(cmd Cmder) *Proxy {
	p := &Proxy{Cmder: cmd}
	p.cmdAble = p.Process
	return p
}

func (p *Proxy) Do() error {
	return p.Run(context.Background(), p.Cmder).Error()
}

func main() {
	cc := NewProxy(&StatusCmd{tag: "StatusCmd"}).Do()
	fmt.Println(cc)
}
