package medtronic

import (
	"time"
)

const (
	PowerControl CommandCode = 0x5D
)

func (pump *Pump) PowerControl(retries int) error {
	cmd := PumpCommand{
		Code:            PowerControl,
		NumRetries:      retries,
		ResponseTimeout: 10 * time.Second,
		ResponseHandler: func(data []byte) interface{} {
			// Return something other than nil
			return true
		},
	}
	_, err := pump.Execute(cmd)
	return err
}
