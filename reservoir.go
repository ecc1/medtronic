package medtronic

const (
	Reservoir Command = 0x73
)

// Reservoir returns the amount of insulin remaining, in milliUnits.
func (pump *Pump) Reservoir() MilliUnits {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	data := pump.Execute(Reservoir)
	if pump.Error() != nil {
		return 0
	}
	if newer {
		if len(data) < 5 || data[0] != 4 {
			pump.BadResponse(Reservoir, data)
			return 0
		}
		return twoByteMilliUnits(data[3:5], true)
	} else {
		if len(data) < 3 || data[0] != 2 {
			pump.BadResponse(Reservoir, data)
			return 0
		}
		return twoByteMilliUnits(data[1:3], false)
	}
}
