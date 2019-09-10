package loader

import (
	"context"
	"sync"
	"time"
)

type Loader struct {
	sync.Mutex
	list     []int64
	wait     chan Response
	MaxBatch int
	Timeout  time.Duration
	Fn       func(baseCtx context.Context, IDs []int64) Response
}

type Response struct {
	Err   error
	Items interface{}
}

func (t *Loader) Get(ctx context.Context, id int64) chan Response {
	t.Lock()
	wait := t.wait
	t.list = append(t.list, id)
	if len(t.list) == 1 {
		wait = make(chan Response)
		t.wait = wait
		go func() {
			time.Sleep(t.Timeout)
			t.Lock()
			defer t.Unlock()
			if t.wait == wait {
				t.consume(ctx)
			}
		}()
	}
	if len(t.list) == t.MaxBatch {
		t.consume(ctx)
	}
	t.Unlock()

	return wait
}

func (t *Loader) consume(ctx context.Context) {
	go func(ids []int64, wait chan Response) {
		resp := t.Fn(ctx, ids)
		for range ids {
			wait <- resp
		}
		close(wait)
	}(append(t.list[:0:0], t.list...), t.wait)

	t.wait = make(chan Response)
	t.list = make([]int64, 0, t.MaxBatch)
}
