package main

import (
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"math/rand"
	"net/http"
	"strconv"
)

func main() {
	fmt.Println("Server started, waiting for connection from client")
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Client connected")
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			fmt.Println("Error starting socket server: " + err.Error())
		}
		go func() {
			defer conn.Close()
			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					fmt.Println("Error receiving data: " + err.Error())
					fmt.Println("Client disconnected")
					return
				}
				fmt.Println("Client message received with random number: " + string(msg))
				randomNumber := strconv.Itoa(rand.Intn(100))
				err = wsutil.WriteServerMessage(conn, op, []byte(randomNumber))
				if err != nil {
					fmt.Println("Error sending data: " + err.Error())
					fmt.Println("Client disconnected")
					return
				}
				fmt.Println("Server message send with random number " + randomNumber)
			}
		}()
	}))
}