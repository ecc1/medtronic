package medtronic

import (
	"time"
)

const (
	Clock    Command = 0x70
	SetClock Command = 0x40
)

func (pump *Pump) Clock() time.Time {
	data := pump.Execute(Clock)
	if pump.Error() != nil {
		return time.Time{}
	}
	if len(data) < 8 && data[0] != 7 {
		pump.BadResponse(Clock, data)
		return time.Time{}
	}
	return time.Date(
		twoByteInt(data[4:6]), // year
		time.Month(data[6]),   // month
		int(data[7]),          // day
		int(data[1]),          // hour
		int(data[2]),          // min
		int(data[3]),          // sec
		0,                     // nsec
		time.Local)
}

func (pump *Pump) SetClock(t time.Time) {
	pump.Execute(SetClock,
		byte(t.Hour()),
		byte(t.Minute()),
		byte(t.Second()),
		byte(t.Year()>>8),
		byte(t.Year()&0xFF),
		byte(t.Month()),
		byte(t.Day()))
}
