package medtronic

const (
	Suspend CommandCode = 0x4D
)

func (pump *Pump) Suspend(suspend bool) {
	var off byte
	if suspend {
		off = 1
	} else {
		off = 0
	}
	pump.Execute(Suspend, nil, off)
}
