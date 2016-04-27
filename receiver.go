package cc1100

import (
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

var (
	// GPIO input connected to GDO0
	gdo0 gpio.InputPin
)

func startReceiver(dev *spi.Device) error {
	// Enable interrupts from GDO0 pin connected to GPIO14
	var err error
	gdo0, err = gpio.Input(14, "both", false)
	if err != nil {
		return err
	}
	return nil
}

func ReceiveMode(dev *spi.Device) error {
	state, err := ReadState(dev)
	if err != nil {
		return err
	}
	if state != STATE_RX {
		err = ChangeState(dev, SRX, STATE_RX)
		if err != nil {
			return err
		}
	}
	return nil
}

func AwaitPacket(dev *spi.Device) error {
	err := ReceiveMode(dev)
	if err != nil {
		return err
	}
	return gdo0.Wait()
}

func PollReceiver(dev *spi.Device) (byte, error) {
	err := ReceiveMode(dev)
	if err != nil {
		return 0, err
	}
	rxbytes, err := ReadRegister(dev, RXBYTES)
	if err != nil {
		return 0, err
	}
	return rxbytes & NUM_RXBYTES_MASK, nil
}

func ReceivePacket(dev *spi.Device) ([]byte, error) {
	var packet []byte
	for {
		n, err := PollReceiver(dev)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			time.Sleep(250 * time.Microsecond)
			continue
		}
		data, err := ReadFifo(dev, n)
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
