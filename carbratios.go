package medtronic

import (
	"time"
)

const (
	CarbRatios Command = 0x8A
)

type CarbRatio struct {
	Start     time.Duration // offset from 00:00:00
	CarbRatio int           // grams or exchanges covered by TEN insulin units
	Units     CarbUnitsType
}

type CarbRatioSchedule struct {
	Schedule []CarbRatio
}

func (pump *Pump) CarbRatios() CarbRatioSchedule {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	data := pump.Execute(CarbRatios)
	if pump.Error() != nil {
		return CarbRatioSchedule{}
	}
	info := []CarbRatio{}
	units := CarbUnitsType(data[1])
	if newer {
		if len(data) < 2 || (data[0]-1)%3 != 0 {
			pump.BadResponse(CarbRatios, data)
			return CarbRatioSchedule{}
		}
		n := (data[0] - 1) / 3
		i := 3
		for n != 0 {
			if data[i] == 0 && len(info) != 0 {
				break
			}
			start := scheduleToDuration(data[i])
			value := twoByteInt(data[i+1 : i+3])
			info = append(info, CarbRatio{
				Start:     start,
				CarbRatio: value,
				Units:     units,
			})
			n--
			i += 3
		}
	} else {
		if len(data) < 2 || (data[0]-1)%2 != 0 {
			pump.BadResponse(CarbRatios, data)
			return CarbRatioSchedule{}
		}
		n := (data[0] - 1) / 2
		i := 2
		for n != 0 {
			start := scheduleToDuration(data[i])
			value := int(data[i+1])
			info = append(info, CarbRatio{
				Start:     start,
				CarbRatio: value,
				Units:     units,
			})
			n--
			i += 2
		}
	}
	return CarbRatioSchedule{Schedule: info}
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
