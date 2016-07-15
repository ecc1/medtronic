package rfm69

import (
	"fmt"
	"log"
	"unsafe"
)

const (
	bitrate   = 16384  // baud
	channelBw = 250000 // Hz
)

func (config *RfConfiguration) Bytes() []byte {
	return (*[RegTemp2 - RegOpMode + 1]byte)(unsafe.Pointer(config))[:]
}

func (r *Radio) ReadConfiguration() *RfConfiguration {
	if r.Error() != nil {
		return nil
	}
	regs := r.hw.ReadBurst(RegOpMode, RegTemp2-RegOpMode+1)
	return (*RfConfiguration)(unsafe.Pointer(&regs[0]))
}

func (r *Radio) WriteConfiguration(config *RfConfiguration) {
	r.hw.WriteBurst(RegOpMode, config.Bytes())
}

func (r *Radio) InitRF(frequency uint32) {
	rf := DefaultRfConfiguration
	fb := frequencyToRegisters(frequency)
	br := bitrateToRegisters(bitrate)
	bw := channelBwToRegister(channelBw)

	rf.RegDataModul = PacketMode | ModulationTypeOOK | 0<<ModulationShapingShift

	rf.RegBitrateMsb = br[0]
	rf.RegBitrateLsb = br[1]

	rf.RegFrfMsb = fb[0]
	rf.RegFrfMid = fb[1]
	rf.RegFrfLsb = fb[2]

	// Use PA1 with 13 dbM output power.
	rf.RegPaLevel = Pa1On | 0x1F<<OutputPowerShift

	// Default != reset value
	rf.RegLna = LnaZin | 1<<LnaCurrentGainShift | 0<<LnaGainSelectShift

	rf.RegRxBw = 2<<DccFreqShift | bw
	rf.RegAfcBw = 4<<DccFreqShift | bw

	// Interrupt when Sync word is seen.
	// Cleared when leaving Rx or FIFO is emptied.
	rf.RegDioMapping1 = 2 << Dio0MappingShift

	// Default != reset value.
	rf.RegDioMapping2 = 7 << ClkOutShift

	// Default != reset value.
	rf.RegRssiThresh = 0xE4

	// Make sure enough preamble bytes are sent.
	rf.RegPreambleMsb = 0x00
	rf.RegPreambleLsb = 0x18

	// Use 4 bytes for Sync word.
	rf.RegSyncConfig = SyncOn | 3<<SyncSizeShift

	// Sync word
	rf.RegSyncValue1 = 0xFF
	rf.RegSyncValue2 = 0x00
	rf.RegSyncValue3 = 0xFF
	rf.RegSyncValue4 = 0x00

	// Use unlimited length packet format (data sheet section 5.5.2.3).
	rf.RegPacketConfig1 = FixedLength
	rf.RegPayloadLength = 0x00
	rf.RegFifoThresh = TxStartFifoNotEmpty | fifoThreshold<<FifoThresholdShift
	rf.RegPacketConfig2 = AutoRxRestartOff

	r.WriteConfiguration(&rf)

	// Default != reset value.
	r.hw.WriteRegister(RegTestDagc, 0x30)
}

func (r *Radio) Frequency() uint32 {
	return registersToFrequency(r.hw.ReadBurst(RegFrfMsb, 3))
}

func registersToFrequency(frf []byte) uint32 {
	f := uint32(frf[0])<<16 + uint32(frf[1])<<8 + uint32(frf[2])
	return uint32(uint64(f) * FXOSC >> 19)
}

func (r *Radio) SetFrequency(freq uint32) {
	r.hw.WriteBurst(RegFrfMsb, frequencyToRegisters(freq))
}

func frequencyToRegisters(freq uint32) []byte {
	f := (uint64(freq)<<19 + FXOSC/2) / FXOSC
	return []byte{byte(f >> 16), byte(f >> 8), byte(f)}
}

