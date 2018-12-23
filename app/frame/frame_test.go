package frame

import (
	"context"
	"testing"
)

func TestContext(t *testing.T) {
	ctx := context.Background()
	fr := FromContext(ctx)
	if fr != nil {
		t.Error("Fetched frame from the background")
	}
	ctx = NewContext(ctx)
	fr = FromContext(ctx)
	if fr == nil {
		t.Error("Frame was not created")
	}
}
