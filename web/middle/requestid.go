package middle

import (
	"net/http"

	"github.com/thoughtdealership/howto/app/frame"

	"github.com/google/uuid"
)

// RequestIDHandler assigns a UUID to the request
func RequestIDHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fr := frame.FromContext(ctx)
		fr.UUID = uuid.New()
		fr.Logger = fr.Logger.With().
			Str("uuid", fr.UUID.String()).
			Logger()
		h.ServeHTTP(w, r)
	})
}
