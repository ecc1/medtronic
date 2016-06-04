package medtronic

const (
	Suspend CommandCode = 0x4D
)

func (pump *Pump) Suspend(suspend bool) error {
	off := 0
	if suspend {
		off = 1
	}
	_, err := pump.Execute(Suspend, nil, byte(off))
	return err
}
