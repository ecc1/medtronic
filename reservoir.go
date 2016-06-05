package medtronic

import (
	"log"
)

const (
	Reservoir CommandCode = 0x73
)

// Reservoir returns the amount of insulin remaining, in milliUnits.
func (pump *Pump) Reservoir() (int, error) {
	// Format of response packet and conversion factor depend on model.
	model, err := pump.Model()
	if err != nil {
		return 0, err
	}
	log.Printf("model %s pump\n", model)
	newer := false
	switch model {
	case "523":
		newer = true
	case "723":
		newer = true
	}
	result, err := pump.Execute(Reservoir, func(data []byte) interface{} {
		if newer {
			if len(data) < 5 || data[0] != 4 {
				return nil
			}
			return twoByteInt(data[3:5]) * 25
		} else {
			if len(data) < 3 || data[0] != 2 {
				return nil
			}
			return twoByteInt(data[1:3]) * 100
		}
	})
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}
