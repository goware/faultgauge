package faultgauge

import (
	"sync"
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
	mu               sync.RWMutex
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
	f.mu.Lock()
	f.sample(true)
	f.mu.Unlock()
}

func (f *FaultGauge) IncrementSuccess() {
	f.mu.Lock()
	f.sample(false)
	f.mu.Unlock()
}

func (f *FaultGauge) FailRate() (float32, float32) {
	f.mu.RLock()
	currFailRate := float32(f.currFailCount) / float32(f.currSuccessCount+f.currFailCount)
	prevFailRate := float32(f.prevFailCount) / float32(f.prevSuccessCount+f.prevFailCount)
	f.mu.RUnlock()
	return currFailRate, prevFailRate
}

func (f *FaultGauge) NumFail() (uint64, uint64) {
	f.mu.RLock()
	currFailCount := f.currFailCount
	prevFailCount := f.prevFailCount
	f.mu.RUnlock()
	return currFailCount, prevFailCount
}

func (f *FaultGauge) NumSuccess() (uint64, uint64) {
	f.mu.RLock()
	currSuccessCount := f.currSuccessCount
	prevSuccessCount := f.prevSuccessCount
	f.mu.RUnlock()
	return currSuccessCount, prevSuccessCount
}

func (f *FaultGauge) Counter() (uint64, uint64) {
	f.mu.RLock()
	currCounter := f.currCounter
	prevCounter := f.prevCounter
	f.mu.RUnlock()
	return currCounter, prevCounter
}

func (f *FaultGauge) sample(fail bool) {
	now := time.Now().UTC()
	currWindow := now.Truncate(f.windowLength)

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

	if fail {
		f.currFailCount += 1
	} else {
		f.currSuccessCount += 1
	}
	f.currCounter += 1
}
