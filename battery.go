package medtronic

const (
	Battery CommandCode = 0x72
)

type BatteryInfo struct {
	MilliVolts int
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
		MilliVolts: twoByteInt(data[2:4]) * 10,
	}
}
