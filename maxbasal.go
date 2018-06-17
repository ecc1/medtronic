package medtronic

import (
	"fmt"
	"log"
)

// SetMaxBasal sets the pump's maximum basal rate.
func (pump *Pump) SetMaxBasal(rate Insulin) {
	if rate < 0 {
		pump.SetError(fmt.Errorf("max basal rate (%d) is negative", rate))
		return
	}
	if rate > maxBasal {
		pump.SetError(fmt.Errorf("max basal rate (%d) is too large", rate))
		return
	}
	m := milliUnitsPerStroke(23)
	strokes := rate / m
	actual := strokes * m
	if actual != rate {
		log.Printf("rounding max basal rate from %v to %v", rate, actual)
	}
	pump.Execute(setMaxBasal, marshalUint16(uint16(strokes))...)
}
