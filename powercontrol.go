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
		ResponseHandler: emptyResponseHandler,
	}
	_, err := pump.Execute(cmd)
	return err
}
