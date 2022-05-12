package main

import (
	"fmt"
	"os"

	"github.com/l2borowski/go-http-rest-server/server"
	"github.com/l2borowski/go-http-rest-server/store"
)

var GlobalStore = make(map[string]string)

func main() {
	var response interface{}
	var err error

	kvs := store.NewStoreData()
	user, _ := os.Hostname()

	fmt.Println("\nHello", user)

	store.Put(kvs, user, "myKey", 123)
	store.Put(kvs, user, "myKey", 1234)
	store.Put(kvs, "Lukasz", "LukaszKey", "Test123")

	//TODO: Write tests instead
	response, err = store.Get(kvs, "myKey")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("response:", response)
	}

	response, err = store.Get(kvs, "myKey2")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("response:", response)
	}

	store.Put(kvs, user, "myKey2", "abc")

	response, err = store.Get(kvs, "myKey2")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("response:", response)
	}

	err = store.Delete(kvs, user, "myKey")
	if err != nil {
		fmt.Println(err.Error())
	}

	err = store.Delete(kvs, user, "myKey")
	if err != nil {
		fmt.Println(err.Error())
	}

	response, err = store.Get(kvs, "myKey")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("response:", response)
	}

	err = store.Put(kvs, "MAC5013", "myKey2", 48)
	if err != nil {
		fmt.Println(err.Error())
	}

	response, err = store.Get(kvs, "myKey2")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("response:", response)
	}

	s := store.ListStore()
	fmt.Println(s)

	server.Listen(kvs, 8000)
}
