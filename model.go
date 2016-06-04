package medtronic

const (
	Model CommandCode = 0x8D
)

func (pump *Pump) Model(retries int, rssi *int) (string, error) {
	cmd := PumpCommand{
		Code:       Model,
		NumRetries: retries,
		ResponseHandler: func(data []byte) interface{} {
			if len(data) >= 2 {
				n := int(data[1])
				if len(data) >= 2+n {
					return string(data[2 : 2+n])
				}
			}
			return nil
		},
		Rssi: rssi,
	}
	result, err := pump.Execute(cmd)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
