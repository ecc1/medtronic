package medtronic

// Generic query returning raw response data.
func (pump *Pump) Query(cmd Command) []byte {
	return pump.Execute(cmd)
}
