package loader

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoader(t *testing.T) {
	l := Loader{
		MaxBatch: 2,
		Timeout:  time.Millisecond,
		Fn: func(ctx context.Context, IDs []int64) Response {
			assert.Equal(t, 1, len(IDs))
			assert.Equal(t, int64(1), IDs[0])
			return Response{
				Items: IDs,
				Err:   nil,
			}
		},
	}

	resp := <-l.Get(context.TODO(), 1)
	assert.Nil(t, resp.Err)
	assert.Equal(t, int64(1), resp.Items.([]int64)[0])
}

func TestLoaderMaxBatch(t *testing.T) {
	l := Loader{
		MaxBatch: 2,
		Timeout:  time.Millisecond,
		Fn: func(ctx context.Context, IDs []int64) Response {
			assert.Equal(t, 2, len(IDs))
			return Response{
				Items: IDs,
				Err:   nil,
			}
		},
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		resp := <-l.Get(context.TODO(), 1)
		assert.Nil(t, resp.Err)
		assert.Equal(t, 2, len(resp.Items.([]int64)))
		found := false
		for _, i := range resp.Items.([]int64) {
			if i == 1 {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
		wg.Done()
	}()
	go func() {
		resp := <-l.Get(context.TODO(), 2)
		assert.Nil(t, resp.Err)
		assert.Equal(t, 2, len(resp.Items.([]int64)))
		found := false
		for _, i := range resp.Items.([]int64) {
			if i == 2 {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
		wg.Done()
	}()

	wg.Wait()
}
