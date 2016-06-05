package medtronic

import (
	"log"
	"time"
)

const (
	Wakeup CommandCode = 0x5D
)

func (pump *Pump) Wakeup() error {
	_, err := pump.Model()
	if err == nil {
		return nil
	}
	log.Printf("waking pump\n")
	const (
		numWakeups = 200
		xmitDelay  = 35 * time.Millisecond
	)
	packet := commandPacket(Wakeup, nil)
	for i := 0; i < numWakeups; i++ {
		pump.Radio.Outgoing() <- packet
		time.Sleep(xmitDelay)
	}
	n := pump.SetRetries(1)
	defer pump.SetRetries(n)
	t := pump.SetTimeout(10 * time.Second)
	defer pump.SetTimeout(t)
	_, err = pump.Execute(Wakeup, nil)
	return err
}
