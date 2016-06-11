package radio

type Packet struct {
	Rssi int
	Data []byte
}

type Counters struct {
	Received int
	Sent     int
}

type Statistics struct {
	Bytes   Counters
	Packets Counters
}

type Interface interface {
	Init() error

	Frequency() (uint32, error)
	SetFrequency(uint32) error

	Incoming() <-chan Packet
	Outgoing() chan<- Packet

	State() string
	Statistics() Statistics
}
