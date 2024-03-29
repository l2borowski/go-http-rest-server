package store

import (
	"fmt"
)

type Entry struct {
	Owner string      `json:"owner"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	//Timestamp time.Time   `json:"timestamp"`
}

var e = make([]Entry, 0)

func AddNewEntry(user, key string, value interface{}) {
	e = append(e, Entry{
		Key:   key,
		Value: value,
		Owner: user,
	})

	// e = append(e, Entry{
	// 	Owner:     user,
	// 	Key:       key,
	// 	Value:     value,
	// 	Timestamp: time.Now(),
	// })
}

func GetEntry(key string) (Entry, error) {
	for i := range e {
		if e[i].Key == key {
			return e[i], nil
		}
	}

	return Entry{}, fmt.Errorf("entry %w for key: %q", ErrNotFound, key)
}

func GetEntryOwner(key string) (string, error) {
	entry, err := GetEntry(key)
	if err != nil {
		return "", err
	}

	return entry.Owner, nil
}

func GetEntryValue(key string) (string, error) {
	entry, err := GetEntry(key)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(entry.Value), nil
}

func UpdateEntryValue(key string, value interface{}) error {
	for i := range e {
		if e[i].Key == key {
			e[i].Value = value
			//eShort[i].Timestamp = time.Now()
			return nil
		}
	}

	return fmt.Errorf("entry not found for key: %q", key)
}

//TODO: Handling short entry - Decide on one
func DeleteEntry(key string) error {
	// Find the index of the entry
	var index int
	var found bool
	for i, v := range e {
		if v.Key == key {
			index = i
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("entry not found for key: %q", key)
	}

	// Remove the entry by index
	e = append(e[:index], e[index+1:]...)

	return nil
}

func GetAllEntries() []Entry {
	return e
}
