package medtronic

import (
	"time"
)

const (
	glucoseTargets Command = 0x9F
)

// GlucoseTarget represents an entry in a glucose target schedule.
type GlucoseTarget struct {
	Start TimeOfDay
	Low   Glucose
	High  Glucose
	Units GlucoseUnitsType
}

// GlucoseTargetSchedule represents a glucose target schedule.
type GlucoseTargetSchedule []GlucoseTarget

func decodeGlucoseTargetSchedule(data []byte, units GlucoseUnitsType) GlucoseTargetSchedule {
	var sched []GlucoseTarget
	for i := 0; i < len(data); i += 3 {
		start := halfHoursToTimeOfDay(data[i])
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

// GlucoseTargets returns the pump's glucose target schedule.
func (pump *Pump) GlucoseTargets() GlucoseTargetSchedule {
	data := pump.Execute(glucoseTargets)
	if pump.Error() != nil {
		return GlucoseTargetSchedule{}
	}
	if len(data) < 2 || (data[0]-1)%3 != 0 {
		pump.BadResponse(glucoseTargets, data)
		return GlucoseTargetSchedule{}
	}
	n := data[0] - 1
	units := GlucoseUnitsType(data[1])
	return decodeGlucoseTargetSchedule(data[2:2+n], units)
}

// GlucoseTargetAt returns the glucose target in effect at the given time.
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
