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

func (pump *Pump) PumpStatus(retries int) (StatusInfo, error) {
	cmd := PumpCommand{
		Code:       PumpStatus,
		NumRetries: retries,
		ResponseHandler: func(data []byte) interface{} {
			if len(data) >= 4 && data[0] == 3 {
				return StatusInfo{
					Normal:    data[1] == 0x03,
					Bolusing:  data[2] == 0x01,
					Suspended: data[3] == 0x01,
				}
			}
			return nil
		},
	}
	result, err := pump.Execute(cmd)
	if err != nil {
		return StatusInfo{}, err
	}
	return result.(StatusInfo), nil
}
