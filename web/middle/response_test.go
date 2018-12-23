package middle

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thoughtdealership/howto/app/exterror"
	"github.com/thoughtdealership/howto/app/frame"
	"github.com/julienschmidt/httprouter"
)

func TestResponseHandlerSuccess(t *testing.T) {
	req, err := http.NewRequest("GET", "/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()
	ctx = frame.NewContext(ctx)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	inner := func(r *http.Request, p httprouter.Params) (string, error) {
		return "success", nil
	}
	handler := ResponseHandler(inner)
	handler(recorder, req, nil)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestResponseHandlerWarning(t *testing.T) {
	req, err := http.NewRequest("GET", "/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()
	ctx = frame.NewContext(ctx)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	inner := func(r *http.Request, p httprouter.Params) (string, error) {
		return "", exterror.Create(http.StatusBadRequest, errors.New("foobar"))
	}
	handler := ResponseHandler(inner)
	handler(recorder, req, nil)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestResponseHandlerError(t *testing.T) {
	req, err := http.NewRequest("GET", "/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()
	ctx = frame.NewContext(ctx)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	inner := func(r *http.Request, p httprouter.Params) (string, error) {
		return "", errors.New("foobar")
	}
	handler := ResponseHandler(inner)
	handler(recorder, req, nil)
	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}
