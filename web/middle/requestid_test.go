package middle

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thoughtdealership/howto/app/frame"

	"github.com/google/uuid"
)

var emptyUUID uuid.UUID

func TestRequestIDHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()
	ctx = frame.NewContext(ctx)
	req = req.WithContext(ctx)
	fr := frame.FromContext(ctx)
	if fr.UUID != emptyUUID {
		t.Error("Request UUID is non-empty")
	}
	recorder := httptest.NewRecorder()
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fr := frame.FromContext(ctx)
		if fr.UUID == emptyUUID {
			t.Error("Request UUID is empty")
		}
	})
	handler := RequestIDHandler(inner)
	handler.ServeHTTP(recorder, req)
}
