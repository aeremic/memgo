package hub

import (
	"context"
	"log"
	"testing"
)

func TestHubBasic(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	h := New("localhost", "3336")
	if err := h.Run(ctx, cancel); err != nil {
		log.Fatal(err)
	}
}
