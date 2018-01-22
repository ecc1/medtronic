package medtronic

const (
	button Command = 0x5B
)

// PumpButton represents a key on the pump keypad.
type PumpButton byte

//go:generate stringer -type PumpButton

// Pump button codes.
const (
	BolusButton PumpButton = 0
	EscButton   PumpButton = 1
	ActButton   PumpButton = 2
	UpButton    PumpButton = 3
	DownButton  PumpButton = 4
)

// Button sends the button-press to the pump.
func (pump *Pump) Button(b PumpButton) {
	n := pump.Retries()
	defer pump.SetRetries(n)
	pump.SetRetries(1)
	pump.Execute(button, byte(b))
}
