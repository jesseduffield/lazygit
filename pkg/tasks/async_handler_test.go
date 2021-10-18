package tasks

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsyncHandler(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	handler := NewAsyncHandler()
	handler.onReject = func() {
		wg.Done()
	}

	result := 0

	wg2 := sync.WaitGroup{}
	wg2.Add(1)

	handler.Do(func() func() {
		wg2.Wait()
		return func() {
			fmt.Println("setting to 1")
			result = 1
		}
	})
	handler.Do(func() func() {
		return func() {
			fmt.Println("setting to 2")
			result = 2
			wg.Done()
			wg2.Done()
		}
	})

	wg.Wait()

	assert.EqualValues(t, 2, result)
}
