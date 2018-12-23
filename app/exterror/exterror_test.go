package exterror

import (
	"context"
	"errors"
	"testing"

	"github.com/thoughtdealership/howto/app/frame"

	"github.com/jonbodner/multierr"
)

type foobar struct{}

func (f foobar) Error() string {
	return "foobar"
}

func TestConvert(t *testing.T) {
	i := ExtError{Status: 404, Err: errors.New("foobar")}
	ctx := frame.NewContext(context.Background())
	o := Convert(ctx, i)
	if o.Status != i.Status {
		t.Error("Expected status 404")
	}
	if o.Err.Error() != i.Err.Error() {
		t.Error("Incorrect error message")
	}
	o = Convert(ctx, foobar{})
	if o.Status != 500 {
		t.Error("Expected status 500")
	}
	if o.Err.Error() != "foobar" {
		t.Error("Incorrect error message")
	}
}

func TestConvertMultiErr(t *testing.T) {
	var err error
	var merr multierr.Error
	e1 := Create(404, nil)
	e2 := Create(401, nil)
	e3 := Create(500, nil)
	out := convertMultiErr(merr)
	if out != 500 {
		t.Error("Expected status 500")
	}
	err = multierr.Append(e1, e2)
	err = multierr.Append(err, e3)
	out = convertMultiErr(err.(multierr.Error))
	if out != 500 {
		t.Error("Expected status 500")
	}
	err = multierr.Append(e1, e2)
	out = convertMultiErr(err.(multierr.Error))
	if out != 400 {
		t.Error("Expected status 400")
	}
	err = multierr.Append(e1, e2)
	err = multierr.Append(err, errors.New("foobar"))
	out = convertMultiErr(err.(multierr.Error))
	if out != 500 {
		t.Error("Expected status 500")
	}
}
