package exterror

import (
	"context"
	"errors"
	"strings"
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

func TestAppend(t *testing.T) {
	ctx := frame.NewContext(context.Background())
	err := Append(ctx, errors.New(""), "baz")
	if err.Error() != "baz" {
		t.Error("Expected plain error")
	}
	prev := ExtError{Status: 404, Err: errors.New("foobar")}
	err = Append(ctx, prev, "baz")
	if prev.Status != err.(ExtError).Status {
		t.Error("Expected status 404")
	}
	if err.Error() != "baz. foobar" {
		t.Errorf("Incorrect error message: %s", err.Error())
	}
	plain := errors.New("foobar")
	err = Append(ctx, plain, "baz")
	if err.Error() != "baz. foobar" {
		t.Errorf("Incorrect error message: %s", err.Error())
	}
	var errs error
	e1 := Create(400, errors.New("foo"))
	e2 := Create(400, errors.New("bar"))
	e3 := Create(400, errors.New("baz"))
	errs = multierr.Append(e1, e2)
	errs = multierr.Append(errs, e3)
	errs = Append(ctx, errs, "prefix")
	if 400 != errs.(ExtError).Status {
		t.Error("Expected status 400")
	}
	if !strings.HasPrefix(errs.Error(), "prefix.") {
		t.Errorf("Incorrect error message: %s", errs.Error())
	}
}

func TestConvertMultiErr(t *testing.T) {
	var err error
	var merr multierr.Error
	e1 := Create(404, nil)
	e2 := Create(401, nil)
	e3 := Create(500, nil)
	out := convertMultiErr(merr)
	if out.Status != 500 {
		t.Error("Expected status 500")
	}
	err = multierr.Append(e1, e2)
	err = multierr.Append(err, e3)
	out = convertMultiErr(err.(multierr.Error))
	if out.Status != 500 {
		t.Error("Expected status 500")
	}
	err = multierr.Append(e1, e2)
	out = convertMultiErr(err.(multierr.Error))
	if out.Status != 400 {
		t.Error("Expected status 400")
	}
	err = multierr.Append(e1, e2)
	err = multierr.Append(err, errors.New("foobar"))
	out = convertMultiErr(err.(multierr.Error))
	if out.Status != 500 {
		t.Error("Expected status 500")
	}
}
