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

func TestGetByKeyAndPath_ShouldReturnName(t *testing.T) {
	s := New()

	s.Set("testKey", `{"a":2, "b":3, "people":{"names": ["Bob", "Alice"]}}`)
	actual := s.GetByKeyAndPath("testKey", ".people.names[0]")

	expected := "Bob"
	if expected != actual {
		t.Fatalf("invalid actual value. got %s instead of %s",
			actual, expected)
	}
}

func TestGetByKeyAndPath_ShouldReturnEmpty(t *testing.T) {
	s := New()

	s.Set("testKey", `{"a":2, "b":3}`)
	actual := s.GetByKeyAndPath("testKey", ".people.names[0]")

	expected := ""
	if expected != actual {
		t.Fatalf("invalid actual value. got %s instead of %s",
			actual, expected)
	}
}

func TestSelect(t *testing.T) {
	s := New()

	s.Set("testKey1", `{"a":2, "b":3, "people":{"names": ["Bob1", "Alice"]}}`)
	s.Set("testKey2", `{"a":2, "b":3, "people":{"names": ["Bob2", "Alice"]}}`)
	s.Set("testKey3", `{"a":2, "b":3}`)

	actual := s.Select(".people.names[0]")

	expected :=
		`{"testKey1":"{"a":2, "b":3, "people":{"names": ["Bob1", "Alice"]}}", "testKey2":"{"a":2, "b":3, "people":{"names": ["Bob2", "Alice"]}}"}`
	if expected != actual {
		t.Fatalf("invalid actual value. got \n%s\ninstead of \n%s\n",
			actual, expected)
	}
}
