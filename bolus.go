package medtronic

import (
	"fmt"
	"log"
)

const (
	maxBolus = 25000 // milliUnits
)

// Bolus delivers the given amount of insulin as a bolus.
func (pump *Pump) Bolus(amount Insulin) {
	if amount < 0 {
		pump.SetError(fmt.Errorf("bolus amount (%d) is negative", amount))
	}
	if amount > maxBolus {
		pump.SetError(fmt.Errorf("bolus amount (%d) is too large", amount))
	}
	family := pump.Family()
	d := milliUnitsPerStroke(family)
	strokes := amount / d
	actual := strokes * d
	if actual != amount {
		log.Printf("rounding bolus from %v to %v", amount, actual)
	}
	if family <= 22 {
		pump.Execute(bolus, uint8(strokes))
	} else {
		pump.Execute(bolus, marshalUint16(uint16(strokes))...)
	}
}
