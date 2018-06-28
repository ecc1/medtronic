package medtronic

import (
	"fmt"
	"time"
)

const (
	maxBasal    = 34000 // milliUnits
	maxDuration = 24 * time.Hour
)

// TempBasalType represents the temp basal type.
type TempBasalType byte

//go:generate stringer -type TempBasalType

const (
	// Absolute represents the pump's use of absolute rates for temporary basals.
	Absolute TempBasalType = 0
	// Percent represents the pump's use of percentage rates for temporary basals.
	Percent TempBasalType = 1
)

// TempBasalInfo represents a temporary basal setting.
type TempBasalInfo struct {
	Duration time.Duration
	Type     TempBasalType
	Rate     *Insulin `json:",omitempty"`
	Percent  *uint8   `json:",omitempty"`
}

func decodeTempBasal(data []byte) (TempBasalInfo, error) {
	if len(data) < 7 || data[0] != 6 {
		return TempBasalInfo{}, BadResponseError{Command: tempBasal, Data: data}
	}
	d := time.Duration(twoByteInt(data[5:7])) * time.Minute
	tempType := TempBasalType(data[1])
	info := TempBasalInfo{Duration: d, Type: tempType}
	switch tempType {
	case Absolute:
		rate := twoByteInsulin(data[3:5], 23)
		info.Rate = &rate
	case Percent:
		percent := data[2]
		info.Percent = &percent
	default:
		return info, BadResponseError{Command: tempBasal, Data: data}
	}
	return info, nil
}

// TempBasal returns the pump's current temporary basal setting.
// If none is in effect, it will have a Duration of 0.
func (pump *Pump) TempBasal() TempBasalInfo {
	data := pump.Execute(tempBasal)
	if pump.Error() != nil {
		return TempBasalInfo{}
	}
	info, err := decodeTempBasal(data)
	pump.SetError(err)
	return info
}

// SetAbsoluteTempBasal sets a temporary basal with the given absolute rate and duration.
func (pump *Pump) SetAbsoluteTempBasal(duration time.Duration, rate Insulin) {
	d := pump.halfHours(duration)
	r, err := encodeBasalRate("temporary basal", rate, pump.Family())
	if err != nil {
		pump.SetError(err)
		return
	}
	args := append(marshalUint16(r), d)
	pump.Execute(setAbsoluteTempBasal, args...)
}

// SetPercentTempBasal sets a temporary basal with the given percent rate and duration.
func (pump *Pump) SetPercentTempBasal(duration time.Duration, percent int) {
	d := pump.halfHours(duration)
	if percent < 0 || 100 < percent {
		pump.SetError(fmt.Errorf("percent temporary basal rate (%d) is not between 0 and 100", percent))
		return
	}
	pump.Execute(setPercentTempBasal, byte(percent), d)
}

func (pump *Pump) halfHours(duration time.Duration) uint8 {
	const halfHour = 30 * time.Minute
	if duration%halfHour != 0 {
		pump.SetError(fmt.Errorf("duration (%v) is not a multiple of 30 minutes", duration))
		return 0
	}
	if duration < 0 {
		pump.SetError(fmt.Errorf("duration (%v) is negative", duration))
		return 0
	}
	if duration > maxDuration {
		pump.SetError(fmt.Errorf("duration (%v) is too large", duration))
		return 0
	}
	return uint8(duration / halfHour)
}
