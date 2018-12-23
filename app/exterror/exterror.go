package exterror

import (
	"context"
	"reflect"

	"github.com/thoughtdealership/howto/app/frame"

	"github.com/jonbodner/multierr"
)

type ExtError struct {
	Status int
	Err    error
}

func (e ExtError) Error() string {
	return e.Err.Error()
}

func Create(status int, err error) ExtError {
	return ExtError{Status: status, Err: err}
}

func Convert(ctx context.Context, err error) ExtError {
	fr := frame.FromContext(ctx)
	switch err := err.(type) {
	case ExtError:
		fr.Logger.Debug().
			Err(err).
			Msg("no conversion necessary for ExtError")
		return err
	case multierr.Error:
		fr.Logger.Debug().
			Err(err).
			Msg("multierr conversion for ExtError")
		return ExtError{Status: convertMultiErr(err), Err: err}
	default:
		fr.Logger.Error().
			Err(err).
			Str("type", reflect.TypeOf(err).String()).
			Msg("automatic promotion to 500 response")
		return ExtError{Status: 500, Err: err}
	}
}

func allExtError(errs multierr.Error) bool {
	if len(errs) == 0 {
		return false
	}
	for _, e := range errs {
		if _, ok := e.(ExtError); !ok {
			return false
		}
	}
	return true
}

func allEqualStatus(errs multierr.Error) bool {
	resp := errs[0].(ExtError).Status
	for _, e := range errs {
		if resp != e.(ExtError).Status {
			return false
		}
	}
	return true
}

func allRangeStatus(errs multierr.Error, low int, high int) bool {
	for _, e := range errs {
		status := e.(ExtError).Status
		if (status < low) || (status >= high) {
			return false
		}
	}
	return true
}

func convertMultiErr(errs multierr.Error) int {
	if !allExtError(errs) {
		return 500
	} else if allEqualStatus(errs) {
		return errs[0].(ExtError).Status
	} else if allRangeStatus(errs, 400, 500) {
		return 400
	}
	return 500
}
