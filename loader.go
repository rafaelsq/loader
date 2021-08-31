package loader

import (
	"sync"
	"time"
)

type Loader struct {
	sync.Mutex
	list     []int64
	wait     chan response
	MaxBatch int
	Timeout  time.Duration
	Fn       func(IDs []int64) (interface{}, error)
}

type response struct {
	Error error
	Value interface{}
}

func (t *Loader) Get(id int64) (interface{}, error) {
	t.Lock()
	wait := t.wait
	t.list = append(t.list, id)
	if len(t.list) == 1 {
		wait = make(chan response)
		t.wait = wait
		go func() {
			time.Sleep(t.Timeout)
			t.Lock()
			defer t.Unlock()
			if t.wait == wait {
				t.consume()
			}
		}()
	}
	if len(t.list) == t.MaxBatch {
		t.consume()
	}
	t.Unlock()

	r := <-wait
	return r.Value, r.Error
}

func (t *Loader) consume() {
	go func(ids []int64, wait chan response) {
		resp := response{}
		resp.Value, resp.Error = t.Fn(ids)
		for range ids {
			wait <- resp
		}
		close(wait)
	}(t.list, t.wait)

	t.wait = make(chan response)
	t.list = make([]int64, 0, t.MaxBatch)
}
