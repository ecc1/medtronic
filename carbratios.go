package medtronic

import (
	"time"
)

const (
	CarbRatios Command = 0x8A
)

type CarbRatio struct {
	Start     time.Duration // offset from 00:00:00
	CarbRatio int           // 10x grams/unit or 100x units/exchange
	Units     CarbUnitsType
}

type CarbRatioSchedule []CarbRatio

func carbRatioStep(newerPump bool) int {
	if newerPump {
		return 3
	} else {
		return 2
	}
}

func decodeCarbRatioSchedule(data []byte, units CarbUnitsType, newerPump bool) CarbRatioSchedule {
	sched := []CarbRatio{}
	step := carbRatioStep(newerPump)
	for i := 0; i < len(data); i += step {
		start := scheduleToDuration(data[i])
		if start == 0 && len(sched) != 0 {
			break
		}
		value := 0
		if newerPump {
			value = twoByteInt(data[i+1 : i+3])
		} else {
			value = int(data[i+1])
		}
		sched = append(sched, CarbRatio{
			Start:     start,
			CarbRatio: value,
			Units:     units,
		})
	}
	return sched
}

func (pump *Pump) CarbRatios() CarbRatioSchedule {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	data := pump.Execute(CarbRatios)
	if pump.Error() != nil {
		return CarbRatioSchedule{}
	}
	if len(data) < 2 {
		pump.BadResponse(CarbRatios, data)
		return CarbRatioSchedule{}
	}
	n := int(data[0]) - 1
	step := carbRatioStep(newer)
	if n%step != 0 {
		pump.BadResponse(CarbRatios, data)
		return CarbRatioSchedule{}
	}
	units := CarbUnitsType(data[1])
	return decodeCarbRatioSchedule(data[step:step+n], units, newer)
}

func (s CarbRatioSchedule) CarbRatioAt(t time.Time) CarbRatio {
	d := sinceMidnight(t)
	last := CarbRatio{}
	for _, v := range s {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
