package multierror

import (
	"errors"
	"testing"
	"time"
)

func TestZeroErrors(t *testing.T) {
	e := New()
	if want, have := 0, e.Len(); want != have {
		t.Errorf("want %v, have %v", want, have)
	}
}

func TestZeroLenIfErrorIsNil(t *testing.T) {
	var e error
	errs := New()
	errs.Add(e)
	errs.Add(e)

	if want, have := 0, errs.Len(); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

	e1 := errors.New("e1")
	errs.Add(e1)
	if want, have := 1, errs.Len(); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

	e2 := errors.New("e1")
	errs.Add(e2)
	if want, have := 1, errs.Len(); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

}

func TestErrorLen(t *testing.T) {
	errs := New()
	e1 := errors.New("e1")
	errs.Add(e1)
	if want, have := 1, errs.Len(); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

	errs.Add(errors.New("e1"))
	if want, have := 1, errs.Len(); want != have {
		t.Errorf("want %v, have %v", want, have)
	}

	if want, have := 2, errs.Count(e1); want != have {
		t.Logf("%v", errs.cnts)
		t.Errorf("want %v, have %v", want, have)
	}

}

func TestErrorReport(t *testing.T) {
	ch := make(chan time.Time)
	tick = func(time.Duration) <-chan time.Time { return ch }
	defer func() { tick = time.Tick }()

	ok := false
	reportFunc := func(e []error, cnt []int, me *Error) {
		ok = true
		t.Logf("e:%v c:%v", e, cnt)
	}
	errs := NewWithReport(time.Second, reportFunc)
	errs.Add(errors.New("ERR1"))

	if want, have := 1, errs.Len(); want != have {
		t.Errorf("want %v,have %v", want, have)
	}

	ch <- time.Now()

	if want, have := true, ok; want != have {
		t.Errorf("want %v,have %v", want, have)
	}
}
