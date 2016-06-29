package medtronic

import (
	"time"
)

const (
	GlucoseTargets Command = 0x9F
)

type GlucoseTarget struct {
	Start time.Duration // offset from 00:00:00
	Low   int           // mg/dL or μmol/L
	High  int           // mg/dL or μmol/L
	Units GlucoseUnitsType
}

type GlucoseTargetSchedule struct {
	Schedule []GlucoseTarget
}

func (pump *Pump) GlucoseTargets() GlucoseTargetSchedule {
	data := pump.Execute(GlucoseTargets)
	if pump.Error() != nil {
		return GlucoseTargetSchedule{}
	}
	if len(data) < 2 || (data[0]-1)%3 != 0 {
		pump.BadResponse(GlucoseTargets, data)
		return GlucoseTargetSchedule{}
	}
	n := (data[0] - 1) / 3
	i := 2
	units := GlucoseUnitsType(data[1])
	info := []GlucoseTarget{}
	for n != 0 {
		start := scheduleToDuration(data[i])
		low := int(data[i+1])
		high := int(data[i+2])
		if units == MmolPerLiter {
			// Convert to μmol/L
			low *= 100
			high *= 100
		}
		info = append(info, GlucoseTarget{
			Start: start,
			Low:   low,
			High:  high,
			Units: units,
		})
		n--
		i += 3
	}
	return GlucoseTargetSchedule{Schedule: info}
}

func (s GlucoseTargetSchedule) GlucoseTargetAt(t time.Time) GlucoseTarget {
	d := sinceMidnight(t)
	last := GlucoseTarget{}
	for _, v := range s.Schedule {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
