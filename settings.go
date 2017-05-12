package medtronic

import (
	"time"
)

const (
	settings Command = 0xC0
)

// SettingsInfo represents the pump's settings.
type SettingsInfo struct {
	AutoOff              time.Duration
	InsulinAction        time.Duration
	InsulinConcentration int // 50 or 100
	MaxBolus             Insulin
	MaxBasal             Insulin
	RFEnabled            bool
	TempBasalType        TempBasalType
	SelectedPattern      int
}

// Settings returns the pump's settings.
func (pump *Pump) Settings() SettingsInfo {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	data := pump.Execute(settings)
	if pump.Error() != nil {
		return SettingsInfo{}
	}
	var info SettingsInfo
	switch newer {
	case true:
		if len(data) < 26 || data[0] != 25 {
			pump.BadResponse(settings, data)
			return info
		}
		info.MaxBolus = byteToInsulin(data[7], false)
		info.MaxBasal = twoByteInsulin(data[8:10], true)
	case false:
		if len(data) < 22 || data[0] != 21 {
			pump.BadResponse(settings, data)
			return info
		}
		info.MaxBolus = byteToInsulin(data[6], false)
		info.MaxBasal = twoByteInsulin(data[7:9], true)
	}
	info.AutoOff = time.Duration(data[1]) * time.Hour
	info.SelectedPattern = int(data[12])
	info.RFEnabled = data[13] == 1
	info.TempBasalType = TempBasalType(data[14])
	info.InsulinAction = time.Duration(data[18]) * time.Hour
	info.InsulinConcentration = pump.insulinConcentration(data)
	return info
}

func (pump *Pump) insulinConcentration(data []byte) int {
	switch data[10] {
	case 0:
		return 100
	case 1:
		return 50
	default:
		pump.BadResponse(settings, data)
		return 0
	}
}
