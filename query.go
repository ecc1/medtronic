package medtronic

// Generic query returning raw response data.
func (pump *Pump) Query(cmd CommandCode) []byte {
	return pump.Execute(cmd)
}
