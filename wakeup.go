package medtronic

import (
	"log"
	"time"
)

const (
	Wakeup CommandCode = 0x5D

	// Older pumps should have RF enabled to increase the
	// frequency with which they listen for wakeups.
	numWakeups = 100
	xmitDelay  = 10 * time.Millisecond
)

func (pump *Pump) Wakeup() {
	pump.Model()
	if pump.Error() == nil {
		return
	}
	_, noResponse := pump.Error().(NoResponseError)
	if !noResponse {
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
	pump.Execute(Wakeup)
	pump.SetError(nil)
	pump.SetTimeout(10 * time.Second)
	pump.Execute(Wakeup)
}
