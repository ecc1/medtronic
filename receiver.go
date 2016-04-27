package cc1100

import (
	"time"
)

func (dev *Device) ReceiveMode() error {
	state, err := dev.ReadState()
	if err != nil {
		return err
	}
	if state != STATE_RX {
		err = dev.ChangeState(SRX, STATE_RX)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dev *Device) AwaitPacket() error {
	err := dev.ReceiveMode()
	if err != nil {
		return err
	}
	return dev.rxGPIO.Wait()
}

func (dev *Device) PollReceiver() (byte, error) {
	err := dev.ReceiveMode()
	if err != nil {
		return 0, err
	}
	rxbytes, err := dev.ReadRegister(RXBYTES)
	if err != nil {
		return 0, err
	}
	return rxbytes & NUM_RXBYTES_MASK, nil
}

func (dev *Device) ReceivePacket() ([]byte, error) {
	var packet []byte
	for {
		n, err := dev.PollReceiver()
		if err != nil {
			return nil, err
		}
		if n == 0 {
			time.Sleep(250 * time.Microsecond)
			continue
		}
		data, err := dev.ReadFifo(n)
		if err != nil {
			return nil, err
		}
		if data[0] != 0 {
			for i := 0; i < len(data); i++ {
				if data[i] == 0 {
					return append(packet, data[:i]...), nil
				}
			}
			packet = append(packet, data...)
		}
	}
}
