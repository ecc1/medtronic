package medtronic

import (
	"log"
	"time"
)

const (
	SetTempBasal CommandCode = 0x4C
)

func (pump *Pump) SetTempBasal(duration time.Duration, milliUnitsPerHour int) error {
	const halfHour = 30 * time.Minute
	if duration%halfHour != 0 {
		log.Panicf("temporary basal duration (%v) is not a multiple of 30 minutes\n", duration)
	}
	if milliUnitsPerHour%25 != 0 {
		log.Panicf("temporary basal rate (%d) is not a multiple of 25 milliUnits per hour\n", milliUnitsPerHour)
	}
	_, err := pump.Execute(SetTempBasal, nil,
		0,
		byte(milliUnitsPerHour/25),
		byte(duration/halfHour))
	return err
}
