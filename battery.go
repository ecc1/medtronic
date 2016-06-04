package medtronic

const (
	Battery CommandCode = 0x72
)

type BatteryInfo struct {
	Millivolts int
	LowBattery bool
}

func (pump *Pump) Battery() (BatteryInfo, error) {
	result, err := pump.Execute(Battery, func(data []byte) interface{} {
		if len(data) >= 4 && data[0] == 3 {
			return BatteryInfo{
				LowBattery: data[1] != 0,
				Millivolts: (int(data[2])<<8 | int(data[3])) * 10,
			}
		}
		return nil
	})
	if err != nil {
		return BatteryInfo{}, err
	}
	return result.(BatteryInfo), nil
}
