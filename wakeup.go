package medtronic

import (
	"log"
	"time"
)

func (pump *Pump) Wakeup() error {
	_, err := pump.Model(3, nil)
	if err == nil {
		log.Printf("pump is awake\n")
		return nil
	}
	const (
		numWakeups = 150
		xmitDelay  = 35 * time.Millisecond
	)
	packet := commandPacket(PumpCommand{Code: PowerControl})
	log.Printf("waking pump\n")
	for i := 0; i < numWakeups; i++ {
		pump.Radio.Outgoing() <- packet
		time.Sleep(xmitDelay)
	}
	return pump.PowerControl(1)
}
