package medtronic

import (
	"fmt"
	"time"
)

const (
	TempBasal    CommandCode = 0x98
	SetTempBasal CommandCode = 0x4C
)

type TempBasalType byte

const (
	Absolute TempBasalType = 0
	Percent  TempBasalType = 1
)

type TempBasalInfo interface {
	Type() TempBasalType
}

type AbsoluteTempBasalInfo struct {
	Duration          time.Duration
	MilliUnitsPerHour int
}

func (info AbsoluteTempBasalInfo) Type() TempBasalType {
	return Absolute
}

type PercentTempBasalInfo struct {
	Duration   time.Duration
	Percentage int
}

func (info PercentTempBasalInfo) Type() TempBasalType {
	return Percent
}

func (pump *Pump) TempBasal() TempBasalInfo {
	data := pump.Execute(TempBasal)
	if pump.Error() != nil {
		return nil
	}
	if len(data) < 7 || data[0] != 6 {
		pump.BadResponse(TempBasal, data)
		return nil
	}
	d := time.Duration(twoByteInt(data[5:7])) * time.Minute
	var info TempBasalInfo
	switch TempBasalType(data[1]) {
	case Absolute:
		info = AbsoluteTempBasalInfo{
			Duration:          d,
			MilliUnitsPerHour: twoByteInt(data[3:5]) * 25,
		}
	case Percent:
		info = PercentTempBasalInfo{
			Duration:   d,
			Percentage: int(data[2]),
		}
	default:
		pump.BadResponse(TempBasal, data)
	}
	return info
}

func (pump *Pump) SetTempBasal(duration time.Duration, milliUnitsPerHour int) {
	const halfHour = 30 * time.Minute
	if duration%halfHour != 0 {
		pump.SetError(fmt.Errorf("temporary basal duration (%v) is not a multiple of 30 minutes", duration))
	}
	if milliUnitsPerHour%25 != 0 {
		pump.SetError(fmt.Errorf("temporary basal rate (%d) is not a multiple of 25 milliUnits per hour", milliUnitsPerHour))
	}
	pump.Execute(SetTempBasal, byte(Absolute), byte(milliUnitsPerHour/25), byte(duration/halfHour))
}
