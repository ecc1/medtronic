package cc1100

var (
	PacketsSent int
)

func (dev *Device) TransmitMode() error {
	state, err := dev.ReadState()
	if err != nil {
		return err
	}
	if state != STATE_TX {
		err = dev.ChangeState(STX, STATE_TX)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dev *Device) TransmitPacket(data []byte) error {
	err := dev.TransmitMode()
	if err != nil {
		return err
	}
	// Terminate packet with zero byte.
	err = dev.WriteFifo(append(data, 0))
	if err != nil {
		return err
	}
	PacketsSent++
	return nil
}
