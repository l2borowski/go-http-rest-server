package store

import (
	"fmt"
	"time"
)

type Request struct {
	Command   string
	Key       string
	Value     interface{}
	Timestamp time.Time
}

type Response struct {
	Value     interface{}
	Timestamp time.Time
}

func (r Request) GetTransaction() error {
	return nil
}

//TODO: Divide into separate send requests
func SendRequest(key, value string, cmd string, rc chan<- Request) error {
	switch cmd {
	case "/get":
	case "/put":
	case "/delete":
	case "/list":
	case "/ping":
		go func() {
			rc <- Request{
				Command:   cmd,
				Key:       key,
				Value:     value,
				Timestamp: time.Now(),
			}
		}()
	default:
		return fmt.Errorf("unknown command")
	}

	return nil
}
