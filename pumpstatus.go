package medtronic

const (
	PumpStatus CommandCode = 0xCE
)

// Possible values other than 3 (normal):
// 0: rewinding
// 1: preparing to prime
// 2: priming

type StatusInfo struct {
	Normal    bool
	Bolusing  bool
	Suspended bool
}

func (pump *Pump) PumpStatus() (StatusInfo, error) {
	result, err := pump.Execute(PumpStatus, func(data []byte) interface{} {
		if len(data) < 4 || data[0] != 3 {
			return nil
		}
		return StatusInfo{
			Normal:    data[1] == 0x03,
			Bolusing:  data[2] == 0x01,
			Suspended: data[3] == 0x01,
		}
	})
	if err != nil {
		return StatusInfo{}, err
	}
	return result.(StatusInfo), nil
}
