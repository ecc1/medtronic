package medtronic

const (
	status Command = 0xCE
)

// StatusInfo represents the pump's status.
type StatusInfo struct {
	Normal    bool
	Bolusing  bool
	Suspended bool
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
	// Observed values for data[1]:
	//   0: rewinding
	//   1: preparing to prime
	//   2: priming
	//   3: normal
	return StatusInfo{
		Normal:    data[1] == 0x03,
		Bolusing:  data[2] == 1,
		Suspended: data[3] == 1,
	}
}
