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
		Fn: func(keys []interface{}) (interface{}, error) {
			IDs := KeysToInt(keys)
			assert.Equal(t, 1, len(IDs))
			assert.Equal(t, 1, IDs[0])
			return IDs, nil
		},
	}

	value, err := l.Load(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, value.([]int)[0])
}

func TestLoaderMaxBatch(t *testing.T) {
	l := Loader{
		MaxBatch: 3,
		Timeout:  time.Millisecond,
		Fn: func(keys []interface{}) (interface{}, error) {
			IDs := KeysToInt64(keys)
			assert.Equal(t, 2, len(IDs))
			return IDs, nil
		},
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		value, err := l.Load(int64(1))
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
		value, err := l.Load(int64(2))
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
	go func() {
		value, err := l.Load(int64(2))
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

func TestLoaderMaxBatchString(t *testing.T) {
	l := Loader{
		MaxBatch: 2,
		Timeout:  time.Millisecond,
		Fn: func(keys []interface{}) (interface{}, error) {
			return KeysToString(keys), nil
		},
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		value, err := l.Load("1")
		assert.Nil(t, err)
		assert.Equal(t, 2, len(value.([]string)))
		found := false
		for _, i := range value.([]string) {
			if i == "1" {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
		wg.Done()
	}()
	go func() {
		time.Sleep(time.Millisecond * 2)

		value, err := l.Load("2")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(value.([]string)))
		found := false
		for _, i := range value.([]string) {
			if i == "2" {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
		wg.Done()
	}()
	go func() {
		value, err := l.Load("3")
		assert.Nil(t, err)
		assert.Equal(t, 2, len(value.([]string)))
		found := false
		for _, i := range value.([]string) {
			if i == "3" {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
		wg.Done()
	}()

	wg.Wait()
}
