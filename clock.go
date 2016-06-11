package medtronic

import (
	"time"
)

const (
	Clock CommandCode = 0x70
)

func (pump *Pump) Clock() time.Time {
	result := pump.Execute(Clock, func(data []byte) interface{} {
		if len(data) < 8 && data[0] != 7 {
			return nil
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
	})
	if pump.Error() != nil {
		return time.Time{}
	}
	return result.(time.Time)
}
