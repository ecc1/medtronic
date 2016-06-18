package medtronic

import (
	"fmt"
)

const (
	CurrentPage CommandCode = 0x9D
	History     CommandCode = 0x80

	historyPageSize = 1024
)

func (pump *Pump) CurrentPage() int {
	result := pump.Execute(CurrentPage, func(data []byte) interface{} {
		if len(data) < 5 || data[0] != 4 {
			return nil
		}
		return fourByteInt(data[1:5])
	})
	if pump.Error() != nil {
		return 0
	}
	return result.(int)
}

func (pump *Pump) History(page int) []byte {
	result := pump.Execute(History, func(data []byte) interface{} {
		return data
	}, byte(page))
	if pump.Error() != nil {
		return nil
	}
	data := result.([]byte)
	prev := byte(0)
	ack := commandPacket(Ack, nil)
	results := []byte{}
	for {
		done := data[0]&0x80 != 0
		data[0] &^= 0x80
		// Skip duplicate responses.
		if data[0] != prev {
			results = append(results, data[1:]...)
			prev = data[0]
		}
		if done {
			break
		}
		pump.Radio.Send(ack)
		next, _ := pump.Radio.Receive(pump.timeout)
		if next == nil {
			pump.SetError(nil)
			continue
		}
		data = pump.DecodePacket(next)
		if pump.Error() != nil {
			pump.SetError(nil)
			continue
		}
		if pump.Error() == nil && !expected(History, data) {
			pump.SetError(BadResponseError{command: History, data: data})
			break
		}
		data = data[5:]
	}
	if len(results) != historyPageSize {
		pump.SetError(fmt.Errorf("unexpected history page size (%d)", len(results)))
		return nil
	}
	dataCrc := twoByteUint(results[historyPageSize-2:])
	results = results[:historyPageSize-2]
	calcCrc := Crc16(results)
	if dataCrc != calcCrc {
		pump.SetError(fmt.Errorf("CRC should be %02X, not %02X", calcCrc, dataCrc))
		return nil
	}
	return results
}
