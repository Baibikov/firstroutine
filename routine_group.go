package firstroutine

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type FirstGroup struct {
	wg 		sync.WaitGroup
	mutex 	sync.RWMutex

	mut 	*Mutable
	errs 	[]error
	_first 	chan struct{}
	_closed bool
}

func New() *FirstGroup {
	fg := FirstGroup{
		wg: sync.WaitGroup{},
		errs: make([]error, 0),
		mut: &Mutable{},
		_first: make(chan struct{}),
	}

	go fg.checkFirst()

	return &fg
}

func (fg *FirstGroup) checkFirst() {
	defer close(fg._first)
	for {
		select {
		case <-fg._first:
			fg._closed = true
			return
		}
	}
}

type Mutable struct {
	_lock bool
}

func (m *Mutable) Set(v *any, value any) {
	if !m._lock {
		*v = value
	}
}

func (fg *FirstGroup) Go(f func(m *Mutable) error) {
	fg.wg.Add(1)
	go func(errs []error) {
		defer fg.wg.Done()

		// before function load
		if fg.closed() {
			return
		}


		err := f(fg.mut)
		if err != nil {
			errs = append(errs, err)
			return
		}

		// after function load
		if fg.closed() {
			return
		}

		fg._first <- struct{}{}

		fg.mutex.Lock()
		fg.mut._lock = true
		fg.mutex.Unlock()


	}(fg.errs)
}

func (fg *FirstGroup) closed() bool {
	fg.mutex.RLock()
	ok := fg._closed
	fg.mutex.RUnlock()

	return ok
}

func (fg *FirstGroup) Wait() error {
	fg.wg.Wait()

	b := strings.Builder{}
	for i, e := range fg.errs {
		if i == 0 {
			b.WriteString(fmt.Sprintf("%v", e.Error()))
			continue
		}

		b.WriteString(fmt.Sprintf(" %v",e.Error()))
	}

	if b.String() != "" {
		return errors.New(b.String())
	}

	return nil
}