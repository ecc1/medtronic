package medtronic

func decodeReservoir(data []byte, family Family) (Insulin, error) {
	if family <= 22 {
		if len(data) < 3 || data[0] != 2 {
			return 0, BadResponseError{Command: reservoir, Data: data}
		}
		return twoByteInsulin(data[1:3], family), nil
	}
	if len(data) < 5 || data[0] != 4 {
		return 0, BadResponseError{Command: reservoir, Data: data}
	}
	return twoByteInsulin(data[3:5], family), nil
}

// Reservoir returns the amount of insulin remaining.
func (pump *Pump) Reservoir() Insulin {
	data := pump.Execute(reservoir)
	if pump.Error() != nil {
		return 0
	}
	// Format of response depends on the pump family.
	i, err := decodeReservoir(data, pump.Family())
	pump.SetError(err)
	return i
}
