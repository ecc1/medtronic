package medtronic

import (
	"time"
)

const (
	SetClock CommandCode = 0x40
)

func (pump *Pump) SetClock(t time.Time, retries int) error {
	cmd := PumpCommand{
		Code: SetClock,
		Params: []byte{
			byte(t.Hour()),
			byte(t.Minute()),
			byte(t.Second()),
			byte(t.Year() >> 8),
			byte(t.Year() & 0xFF),
			byte(t.Month()),
			byte(t.Day()),
		},
		NumRetries: retries,
	}
	_, err := pump.Execute(cmd)
	if err != nil {
		return err
	}
	return nil
}
