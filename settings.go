package medtronic

import (
	"time"
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

func decodeSettings(data []byte, family Family) (SettingsInfo, error) {
	var info SettingsInfo
	if family <= 12 {
		if len(data) < 19 || data[0] != 18 {
			return info, BadResponseError{Command: settings, Data: data}
		}
		info.MaxBolus = byteToInsulin(data[6], 22)
		info.MaxBasal = twoByteInsulin(data[7:9], 23)
	} else if family <= 22 {
		if len(data) < 22 || data[0] != 21 {
			return info, BadResponseError{Command: settings, Data: data}
		}
		info.MaxBolus = byteToInsulin(data[6], 22)
		info.MaxBasal = twoByteInsulin(data[7:9], 23)
		info.InsulinAction = time.Duration(data[18]) * time.Hour
	} else {
		if len(data) < 26 || data[0] != 25 {
			return info, BadResponseError{Command: settings, Data: data}
		}
		info.MaxBolus = byteToInsulin(data[7], 22)
		info.MaxBasal = twoByteInsulin(data[8:10], 23)
		info.InsulinAction = time.Duration(data[18]) * time.Hour
	}
	info.AutoOff = time.Duration(data[1]) * time.Hour
	info.SelectedPattern = int(data[12])
	info.RFEnabled = data[13] == 1
	info.TempBasalType = TempBasalType(data[14])
	var err error
	info.InsulinConcentration, err = insulinConcentration(data)
	return info, err
}

// Settings returns the pump's settings.
func (pump *Pump) Settings() SettingsInfo {
	// Command opcode and format of response depend on the pump family.
	family := pump.Family()
	var cmd Command
	if family <= 12 {
		cmd = settings512
	} else {
		cmd = settings
	}
	data := pump.Execute(cmd)
	if pump.Error() != nil {
		return SettingsInfo{}
	}
	i, err := decodeSettings(data, family)
	pump.SetError(err)
	return i
}

func insulinConcentration(data []byte) (int, error) {
	switch data[10] {
	case 0:
		return 100, nil
	case 1:
		return 50, nil
	default:
		return 0, BadResponseError{Command: settings, Data: data}
	}
}
