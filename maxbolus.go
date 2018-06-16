package medtronic

import (
	"fmt"
	"log"
)

// SetMaxBolus sets the pump's maximum bolus.
func (pump *Pump) SetMaxBolus(amount Insulin) {
	if amount < 0 {
		pump.SetError(fmt.Errorf("bolus amount (%d) is negative", amount))
	}
	if amount > maxBolus {
		pump.SetError(fmt.Errorf("bolus amount (%d) is too large", amount))
	}
	if pump.Error() != nil {
		return
	}
	m := milliUnitsPerStroke(22)
	strokes := amount / m
	actual := strokes * m
	if actual != amount {
		log.Printf("rounding max bolus from %v to %v", amount, actual)
	}
	pump.Execute(setMaxBolus, uint8(strokes))
}
