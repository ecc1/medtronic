package medtronic

// PumpID returns the pump's ID.
func (pump *Pump) PumpID() string {
	data := pump.Execute(pumpID)
	if pump.Error() != nil {
		return ""
	}
	if len(data) == 0 {
		pump.BadResponse(pumpID, data)
		return ""
	}
	n := int(data[0])
	if len(data) < 1+n {
		pump.BadResponse(pumpID, data)
		return ""
	}
	return string(data[1 : 1+n])
}
