package radio

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

const (
	interruptPin       = 46 // Intel Edison GPIO for receive interrupts
	writeUsingTransfer = false
)

type HwFlavor interface {
	Name() string
	Speed() int
	ReadSingleAddress(byte) byte
	ReadBurstAddress(byte) byte
	WriteSingleAddress(byte) byte
	WriteBurstAddress(byte) byte
}

type Hardware struct {
	device       *spi.Device
	hf           HwFlavor
	err          error
	interruptPin gpio.InputPin
}

func (hw *Hardware) Name() string {
	return hw.hf.Name()
}

func (hw *Hardware) Error() error {
	return hw.err
}

func (hw *Hardware) SetError(err error) {
	hw.err = err
}

func (hw *Hardware) AwaitInterrupt(timeout time.Duration) {
	hw.err = hw.interruptPin.Wait(timeout)
}

func Open(hf HwFlavor) *Hardware {
	hw := &Hardware{hf: hf}
	hw.device, hw.err = spi.Open(hf.Speed())
	if hw.Error() != nil {
		return hw
	}
	hw.err = hw.device.SetMaxSpeed(hf.Speed())
	if hw.Error() != nil {
		hw.Close()
		return hw
	}
	hw.interruptPin, hw.err = gpio.Input(interruptPin, "rising", false)
	if hw.Error() != nil {
		hw.Close()
		return hw
	}
	return hw
}

func (hw *Hardware) Close() {
	hw.device.Close()
}

func (hw *Hardware) ReadRegister(addr byte) byte {
	if hw.Error() != nil {
		return 0
	}
	buf := []byte{hw.hf.ReadSingleAddress(addr), 0}
	hw.err = hw.device.Transfer(buf)
	return buf[1]
}

func (hw *Hardware) ReadBurst(addr byte, n int) []byte {
	if hw.Error() != nil {
		return nil
	}
	buf := make([]byte, n+1)
	buf[0] = hw.hf.ReadBurstAddress(addr)
	hw.err = hw.device.Transfer(buf)
	return buf[1:]
}

func (hw *Hardware) writeData(data []byte) {
	if hw.Error() != nil {
		return
	}
	if writeUsingTransfer {
		hw.err = hw.device.Transfer(data)
	} else {
		hw.err = hw.device.Write(data)
	}
}

func (hw *Hardware) WriteRegister(addr byte, value byte) {
	hw.writeData([]byte{hw.hf.WriteSingleAddress(addr), value})
}

func (hw *Hardware) WriteBurst(addr byte, data []byte) {
	hw.writeData(append([]byte{hw.hf.WriteBurstAddress(addr)}, data...))
}

func (hw *Hardware) WriteEach(data []byte) {
	n := len(data)
	if n%2 != 0 {
		log.Panicf("odd data length (%d)", n)
	}
	for i := 0; i < n; i += 2 {
		hw.WriteRegister(data[i], data[i+1])
	}
}

func (hw *Hardware) SpiDevice() *spi.Device {
	return hw.device
}

type HardwareVersionError struct {
	Actual   uint16
	Expected uint16
}

func (e HardwareVersionError) Error() string {
	return fmt.Sprintf("unexpected hardware version %04X (should be %04X)", e.Actual, e.Expected)
}
