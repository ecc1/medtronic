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
	family := pump.Family()
	n, err := encodeBolus(amount, family)
	if err != nil {
		pump.SetError(err)
		return
	}
	if family <= 22 {
		pump.Execute(bolus, uint8(n))
	} else {
		pump.Execute(bolus, marshalUint16(n)...)
	}
}

func encodeBolus(amount Insulin, family Family) (uint16, error) {
	if amount < 0 {
		return 0, fmt.Errorf("bolus amount (%d) is negative", amount)
	}
	if amount > maxBasal {
		return 0, fmt.Errorf("bolus amount (%d) is too large", amount)
	}
	// Round the amount to the pump's delivery resolution.
	var res Insulin
	if family <= 22 {
		res = 100
	} else if amount < 1000 {
		res = 25
	} else if amount < 10000 {
		res = 50
	} else {
		res = 100
	}
	actual := (amount / res) * res
	if actual != amount {
		log.Printf("rounding bolus from %v to %v", amount, actual)
	}
	// Encode the rounded value using the family-specific units/stroke.
	m := milliUnitsPerStroke(family)
	return uint16(actual / m), nil
}
