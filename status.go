package medtronic

const (
	status Command = 0xCE
)

// StatusInfo represents the pump's status.
type StatusInfo struct {
	Code      byte
	Bolusing  bool
	Suspended bool
}

// Normal returns true if the status code indicates normal pump operation.
// Observed values:
//   0: rewinding
//   1: preparing to prime
//   2: priming
//   3: normal
func (s StatusInfo) Normal() bool {
	return s.Code == 0x03
}

// Status returns the pump's status.
func (pump *Pump) Status() StatusInfo {
	data := pump.Execute(status)
	if pump.Error() != nil {
		return StatusInfo{}
	}
	if len(data) < 4 || data[0] != 3 {
		pump.BadResponse(status, data)
		return StatusInfo{}
	}
	return StatusInfo{
		Code:      data[1],
		Bolusing:  data[2] == 1,
		Suspended: data[3] == 1,
	}
}
