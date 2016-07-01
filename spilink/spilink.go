package main

import (
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
		Data    []byte // base64-encoded by json.Marshal
		Repeat  int
		Timeout int // microseconds
	}

	SpiLinkResult struct {
		Data  []byte // base64-encoded by json.Marshal
		Rssi  int
		Error bool
	}
)

var (
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

func send(data []byte, repeat int) SpiLinkResult {
	packet := medtronic.EncodePacket(data)
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
	result.Data = data
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
