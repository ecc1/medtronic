package medtronic

const (
	PumpId Command = 0x71
)

func (pump *Pump) PumpId() string {
	data := pump.Execute(PumpId)
	if pump.Error() != nil {
		return ""
	}
	if len(data) == 0 {
		pump.BadResponse(PumpId, data)
		return ""
	}
	n := int(data[0])
	if len(data) < 1+n {
		pump.BadResponse(PumpId, data)
		return ""
	}
	return string(data[1 : 1+n])
}
