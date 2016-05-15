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

func (dev *Device) ReceiveMode() error {
	if !dev.receiverStarted {
		dev.receiverStarted = true
		dev.receivedPackets = make(chan []byte, 10)
		go dev.receiver()
	}
	return dev.ChangeState(SRX, STATE_RX)
}

func (dev *Device) AwaitPacket() error {
	for {
		n, err := dev.ReadNumRxBytes()
		if err != nil {
			return err
		}
		if n != 0 {
			return nil
		}
		if !usePolling {
			return dev.interruptPin.Wait()
		}
		time.Sleep(pollInterval)
	}
}

func (dev *Device) IncomingPackets() <-chan []byte {
	return dev.receivedPackets
}

func (dev *Device) flushRxFifo() {
	dev.ChangeState(SFRX, STATE_IDLE)
	dev.ChangeState(SRX, STATE_RX)
}

func (dev *Device) receiver() {
	var packet bytes.Buffer
	for {
		numBytes, err := dev.ReadNumRxBytes()
		if err == RxFifoOverflow {
			dev.flushRxFifo()
			continue
		} else if err != nil {
			log.Fatal(err)
		}
		if numBytes == 0 {
			dev.AwaitPacket()
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
			dev.packetsReceived++
			p := make([]byte, size)
			_, err := packet.Read(p)
			if err != nil {
				log.Fatal(err)
			}
			dev.receivedPackets <- p
		}
		packet.Reset()
		if numBytes > 1 {
			if Verbose {
				log.Printf("flushing %d bytes\n", numBytes-1)
			}
			dev.ChangeState(SIDLE, STATE_IDLE)
			dev.flushRxFifo()
		}
	}
}
