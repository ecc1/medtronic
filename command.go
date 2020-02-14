package medtronic

import (
	"bytes"
	"fmt"
	"log"

	"github.com/ecc1/medtronic/packet"
)

// Command represents a pump command.
type Command byte

//go:generate stringer -type Command

const (
	ack                  Command = 0x06
	nak                  Command = 0x15
	cgmWriteTimestamp    Command = 0x28
	setBasalPatternA     Command = 0x30
	setBasalPatternB     Command = 0x31
	setClock             Command = 0x40
	setMaxBolus          Command = 0x41
	bolus                Command = 0x42
	selectBasalPattern   Command = 0x4A
	setAbsoluteTempBasal Command = 0x4C
	suspend              Command = 0x4D
	button               Command = 0x5B
	wakeup               Command = 0x5D
	setPercentTempBasal  Command = 0x69
	setMaxBasal          Command = 0x6E
	setBasalRates        Command = 0x6F
	clock                Command = 0x70
	pumpID               Command = 0x71
	battery              Command = 0x72
	reservoir            Command = 0x73
	firmwareVersion      Command = 0x74
	errorStatus          Command = 0x75
	historyPage          Command = 0x80
	carbUnits            Command = 0x88
	glucoseUnits         Command = 0x89
	carbRatios           Command = 0x8A
	insulinSensitivities Command = 0x8B
	glucoseTargets512    Command = 0x8C
	model                Command = 0x8D
	settings512          Command = 0x91
	basalRates           Command = 0x92
	basalPatternA        Command = 0x93
	basalPatternB        Command = 0x94
	tempBasal            Command = 0x98
	glucosePage          Command = 0x9A
	isigPage             Command = 0x9B
	calibrationFactor    Command = 0x9C
	lastHistoryPage      Command = 0x9D
	glucoseTargets       Command = 0x9F
	settings             Command = 0xC0
	cgmPageCount         Command = 0xCD
	status               Command = 0xCE
	vcntrPage            Command = 0xD5
)

// NoResponseError indicates that no response to a command was received.
type NoResponseError Command

func (e NoResponseError) Error() string {
	return fmt.Sprintf("no response to %v", Command(e))
}

// NoResponse checks whether the pump has a NoResponseError.
func (pump *Pump) NoResponse() bool {
	_, ok := pump.Error().(NoResponseError)
	return ok
}

// InvalidCommandError indicates that the pump rejected a command as invalid.
type InvalidCommandError struct {
	Command   Command
	PumpError PumpError
}

// PumpError represents an error response from the pump.
type PumpError byte

//go:generate stringer -type PumpError

// Pump error codes.
const (
	CommandRefused           PumpError = 0x08
	SettingOutOfRange        PumpError = 0x09
	BolusInProgress          PumpError = 0x0C
	InvalidHistoryPageNumber PumpError = 0x0D
)

func (e InvalidCommandError) Error() string {
	return fmt.Sprintf("%v error: %v", e.Command, e.PumpError)
}

// BadResponseError indicates an unexpected response to a command.
type BadResponseError struct {
	Command Command
	Data    []byte
}

func (e BadResponseError) Error() string {
	return fmt.Sprintf("unexpected response to %v: % X", e.Command, e.Data)
}

// BadResponse sets the pump's error state to a BadResponseError.
func (pump *Pump) BadResponse(cmd Command, data []byte) {
	pump.SetError(BadResponseError{Command: cmd, Data: data})
}

const (
	shortPacketLength       = 6  // excluding CRC byte
	longPacketLength        = 70 // excluding CRC byte
	encodedLongPacketLength = 107
	payloadLength           = 64
	fragmentLength          = payloadLength + 1 // including sequence number
	doneBit                 = 1 << 7
	maxNAKs                 = 10
)

var (
	shortPacket = make([]byte, shortPacketLength)
	longPacket  = make([]byte, longPacketLength)
	ackPacket   []byte
)

