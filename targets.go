package medtronic

import (
	"time"
)

const (
	GlucoseTargets Command = 0x9F
)

type GlucoseTarget struct {
	Start time.Duration // offset from 00:00:00
	Low   Glucose
	High  Glucose
	Units GlucoseUnitsType
}

type GlucoseTargetSchedule []GlucoseTarget

func decodeGlucoseTargetSchedule(data []byte, units GlucoseUnitsType) GlucoseTargetSchedule {
	sched := []GlucoseTarget{}
	for i := 0; i < len(data); i += 3 {
		start := scheduleToDuration(data[i])
		if start == 0 && len(sched) != 0 {
			break
		}
		sched = append(sched, GlucoseTarget{
			Start: start,
			Low:   byteToGlucose(data[i+1], units),
			High:  byteToGlucose(data[i+2], units),
			Units: units,
		})
	}
	return sched
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
	n := data[0] - 1
	units := GlucoseUnitsType(data[1])
	return decodeGlucoseTargetSchedule(data[2:2+n], units)
}

func (s GlucoseTargetSchedule) GlucoseTargetAt(t time.Time) GlucoseTarget {
	d := sinceMidnight(t)
	last := GlucoseTarget{}
	for _, v := range s {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
