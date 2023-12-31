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

	// Stats for current and previous windows respectively.
	NumFail() (uint64, uint64)
	NumSuccess() (uint64, uint64)
	Counter() (uint64, uint64)
}

type FaultGauge struct {
	windowLength     time.Duration
	currWindow       time.Time
	prevWindow       time.Time
	currFailCount    uint64
	prevFailCount    uint64
	currSuccessCount uint64
	prevSuccessCount uint64
	currCounter      uint64
	prevCounter      uint64
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
	f.mu.Lock()
	defer f.mu.Unlock()
	currFailRate := float32(f.currFailCount) / float32(f.currSuccessCount+f.currFailCount)
	prevFailRate := float32(f.prevFailCount) / float32(f.prevSuccessCount+f.prevFailCount)
	return currFailRate, prevFailRate
}

func (f *FaultGauge) NumFail() (uint64, uint64) {
	return atomic.LoadUint64(&f.currFailCount), atomic.LoadUint64(&f.prevFailCount)
}

func (f *FaultGauge) NumSuccess() (uint64, uint64) {
	return atomic.LoadUint64(&f.currSuccessCount), atomic.LoadUint64(&f.prevSuccessCount)
}

func (f *FaultGauge) Counter() (uint64, uint64) {
	return atomic.LoadUint64(&f.currCounter), atomic.LoadUint64(&f.prevCounter)
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
		f.prevCounter = f.currCounter
		f.currCounter = 0
	}
	f.mu.Unlock()

	if fail {
		atomic.AddUint64(&f.currFailCount, 1)
	} else {
		atomic.AddUint64(&f.currSuccessCount, 1)
	}
	atomic.AddUint64(&f.currCounter, 1)
}
