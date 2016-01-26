package multierror

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	tick = time.Tick
)

type ReportFunc func([]error, []int, *Error)

type Error struct {
	errs       []error
	cnts       []int
	reportFunc ReportFunc
	mutex      sync.RWMutex
}

func New() *Error {
	e := &Error{
		errs: make([]error, 0),
		cnts: make([]int, 0),
	}
	return e
}

func NewWithReport(d time.Duration, reportFunc ReportFunc) *Error {
	e := New()
	e.reportFunc = reportFunc
	go e.fwd(d)
	return e
}

func (e *Error) Add(err error) {
	if err == nil {
		return
	}

	e.mutex.Lock()
	ok := false
	for i, _err := range e.errs {
		if err.Error() == _err.Error() {
			e.cnts[i]++
			ok = true
			break
		}
	}
	if !ok {
		e.errs = append(e.errs, err)
		e.cnts = append(e.cnts, 1)

		if e.reportFunc != nil {
			e.reportFunc(e.errs, e.cnts, e)
		}
	}
	e.mutex.Unlock()
}

func (e *Error) Len() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return len(e.errs)
}

func (e *Error) Count(err error) int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	for i, _err := range e.errs {
		if _err.Error() == err.Error() && i < len(e.cnts) {
			return e.cnts[i]
		}
	}
	return 0
}

func (e *Error) Error() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	msgs := make([]string, len(e.errs))
	for i, err := range e.errs {
		msgs[i] = fmt.Sprintf("%s", err)
	}

	return strings.Join(msgs, "\n")
}

func (e *Error) Errors() []error {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return e.errs
}

func (e *Error) Reset() {
	e.mutex.Lock()
	e.errs = e.errs[:0]
	e.cnts = e.cnts[:0]
	e.mutex.Unlock()
}

func (e *Error) fwd(d time.Duration) {
	tick := tick(d)
	for {
		<-tick
		e.mutex.RLock()
		if len(e.errs) > 0 {
			e.reportFunc(e.errs, e.cnts, e)
		}
		e.mutex.RUnlock()
	}
}
