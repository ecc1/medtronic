package medtronic

import (
	"bytes"
)

// FirmwareVersion returns the pump's firmware version.
func (pump *Pump) FirmwareVersion() string {
	data := pump.Execute(firmwareVersion)
	if pump.Error() != nil {
		return ""
	}
	if len(data) == 0 {
		pump.BadResponse(firmwareVersion, data)
		return ""
	}
	n := int(data[0])
	if len(data) < 1+n {
		pump.BadResponse(firmwareVersion, data)
		return ""
	}
	return string(bytes.TrimRight(data[1:1+n], " \x0B"))
}
