package medtronic

import (
	"time"
)

const (
	TempBasal CommandCode = 0x98
)

type TempBasalInfo interface {
	tagTempBasalInfo()
}

type AbsoluteTempBasalInfo struct {
	Duration          time.Duration
	MilliUnitsPerHour int
}

func (info AbsoluteTempBasalInfo) tagTempBasalInfo() {
}

type PercentTempBasalInfo struct {
	Duration   time.Duration
	Percentage int
}

func (info PercentTempBasalInfo) tagTempBasalInfo() {
}

func (pump *Pump) TempBasal() (TempBasalInfo, error) {
	result, err := pump.Execute(TempBasal, func(data []byte) interface{} {
		if len(data) < 7 || data[0] != 6 {
			return nil
		}
		d := time.Duration(twoByteInt(data[5:7])) * time.Minute
		var info TempBasalInfo
		switch data[1] {
		case 0:
			info = AbsoluteTempBasalInfo{
				Duration:          d,
				MilliUnitsPerHour: twoByteInt(data[3:5]) * 25,
			}
		case 1:
			info = PercentTempBasalInfo{
				Duration:   d,
				Percentage: int(data[2]),
			}
		default:
			return nil
		}
		return info
	})
	if err != nil {
		return nil, err
	}
	return result.(TempBasalInfo), nil
}
