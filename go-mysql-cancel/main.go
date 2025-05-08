package main

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"log"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(e *canal.RowsEvent) error {
	log.Printf("%s %v\n", e.Action, e.Rows)
	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func main() {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = "127.0.0.1:3306"
	cfg.User = "root"
	cfg.Password = "meimima123"
	cfg.Dump.ExecutionPath = ""
	// We only care table canal_test in test db
	//cfg.Dump.TableDB = "test"
	cfg.Dump.SkipMasterData = true
	//cfg.Dump.Tables = []string{"canal_test"}
	// 关键配置：确保使用ROW模式
	//cfg.Flavor = "mysql"

	c, err := canal.NewCanal(cfg)
	c.SetEventHandler(&MyEventHandler{})

	n, err := c.GetMasterPos()
	c.RunFrom(mysql.Position{
		Name: n.Name,
		Pos:  n.Pos,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Register a handler to handle RowsEvent

	// Start canal
	c.Run()
}
