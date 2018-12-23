package transport

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/thoughtdealership/howto/app/frame"
)

type LoggingTransport struct {
	Inner http.RoundTripper
}

func CreateLoggingTransport() http.RoundTripper {
	return LoggingTransport{
		Inner: http.DefaultTransport,
	}
}

// RoundTrip implements the RoundTripper interface.
func (t LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		var err error
		var buf bytes.Buffer
		tee := io.TeeReader(req.Body, &buf)
		body, err = ioutil.ReadAll(tee)
		if err != nil {
			return nil, err
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
	}
	ctx := req.Context()
	fr := frame.FromContext(ctx)
	fr.Logger.Info().
		Bytes("body", body).
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Msg("remote request")
	return t.Inner.RoundTrip(req)
}
