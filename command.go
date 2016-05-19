package medtronic

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ecc1/cc1100"
)

const (
	PumpDevice             = 0xA7
	defaultResponseTimeout = 100 * time.Millisecond
)

type CommandCode byte

//go:generate stringer -type=CommandCode

const (
	Ack          CommandCode = 0x06
	SetClock     CommandCode = 0x40
	GetClock     CommandCode = 0x70
	PowerControl CommandCode = 0x5D
	GetID        CommandCode = 0x71
	GetBattery   CommandCode = 0x72
	GetModel     CommandCode = 0x8D
)

func noResponse(code CommandCode) error {
	return fmt.Errorf("no response to %s", code.String())
}

func unexpectedResponse(code CommandCode, data []byte) error {
	return fmt.Errorf("unexpected response to %s: % X", code.String(), data)
}

type PumpCommand struct {
	Code            CommandCode
	Params          []byte
	ResponseHandler func([]byte) interface{}
	ResponseTimeout time.Duration
	NumRetries      int
	Rssi            *int
}

var commandPrefix = []byte{
	PumpDevice,
	pumpID[0]<<4 | pumpID[1],
	pumpID[2]<<4 | pumpID[3],
	pumpID[4]<<4 | pumpID[5],
}

func commandPacket(cmd PumpCommand) cc1100.Packet {
	data := append(commandPrefix, byte(cmd.Code), byte(len(cmd.Params)))
	if len(cmd.Params) != 0 {
		data = append(data, cmd.Params...)
	}
	return EncodePacket(data)
}

func (pump *Pump) Execute(cmd PumpCommand) (interface{}, error) {
	packet := commandPacket(cmd)
	responseTimeout := defaultResponseTimeout
	if cmd.ResponseTimeout != 0 {
		responseTimeout = cmd.ResponseTimeout
	}
	for tries := 0; tries < cmd.NumRetries || cmd.NumRetries == 0; tries++ {
		pump.Radio.Outgoing() <- packet
		timeout := time.After(responseTimeout)
		var response cc1100.Packet
		select {
		case response = <-pump.Radio.Incoming():
			break
		case <-timeout:
			continue
		}
		data, err := pump.DecodePacket(response)
		if err != nil {
			continue
		}
		if !expected(cmd.Code, data) {
			return nil, unexpectedResponse(cmd.Code, data)
		}
		if cmd.Rssi != nil {
			*cmd.Rssi = response.Rssi
		}
		result := cmd.ResponseHandler(data[5:])
		if result == nil {
			return nil, unexpectedResponse(cmd.Code, data)
		}
		return result, nil

	}
	return nil, noResponse(cmd.Code)
}

func expected(code CommandCode, data []byte) bool {
	if len(data) < 5 {
		return false
	}
	if !bytes.Equal(data[:len(commandPrefix)], commandPrefix) {
		return false
	}
	if code == PowerControl {
		return data[4] == byte(Ack)
	} else {
		return data[4] == byte(code)
	}
}

func (pump *Pump) ID(retries int) (string, error) {
	cmd := PumpCommand{
		Code:       GetID,
		NumRetries: retries,
		ResponseHandler: func(data []byte) interface{} {
			if len(data) >= 1 {
				n := int(data[0])
				if len(data) >= 1+n {
					return string(data[1 : 1+n])
				}
			}
			return nil
		},
	}
	result, err := pump.Execute(cmd)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

type BatteryInfo struct {
	Millivolts int
	LowBattery bool
}

func (pump *Pump) Battery(retries int) (BatteryInfo, error) {
	cmd := PumpCommand{
		Code:       GetBattery,
		NumRetries: retries,
		ResponseHandler: func(data []byte) interface{} {
			if len(data) >= 4 && data[0] == 3 {
				return &BatteryInfo{
					LowBattery: data[1] != 0,
					Millivolts: (int(data[2])<<8 | int(data[3])) * 10,
				}
			}
			return nil
		},
	}
	result, err := pump.Execute(cmd)
	if err != nil {
		return BatteryInfo{}, err
	}
	return result.(BatteryInfo), nil
}

type ClockInfo struct {
	Hour, Minute, Second int
	Year, Month, Day     int
}

func (pump *Pump) Clock(retries int) (time.Time, error) {
	cmd := PumpCommand{
		Code:       GetClock,
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

func (pump *Pump) Model(retries int, rssi *int) (string, error) {
	cmd := PumpCommand{
		Code:       GetModel,
		NumRetries: retries,
		ResponseHandler: func(data []byte) interface{} {
			if len(data) >= 2 {
				n := int(data[1])
				if len(data) >= 2+n {
					return string(data[2 : 2+n])
				}
			}
			return nil
		},
		Rssi: rssi,
	}
	result, err := pump.Execute(cmd)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (pump *Pump) PowerControl(retries int) error {
	cmd := PumpCommand{
		Code:            PowerControl,
		NumRetries:      retries,
		ResponseTimeout: 10 * time.Second,
		ResponseHandler: func(data []byte) interface{} {
			// Return something other than nil
			return true
		},
	}
	_, err := pump.Execute(cmd)
	return err
}

func (pump *Pump) Wakeup() error {
	const (
		numWakeups = 125
		xmitDelay  = 35 * time.Millisecond
	)
	packet := commandPacket(PumpCommand{Code: PowerControl})
	for i := 0; i < numWakeups; i++ {
		pump.Radio.Outgoing() <- packet
		time.Sleep(xmitDelay)
	}
	return pump.PowerControl(2)
}
