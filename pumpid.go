package medtronic

const (
	PumpId CommandCode = 0x71
)

func (pump *Pump) PumpId() string {
	result := pump.Execute(PumpId, func(data []byte) interface{} {
		if len(data) == 0 {
			return nil
		}
		n := int(data[0])
		if len(data) < 1+n {
			return nil
		}
		return string(data[1 : 1+n])
	})
	if pump.Error() != nil {
		return ""
	}
	return result.(string)
}
