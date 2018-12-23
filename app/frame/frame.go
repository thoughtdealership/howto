package frame

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type frameCtx struct{}

var frameKey frameCtx

type Frame struct {
	UUID   uuid.UUID
	Logger zerolog.Logger

	Foo string
	Bar bool
	Baz struct {
		A int
		B byte
		C string
	}
}

func NewContext(ctx context.Context) context.Context {
	fr := new(Frame)
	fr.Logger = log.Logger
	return context.WithValue(ctx, frameKey, fr)
}

func FromContext(ctx context.Context) *Frame {
	fr := ctx.Value(frameKey)
	if fr == nil {
		return nil
	}
	return fr.(*Frame)
}
