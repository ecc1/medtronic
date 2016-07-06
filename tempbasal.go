package medtronic

import (
	"fmt"
	"time"
)

const (
	TempBasal            Command = 0x98
	SetAbsoluteTempBasal Command = 0x4C
	SetPercentTempBasal  Command = 0x69
)

type TempBasalType byte

//go:generate stringer -type TempBasalType

const (
	Absolute TempBasalType = 0
	Percent  TempBasalType = 1
)

type TempBasalInfo struct {
	Duration time.Duration
	Type     TempBasalType
	Rate     *Insulin `json:",omitempty"`
	Percent  *uint8   `json:",omitempty"`
}

func (pump *Pump) TempBasal() TempBasalInfo {
	data := pump.Execute(TempBasal)
	if pump.Error() != nil {
		return TempBasalInfo{}
	}
	if len(data) < 7 || data[0] != 6 {
		pump.BadResponse(TempBasal, data)
		return TempBasalInfo{}
	}
	d := time.Duration(twoByteInt(data[5:7])) * time.Minute
	tempType := TempBasalType(data[1])
	info := TempBasalInfo{Duration: d, Type: tempType}
	switch TempBasalType(data[1]) {
	case Absolute:
		rate := twoByteInsulin(data[3:5], true)
		info.Rate = &rate
	case Percent:
		percent := data[2]
		info.Percent = &percent
	default:
		pump.BadResponse(TempBasal, data)
	}
	return info
}

func (pump *Pump) SetAbsoluteTempBasal(duration time.Duration, rate Insulin) {
	d := pump.halfHours(duration)
	if rate%25 != 0 {
		pump.SetError(fmt.Errorf("absolute temporary basal rate (%d) is not a multiple of 25 milliUnits per hour", rate))
	}
	pump.Execute(SetAbsoluteTempBasal, 0, byte(rate/25), d)
}

func (pump *Pump) SetPercentTempBasal(duration time.Duration, percent int) {
	d := pump.halfHours(duration)
	if percent < 0 || 100 < percent {
		pump.SetError(fmt.Errorf("percent temporary basal rate (%d) is not between 0 and 100", percent))
	}
	pump.Execute(SetPercentTempBasal, byte(percent), d)
}

func (pump *Pump) halfHours(duration time.Duration) uint8 {
	const halfHour = 30 * time.Minute
	if duration%halfHour != 0 {
		pump.SetError(fmt.Errorf("duration (%v) is not a multiple of 30 minutes", duration))
	}
	return uint8(duration / halfHour)
}