func precomputePackets() {
	addr := PumpAddress()

	shortPacket[0] = packet.Pump
	copy(shortPacket[1:4], addr)

	longPacket[0] = packet.Pump
	copy(longPacket[1:4], addr)

	ackPacket = shortPumpPacket(ack)
}

// shortPumpPacket constructs a 7-byte packet with the specified command code:
//   device type (0xA7)
//   3 bytes of pump ID
//   command code
//   length of parameters (0)
//   CRC-8 (added by packet.Encode)
func shortPumpPacket(cmd Command) []byte {
	p := shortPacket
	p[4] = byte(cmd)
	p[5] = 0
	return packet.Encode(p)
}

// longPumpPacket constructs a 71-byte packet with
// the specified command code and parameters:
//   device type (0xA7)
//   3 bytes of pump ID
//   command code
//   length of parameters (or fragment number if non-zero)
//   64 bytes of parameters plus zero padding
//   CRC-8 (added by packet.Encode)
func longPumpPacket(cmd Command, fragNum int, params []byte) []byte {
	p := longPacket
	p[4] = byte(cmd)
	if fragNum == 0 {
		p[5] = byte(len(params))
	} else {
		// Use a fragment number instead of the length.
		p[5] = uint8(fragNum)
	}
	copy(p[6:], params)
	// Zero-pad the remainder of the packet.
	for i := 6 + len(params); i < longPacketLength; i++ {
		p[i] = 0
	}
	return packet.Encode(p)
}

// Execute sends a command and parameters to the pump and returns its response.
// Commands with parameters require an initial exchange with no parameters,
// followed by an exchange with the actual arguments.
func (pump *Pump) Execute(cmd Command, params ...byte) []byte {
	if len(params) == 0 {
		return pump.perform(cmd, cmd, shortPumpPacket(cmd))
	}
	pump.perform(cmd, ack, shortPumpPacket(cmd))
	if pump.NoResponse() {
		pump.SetError(fmt.Errorf("%v command not performed", cmd))
		return nil
	}
	t := pump.Timeout()
	defer pump.SetTimeout(t)
	pump.SetTimeout(2 * t)
	return pump.perform(cmd, ack, longPumpPacket(cmd, 0, params))
}

// ExtendedRequest sends a command and a sequence of parameter packets
// to the pump and returns its response.
func (pump *Pump) ExtendedRequest(cmd Command, params ...byte) []byte {
	seqNum := 1
	i := 0
	var result []byte
	done := false
	for !done && pump.Error() == nil {
		j := i + payloadLength
		if j >= len(params) {
			done = true
			j = len(params)
		}
		if seqNum == 1 {
			pump.perform(cmd, ack, shortPumpPacket(cmd))
			if pump.NoResponse() {
				pump.SetError(fmt.Errorf("%v command not performed", cmd))
				break
			}
		}
		p := longPumpPacket(cmd, seqNum, params[i:j])
		data := pump.perform(cmd, ack, p)
		result = append(result, data...)
		seqNum++
		i = j
	}
	if done {
		t := pump.Timeout()
		defer pump.SetTimeout(t)
		pump.SetTimeout(2 * t)
		p := longPumpPacket(cmd, seqNum|doneBit, nil)
		data := pump.perform(cmd, ack, p)
		result = append(result, data...)
	}
	return result
}

// ExtendedResponse sends a command and parameters to the pump and
// collects the sequence of packets that make up its response.
func (pump *Pump) ExtendedResponse(cmd Command, params ...byte) []byte {
	var result []byte
	data := pump.Execute(cmd, params...)
	expected := 1
	retries := pump.Retries()
	defer pump.SetRetries(retries)
	pump.SetRetries(1)
	for pump.Error() == nil {
		if len(data) != fragmentLength {
			pump.SetError(fmt.Errorf("%v: received %d-byte response", cmd, len(data)))
			break
		}
		seqNum := int(data[0] &^ doneBit)
		if seqNum != expected {
			pump.SetError(fmt.Errorf("%v: received response %d instead of %d", cmd, seqNum, expected))
			break
		}
		result = append(result, data[1:]...)
		if data[0]&doneBit != 0 {
			break
		}
		// Acknowledge this fragment.
		data = pump.perform(ack, cmd, ackPacket)
		expected++
	}
	return result
}

