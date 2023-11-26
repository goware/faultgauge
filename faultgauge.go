package faultgauge

import (
	"sync"
	"sync/atomic"
	"time"
)

// FailRate returns the current fail rate and the previous fail rate.
type FailRate interface {
	FailRate() (float32, float32)
}

// Controller controls the gauge by incrementing the success or fail count.
type Controller interface {
	IncrementFail()
	IncrementSuccess()

	NumFail() uint64
	NumSuccess() uint64
}

type FaultGauge struct {
	windowLength     time.Duration
	currWindow       time.Time
	prevWindow       time.Time
	currFailCount    uint64
	prevFailCount    uint64
	currSuccessCount uint64
	prevSuccessCount uint64
	mu               sync.Mutex
}

// NewFaultGauge creates a new gauge which samples the fail rate over the
// window length. Be sure to call IncrementFail() or IncrementSuccess() to
// inform the gauge of the success or failure. A windowLength of 10 seconds
// or more is recommended.
func NewFaultGauge(windowLength time.Duration) *FaultGauge {
	return &FaultGauge{
		windowLength: windowLength,
	}
}

var _ FailRate = (*FaultGauge)(nil)

var _ Controller = (*FaultGauge)(nil)

func (f *FaultGauge) IncrementFail() {
	f.sample(true)
}

func (f *FaultGauge) IncrementSuccess() {
	f.sample(false)
}

func (f *FaultGauge) FailRate() (float32, float32) {
	currFailRate := float32(f.currFailCount) / float32(f.currSuccessCount+f.currFailCount)
	prevFailRate := float32(f.prevFailCount) / float32(f.prevSuccessCount+f.prevFailCount)
	return currFailRate, prevFailRate
}

func (f *FaultGauge) NumFail() uint64 {
	return f.currFailCount
}

func (f *FaultGauge) NumSuccess() uint64 {
	return f.currSuccessCount
}

func (f *FaultGauge) sample(fail bool) {
	now := time.Now().UTC()
	currWindow := now.Truncate(f.windowLength)

	f.mu.Lock()
	if f.currWindow != currWindow {
		f.prevWindow = f.currWindow
		f.currWindow = currWindow
		f.prevFailCount = f.currFailCount
		f.currFailCount = 0
		f.prevSuccessCount = f.currSuccessCount
		f.currSuccessCount = 0
	}
	f.mu.Unlock()

	if fail {
		atomic.AddUint64(&f.currFailCount, 1)
	} else {
		atomic.AddUint64(&f.currSuccessCount, 1)
	}
}
