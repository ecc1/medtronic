package medtronic

const (
	bolus Command = 0x42
)

// Bolus delivers the given amount of insulin as a bolus.
func (pump *Pump) Bolus(amount Insulin) {
	newer := pump.Family() >= 23
	if newer {
		panic("unimplemented")
	}
	n := byte(amount / 100)
	pump.Execute(bolus, n)
}
