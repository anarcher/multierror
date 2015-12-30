package multierror

import (
	"errors"
	"testing"
)

func TestZeroErrors(t *testing.T) {
	e := New()
	if e.err != nil {
		t.Errorf("want %q, have %q", nil, e.err)
	}
}

func TestZeroCountIfErrorIsNil(t *testing.T) {
	var e error
	errs := New()
	errs.Add(e)
	errs.Add(e)

	if want, have := 0, errs.Count(); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

	e1 := errors.New("e1")
	if want, have := true, errs.Add(e1); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

	e2 := errors.New("e1")
	if want, have := true, errs.Add(e2); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

}

func TestErrorCount(t *testing.T) {
	errs := New()
	e1 := errors.New("e1")
	if want, have := true, errs.Add(e1); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

	e2 := errors.New("e1")
	if want, have := true, errs.Add(e2); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

	if want, have := 2, errs.Count(); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

}
