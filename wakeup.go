package medtronic

import (
	"log"
	"time"
)

const (
	Wakeup CommandCode = 0x5D
)

func (pump *Pump) Wakeup() error {
	_, err := pump.Model(3, nil)
	if err == nil {
		log.Printf("pump is awake\n")
		return nil
	}
	log.Printf("waking pump\n")
	cmd := PumpCommand{
		Code:            Wakeup,
		ResponseTimeout: 10 * time.Second,
		NumRetries:      1,
	}
	packet := commandPacket(cmd)
	const (
		numWakeups = 200
		xmitDelay  = 35 * time.Millisecond
	)
	for i := 0; i < numWakeups; i++ {
		pump.Radio.Outgoing() <- packet
		time.Sleep(xmitDelay)
	}
	_, err = pump.Execute(cmd)
	return err
}
