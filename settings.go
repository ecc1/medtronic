package medtronic

import (
	"time"
)

const (
	Settings CommandCode = 0xC0
)

type SettingsInfo struct {
	AutoOff              time.Duration
	MaxBolus             int // milliUnits
	MaxBasal             int // milliUnits
	TempBasalPercent     int
	InsulinAction        time.Duration
	InsulinConcentration int // 50 or 100
}

func (pump *Pump) Settings() (SettingsInfo, error) {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	result, err := pump.Execute(Settings, func(data []byte) interface{} {
		if newer {
			if len(data) < 26 || data[0] != 25 {
				return nil
			}
		} else {
			if len(data) < 22 || data[0] != 21 {
				return nil
			}
		}
		info := SettingsInfo{
			AutoOff:          time.Duration(data[1]) * time.Hour,
			MaxBolus:         int(data[6]) * 100,
			TempBasalPercent: int(data[16]),
			InsulinAction:    time.Duration(data[18]) * time.Hour,
		}
		switch data[10] {
		case 0:
			info.InsulinConcentration = 100
		case 1:
			info.InsulinConcentration = 50
		default:
			return nil
		}
		if newer {
			info.MaxBasal = twoByteInt(data[8:10]) * 25
		} else {
			info.MaxBasal = twoByteInt(data[7:9]) * 100
		}
		return info
	})
	if err != nil {
		return SettingsInfo{}, err
	}
	return result.(SettingsInfo), nil
}
