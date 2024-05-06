package cache

import (
	"reflect"
	"testing"
)

type TestType struct {
	Value string
}

func Test_Cache(t *testing.T) {
	initialType := TestType{
		Value: "Hello, world",
	}
	err := Insert("a", "b", initialType)
	if err != nil {
		t.Error("unexpected error.")
	}
	err = Insert("a", "b", initialType)
	if err == nil || err != ErrAlreadyCached {
		t.Error("expected error.")
	}
	retrieved := Retrieve("a", "b")
	if retrieved == nil {
		t.Error("Unexpected value")
	}
	if !reflect.DeepEqual(initialType, retrieved) {
		t.Error("Error comparing")
	}
	if nil != Retrieve("a", "c") {
		t.Error("Unexpected value")
	}
}
