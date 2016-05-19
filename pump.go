package medtronic

import (
	"github.com/ecc1/cc1100"
)

type Pump struct {
	Radio *cc1100.Radio

	DecodingErrors int
	CrcErrors      int
}

func Open() (*Pump, error) {
	r, err := cc1100.Open()
	if err != nil {
		return nil, err
	}
	err = r.Init()
	if err != nil {
		return nil, err
	}
	return &Pump{Radio: r}, nil
}
