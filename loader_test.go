package loader

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoader(t *testing.T) {
	l := Loader{
		MaxBatch: 2,
		Timeout:  time.Millisecond,
		Fn: func(IDs []int64) (interface{}, error) {
			assert.Equal(t, 1, len(IDs))
			assert.Equal(t, int64(1), IDs[0])
			return IDs, nil
		},
	}

	value, err := l.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value.([]int64)[0])
}

func TestLoaderMaxBatch(t *testing.T) {
	l := Loader{
		MaxBatch: 2,
		Timeout:  time.Millisecond,
		Fn: func(IDs []int64) (interface{}, error) {
			assert.Equal(t, 2, len(IDs))
			return IDs, nil
		},
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		value, err := l.Get(1)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(value.([]int64)))
		found := false
		for _, i := range value.([]int64) {
			if i == 1 {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
		wg.Done()
	}()
	go func() {
		value, err := l.Get(2)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(value.([]int64)))
		found := false
		for _, i := range value.([]int64) {
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
