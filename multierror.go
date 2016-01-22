package multierror

import (
	"time"
)

var (
	tick = time.Tick
)

type ReportFunc func(error, int, *MultiError)

type MultiError struct {
	err        error
	cnt        int
	reportFunc ReportFunc
}

func New() *MultiError {
	e := &MultiError{}
	return e
}

func NewWithReport(d time.Duration, reportFunc ReportFunc) *MultiError {
	e := New()
	e.reportFunc = reportFunc
	go e.fwd(d)
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

	if err.Error() != e.err.Error() {
		if e.reportFunc != nil {
			e.reportFunc(e.err, e.cnt, e)
			e.err = err
			e.cnt = 1
		}
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

func (e *MultiError) Reset() {
	e.cnt = 0
}

func (e *MultiError) fwd(d time.Duration) {
	tick := tick(d)
	for {
		<-tick
		if e.cnt > 0 {
			e.reportFunc(e.err, e.cnt, e)
		}
	}
}
