// +build !cc1101,!cc111x

package medtronic

import "github.com/ecc1/rfm69"

var radioInterface = rfm69.Open
