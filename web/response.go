package web

import (
	"errors"
	"net/http"

	"github.com/thoughtdealership/howto/app/exterror"
	"github.com/thoughtdealership/howto/app/frame"
	"github.com/thoughtdealership/howto/web/middle"
)

func UnauthorizedError(w http.ResponseWriter, r *http.Request, action string, hidden error) {
	ctx := r.Context()
	fr := frame.FromContext(ctx)
	fr.Logger.Warn().
		Err(hidden).
		Msg(action)
	err := exterror.Create(http.StatusUnauthorized, errors.New("unauthorized request"))
	middle.HandleError(w, r, err)
}

func UserError(w http.ResponseWriter, r *http.Request, err error) {
	err = exterror.Create(http.StatusBadRequest, err)
	middle.HandleError(w, r, err)
}

func ServerError(w http.ResponseWriter, r *http.Request, action string, hidden error) {
	ctx := r.Context()
	fr := frame.FromContext(ctx)
	fr.Logger.Error().
		Err(hidden).
		Msg(action)
	err := exterror.Create(http.StatusInternalServerError, errors.New("internal server error"))
	middle.HandleError(w, r, err)
}
