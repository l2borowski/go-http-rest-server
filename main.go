package main

import (
	"github.com/l2borowski/go-http-rest-server/server"
	"github.com/l2borowski/go-http-rest-server/store"
)

func main() {
	kvs := store.NewStoreData()
	server.Listen(kvs, 8000)
}
