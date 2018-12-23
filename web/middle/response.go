package middle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/thoughtdealership/howto/app/envvars"
	"github.com/thoughtdealership/howto/app/exterror"
	"github.com/thoughtdealership/howto/app/frame"

	"github.com/julienschmidt/httprouter"
)

type respHandle func(*http.Request, httprouter.Params) (string, error)

const genericError = `
{
	"error": "internal server error"
}
`

type Response struct {
	UUID    string `json:"uuid"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func createResponse(r *http.Request) Response {
	var resp Response
	ctx := r.Context()
	fr := frame.FromContext(ctx)
	resp.UUID = fr.UUID.String()
	return resp
}

func ResponseHandler(handle respHandle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		msg, err := handle(r, p)
		if err != nil {
			HandleError(w, r, err)
		} else {
			HandleResult(w, r, msg)
		}
	}
}

const messageTemplate = `
<html>
    <head>
        <title>{{.Title}}</title>
    </head>
    <body>{{.Body}}</body>
</html>
`

func WriteMessage(w http.ResponseWriter, r *http.Request, statusCode int, body string) {
	var buf bytes.Buffer
	t, err := template.New("message").Parse(messageTemplate)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	data := struct {
		Title string
		Body  string
	}{
		envvars.Env.Info.Name,
		body,
	}
	err = t.Execute(&buf, data)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/html")
	w.Write(buf.Bytes())
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	fr := frame.FromContext(ctx)
	exterr := exterror.Convert(ctx, err)
	if exterr.Status < 500 {
		fr.Logger.Warn().
			Str("error", exterr.Error()).
			Int("status", exterr.Status).
			Msg("user error reported")
	} else {
		fr.Logger.Error().
			Str("error", exterr.Error()).
			Int("status", exterr.Status).
			Msg("server error reported")
	}
	w.WriteHeader(exterr.Status)
	resp := createResponse(r)
	resp.Error = exterr.Error()
	writeResponse(w, r, resp)
}

func HandleResult(w http.ResponseWriter, r *http.Request, msg string) {
	ctx := r.Context()
	fr := frame.FromContext(ctx)
	fr.Logger.Info().Msg(msg)
	resp := createResponse(r)
	resp.Message = msg
	writeResponse(w, r, resp)
}

func writeResponse(w http.ResponseWriter, r *http.Request, resp Response) {
	ctx := r.Context()
	fr := frame.FromContext(ctx)
	body, err := json.Marshal(resp)
	if err != nil {
		fr.Logger.Error().
			Str("error", err.Error()).
			Msg("unable to convert response to JSON")
		body = []byte(genericError)
	}
	_, err = w.Write(body)
	if err != nil {
		fr.Logger.Error().
			Str("error", err.Error()).
			Str("body", fmt.Sprintf("%s", string(body))).
			Msg("unable to write response body")
		w.WriteHeader(500)
	}
}
