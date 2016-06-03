package medtronic

const (
	Suspend CommandCode = 0x4D
)

func (pump *Pump) Suspend(suspend bool, retries int) error {
	code := byte(0)
	if suspend {
		code = byte(1)
	}
	cmd := PumpCommand{
		Code:       Suspend,
		Params:     []byte{code},
		NumRetries: retries,
	}
	_, err := pump.Execute(cmd)
	if err != nil {
		return err
	}
	return nil
}
