package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/ecc1/medtronic"
	"github.com/ecc1/medtronic/packet"
)

const (
	verbose = false
)

type (
	// SPILinkCommand represents a command sent by the client.
	SPILinkCommand struct {
		Command string
		Data    []byte // base64-encoded by json.Marshal
		Repeat  int
		Timeout int // microseconds
	}

	// SPILinkResult represents a result returned to the client.
	SPILinkResult struct {
		Data  []byte // base64-encoded by json.Marshal
		RSSI  int
		Error bool
	}
)

var (
	input  = json.NewDecoder(os.Stdin)
	output = json.NewEncoder(os.Stdout)

	radio = medtronic.Open().Radio
)

func main() {
	if radio.Error() != nil {
		log.Fatal(radio.Error())
	}
	for {
		cmd := readCommand()
		result := cmd.perform()
		err := output.Encode(result)
		radio.SetError(err)
		if radio.Error() != nil {
			log.Print(radio.Error())
			radio.SetError(nil)
		}
	}
}

func readCommand() SPILinkCommand {
	cmd := SPILinkCommand{}
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

func (cmd SPILinkCommand) perform() SPILinkResult {
	if verbose {
		log.Printf("received %s command", cmd.Command)
	}
	timeout := time.Duration(cmd.Timeout) * time.Microsecond
	result := SPILinkResult{}
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
		log.Printf("returning %d-byte result", len(result.Data))
	}
	return result
}

func send(data []byte, repeat int) SPILinkResult {
	p := packet.Encode(data)
	if repeat == 0 {
		repeat = 1
	}
	if verbose {
		if repeat == 1 {
			log.Printf("sending %d-byte packet", len(p))
		} else {
			log.Printf("sending %d-byte packet %d times", len(p), repeat)
		}
	}
	for i := 0; i < repeat; i++ {
		radio.Send(p)
	}
	return SPILinkResult{}
}

func receive(timeout time.Duration) SPILinkResult {
	if verbose {
		log.Printf("receiving with timeout = %v", timeout)
	}
	result := SPILinkResult{}
	p, rssi := radio.Receive(timeout)
	data, err := packet.Decode6b4b(p)
	radio.SetError(err)
	result.Data = data
	if radio.Error() != nil {
		log.Print(radio.Error())
		radio.SetError(nil)
		result.Error = true
		return result
	}
	if verbose {
		log.Printf("received %d-byte packet (RSSI = %d)", len(p), rssi)
	}
	return result
}
