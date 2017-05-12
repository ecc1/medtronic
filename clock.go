package medtronic

import (
	"time"
)

const (
	clock    Command = 0x70
	setClock Command = 0x40
)

func decodeClock(data []byte) time.Time {
	hour := int(data[1])
	min := int(data[2])
	sec := int(data[3])
	year := twoByteInt(data[4:6])
	month := time.Month(data[6])
	day := int(data[7])
	return time.Date(year, month, day, hour, min, sec, 0, time.Local)
}

// Clock returns the time according to the pump's clock.
func (pump *Pump) Clock() time.Time {
	data := pump.Execute(clock)
	if pump.Error() != nil {
		return time.Time{}
	}
	if len(data) < 8 && data[0] != 7 {
		pump.BadResponse(clock, data)
		return time.Time{}
	}
	return decodeClock(data)
}

// SetClock sets the pump's clock to the given time.
func (pump *Pump) SetClock(t time.Time) {
	pump.Execute(setClock,
		byte(t.Hour()),
		byte(t.Minute()),
		byte(t.Second()),
		byte(t.Year()>>8),
		byte(t.Year()&0xFF),
		byte(t.Month()),
		byte(t.Day()))
}
