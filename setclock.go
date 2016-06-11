package medtronic

import (
	"time"
)

const (
	SetClock CommandCode = 0x40
)

func (pump *Pump) SetClock(t time.Time) {
	pump.Execute(SetClock, nil,
		byte(t.Hour()),
		byte(t.Minute()),
		byte(t.Second()),
		byte(t.Year()>>8),
		byte(t.Year()&0xFF),
		byte(t.Month()),
		byte(t.Day()))
}
