package medtronic

const (
	reservoir Command = 0x73
)

func decodeReservoir(data []byte, newerPump bool) (Insulin, error) {
	switch newerPump {
	case true:
		if len(data) < 5 || data[0] != 4 {
			return 0, BadResponseError{Command: reservoir, Data: data}
		}
		return twoByteInsulin(data[3:5], true), nil
	case false:
		if len(data) < 3 || data[0] != 2 {
			return 0, BadResponseError{Command: reservoir, Data: data}
		}
		return twoByteInsulin(data[1:3], false), nil
	}
	panic("unreachable")
}

// Reservoir returns the amount of insulin remaining.
func (pump *Pump) Reservoir() Insulin {
	// Format of response depends on the pump family.
	newer := pump.Family() >= 23
	data := pump.Execute(reservoir)
	if pump.Error() != nil {
		return 0
	}
	i, err := decodeReservoir(data, newer)
	pump.SetError(err)
	return i
}
