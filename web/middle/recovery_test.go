package middle

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoveryHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("at the disco")
	})
	handler := RecoveryHandler(inner)
	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected := `internal server error`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			recorder.Body.String(), expected)
	}
}
