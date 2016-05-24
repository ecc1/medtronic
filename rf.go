package rfm69

import (
	"fmt"
	"log"
)

func (r *Radio) InitRF() error {
	err := r.WriteEach([]byte{
		RegDataModul, PacketMode | ModulationTypeOOK | 0<<ModulationShapingShift,

		// FxOsc / BitRate = 16385 baud
		RegBitrateMsb, 0x07,
		RegBitrateLsb, 0xA1,

		// Use PA1 with 13 dbM output power
		RegPaLevel, Pa1On | 0x1F<<OutputPowerShift,

		// Default != reset value
		RegLna, LnaZin | 1<<LnaCurrentGainShift | 0<<LnaGainSelectShift,

		// FXOSC / (RxBwMant * 2^(RxBwExp + 3)) = 200 kHz
		RegRxBw, 2<<DccFreqShift | RxBwMant20 | 0<<RxBwExpShift,
		RegAfcBw, 4<<DccFreqShift | RxBwMant20 | 0<<RxBwExpShift,

		// Interrupt when Sync word is seen
		RegDioMapping1, 2 << Dio0MappingShift,

		// Default != reset value
		RegDioMapping2, 7 << ClkOutShift,

		// Default != reset value
		RegRssiThresh, 0xE4,

		// Use 4 bytes for Sync word
		RegSyncConfig, SyncOn | 3<<SyncSizeShift,

		// Sync word
		RegSyncValue1, 0xFF,
		RegSyncValue2, 0x00,
		RegSyncValue3, 0xFF,
		RegSyncValue4, 0x00,

		//XXX		RegPacketConfig1, VariableLength,
		//XXX		RegPayloadLength, 0xFF,
		RegFifoThresh, TxStartFifoNotEmpty | fifoThreshold<<FifoThresholdShift,
		RegPacketConfig2, AutoRxRestartOff,

		// Default != reset value
		RegTestDagc, 0x30,
	})
	return err
}

func (r *Radio) Frequency() (uint32, error) {
	frf, err := r.ReadBurst(RegFrfMsb, 3)
	if err != nil {
		return 0, err
	}
	f := uint32(frf[0])<<16 + uint32(frf[1])<<8 + uint32(frf[2])
	return uint32(uint64(f) * FXOSC >> 19), nil
}

func (r *Radio) SetFrequency(freq uint32) error {
	f := (uint64(freq)<<19 + FXOSC/2) / FXOSC
	return r.WriteBurst(RegFrfMsb, []byte{
		byte(f >> 16),
		byte(f >> 8),
		byte(f),
	})
}

func (r *Radio) ReadRSSI() (int, error) {
	rssi, err := r.ReadRegister(RegRssiValue)
	if err != nil {
		return 0, err
	}
	return -int(rssi) / 2, nil
}

func (r *Radio) ReadBitrate() (uint32, error) {
	br, err := r.ReadBurst(RegBitrateMsb, 2)
	if err != nil {
		return 0, err
	}
	d := uint32(br[0])<<8 + uint32(br[1])
	return (FXOSC + d/2) / d, nil
}

func (r *Radio) ReadModulationType() (byte, error) {
	m, err := r.ReadRegister(RegDataModul)
	if err != nil {
		return 0, err
	}
	return m & ModulationTypeMask, nil
}

func (r *Radio) ReadChannelBw() (uint32, error) {
	bw, err := r.ReadRegister(RegRxBw)
	if err != nil {
		return 0, err
	}
	mant := 0
	switch bw & RxBwMantMask {
	case RxBwMant16:
		mant = 16
	case RxBwMant20:
		mant = 20
	case RxBwMant24:
		mant = 24
	default:
		panic(fmt.Sprintf("unknown RX bandwidth mantissa (%X)", bw&RxBwMantMask))
	}
	e := bw & RxBwExpMask
	m, err := r.ReadModulationType()
	if err != nil {
		return 0, err
	}
	switch m {
	case ModulationTypeFSK:
		return uint32(FXOSC) / (uint32(mant) << (e + 2)), nil
	case ModulationTypeOOK:
		return uint32(FXOSC) / (uint32(mant) << (e + 3)), nil
	default:
		panic(fmt.Sprintf("unknown modulation mode (%X)", m))
	}
}

func (r *Radio) mode() (uint8, error) {
	cur, err := r.ReadRegister(RegOpMode)
	if err != nil {
		return 0, err
	}
	return cur & ModeMask, nil
}

func (r *Radio) setMode(mode uint8) error {
	//XXX
	flags, _ := r.ReadRegister(RegIrqFlags2)
	if flags&FifoOverrun != 0 {
		fmt.Printf("FIFO overrun!\n")
	}
	//XXX
	old, err := r.ReadRegister(RegOpMode)
	if err != nil {
		return err
	}
	if old&ModeMask == mode {
		return nil
	}
	if verbose {
		log.Printf("change from %s to %s\n", stateName(old&ModeMask), stateName(mode))
	}
	new := old&^ModeMask | mode
	return r.WriteRegister(RegOpMode, new)
}

func (r *Radio) Sleep() error {
	return r.setMode(SleepMode)
}

func stateName(mode uint8) string {
	switch mode {
	case SleepMode:
		return "Sleep"
	case StandbyMode:
		return "Standby"
	case FreqSynthMode:
		return "Frequency Synthesizer"
	case TransmitterMode:
		return "Transmitter"
	case ReceiverMode:
		return "Receiver"
	default:
		return fmt.Sprintf("Unknown operating mode (%X)", mode)
	}
}

func (r *Radio) State() string {
	mode, err := r.mode()
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return stateName(mode)
}
