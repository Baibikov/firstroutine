package firstroutine

import (
	"errors"
	"testing"
	"time"
)

func TestFirstGroup_Go(t *testing.T) {
	// don't fast test

	for i := 0; i < 10; i++ {
		var v any

		fg := New()
		fg.Go(func(m *Mutable) error {
			time.Sleep(time.Second * 2)

			m.Set(&v, "test_old")
			return nil
		})

		fg.Go(func(m *Mutable) error {
			time.Sleep(time.Second * 3)

			m.Set(&v, "test_old_1")
			return nil
		})

		fg.Go(func(m *Mutable) error {
			time.Sleep(time.Second * 1)

			m.Set(&v, "test_old_2")
			return nil
		})

		fg.Go(func(m *Mutable) error {
			m.Set(&v, "test_fast")
			return nil
		})


		err := fg.Wait()
		if err != nil {
			t.Error(err)
		}

		test, _ := v.(string)

		if test != "test_fast" {
			t.Error(errors.New("undefined state"))
		}
	}
}