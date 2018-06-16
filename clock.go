package medtronic

import (
	"time"
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
	year := marshalUint16(uint16(t.Year()))
	pump.Execute(setClock,
		byte(t.Hour()),
		byte(t.Minute()),
		byte(t.Second()),
		year[0], year[1],
		byte(t.Month()),
		byte(t.Day()))
}
