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

//If The return of the ReportFunc is true,The Error resets the internal errors
type ReportFunc func([]*ErrorItem, *Error) bool

type ErrorItem struct {
	err error
	cnt int
}

func (e *ErrorItem) Error() string {
	return e.err.Error()
}

func (e ErrorItem) Get() error {
	return e.err
}

func (e ErrorItem) Count() int {
	return e.cnt
}

type Error struct {
	errs       []*ErrorItem
	reportFunc ReportFunc
	lastReport time.Time
	mutex      sync.RWMutex
}

func New() *Error {
	e := &Error{
		errs: make([]*ErrorItem, 0),
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
	for _, _err := range e.errs {
		if err.Error() == _err.Error() {
			_err.cnt++
			ok = true
			break
		}
	}
	if !ok {
		errItem := &ErrorItem{
			err: err,
			cnt: 1,
		}
		e.errs = append(e.errs, errItem)
		e.firstReport()
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

	for _, _err := range e.errs {
		if _err.Error() == err.Error() {
			return _err.cnt
		}
	}
	return 0
}

func (e *Error) Error() string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	msgs := make([]string, len(e.errs))
	for i, err := range e.errs {
		msgs[i] = fmt.Sprintf("%s (%v)", err, err.Count())
	}

	return strings.Join(msgs, "\n")
}

func (e *Error) Errors() []error {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	errs := make([]error, len(e.errs))
	for _, err := range e.errs {
		errs = append(errs, err)
	}

	return errs
}

func (e *Error) Reset() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.reset()
}

func (e *Error) reset() {
	e.errs = e.errs[:0]
}

func (e *Error) fwd(d time.Duration) {
	tick := tick(d)
	for {
		<-tick
		e.mutex.RLock()
		cnt := len(e.errs)
		e.mutex.RUnlock()
		if cnt > 0 {
			e.reportWithReset()
		}
	}
}

func (e *Error) firstReport() {
	if time.Since(e.lastReport) >= 1*time.Minute {
		e.report(false)
	}
}

func (e *Error) reportWithReset() {
	e.report(true)
}

func (e *Error) report(reset bool) {
	if e.reportFunc == nil {
		return
	}
	if e.reportFunc(e.errs, e) == true && reset == true {
		e.Reset()
	}
	e.lastReport = time.Now()
}
