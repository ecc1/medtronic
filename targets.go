package medtronic

import (
	"time"
)

const (
	GlucoseTargets CommandCode = 0x9F
)

type GlucoseTarget struct {
	Start time.Duration // offset from 00:00:00
	Low   int           // mg/dL or μmol/L
	High  int           // mg/dL or μmol/L
}

type GlucoseTargetSchedule struct {
	Schedule []GlucoseTarget
}

func (pump *Pump) GlucoseTargets() GlucoseTargetSchedule {
	result := pump.Execute(GlucoseTargets, func(data []byte) interface{} {
		if len(data) < 2 || (data[0]-1)%3 != 0 {
			return nil
		}
		n := (data[0] - 1) / 3
		i := 2
		units := GlucoseUnitsInfo(data[1])
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
			})
			n--
			i += 3
		}
		return GlucoseTargetSchedule{Schedule: info}
	})
	if pump.Error() != nil {
		return GlucoseTargetSchedule{}
	}
	return result.(GlucoseTargetSchedule)
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
