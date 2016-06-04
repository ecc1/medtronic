package medtronic

const (
	PumpId CommandCode = 0x71
)

func (pump *Pump) PumpId(retries int) (string, error) {
	cmd := PumpCommand{
		Code:       PumpId,
		NumRetries: retries,
		ResponseHandler: func(data []byte) interface{} {
			if len(data) >= 1 {
				n := int(data[0])
				if len(data) >= 1+n {
					return string(data[1 : 1+n])
				}
			}
			return nil
		},
	}
	result, err := pump.Execute(cmd)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
