package medtronic

import (
	"time"
)

const (
	InsulinSensitivities Command = 0x8B
)

type InsulinSensitivity struct {
	Start       TimeOfDay
	Sensitivity Glucose // glucose reduction per insulin unit
	Units       GlucoseUnitsType
}

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

func (pump *Pump) InsulinSensitivities() InsulinSensitivitySchedule {
	data := pump.Execute(InsulinSensitivities)
	if pump.Error() != nil {
		return InsulinSensitivitySchedule{}
	}
	if len(data) < 2 || (data[0]-1)%2 != 0 {
		pump.BadResponse(InsulinSensitivities, data)
		return InsulinSensitivitySchedule{}
	}
	n := int(data[0]) - 1
	units := GlucoseUnitsType(data[1])
	return decodeInsulinSensitivitySchedule(data[2:2+n], units)
}

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
