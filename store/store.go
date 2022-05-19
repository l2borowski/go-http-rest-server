package store

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Operation struct {
	action  string
	user    string
	key     string
	value   interface{}
	channel chan Response
}

type Response struct {
	Action   string
	User     string
	Response []byte
	Value    string
	Err      error
}

type StoreData struct {
	Busy bool
	Data map[string]interface{}
}

var requests chan Operation = make(chan Operation)

var done chan struct{} = make(chan struct{})

var ErrNotOwner = errors.New("not owner")
var ErrNotFound = errors.New("not found")

func NewStoreData() *StoreData {
	return &StoreData{
		Data: make(map[string]interface{}),
	}
}

func Start(kvs *StoreData) {
	go monitorRequests(kvs)
}

func Stop() {
	close(requests)
	<-done
}

func monitorRequests(kvs *StoreData) {
	for op := range requests {
		switch op.action {
		case "Get":
			response, err := Get(kvs, op.key)
			if err != nil {
				fmt.Println(err.Error())
			}
			op.channel <- Response{
				Action:   op.action,
				Response: []byte(response),
				Err:      err,
			}

		case "Put":
			err := Put(kvs, op.user, op.key, op.value)
			if err != nil {
				fmt.Println(err.Error())
			}
			value := fmt.Sprint(op.value)
			op.channel <- Response{
				Action:   op.action,
				Value:    value,
				Response: []byte(value),
				Err:      err,
			}

		case "Delete":
			err := Delete(kvs, op.user, op.key)
			if err != nil {
				fmt.Println(err.Error())
			}
			op.channel <- Response{
				Action: op.action,
				Err:    err,
			}

		case "ListStore":
			response := ListStore()
			op.channel <- Response{
				Action:   op.action,
				Response: response,
			}

		case "ListKey":
			response, err := ListKey(op.key)
			op.channel <- Response{
				Action:   op.action,
				Response: response,
				Err:      err,
			}

		default:
		}
	}
	close(done)
}

func GetRequest(key string, c chan Response) {
	op := Operation{
		action:  "Get",
		key:     key,
		channel: c,
	}

	requests <- op
}

func PutRequest(user, key string, value interface{}, c chan Response) {
	op := Operation{
		action:  "Put",
		user:    user,
		key:     key,
		value:   value,
		channel: c,
	}

	requests <- op
}

func DeleteRequest(user, key string, c chan Response) {
	op := Operation{
		action:  "Delete",
		user:    user,
		key:     key,
		channel: c,
	}

	requests <- op
}

func ListStoreRequest(c chan Response) {
	op := Operation{
		action:  "ListStore",
		channel: c,
	}

	requests <- op
}

func ListKeyRequest(key string, c chan Response) {
	op := Operation{
		action:  "ListKey",
		key:     key,
		channel: c,
	}

	requests <- op
}

func Get(kvs *StoreData, key string) (string, error) {
	kvs.Busy = true
	defer func() { kvs.Busy = false }()

	element, ok := kvs.Data[key]
	if !ok {
		return "", fmt.Errorf("GET: key %q: %w", key, ErrNotFound)
	}

	value := fmt.Sprint(element)

	return value, nil
}

func Put(kvs *StoreData, user, key string, value interface{}) error {
	kvs.Busy = true
	defer func() { kvs.Busy = false }()

	_, ok := kvs.Data[key]
	if ok {
		// Get corresponding entry
		entry, err := GetEntry(key)
		if err != nil {
			panic("Something went wrong - Missing entry")
		}

		// Check for permission
		if !Authorised(user, entry.Owner) {
			return fmt.Errorf("PUT: %q is %w of the %q", user, ErrNotOwner, key)
		}

		// Overwrite value for a given key
		kvs.Data[key] = value

		// Update entries
		err = UpdateEntryValue(key, value)
		if err != nil {
			panic("Something went wrong - Missing entry")
		}

	} else {
		// Create value for a given key
		kvs.Data[key] = value

		// Add entry
		AddNewEntry(user, key, value)
	}

	return nil
}

func Delete(kvs *StoreData, user, key string) error {
	kvs.Busy = true
	defer func() { kvs.Busy = false }()

	_, ok := kvs.Data[key]
	if !ok {
		return fmt.Errorf("DELETE: key %q: %w", key, ErrNotFound)
	}

	// Get corresponding entry
	entry, err := GetEntry(key)
	if err != nil {
		panic("Something went wrong - Missing entry")
	}

	// Check for permission
	if !Authorised(user, entry.Owner) {
		return fmt.Errorf("DELETE: %q is %w of the %q", user, ErrNotOwner, key)
	}

	// Delete key and its value
	delete(kvs.Data, key)

	// Delete entry
	DeleteEntry(key)
	if err != nil {
		panic("Something went wrong - Missing entry")
	}

	return nil
}

func ListStore() []byte {
	entries, err := json.Marshal(GetAllEntries())
	if err != nil {
		fmt.Println(err.Error())
	}

	return entries
}

func ListKey(key string) ([]byte, error) {
	entry, err := GetEntry(key)
	response, e := json.Marshal(entry)
	if e != nil {
		fmt.Println(e.Error())
	}

	return response, err
}

func Authorised(user, owner string) bool {
	if user == "admin" {
		return true
	}

	return user == owner
}
