package medtronic

import (
	"time"
)

const (
	insulinSensitivities Command = 0x8B
)

// InsulinSensitivity represents an entry in an insulin sensitivity schedule.
type InsulinSensitivity struct {
	Start       TimeOfDay
	Sensitivity Glucose // glucose reduction per insulin unit
	Units       GlucoseUnitsType
}

// InsulinSensitivitySchedule represents an insulin sensitivity schedule.
type InsulinSensitivitySchedule []InsulinSensitivity

func decodeInsulinSensitivitySchedule(data []byte, units GlucoseUnitsType) InsulinSensitivitySchedule {
	var sched []InsulinSensitivity
	for i := 0; i < len(data); i += 2 {
		n := data[i]
		start := halfHoursToTimeOfDay(n & 0x1F)
		if start == 0 && len(sched) != 0 {
			break
		}
		s := int((n>>6)&0x1)<<8 | int(data[i+1])
		sched = append(sched, InsulinSensitivity{
			Start:       start,
			Sensitivity: intToGlucose(s, units),
			Units:       units,
		})
	}
	return sched
}

// InsulinSensitivities returns the pump's insulin sensitivity schedule.
func (pump *Pump) InsulinSensitivities() InsulinSensitivitySchedule {
	data := pump.Execute(insulinSensitivities)
	if pump.Error() != nil {
		return InsulinSensitivitySchedule{}
	}
	if len(data) < 2 || (data[0]-1)%2 != 0 {
		pump.BadResponse(insulinSensitivities, data)
		return InsulinSensitivitySchedule{}
	}
	n := int(data[0]) - 1
	units := GlucoseUnitsType(data[1])
	return decodeInsulinSensitivitySchedule(data[2:2+n], units)
}

// InsulinSensitivityAt returns the insulin sensitivity in effect at the given time.
func (s InsulinSensitivitySchedule) InsulinSensitivityAt(t time.Time) InsulinSensitivity {
	d := sinceMidnight(t)
	last := InsulinSensitivity{}
	for _, v := range s {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
