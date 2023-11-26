package faultgauge

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFaultGauge(t *testing.T) {
	windowLength := 2000 * time.Millisecond
	fg := NewFaultGauge(windowLength)

	syncTime()

	{
		fg.IncrementSuccess()
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(0), fg.currFailCount)
		assert.Equal(t, uint64(1), fg.currSuccessCount)
		assert.Equal(t, uint64(0), fg.prevFailCount)
		assert.Equal(t, uint64(0), fg.prevSuccessCount)
	}

	{
		fg.IncrementSuccess()
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(0), fg.currFailCount)
		assert.Equal(t, uint64(2), fg.currSuccessCount)
		assert.Equal(t, uint64(0), fg.prevFailCount)
		assert.Equal(t, uint64(0), fg.prevSuccessCount)
	}

	{
		fg.IncrementFail() // !!!
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(1), fg.currFailCount)
		assert.Equal(t, uint64(2), fg.currSuccessCount)
		assert.Equal(t, uint64(0), fg.prevFailCount)
		assert.Equal(t, uint64(0), fg.prevSuccessCount)
	}

	{
		fg.IncrementFail() // !!!
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(2), fg.currFailCount)
		assert.Equal(t, uint64(2), fg.currSuccessCount)
		assert.Equal(t, uint64(0), fg.prevFailCount)
		assert.Equal(t, uint64(0), fg.prevSuccessCount)
	}

	{
		fg.IncrementFail() // !!!
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(3), fg.currFailCount)
		assert.Equal(t, uint64(2), fg.currSuccessCount)
		assert.Equal(t, uint64(0), fg.prevFailCount)
		assert.Equal(t, uint64(0), fg.prevSuccessCount)
	}

	{
		fg.IncrementSuccess()
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(3), fg.currFailCount)
		assert.Equal(t, uint64(3), fg.currSuccessCount)
		assert.Equal(t, uint64(0), fg.prevFailCount)
		assert.Equal(t, uint64(0), fg.prevSuccessCount)
	}

	//-- new window, at 4 seconds
	time.Sleep(windowLength + (1 * time.Second))
	syncTime()

	{
		fg.IncrementSuccess()
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(0), fg.currFailCount)
		assert.Equal(t, uint64(1), fg.currSuccessCount)
		assert.Equal(t, uint64(3), fg.prevFailCount)
		assert.Equal(t, uint64(3), fg.prevSuccessCount)
	}

	// os.Exit(1)
	{
		fg.IncrementSuccess()
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(0), fg.currFailCount)
		assert.Equal(t, uint64(2), fg.currSuccessCount)
		assert.Equal(t, uint64(3), fg.prevFailCount)
		assert.Equal(t, uint64(3), fg.prevSuccessCount)
	}

	{
		fg.IncrementFail() // !!!
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(1), fg.currFailCount)
		assert.Equal(t, uint64(2), fg.currSuccessCount)
		assert.Equal(t, uint64(3), fg.prevFailCount)
		assert.Equal(t, uint64(3), fg.prevSuccessCount)
	}

	{
		fg.IncrementFail() // !!!
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(2), fg.currFailCount)
		assert.Equal(t, uint64(2), fg.currSuccessCount)
		assert.Equal(t, uint64(3), fg.prevFailCount)
		assert.Equal(t, uint64(3), fg.prevSuccessCount)
	}

	{
		fg.IncrementFail() // !!!
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(3), fg.currFailCount)
		assert.Equal(t, uint64(2), fg.currSuccessCount)
		assert.Equal(t, uint64(3), fg.prevFailCount)
		assert.Equal(t, uint64(3), fg.prevSuccessCount)
	}

	{
		fg.IncrementSuccess()
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, uint64(3), fg.currFailCount)
		assert.Equal(t, uint64(3), fg.currSuccessCount)
		assert.Equal(t, uint64(3), fg.prevFailCount)
		assert.Equal(t, uint64(3), fg.prevSuccessCount)
	}

	currFailRate, prevFailRate := fg.FailRate()

	assert.Equal(t, float32(0.5), currFailRate)
	assert.Equal(t, float32(0.5), prevFailRate)
}

// sync time to the nearest second
func syncTime() {
	t := time.Now()
	curr := t.Truncate(1 * time.Second).Add(1 * time.Second)
	time.Sleep(curr.Sub(t))
}
