package middle

import (
	"net/http"

	"github.com/thoughtdealership/howto/app/frame"
)

// FrameHandler creates a new request context.
// The new context contains an empty frame.
func FrameHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = frame.NewContext(ctx)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
