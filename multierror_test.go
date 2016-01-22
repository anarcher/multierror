package multierror

import (
	"errors"
	"testing"
	"time"
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

func TestErrorReport(t *testing.T) {
	ch := make(chan time.Time)
	tick = func(time.Duration) <-chan time.Time { return ch }
	defer func() { tick = time.Tick }()

	ok := false
	reportFunc := func(e error, cnt int, me *MultiError) {
		ok = true
		t.Logf("e:%v c:%v", e, cnt)
	}
	errs := NewWithReport(time.Second, reportFunc)

	if errs.Add(errors.New("ERR1")) == false {
		t.Errorf("want true,have false")
	}

	ch <- time.Now()

	if want, have := true, ok; want != have {
		t.Errorf("want %v,have %v", want, have)
	}
}
