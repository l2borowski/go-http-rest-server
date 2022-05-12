package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/l2borowski/go-http-rest-server/store"
)

var user, _ = os.Hostname()

func Listen(kvs *store.StoreData, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpHandler(kvs))

	fmt.Printf("Starting server on port %d\n\n", port)
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func httpHandler(kvs *store.StoreData) func(http.ResponseWriter, *http.Request) {
	if kvs == nil {
		panic("nil key value store!")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		param := strings.Split(path, "/")
		//fmt.Println("URL:", path, "Param:", param)

		switch path {
		case "/list":
			if r.Method == "GET" {
				listStore(w, r)
				return
			}
		case "/ping":
			if r.Method == "GET" {
				ping(w, r)
				return
			}
		case "/shutdown":
			if r.Method == "GET" {
				shutdown(w, r)
			}
		}

		if len(param) >= 3 {
			fmt.Println("Param:", param[2], "Method:", r.Method, "Path:", path)
			switch path {
			case "/store/" + param[2]:
				if r.Method == "GET" {
					get(kvs, w, param[2])
				} else if r.Method == "PUT" {
					put(kvs, w, r)
				} else if r.Method == "DELETE" {
					delete(kvs, w, r)
				}
			case "/list/%v", param[2]:
				if r.Method == "GET" {
					listKey(w, r)
				}
			}
		}
	}
}

func get(kvs *store.StoreData, w http.ResponseWriter, param string) {
	fmt.Println("GET:", param)
	response, err := store.Get(kvs, param)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("GET: ", param)
	//w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func put(kvs *store.StoreData, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("URL Param 'key' is missing")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	key := keys[0]
	values, ok := r.URL.Query()["value"]
	if !ok || len(values[0]) < 1 {
		fmt.Println("URL Param 'value' is missing")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	value := values[0]
	err := store.Put(kvs, user, key, value)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, store.ErrNotOwner) {
			w.WriteHeader(http.StatusForbidden)
		} else if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}

	//w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func delete(kvs *store.StoreData, w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		fmt.Println("URL Param 'key' is missing")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	key := keys[0]
	err := store.Delete(kvs, user, key)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, store.ErrNotOwner) {
			w.WriteHeader(http.StatusForbidden)
		} else if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func listStore(w http.ResponseWriter, r *http.Request) {
	response := store.ListStore()
	for _, v := range response {
		w.Write([]byte(fmt.Sprint(v) + "\n"))
	}

	w.WriteHeader(http.StatusOK)
}

func listKey(w http.ResponseWriter, r *http.Request) {

}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func shutdown(w http.ResponseWriter, r *http.Request) {
	//TODO: Get actual user name - Overwriting user to admin as a temporary solution
	user = "admin"
	if user == "admin" {
		w.WriteHeader(http.StatusOK)

		go func() {
			time.Sleep(time.Millisecond)
			os.Exit(0)
		}()
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
