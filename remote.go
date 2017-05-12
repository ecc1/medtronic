package medtronic

import (
	"log"
	"os"
	"time"

	"github.com/ecc1/medtronic/packet"
)

// nolint
const (
	rfRemoteEnvVar = "MEDTRONIC_REMOTE_ID"
	rfRemoteDevice = 0xA6

	rfRemoteS   Command = 0x81
	rfRemoteACT Command = 0x86
	rfRemoteB   Command = 0x88
)

var (
	rfRemotePrefix []byte
)

func initRFRemotePrefix() {
	if len(rfRemotePrefix) != 0 {
		return
	}
	id := os.Getenv(rfRemoteEnvVar)
	if len(id) == 0 {
		log.Fatalf("%s environment variable is not set", rfRemoteEnvVar)
	}
	if len(id) != 6 {
		log.Fatalf("%s environment variable must be 6 digits", rfRemoteEnvVar)
	}
	rfRemotePrefix = append([]byte{rfRemoteDevice}, marshalDeviceID(id)...)
}

// rfRemotePacket constructs a packet with the specified command code:
//   device type (0xA6)
//   3 bytes of RF remote ID
//   command code
//   sequence number
//   CRC-8
func rfRemotePacket(cmd Command, seq uint8) []byte {
	initRFRemotePrefix()
	data := make([]byte, 7)
	copy(data, rfRemotePrefix)
	data[4] = byte(cmd)
	data[5] = seq
	return packet.Encode(data)
}

// RFRemote sends an RF remote command with the given sequence number to the pump.
func (pump *Pump) RFRemote(cmd Command, seq uint8) {
	if pump.Error() != nil {
		return
	}
	p := rfRemotePacket(cmd, seq)
	for tries := 0; tries < pump.retries || pump.retries == 0; tries++ {
		pump.Radio.Send(p)
		time.Sleep(pump.Timeout())
	}
}
