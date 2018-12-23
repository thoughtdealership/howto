package main

import (
	"net/http"

	"github.com/thoughtdealership/howto/app/envvars"
	"github.com/thoughtdealership/howto/web"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.SetGlobalLevel(envvars.Env.Log.Lvl)
}

func main() {
	router := web.Create()

	// Use github.com/akrylysov/algnhsa to run on AWS Lambda/API Gateway
	// without changing the existing HTTP handlers.
	// algnhsa.ListenAndServe(router, nil)
	http.ListenAndServe(":8080", router)
}
