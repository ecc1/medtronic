package medtronic

import (
	"fmt"
)

const (
	Battery Command = 0x72
)

// Battery voltage is represented in milliVolts.
type Voltage int

func (r Voltage) String() string {
	return fmt.Sprintf("%g", float64(r)/1000)
}

type BatteryInfo struct {
	Voltage    Voltage
	LowBattery bool
}

func (pump *Pump) Battery() BatteryInfo {
	data := pump.Execute(Battery)
	if pump.Error() != nil {
		return BatteryInfo{}
	}
	if len(data) < 4 || data[0] != 3 {
		pump.BadResponse(Battery, data)
		return BatteryInfo{}
	}
	return BatteryInfo{
		LowBattery: data[1] != 0,
		Voltage:    Voltage(twoByteInt(data[2:4]) * 10),
	}
}
