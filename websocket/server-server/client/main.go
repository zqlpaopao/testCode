package main

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Client started")
	for {
		conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "ws://127.0.0.1:8080/")
		if err != nil {
			fmt.Println("Cannot connect: " + err.Error())
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		fmt.Println("Connected to server")
		for i := 0; i < 10; i++ {
			randomNumber := strconv.Itoa(rand.Intn(100))
			msg := []byte(randomNumber)
			err = wsutil.WriteClientMessage(conn, ws.OpText, msg)
			if err != nil {
				fmt.Println("Cannot send: " + err.Error())
				continue
			}
			fmt.Println("Client message send with random number " + randomNumber)
			msg, _, err := wsutil.ReadServerData(conn)
			if err != nil {
				fmt.Println("Cannot receive data: " + err.Error())
				continue
			}
			fmt.Println("Server message received with random number: " + string(msg))
			time.Sleep(time.Duration(5) * time.Second)
		}
		err = conn.Close()
		if err != nil {
			fmt.Println("Cannot close the connection: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("Disconnected from server")
	}
}