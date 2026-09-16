package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRequestQueueHandsBackWhatNoTaskWillServe(t *testing.T) {
	queue := newReadRequestQueue()

	// Nothing has begun serving, so the request comes straight back to its caller.
	assert.False(t, queue.enqueue(LinesToRead{Total: 1}))

	queue.beginServing()
	assert.True(t, queue.enqueue(LinesToRead{Total: 1}))
	assert.True(t, queue.enqueue(LinesToRead{Total: 2}))

	request, ok := queue.dequeue()
	assert.True(t, ok)
	assert.Equal(t, 1, request.Total)

	// What the task never got to comes back when it stops, and nothing is taken
	// from a caller after that.
	unanswered := queue.stopServing()
	assert.Len(t, unanswered, 1)
	assert.Equal(t, 2, unanswered[0].Total)

	assert.False(t, queue.enqueue(LinesToRead{Total: 3}))
	_, ok = queue.dequeue()
	assert.False(t, ok)
}

func TestReadRequestQueueRingsTheDoorbell(t *testing.T) {
	queue := newReadRequestQueue()
	doorbell := queue.beginServing()

	select {
	case <-doorbell:
		assert.Fail(t, "the doorbell rang before anything was queued")
	default:
	}

	// A burst leaves one token, which stands for everything queued.
	queue.enqueue(LinesToRead{Total: 1})
	queue.enqueue(LinesToRead{Total: 2})

	select {
	case <-doorbell:
	default:
		assert.Fail(t, "the doorbell didn't ring")
	}
	select {
	case <-doorbell:
		assert.Fail(t, "the doorbell rang twice for one wake")
	default:
	}
}
