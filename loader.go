package loader

import (
	"sync"
	"time"
)

type Loader struct {
	mutex    sync.Mutex
	list     []interface{}
	wait     chan response
	MaxBatch int
	Timeout  time.Duration
	Fn       func(keys []interface{}) (interface{}, error)
}

type response struct {
	Error error
	Value interface{}
}

func (t *Loader) Load(key interface{}) (interface{}, error) {
	t.mutex.Lock()
	wait := t.wait
	t.list = append(t.list, key)
	if len(t.list) == 1 {
		wait = make(chan response)
		t.wait = wait
		go func() {
			time.Sleep(t.Timeout)
			t.mutex.Lock()
			defer t.mutex.Unlock()
			if t.wait == wait {
				t.consume()
			}
		}()
	}
	if len(t.list) == t.MaxBatch {
		t.consume()
	}
	t.mutex.Unlock()

	r := <-wait
	return r.Value, r.Error
}

func (t *Loader) consume() {
	go func(ids []interface{}, wait chan response) {
		resp := response{}
		resp.Value, resp.Error = t.Fn(ids)
		for range ids {
			wait <- resp
		}
		close(wait)
	}(t.list, t.wait)

	t.wait = make(chan response)
	t.list = make([]interface{}, 0, t.MaxBatch)
}
