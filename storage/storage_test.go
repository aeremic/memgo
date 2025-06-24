package storage

import (
	"testing"
)

func TestSettingStorageValuesBasic(t *testing.T) {
	s := New()

	s.Set("testKey", "testValue")

	actual := s.Get("testKey")

	expected := Value{Data: "testValue"}
	if expected != actual {
		t.Fatalf("invalid actual value. got %s instead of %s",
			actual, expected)
	}
}
