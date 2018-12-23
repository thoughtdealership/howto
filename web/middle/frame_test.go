package middle

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thoughtdealership/howto/app/frame"
)

func TestFrameHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()
	fr := frame.FromContext(ctx)
	if fr != nil {
		t.Error("Frame is non-nil")
	}
	recorder := httptest.NewRecorder()
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fr := frame.FromContext(ctx)
		if fr == nil {
			t.Error("Frame is nil")
		}
	})
	handler := FrameHandler(inner)
	handler.ServeHTTP(recorder, req)
}
