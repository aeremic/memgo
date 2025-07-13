package storage

import (
	"testing"
)

func TestGetSet(t *testing.T) {
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

	expected := `{"testKey1":"testValue1", "testKey2":"testValue2"}`
	if expected != actual {
		t.Fatalf("invalid actual value. got %s instead of %s",
			actual, expected)
	}
}

func TestDelete(t *testing.T) {
	s := New()

	s.Set("testKey", "testValue")

	s.Delete("testKey")

	actual := s.Get("testKey")
	expected := Value{Data: ""}
	if expected.Data != actual {
		t.Fatalf("invalid actual value. got %s instead of %s",
			actual, expected.Data)
	}
}

func TestDeleteAll(t *testing.T) {
	s := New()

	s.Set("testKey", "testValue")

	s.DeleteAll()

	actual := s.GetAll()
	expected := "{}"
	if expected != actual {
		t.Fatalf("invalid actual value. got %s instead of %s",
			actual, expected)
	}
}
