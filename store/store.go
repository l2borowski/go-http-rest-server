package store

import (
	"errors"
	"fmt"
)

type StoreData struct {
	Busy bool

	KeyOwner map[string]string
	Data     map[string]interface{}
}

var ErrNotOwner = errors.New("not owner")
var ErrNotFound = errors.New("not found")

func NewStoreData() *StoreData {
	return &StoreData{
		Data: make(map[string]interface{}),
	}
}

func Get(s *StoreData, key string) (interface{}, error) {
	s.Busy = true
	defer func() { s.Busy = false }()

	element, ok := s.Data[key]
	if !ok {
		return nil, fmt.Errorf("get: key %q: %w", key, ErrNotFound)
	}

	return element, nil
}

func Put(s *StoreData, user, key string, value interface{}) error {
	s.Busy = true
	defer func() { s.Busy = false }()

	_, ok := s.Data[key]
	if ok {
		// Get corresponding entry
		entry, err := GetEntryByKey(key)
		if err != nil {
			panic("Something went wrong - Missing entry")
		}

		// Check for permission
		if !Authorised(user, entry.User) {
			return fmt.Errorf("put: %q is %w of the %q", user, ErrNotOwner, key)
		}

		// Overwrite value for a given key
		s.Data[key] = value

		// Update entries
		err = UpdateEntryValue(key, value)
		if err != nil {
			panic("Something went wrong - Missing entry")
		}

	} else {
		// Create value for a given key
		s.Data[key] = value

		// Add entry
		AddNewEntry(user, key, value)
	}

	return nil
}

func Delete(s *StoreData, user, key string) error {
	s.Busy = true
	defer func() { s.Busy = false }()

	_, ok := s.Data[key]
	if !ok {
		return fmt.Errorf("delete: key %q: %w", key, ErrNotFound)
	}

	// Get corresponding entry
	entry, err := GetEntryByKey(key)
	if err != nil {
		panic("Something went wrong - Missing entry")
	}

	// Check for permission
	if !Authorised(user, entry.User) {
		return fmt.Errorf("delete: %q is %w of the %q", user, ErrNotOwner, key)
	}

	// Delete key and its value
	delete(s.Data, key)

	// Delete entry
	DeleteEntry(key)
	if err != nil {
		panic("Something went wrong - Missing entry")
	}

	return nil
}

func ListAll() {}

func Authorised(user, owner string) bool {
	if user == "admin" {
		return true
	}

	return user == owner
}
