// +build !cc1101,!rfm69,!rfm95

package medtronic

import "github.com/ecc1/cc111x"

var radioInterface = cc111x.Open
