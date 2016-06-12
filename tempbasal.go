package medtronic

import (
	"fmt"
	"time"
)

const (
	TempBasal    CommandCode = 0x98
	SetTempBasal CommandCode = 0x4C
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

func (pump *Pump) TempBasal() TempBasalInfo {
	result := pump.Execute(TempBasal, func(data []byte) interface{} {
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
	if pump.Error() != nil {
		return nil
	}
	return result.(TempBasalInfo)
}

func (pump *Pump) SetTempBasal(duration time.Duration, milliUnitsPerHour int) {
	const halfHour = 30 * time.Minute
	if duration%halfHour != 0 {
		pump.err = fmt.Errorf("temporary basal duration (%v) is not a multiple of 30 minutes", duration)
	}
	if milliUnitsPerHour%25 != 0 {
		pump.err = fmt.Errorf("temporary basal rate (%d) is not a multiple of 25 milliUnits per hour", milliUnitsPerHour)
	}
	pump.Execute(SetTempBasal, nil,
		0,
		byte(milliUnitsPerHour/25),
		byte(duration/halfHour))
}