// History pages are returned as a series of 65-byte fragments:
//   sequence number (1 to numFragments)
//   64 bytes of payload
// The caller must send an ACK to receive the next fragment
// or a NAK to have the current one retransmitted.
// The 0x80 bit is set in the sequence number of the final fragment.
// The page consists of the concatenated payloads.
// The final 2 bytes are the CRC-16 of the preceding data.

type pageStructure struct {
	paramBytes   int // 1 or 4
	numFragments int // 16 or 32 fragments of 64 bytes each
}

var pageData = map[Command]pageStructure{
	historyPage: {
		paramBytes:   1,
		numFragments: 16,
	},
	glucosePage: {
		paramBytes:   4,
		numFragments: 16,
	},
	isigPage: {
		paramBytes:   4,
		numFragments: 32,
	},
	vcntrPage: {
		paramBytes:   1,
		numFragments: 16,
	},
}

// Download requests the given history page from the pump.
func (pump *Pump) Download(cmd Command, page int) []byte {
	maxTries := pump.Retries()
	defer pump.SetRetries(maxTries)
	pump.SetRetries(1)
	for tries := 0; tries < maxTries; tries++ {
		pump.SetError(nil)
		data := pump.tryDownload(cmd, page)
		if pump.Error() == nil {
			logTries(cmd, tries)
			return data
		}
	}
	return nil
}

func (pump *Pump) tryDownload(cmd Command, page int) []byte {
	data := pump.execPage(cmd, page)
	if pump.Error() != nil {
		return nil
	}
	numFragments := pageData[cmd].numFragments
	results := make([]byte, 0, numFragments*payloadLength)
	seq := 1
	for {
		payload, n := pump.checkFragment(page, data, seq, numFragments)
		if pump.Error() != nil {
			return nil
		}
		if n == seq {
			results = append(results, payload...)
			seq++
		}
		if n == numFragments {
			return pump.checkPageCRC(page, results)
		}
		// Acknowledge the current fragment and receive the next.
		next := pump.perform(ack, cmd, ackPacket)
		if pump.Error() != nil {
			if !pump.NoResponse() {
				return nil
			}
			next = pump.handleNoResponse(cmd, page, seq)
		}
		data = next
	}
}

func (pump *Pump) execPage(cmd Command, page int) []byte {
	n := pageData[cmd].paramBytes
	switch n {
	case 1:
		return pump.Execute(cmd, byte(page))
	case 4:
		return pump.Execute(cmd, marshalUint32(uint32(page))...)
	default:
		log.Panicf("%v: unexpected parameter size (%d bytes)", cmd, n)
	}
	panic("unreachable")
}

// checkFragment verifies that a fragment has the expected sequence number
// and returns the payload and sequence number.
func (pump *Pump) checkFragment(page int, data []byte, expected int, numFragments int) ([]byte, int) {
	if len(data) != fragmentLength {
		pump.SetError(fmt.Errorf("history page %d: unexpected fragment length (%d)", page, len(data)))
		return nil, 0
	}
	seqNum := int(data[0] &^ doneBit)
	if seqNum > expected {
		// Missed fragment.
		pump.SetError(fmt.Errorf("history page %d: received fragment %d instead of %d", page, seqNum, expected))
		return nil, 0
	}
	if seqNum < expected {
		// Skip duplicate responses.
		return nil, seqNum
	}
	// This is the next fragment.
	done := data[0]&doneBit != 0
	if (done && seqNum != numFragments) || (!done && seqNum == numFragments) {
		pump.SetError(fmt.Errorf("history page %d: unexpected final sequence number (%d)", page, seqNum))
		return nil, seqNum
	}
	return data[1:], seqNum
}

