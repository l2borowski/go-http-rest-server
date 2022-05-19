package store

import (
	"errors"
	"fmt"
	"testing"
)

func TestPutStoreGet(t *testing.T) {
	kvs := NewStoreData()
	user := "User_1"
	key := "TestKey"
	value := 123
	expected := "123"

	err := Put(kvs, user, key, value)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Put value: %q to key: %q FAILED.", value, key)
	}

	response, err := Get(kvs, key)
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("Get value from key: %q FAILED.", key)
	}

	if response != "123" {
		t.Errorf("Get value: %d from key: %q FAILED. Expected: %q Received: %q", value, key, expected, response)
	} else {
		t.Logf("Get value: %d from key: %q PASSED. Expected: %q Received: %q", value, key, expected, response)
	}
}

func TestDeleteForbidden(t *testing.T) {
	kvs := NewStoreData()
	user := "User_1"

	adminUser := "admin"
	adminKey := "secure"
	adminValue := "48"

	err := Put(kvs, "admin", "secure", "48")
	if err != nil {
		fmt.Println(err.Error())
		t.Errorf("User: %q Put value: %q to key: %q FAILED.", adminUser, adminValue, adminKey)
	}

	err = Delete(kvs, user, adminKey)
	if err != nil {
		if errors.Is(err, ErrNotOwner) {
			t.Logf("User: %q not authorized to delete key: %q PASSED. Expected error: %q", user, adminKey, ErrNotOwner)
			return
		}
	}

	t.Errorf("User: %q not authorized to delete key: %q FAILED.", user, adminKey)
}
