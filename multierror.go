package multierror

import (
	"time"
)

var (
	tick = time.Tick
)

type ReportFunc func(error, int)

type MultiError struct {
	err error
	cnt int
}

func New() *MultiError {
	e := &MultiError{}
	return e
}

func NewWithReport(d time.Duration, reportFunc ReportFunc) *MultiError {
	e := New()
	go e.fwd(d, reportFunc)
	return e
}

func (e *MultiError) Add(err error) bool {
	if err == nil {
		return true
	}

	if e.err == nil {
		e.err = err
		e.cnt++

		return true
	}

	if err != e.err {
		return false
	}

	e.cnt++

	return true
}

func (e *MultiError) Error() string {
	return e.err.Error()
}

func (e *MultiError) Count() int {
	return e.cnt
}

func (e *MultiError) fwd(d time.Duration, reportFunc ReportFunc) {
	tick := tick(d)
	for {
		<-tick
		reportFunc(e.err, e.cnt)
	}
}
