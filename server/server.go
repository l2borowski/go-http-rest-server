package server

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func Listen() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/shutdown", shutdownHandler)

	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received: ping request")
	fmt.Println("Response: pong")

}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received: shutdown request")
	fmt.Println("Response: shutting down server...")
	time.Sleep(time.Millisecond)
	os.Exit(0)
}
