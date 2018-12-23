package middle

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/thoughtdealership/howto/app/frame"
)

func BodyHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		ctx := r.Context()
		fr := frame.FromContext(ctx)
		if r.Body != nil {
			var err error
			var buf bytes.Buffer
			tee := io.TeeReader(r.Body, &buf)
			body, err = ioutil.ReadAll(tee)
			if err != nil {
				fr.Logger.Error().Err(err).Msg("unable to read request body")
			} else {
				r.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
				fr.Logger = fr.Logger.With().
					Bytes("body", body).
					Logger()

			}
		}
		h.ServeHTTP(w, r)
	})
}
