package store

import (
	"errors"
	"fmt"
)

type StoreData struct {
	Busy bool
	Data map[string]interface{}
}

var ErrNotOwner = errors.New("not owner")
var ErrNotFound = errors.New("not found")

func NewStoreData() *StoreData {
	return &StoreData{
		Data: make(map[string]interface{}),
	}
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

//TODO: Handle not found error
func ListStore() []Entry {
	return GetAllEntries()
}

//TODO: Returning short entry - Decide on one
//TODO: Handle not found error
func ListKey(key string) Entry {
	entry, err := GetEntry(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	return entry
}

func Authorised(user, owner string) bool {
	if user == "admin" {
		return true
	}

	return user == owner
}