// handleNoResponse sends NAKs to request retransmission of the expected fragment.
func (pump *Pump) handleNoResponse(cmd Command, page int, expected int) []byte {
	for count := 0; count < maxNAKs; count++ {
		pump.SetError(nil)
		data := pump.perform(nak, cmd, shortPumpPacket(nak))
		if pump.Error() == nil {
			seqNum := int(data[0] &^ doneBit)
			format := "history page %d: received fragment %d after %d NAK"
			if count != 0 {
				format += "s"
			}
			log.Printf(format, page, seqNum, count+1)
			return data
		}
		if !pump.NoResponse() {
			return nil
		}
	}
	pump.SetError(fmt.Errorf("history page %d: lost fragment %d", page, expected))
	return nil
}

// checkPageCRC verifies the history page CRC and returns the page data with the CRC removed.
// In a 2048-byte ISIG page, the CRC-16 is stored in the last 4 bytes: [high 0 low 0]
func (pump *Pump) checkPageCRC(page int, data []byte) []byte {
	if len(data) != cap(data) {
		pump.SetError(fmt.Errorf("history page %d: unexpected size (%d)", page, len(data)))
		return nil
	}
	var dataCRC uint16
	switch cap(data) {
	case 1024:
		dataCRC = twoByteUint(data[1022:])
		data = data[:1022]
	case 2048:
		dataCRC = uint16(data[2044])<<8 | uint16(data[2046])
		data = data[:2044]
	default:
		log.Panicf("unexpected history page size (%d)", cap(data))
	}
	calcCRC := packet.CRC16(data)
	if calcCRC != dataCRC {
		pump.SetError(fmt.Errorf("history page %d: computed CRC %04X but received %04X", page, calcCRC, dataCRC))
		return nil
	}
	return data
}

func (pump *Pump) perform(cmd Command, resp Command, p []byte) []byte {
	if pump.Error() != nil {
		return nil
	}
	maxTries := pump.retries
	if len(p) == encodedLongPacketLength {
		// Don't attempt state-changing commands more than once.
		maxTries = 1
	}
	for tries := 0; tries < maxTries; tries++ {
		pump.SetError(nil)
		response, rssi := pump.Radio.SendAndReceive(p, pump.Timeout())
		if pump.Error() != nil {
			continue
		}
		if len(response) == 0 {
			pump.SetError(NoResponseError(cmd))
			continue
		}
		data, err := packet.Decode(response)
		if err != nil {
			pump.SetError(err)
			continue
		}
		if pump.unexpected(cmd, resp, data) {
			return nil
		}
		logTries(cmd, tries)
		pump.rssi = rssi
		return data[5:]
	}
	if pump.Error() == nil {
		panic("perform")
	}
	return nil
}

func logTries(cmd Command, tries int) {
	if tries == 0 {
		return
	}
	r := "retries"
	if tries == 1 {
		r = "retry"
	}
	log.Printf("%v command required %d %s", cmd, tries, r)
}

func (pump *Pump) unexpected(cmd Command, resp Command, data []byte) bool {
	if len(data) < 6 {
		pump.BadResponse(cmd, data)
		return true
	}
	if !bytes.Equal(data[:4], shortPacket[:4]) {
		pump.BadResponse(cmd, data)
		return true
	}
	switch Command(data[4]) {
	case cmd:
		return false
	case resp:
		return false
	case ack:
		if cmd == cgmWriteTimestamp || cmd == wakeup {
			return false
		}
		pump.BadResponse(cmd, data)
		return true
	case nak:
		pump.SetError(InvalidCommandError{
			Command:   cmd,
			PumpError: PumpError(data[5]),
		})
		return true
	default:
		pump.BadResponse(cmd, data)
		return true
	}
}
