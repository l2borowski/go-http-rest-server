package server

import (
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
	http.HandleFunc("/", httpHandler(kvs))
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Starting server on port %d\n\n", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func httpHandler(kvs *store.StoreData) func(http.ResponseWriter, *http.Request) {
	if kvs == nil {
		panic("nil key value store!")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialise channel for handling HTTP responses
		var resChan = make(chan store.Response)

		// Get URL path and parameters
		path := r.URL.Path
		param := strings.Split(path, "/")

		// Get username
		user := ""
		if len(r.Header["Authorization"]) > 0 {
			user = r.Header["Authorization"][0]
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
				if r.Method == http.MethodGet {
					//fmt.Println("GET:", r.URL.String(), "key:", key)
					go store.GetRequest(key, resChan)
				} else if r.Method == http.MethodPut {
					//fmt.Println("PUT:", r.URL.String(), "value:", value)
					go store.PutRequest(user, key, value, resChan)
				} else if r.Method == http.MethodDelete {
					//fmt.Println("DELETE:", r.URL.String(), "key:", key)
					go store.DeleteRequest(user, key, resChan)
				}
			case "/list/" + key:
				if r.Method == http.MethodGet {
					//fmt.Println("GET:", r.URL.String(), "key:", key)
					go store.ListKeyRequest(key, resChan)
				}
			}
		} else {
			switch path {
			case "/list":
				if r.Method == http.MethodGet {
					//fmt.Println("GET:", r.URL.String())
					go store.ListStoreRequest(resChan)
				}
			case "/ping":
				if r.Method == http.MethodGet {
					//fmt.Println("GET:", r.URL.String())
					PingResponse(w)
					return
				}
			case "/shutdown":
				if r.Method == http.MethodGet {
					//fmt.Println("GET:", r.URL.String(), "user:", user)
					ShutdownResponse(user, w)
					return
				}
			}
		}

		// Handle HTTP responses
		for res := range resChan {
			switch res.Action {
			case "Get":
				GetResponse(res.Response, res.Err, w)
				return
			case "Put":
				PutResponse(res.Value, res.Err, w)
				return
			case "Delete":
				DeleteResponse(res.Err, w)
				return
			case "ListStore":
				ListStoreResponse(res.Response, w)
				return
			case "ListKey":
				ListKeyResponse(res.Response, res.Err, w)
				return
			case "Ping":
				fmt.Println("Pinging...")
				PingResponse(w)
				return
			case "Shutdown":
				fmt.Println("Shutingdown...")
				ShutdownResponse(res.User, w)
				return
			}
		}
	}
}

func GetResponse(response []byte, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Write([]byte(response))
}

func PutResponse(value string, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")

	if err != nil {
		if errors.Is(err, store.ErrNotOwner) {
			w.WriteHeader(http.StatusForbidden)
		} else if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}

	w.Write([]byte(value))
}

func DeleteResponse(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")

	if err != nil {
		if errors.Is(err, store.ErrNotOwner) {
			w.WriteHeader(http.StatusForbidden)
		} else if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func ListStoreResponse(response []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ListKeyResponse(response []byte, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func PingResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func ShutdownResponse(u string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
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
