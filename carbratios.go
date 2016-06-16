package medtronic

import (
	"time"
)

const (
	CarbRatios CommandCode = 0x8A
)

type CarbRatio struct {
	Start time.Duration // offset from 00:00:00
	Carbs int           // grams or exchanges covered by one insulin unit
}

type CarbRatioSchedule struct {
	Schedule []CarbRatio
}

func (pump *Pump) CarbRatios() CarbRatioSchedule {
	result := pump.Execute(CarbRatios, func(data []byte) interface{} {
		if len(data) < 2 || (data[0]-1)%2 != 0 {
			return nil
		}
		n := (data[0] - 1) / 2
		i := 2
		// units := CarbUnitsInfo(data[1])
		info := []CarbRatio{}
		for n != 0 {
			start := scheduleToDuration(data[i])
			value := int(data[i+1])
			info = append(info, CarbRatio{
				Start: start,
				Carbs: value,
			})
			n--
			i += 2
		}
		return CarbRatioSchedule{Schedule: info}
	})
	if pump.Error() != nil {
		return CarbRatioSchedule{}
	}
	return result.(CarbRatioSchedule)
}

func (s CarbRatioSchedule) CarbRatioAt(t time.Time) CarbRatio {
	d := sinceMidnight(t)
	last := CarbRatio{}
	for _, v := range s.Schedule {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
