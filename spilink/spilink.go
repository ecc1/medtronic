package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/ecc1/medtronic"
)

const (
	verbose = false
)

type (
	SpiLinkCommand struct {
		Command string
		Data    string // hex-encoded
		Repeat  int
		Timeout int // microseconds
	}

	SpiLinkResult struct {
		Data  string // hex-encoded
		Rssi  int
		Error bool
	}
)

var (
	version = "spilink 0.1"

	input  = json.NewDecoder(os.Stdin)
	output = json.NewEncoder(os.Stdout)

	pump = medtronic.Open()
)

func main() {
	if pump.Error() != nil {
		log.Fatal(pump.Error())
	}
	for {
		cmd := readCommand()
		result := cmd.perform()
		err := output.Encode(result)
		pump.SetError(err)
		if pump.Error() != nil {
			log.Println(pump.Error())
			pump.SetError(nil)
		}
	}
}

func readCommand() SpiLinkCommand {
	cmd := SpiLinkCommand{}
	err := input.Decode(&cmd)
	if err == io.EOF {
		if verbose {
			log.Printf("EOF: exiting")
		}
		os.Exit(0)
	}
	if err != nil {
		log.Print(err)
	}
	return cmd
}

func (cmd SpiLinkCommand) perform() SpiLinkResult {
	if verbose {
		log.Printf("received %s command", cmd.Command)
	}
	timeout := time.Duration(cmd.Timeout) * time.Microsecond
	result := SpiLinkResult{}
	switch cmd.Command {
	case "send_packet":
		result = send(cmd.Data, cmd.Repeat)
	case "get_packet":
		result = receive(timeout)
	case "send_and_listen":
		send(cmd.Data, cmd.Repeat)
		result = receive(timeout)
	default:
		log.Printf("unknown spilink command: %+v", cmd)
		result.Error = true
	}
	if verbose {
		log.Printf("returning %d-byte result", len(result.Data)/2)
	}
	return result
}

func send(data string, repeat int) SpiLinkResult {
	packet := medtronic.EncodePacket(hexDecode(data))
	if repeat == 0 {
		repeat = 1
	}
	if verbose {
		if repeat == 1 {
			log.Printf("sending %d-byte packet", len(packet))
		} else {
			log.Printf("sending %d-byte packet %d times", len(packet), repeat)
		}
	}
	for i := 0; i < repeat; i++ {
		pump.Radio.Send(packet)
	}
	return SpiLinkResult{}
}

func receive(timeout time.Duration) SpiLinkResult {
	if verbose {
		log.Printf("receiving with timeout = %v", timeout)
	}
	result := SpiLinkResult{}
	packet, rssi := pump.Radio.Receive(timeout)
	data, err := medtronic.Decode6b4b(packet)
	pump.SetError(err)
	result.Data = hexEncode(data)
	if pump.Error() != nil {
		log.Printf("%v", pump.Error())
		pump.SetError(nil)
		result.Error = true
		return result
	}
	if verbose {
		log.Printf("received %d-byte packet (RSSI = %d)", len(packet), rssi)
	}
	return result
}

// The hexEncode and hexDecode functions correspond to
// the Python encode('hex') and decode('hex') string methods,
// so that arbitrary byte arrays can be transferred as JSON strings.

var (
	hexChar = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
	}
	hexDigit = map[byte]byte{}
)

func init() {
	// Set up the inverse of the hexChar mapping for decoding.
	for i, c := range hexChar {
		hexDigit[c] = byte(i)
	}

}

func hexEncode(data []byte) string {
	var buf bytes.Buffer
	for _, b := range data {
		buf.WriteByte(hexChar[b>>4])
		buf.WriteByte(hexChar[b&0xF])
	}
	return buf.String()
}

func hexDecode(str string) []byte {
	var buf bytes.Buffer
	n := len(str)
	if n%2 != 0 {
		log.Panicf("odd-length hex-encoded string (%s)", str)
	}
	for i := 0; i < n; i += 2 {
		buf.WriteByte(hexDigit[str[i]]<<4 | hexDigit[str[i+1]])
	}
	return buf.Bytes()
}
