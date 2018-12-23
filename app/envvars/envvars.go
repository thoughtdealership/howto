package envvars

import (
	"github.com/ianschenck/envflag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EnvValues struct {
	Info struct {
		Name string
	}
	Log struct {
		LvlAsText string
		Lvl       zerolog.Level
	}
}

var Env EnvValues

func init() {
	configure()
}

func loglevel() {
	lvl, err := zerolog.ParseLevel(Env.Log.LvlAsText)
	if err != nil {
		log.Fatal().
			Str("envvar", "LOG_LEVEL").
			Err(err).
			Msg("unable to parse log level")
	}
	Env.Log.Lvl = lvl
}

func configure() {
	envflag.StringVar(&Env.Info.Name, "APPLICATION_NAME", "", "Application name")
	envflag.StringVar(&Env.Log.LvlAsText, "LOG_LEVEL", "debug", "Log level")
	envflag.Parse()

	loglevel()
}
