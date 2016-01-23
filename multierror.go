package multierror

import (
	"sync"
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
	mutex      sync.RWMutex
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
		e.mutex.Lock()
		e.err = err
		e.cnt++
		e.mutex.Unlock()
		if e.reportFunc != nil {
			e.reportFunc(e.err, e.cnt, e)
		}
		return true
	}

	if err.Error() != e.err.Error() {
		if e.reportFunc != nil {
			e.reportFunc(e.err, e.cnt, e)
			e.mutex.Lock()
			e.err = err
			e.cnt = 1
			e.mutex.Unlock()
		}
		return false
	}

	e.mutex.Lock()
	e.cnt++
	e.mutex.Unlock()

	return true
}

func (e *MultiError) Error() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.err.Error()
}

func (e *MultiError) Count() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return e.cnt
}

func (e *MultiError) Reset() {
	e.mutex.Lock()
	e.cnt = 0
	e.mutex.Unlock()
}

func (e *MultiError) fwd(d time.Duration) {
	tick := tick(d)
	for {
		<-tick
		e.mutex.RLock()
		if e.cnt > 0 {
			e.reportFunc(e.err, e.cnt, e)
		}
		e.mutex.RUnlock()
	}
}
