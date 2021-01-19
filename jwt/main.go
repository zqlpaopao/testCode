package main

import (
	"log"
	"net/http"
	"test/jwt/api"
)

func main() {

	// "Signin" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/signin", api.Signin)
	http.HandleFunc("/welcome", api.Welcome)
	//httpHandler.HandleFunc("/refresh", Refresh)

	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}
