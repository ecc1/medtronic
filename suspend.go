package medtronic

const (
	suspend Command = 0x4D
)

// Suspend suspends or resumes the pump.
func (pump *Pump) Suspend(yes bool) {
	if yes {
		pump.Execute(suspend, 1)
	} else {
		pump.Execute(suspend, 0)
	}
}
