package rfm69

import (
	"fmt"
	"log"
	"unsafe"
)

func (config *RfConfiguration) Bytes() []byte {
	return (*[RegTemp2 - RegOpMode + 1]byte)(unsafe.Pointer(config))[:]
}

func (r *Radio) ReadConfiguration() (*RfConfiguration, error) {
	regs, err := r.ReadBurst(RegOpMode, RegTemp2-RegOpMode+1)
	return (*RfConfiguration)(unsafe.Pointer(&regs[0])), err
}

func (r *Radio) WriteConfiguration(config *RfConfiguration) error {
	return r.WriteBurst(RegOpMode, config.Bytes())
}

func (r *Radio) InitRF(frequency uint32) error {
	rf := DefaultRfConfiguration
	fb := frequencyBytes(frequency)

	rf.RegDataModul = PacketMode | ModulationTypeOOK | 2<<ModulationShapingShift

	// FxOsc / BitRate = 16385 baud
	rf.RegBitrateMsb = 0x07
	rf.RegBitrateLsb = 0xA1

	rf.RegFrfMsb = fb[0]
	rf.RegFrfMid = fb[1]
	rf.RegFrfLsb = fb[2]

	// Use PA1 with 13 dbM output power
	rf.RegPaLevel = Pa1On | 0x1F<<OutputPowerShift

	// Default != reset value
	rf.RegLna = LnaZin | 1<<LnaCurrentGainShift | 0<<LnaGainSelectShift

	// FXOSC / (RxBwMant * 2^(RxBwExp + 3)) = 200 kHz
	rf.RegRxBw = 2<<DccFreqShift | RxBwMant20 | 0<<RxBwExpShift
	rf.RegAfcBw = 4<<DccFreqShift | RxBwMant20 | 0<<RxBwExpShift

	// Interrupt when Sync word is seen
	rf.RegDioMapping1 = 2 << Dio0MappingShift

	// Default != reset value
	rf.RegDioMapping2 = 7 << ClkOutShift

	// Default != reset value
	rf.RegRssiThresh = 0xE4

	// Use 4 bytes for Sync word
	rf.RegSyncConfig = SyncOn | 3<<SyncSizeShift

	// Sync word
	rf.RegSyncValue1 = 0xFF
	rf.RegSyncValue2 = 0x00
	rf.RegSyncValue3 = 0xFF
	rf.RegSyncValue4 = 0x00

	rf.RegPacketConfig1 = VariableLength
	rf.RegPayloadLength = 0xFF
	rf.RegFifoThresh = TxStartFifoNotEmpty | fifoThreshold<<FifoThresholdShift
	rf.RegPacketConfig2 = AutoRxRestartOff

	err := r.WriteConfiguration(&rf)
	if err != nil {
		return err
	}

	// Default != reset value
	err = r.WriteRegister(RegTestDagc, 0x30)
	if err != nil {
		return err
	}

	return nil
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
	return r.WriteBurst(RegFrfMsb, frequencyBytes(freq))
}

func frequencyBytes(freq uint32) []byte {
	f := (uint64(freq)<<19 + FXOSC/2) / FXOSC
	return []byte{
		byte(f >> 16),
		byte(f >> 8),
		byte(f),
	}
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
