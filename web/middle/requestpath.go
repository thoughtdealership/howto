package middle

import (
	"net/http"

	"github.com/thoughtdealership/howto/app/frame"
)

// RequestPathHandler adds the request path to the logging context
func RequestPathHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fr := frame.FromContext(ctx)
		fr.Logger = fr.Logger.With().
			Str("request_path", r.URL.Path).
			Logger()
		h.ServeHTTP(w, r)
	})
}
