package cc1100

import (
	"bytes"
	"log"
	"time"
)

const (
	usePolling   = false
	pollInterval = time.Millisecond
)

var (
	PacketsReceived = 0
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

func (dev *Device) AwaitPacket(timeout time.Duration) (bool, error) {
	err := dev.ReceiveMode()
	if err != nil {
		return false, err
	}
	if usePolling {
		poll := time.Tick(pollInterval)
		var t <-chan time.Time
		if timeout != 0 {
			t = time.After(timeout)
		}
		for {
			select {
			case <-poll:
				n, err := dev.ReadNumRxBytes()
				if err != nil {
					return false, err
				}
				if n != 0 {
					return true, nil
				}
			case <-t:
				return false, nil
			}
		}
	}
	if timeout == 0 {
		return true, dev.rxGPIO.Wait()
	} else {
		return dev.rxGPIO.TimedWait(timeout)
	}
}

var (
	receiverStarted = false
	receivedPackets = make(chan []byte, 10)
)

func (dev *Device) ReceivePacket() ([]byte, error) {
	if !receiverStarted {
		go dev.receiver()
		receiverStarted = true
	}
	return <-receivedPackets, nil
}

func (dev *Device) receiver() {
	var packet bytes.Buffer
	for {
		numBytes, err := dev.ReadNumRxBytes()
		if err == RxFifoOverflow {
			dev.Strobe(SIDLE)
			dev.Strobe(SFRX)
			dev.AwaitPacket(0)
			continue
		} else if err != nil {
			log.Fatal(err)
		}
		if numBytes == 0 {
			dev.AwaitPacket(0)
			continue
		}
		c, err := dev.ReadRegister(RXFIFO)
		if err != nil {
			log.Fatal(err)
		}
		if c != 0 {
			err = packet.WriteByte(c)
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		// End of packet.
		size := packet.Len()
		if size != 0 {
			PacketsReceived++
			p := make([]byte, size)
			_, err := packet.Read(p)
			if err != nil {
				log.Fatal(err)
			}
			receivedPackets <- p
		}
		packet.Reset()
		if numBytes > 1 {
			dev.Strobe(SIDLE)
			dev.Strobe(SFRX)
		}
	}
}
