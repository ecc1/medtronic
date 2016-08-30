package medtronic

import (
	"github.com/ecc1/cc1101"
)

func init() {
	radios = append(radios, cc1101.Open)
}
