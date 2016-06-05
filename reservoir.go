package medtronic

const (
	Reservoir CommandCode = 0x73
)

// Reservoir returns the amount of insulin remaining, in milliUnits.
func (pump *Pump) Reservoir() (int, error) {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	result, err := pump.Execute(Reservoir, func(data []byte) interface{} {
		if newer {
			if len(data) < 5 || data[0] != 4 {
				return nil
			}
			return twoByteInt(data[3:5]) * 25
		} else {
			if len(data) < 3 || data[0] != 2 {
				return nil
			}
			return twoByteInt(data[1:3]) * 100
		}
	})
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}
