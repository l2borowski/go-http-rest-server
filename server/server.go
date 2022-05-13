package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/l2borowski/go-http-rest-server/store"
)

func Listen(kvs *store.StoreData, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpHandler(kvs))

	fmt.Printf("Starting server on port %d\n\n", port)
	err := http.ListenAndServe(":8000", mux) //TODO: Pass port value to the function
	if err != nil {
		fmt.Println(err.Error())
	}
}

func httpHandler(kvs *store.StoreData) func(http.ResponseWriter, *http.Request) {
	if kvs == nil {
		panic("nil key value store!")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Get URL path and parameters
		path := r.URL.Path
		param := strings.Split(path, "/")

		// Get username
		user := ""
		if len(r.Header["Authorization"]) > 0 {
			user = r.Header["Authorization"][0]
		}

		switch path {
		case "/list":
			if r.Method == "GET" {
				fmt.Println("GET:", r.URL.String())
				listStore(w, r)
				return
			}
		case "/ping":
			if r.Method == "GET" {
				fmt.Println("GET:", r.URL.String())
				ping(w, r)
				return
			}
		case "/shutdown":
			if r.Method == "GET" {
				fmt.Println("GET:", r.URL.String(), "user:", user)
				shutdown(user, w, r)
			}
		}

		// Check if any parameters are passed
		if len(param) >= 3 {
			key := param[2]

			// Convert response body to string
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println(err.Error())
			}
			value := string(bodyBytes)

			switch path {
			case "/store/" + key:
				if r.Method == "GET" {
					fmt.Println("GET:", r.URL.String(), "key:", key)
					get(kvs, key, w, r)
				} else if r.Method == "PUT" {
					fmt.Println("PUT:", r.URL.String(), "value:", value)
					put(kvs, user, key, value, w, r)
				} else if r.Method == "DELETE" {
					fmt.Println("DELETE:", r.URL.String(), "key:", key)
					delete(kvs, user, key, w, r)
				}
			case "/list/" + key:
				if r.Method == "GET" {
					fmt.Println("GET:", r.URL.String(), "key:", key)
					listKey(key, w, r)
				}
			}
		}
	}
}

func get(kvs *store.StoreData, k string, w http.ResponseWriter, r *http.Request) {
	response, err := store.Get(kvs, k)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(response))
}

func put(kvs *store.StoreData, u, k, v string, w http.ResponseWriter, r *http.Request) {
	err := store.Put(kvs, u, k, v)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, store.ErrNotOwner) {
			w.WriteHeader(http.StatusForbidden)
		} else if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}

	w.Write([]byte(v))
}

func delete(kvs *store.StoreData, u, k string, w http.ResponseWriter, r *http.Request) {
	err := store.Delete(kvs, u, k)
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
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Write(responseBytes)
	w.WriteHeader(http.StatusOK)
}

func listKey(k string, w http.ResponseWriter, r *http.Request) {
	response := store.ListKey(k)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Write(responseBytes)
	w.WriteHeader(http.StatusOK)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func shutdown(u string, w http.ResponseWriter, r *http.Request) {
	if u == "admin" {
		w.WriteHeader(http.StatusOK)

		go func() {
			time.Sleep(time.Millisecond)
			os.Exit(0)
		}()
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
