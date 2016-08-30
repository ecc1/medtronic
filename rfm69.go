package medtronic

import (
	"github.com/ecc1/rfm69"
)

func init() {
	radios = append(radios, rfm69.Open)
}
