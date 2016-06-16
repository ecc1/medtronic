package medtronic

import (
	"time"
)

const (
	InsulinSensitivities CommandCode = 0x8B
)

type InsulinSensitivity struct {
	Start       time.Duration // offset from 00:00:00
	Sensitivity int           // mg/dL or μmol/L reduction per insulin unit
}

type InsulinSensitivitySchedule struct {
	Schedule []InsulinSensitivity
}

func (pump *Pump) InsulinSensitivities() InsulinSensitivitySchedule {
	result := pump.Execute(InsulinSensitivities, func(data []byte) interface{} {
		if len(data) < 2 || (data[0]-1)%2 != 0 {
			return nil
		}
		n := (data[0] - 1) / 2
		i := 2
		units := GlucoseUnitsInfo(data[1])
		info := []InsulinSensitivity{}
		for n != 0 {
			start := scheduleToDuration(data[i])
			value := int(data[i+1])
			if units == MmolPerLiter {
				// Convert to μmol/L
				value *= 100
			}
			info = append(info, InsulinSensitivity{
				Start:       start,
				Sensitivity: value,
			})
			n--
			i += 2
		}
		return InsulinSensitivitySchedule{Schedule: info}
	})
	if pump.Error() != nil {
		return InsulinSensitivitySchedule{}
	}
	return result.(InsulinSensitivitySchedule)
}

func (s InsulinSensitivitySchedule) InsulinSensitivityAt(t time.Time) InsulinSensitivity {
	d := sinceMidnight(t)
	last := InsulinSensitivity{}
	for _, v := range s.Schedule {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
