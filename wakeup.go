package medtronic

import (
	"log"
	"time"
)

const (
	wakeup Command = 0x5D

	// Older pumps should have RF enabled to increase the
	// frequency with which they listen for wakeups.
	numWakeups    = 100
	xmitDelay     = 10 * time.Millisecond
	wakeupTimeout = 10 * time.Second
)

// Wakeup wakes up the pump.
// It first attempts a model command, which will succeed quickly if
// the pump is already awake.  If that times out, it will repeatedly
// send wakeup commands.
func (pump *Pump) Wakeup() {
	pump.Model()
	if pump.Error() == nil {
		return
	}
	if !pump.NoResponse() {
		return
	}
	pump.SetError(nil)
	log.Printf("waking pump")
	n := pump.Retries()
	defer pump.SetRetries(n)
	t := pump.Timeout()
	defer pump.SetTimeout(t)
	pump.SetRetries(numWakeups)
	pump.SetTimeout(xmitDelay)
	pump.Execute(wakeup)
	pump.SetError(nil)
	pump.SetRetries(1)
	pump.SetTimeout(wakeupTimeout)
	pump.Execute(wakeup)
}