func (r *Radio) ReadRSSI() int {
	rssi := r.hw.ReadRegister(RegRssiValue)
	return -int(rssi) / 2
}

func (r *Radio) Bitrate() uint32 {
	return registersToBitrate(r.hw.ReadBurst(RegBitrateMsb, 2))
}

// See data sheet section 3.3.2 and table 9.
func registersToBitrate(br []byte) uint32 {
	d := uint32(br[0])<<8 + uint32(br[1])
	return (FXOSC + d/2) / d
}

func bitrateToRegisters(br uint32) []byte {
	b := (FXOSC + br/2) / br
	return []byte{byte(b >> 8), byte(b)}
}

func (r *Radio) ReadModulationType() byte {
	return r.hw.ReadRegister(RegDataModul) & ModulationTypeMask
}

func (r *Radio) ChannelBw() uint32 {
	bw := r.hw.ReadRegister(RegRxBw)
	m := r.ReadModulationType()
	return registerToChannelBw(bw, m)
}

func registerToChannelBw(bw byte, modType byte) uint32 {
	mant := 0
	switch bw & RxBwMantMask {
	case RxBwMant16:
		mant = 16
	case RxBwMant20:
		mant = 20
	case RxBwMant24:
		mant = 24
	default:
		log.Panicf("unknown RX bandwidth mantissa (%X)", bw&RxBwMantMask)
	}
	e := bw & RxBwExpMask
	switch modType {
	case ModulationTypeFSK:
		return uint32(FXOSC) / (uint32(mant) << (e + 2))
	case ModulationTypeOOK:
		return uint32(FXOSC) / (uint32(mant) << (e + 3))
	default:
		log.Panicf("unknown modulation mode (%X)", modType)
	}
	panic("unreachable")
}

func (r *Radio) SetChannelBw(bw uint32) {
	v := channelBwToRegister(bw)
	r.hw.WriteRegister(RegRxBw, 2<<DccFreqShift|v)
	r.hw.WriteRegister(RegAfcBw, 4<<DccFreqShift|v)
}

// Channel BW = FXOSC / (RxBwMant * 2^(RxBwExp + 3), assuming OOK modulation.
// The caller must add the desired DccFreq field to the result.
func channelBwToRegister(bw uint32) byte {
	bb := uint32(1302) // lowest possible channel bandwidth
	rr := byte(RxBwMant24 | 7<<RxBwExpShift)
	if bw < bb {
		return rr
	}
	for i := 0; i < 8; i++ {
		e := byte(7 - i)
		for j := 0; j < 3; j++ {
			m := byte((6 - j) * 4)
			b := uint32(FXOSC) / (uint32(m) << (e + 3))
			r := byte(2-j)<<RxBwMantShift | e<<RxBwExpShift
			if b >= bw {
				if b-bw < bw-bb {
					return r
				} else {
					return rr
				}
			}
			bb = b
			rr = r
		}
	}
	return rr
}

func (r *Radio) mode() byte {
	return r.hw.ReadRegister(RegOpMode) & ModeMask
}

func (r *Radio) setMode(mode uint8) {
	r.SetError(nil)
	cur := r.hw.ReadRegister(RegOpMode)
	if cur&ModeMask == mode {
		return
	}
	if verbose {
		log.Printf("change from %s to %s", stateName(cur&ModeMask), stateName(mode))
	}
	r.hw.WriteRegister(RegOpMode, cur&^ModeMask|mode)
	for r.Error() == nil {
		s := r.mode()
		if s == mode && r.modeReady() {
			break
		}
		if verbose {
			log.Printf("  %s", stateName(s))
		}
	}
}

func (r *Radio) modeReady() bool {
	return r.hw.ReadRegister(RegIrqFlags1)&ModeReady != 0
}

func (r *Radio) Sleep() {
	r.setMode(SleepMode)
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
		return fmt.Sprintf("Unknown Mode (%X)", mode)
	}
}

func (r *Radio) State() string {
	return stateName(r.mode())
}
