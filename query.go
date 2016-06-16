package medtronic

// Generic query returning raw response data.
func (pump *Pump) Query(cmd CommandCode) []byte {
	result := pump.Execute(cmd, func(data []byte) interface{} {
		return data
	})
	if pump.Error() != nil {
		return nil
	}
	return result.([]byte)
}
