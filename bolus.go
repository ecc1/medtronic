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
	if amount < 0 {
		pump.SetError(fmt.Errorf("bolus amount (%d) is negative", amount))
	}
	if amount > maxBolus {
		pump.SetError(fmt.Errorf("bolus amount (%d) is too large", amount))
	}
	newer := pump.Family() >= 23
	strokes := int(amount / milliUnitsPerStroke(newer))
	n := pump.Retries()
	defer pump.SetRetries(n)
	pump.SetRetries(1)
	switch newer {
	case true:
		pump.Execute(bolus, marshalUint16(uint16(strokes))...)
	case false:
		pump.Execute(bolus, uint8(strokes))
	}
}
