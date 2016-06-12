package medtronic

import (
	"log"
	"time"
)

const (
	Wakeup CommandCode = 0x5D
)

func (pump *Pump) Wakeup() {
	pump.Model()
	if pump.Error() == nil {
		return
	}
	pump.SetError(nil)
	log.Printf("waking pump")
	const (
		// Older pumps should have RF enabled to increase the
		// frequency with which they listen for wakeups.
		numWakeups = 100
		xmitDelay  = 10 * time.Millisecond
	)
	packet := commandPacket(Wakeup, nil)
	for i := 0; i < numWakeups; i++ {
		pump.Radio.Send(packet)
		time.Sleep(xmitDelay)
	}
	n := pump.Retries()
	pump.SetRetries(1)
	defer pump.SetRetries(n)
	t := pump.Timeout()
	pump.SetTimeout(10 * time.Second)
	defer pump.SetTimeout(t)
	pump.Execute(Wakeup, nil)
}
