// +build !cc1101,!rfm69

package medtronic

import "github.com/ecc1/radio"

// Catch builds that do not specify which radio to use.

var radioInterface = func() radio.Interface {
	panic("no radio was specified at build time")
}
