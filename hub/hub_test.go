package hub

import "testing"

func TestHubBasic(t *testing.T) {
	h := New("localhost", "3336")
	h.Run()
}
