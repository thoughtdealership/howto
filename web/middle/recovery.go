package middle

import (
	"fmt"
	"net/http"

	"github.com/thoughtdealership/howto/app/frame"

	"github.com/go-stack/stack"
)

// RecoveryHandler adds panic recovery and logs the stack trace
func RecoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			e := recover()
			if e != nil {
				ctx := r.Context()
				fr := frame.FromContext(ctx)
				// initialize empty frame if none is found
				if fr == nil {
					ctx = frame.NewContext(ctx)
					fr = frame.FromContext(ctx)
				}
				trace := stack.Trace().TrimRuntime()
				fr.Logger.Error().
					Str("error", fmt.Sprintf("%s", e)).
					Str("stacktrace", fmt.Sprintf("%+v", trace)).
					Msg("panic recovery")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
			}
		}()
		h.ServeHTTP(w, r)
	})
}
