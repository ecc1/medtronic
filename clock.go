package medtronic

import (
	"time"
)

const (
	Clock CommandCode = 0x70
)

func (pump *Pump) Clock(retries int) (time.Time, error) {
	cmd := PumpCommand{
		Code:       Clock,
		NumRetries: retries,
		ResponseHandler: func(data []byte) interface{} {
			if len(data) >= 8 && data[0] == 7 {
				return time.Date(
					int(data[4])<<8|int(data[5]), // year
					time.Month(data[6]),          // month
					int(data[7]),                 // day
					int(data[1]),                 // hour
					int(data[2]),                 // min
					int(data[3]),                 // sec
					0,                            // nsec
					time.Local)
			}
			return nil
		},
	}
	result, err := pump.Execute(cmd)
	if err != nil {
		return time.Time{}, err
	}
	return result.(time.Time), nil
}
