package medtronic

import (
	"fmt"
)

const (
	battery Command = 0x72
)

// Voltage represents the battery voltage in milliVolts.
type Voltage int

func (r Voltage) String() string {
	return fmt.Sprintf("%g", float64(r)/1000)
}

// BatteryInfo represents the pump's battery voltage and low-battery state.
type BatteryInfo struct {
	Voltage    Voltage
	LowBattery bool
}

func decodeBatteryInfo(data []byte) BatteryInfo {
	return BatteryInfo{
		LowBattery: data[1] != 0,
		Voltage:    Voltage(twoByteInt(data[2:4]) * 10),
	}
}

// Battery returns the pump's battery information.
func (pump *Pump) Battery() BatteryInfo {
	data := pump.Execute(battery)
	if pump.Error() != nil {
		return BatteryInfo{}
	}
	if len(data) < 4 || data[0] != 3 {
		pump.BadResponse(battery, data)
		return BatteryInfo{}
	}
	return decodeBatteryInfo(data)
}
