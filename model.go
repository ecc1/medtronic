package medtronic

const (
	Model CommandCode = 0x8D
)

func (pump *Pump) Model() (string, error) {
	result, err := pump.Execute(Model, func(data []byte) interface{} {
		if len(data) < 2 {
			return nil
		}
		n := int(data[1])
		if len(data) < 2+n {
			return nil
		}
		return string(data[2 : 2+n])
	})
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
