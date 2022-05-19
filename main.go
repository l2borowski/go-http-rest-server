package main

import (
	"flag"

	"github.com/l2borowski/go-http-rest-server/server"
	"github.com/l2borowski/go-http-rest-server/store"
)

func main() {
	kvs := store.NewStoreData()
	store.Start(kvs)

	var port int
	flag.IntVar(&port, "port", 8000, "port to listen on")
	flag.Parse()

	server.Listen(kvs, port)
}
