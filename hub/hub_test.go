package hub

import (
	"context"
	"log"
	"memgo/storage"
	"testing"
)

func TestHubBasic(t *testing.T) {
	t.Skip("Skipping this test..")

	ctx, cancel := context.WithCancel(context.Background())

	s := storage.New()
	h := New("localhost", "3336")
	if err := h.Run(ctx, cancel, s); err != nil {
		log.Fatal(err)
	}
}
