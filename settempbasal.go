package medtronic

import (
	"fmt"
	"time"
)

const (
	SetTempBasal CommandCode = 0x4C
)

func (pump *Pump) SetTempBasal(duration time.Duration, milliUnitsPerHour int) {
	const halfHour = 30 * time.Minute
	if duration%halfHour != 0 {
		pump.err = fmt.Errorf("temporary basal duration (%v) is not a multiple of 30 minutes", duration)
	}
	if milliUnitsPerHour%25 != 0 {
		pump.err = fmt.Errorf("temporary basal rate (%d) is not a multiple of 25 milliUnits per hour", milliUnitsPerHour)
	}
	pump.Execute(SetTempBasal, nil,
		0,
		byte(milliUnitsPerHour/25),
		byte(duration/halfHour))
}
