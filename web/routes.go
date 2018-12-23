package web

import (
	"errors"
	"net/http"

	"github.com/jonbodner/multierr"
	"github.com/julienschmidt/httprouter"
	"github.com/thoughtdealership/howto/app/exterror"
	"github.com/thoughtdealership/howto/app/version"
)

// Hello generates a successful response
func Hello(r *http.Request, p httprouter.Params) (string, error) {
	return "world", nil
}

// ErrorRoute generates a 500-level error response
func ErrorRoute(r *http.Request, p httprouter.Params) (string, error) {
	return "", errors.New("error message")
}

// UserErrorRoute generates a 400-level error response.
func UserErrorRoute(r *http.Request, p httprouter.Params) (string, error) {
	err := exterror.Create(http.StatusBadRequest, errors.New("user error"))
	return "", err
}

// MultiErrorRoute combines a 400-level error respones and a 500-level error response
// to generate a 500-level error response
func MultiErrorRoute(r *http.Request, p httprouter.Params) (string, error) {
	err1 := exterror.Create(http.StatusBadRequest, errors.New("user error"))
	err2 := errors.New("error message")
	err := multierr.Append(err1, err2)
	return "", err
}

// Panic is an example of a route that panics
func Panic(r *http.Request, p httprouter.Params) (string, error) {
	panic("at the disco")
}

// Version returns the application version
func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte(version.Version))
}
