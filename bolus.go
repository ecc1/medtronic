package medtronic

import (
	"fmt"
)

const (
	bolus Command = 0x42

	maxBolus = 25000 // milliUnits
)

// Bolus delivers the given amount of insulin as a bolus.
// For safety, this command is not attempted more than once.
func (pump *Pump) Bolus(amount Insulin) {
	newer := pump.Family() >= 23
	if newer {
		panic("unimplemented")
	}
	if amount < 0 {
		pump.SetError(fmt.Errorf("bolus amount (%d) is negative", amount))
	}
	if amount > maxBolus {
		pump.SetError(fmt.Errorf("bolus amount (%d) is too large", amount))
	}
	b := byte(amount / 100)
	n := pump.Retries()
	defer pump.SetRetries(n)
	pump.SetRetries(1)
	pump.Execute(bolus, b)
}
