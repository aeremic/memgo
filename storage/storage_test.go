package storage

import (
	"testing"
)

func TestGet(t *testing.T) {
	s := New()

	s.Set("testKey", "testValue")

	actual := s.Get("testKey")

	expected := Value{Data: "testValue"}
	if expected.Data != actual {
		t.Fatalf("invalid actual value. got %s instead of %s",
			actual, expected.Data)
	}
}

func TestGetAllFormatted(t *testing.T) {
	s := New()

	s.Set("testKey1", "testValue1")
	s.Set("testKey2", "testValue2")

	actual := s.GetAll()

	expected := `{"testKey2":"testValue2", "testKey1":"testValue1"}`
	if expected != actual {
		t.Fatalf("invalid actual value. got %s instead of %s",
			actual, expected)
	}
}
