package cc1101

import (
	"errors"
	"unsafe"
)

var (
	RxFifoOverflow  = errors.New("RXFIFO overflow")
	TxFifoUnderflow = errors.New("TXFIFO underflow")
)

func (config *RfConfiguration) Bytes() []byte {
	return (*[TEST0 - IOCFG2 + 1]byte)(unsafe.Pointer(config))[:]
}

func (r *Radio) ReadConfiguration() *RfConfiguration {
	if r.Error() != nil {
		return nil
	}
	regs := r.ReadBurst(IOCFG2, TEST0-IOCFG2+1)
	return (*RfConfiguration)(unsafe.Pointer(&regs[0]))
}

func (r *Radio) WriteConfiguration(config *RfConfiguration) {
	r.WriteBurst(IOCFG2, config.Bytes())
}

func (r *Radio) InitRF(frequency uint32) {
	rf := ResetRfConfiguration
	fb := frequencyToRegisters(frequency)

	rf.IOCFG2 = 0x2F
	rf.IOCFG1 = 0x2F

	// Assert when sync word has been sent/received
	rf.IOCFG0 = 0x06

	// 4 bytes in RX FIFO, 61 bytes in TX FIFO
	rf.FIFOTHR = 0x00

	rf.SYNC1 = 0xFF
	rf.SYNC0 = 0x00

	rf.PKTCTRL1 = 4 << PKTCTRL1_PQT_SHIFT
	rf.PKTCTRL0 = PKTCTRL0_LENGTH_CONFIG_INFINITE

	// Intermediate frequency
	// 0x06 * 24 MHz / 2^10 == 140625 Hz
	rf.FSCTRL1 = 0x06

	rf.FREQ2 = fb[0]
	rf.FREQ1 = fb[1]
	rf.FREQ0 = fb[2]

	// CHANBW_E = 1, CHANBW_M = 1, DRATE_E = 9
	// Channel BW = 24 MHz / (8 * (4 + CHANBW_M) * 2^CHANBW_E) == 300 kHz
	rf.MDMCFG4 = 1<<MDMCFG4_CHANBW_E_SHIFT |
		1<<MDMCFG4_CHANBW_M_SHIFT |
		9<<MDMCFG4_DRATE_E_SHIFT

	// DRATE_M = 102 (0x66)
	// Data rate = (256 + DRATE_M) * 2^DRATE_E * 24 MHz / 2^28 == 16388 Baud
	rf.MDMCFG3 = 0x66

	rf.MDMCFG2 = MDMCFG2_DEM_DCFILT_ON |
		MDMCFG2_MOD_FORMAT_ASK_OOK |
		MDMCFG2_SYNC_MODE_30_32_THRES

		// CHANSPC_E = 2
	rf.MDMCFG1 = MDMCFG1_FEC_DIS |
		MDMCFG1_NUM_PREAMBLE_24 |
		2<<MDMCFG1_CHANSPC_E_SHIFT

	// CHANSPC_M = 26 (0x1A)
	// Channel spacing = (256 + CHANSPC_M) * 2^CHANSPC_E * 24 MHz / 2^18 == 103271 Hz
	rf.MDMCFG0 = 0x1A

	rf.MCSM2 = MCSM2_RX_TIME_END_OF_PACKET

	rf.MCSM1 = MCSM1_CCA_MODE_RSSI_BELOW_UNLESS_RECEIVING |
		MCSM1_RXOFF_MODE_IDLE |
		MCSM1_TXOFF_MODE_IDLE

	rf.MCSM0 = MCSM0_FS_AUTOCAL_FROM_IDLE

	rf.FOCCFG = FOCCFG_FOC_PRE_K_3K |
		FOCCFG_FOC_POST_K_PRE_K_OVER_2 |
		FOCCFG_FOC_LIMIT_BW_OVER_2

	rf.BSCFG = BSCFG_BS_PRE_K_2K |
		BSCFG_BS_PRE_KP_3KP |
		BSCFG_BS_POST_KI_PRE_KI_OVER_2 |
		BSCFG_BS_LIMIT_0

	rf.AGCCTRL2 = AGCCTRL2_MAX_DVGA_GAIN_ALL |
		AGCCTRL2_MAX_LNA_GAIN_0 |
		AGCCTRL2_MAGN_TARGET_38dB

	rf.AGCCTRL1 = AGCCTRL1_AGC_LNA_PRIORITY_0 |
		AGCCTRL1_CARRIER_SENSE_REL_THR_DISABLE |
		AGCCTRL1_CARRIER_SENSE_ABS_THR_0DB

	rf.AGCCTRL0 = AGCCTRL0_HYST_LEVEL_MEDIUM |
		AGCCTRL0_WAIT_TIME_16 |
		AGCCTRL0_AGC_FREEZE_NORMAL |
		AGCCTRL0_FILTER_LENGTH_32

	rf.FREND1 = 1<<FREND1_LNA_CURRENT_SHIFT |
		1<<FREND1_LNA2MIX_CURRENT_SHIFT |
		1<<FREND1_LODIV_BUF_CURRENT_RX_SHIFT |
		2<<FREND1_MIX_CURRENT_SHIFT

	// Use PA_TABLE 1 for transmitting '1' in ASK
	// (PA_TABLE 0 is always used for '0')
	rf.FREND0 = 1<<FREND0_LODIV_BUF_CURRENT_TX_SHIFT |
		1<<FREND0_PA_POWER_SHIFT

	rf.FSCAL3 = 3<<6 | 2<<4 | 0x09
	rf.FSCAL2 = 1<<5 | 0x0A // VCO high
	rf.FSCAL1 = 0x00
	rf.FSCAL0 = 0x1F

	rf.TEST2 = TEST2_RX_LOW_DATA_RATE_MAGIC
	rf.TEST1 = TEST1_RX_LOW_DATA_RATE_MAGIC
	rf.TEST0 = 2<<2 | 1 // disable VCO selection calibration

	r.WriteConfiguration(&rf)

	// Power amplifier output settings (see section 24 of the data sheet)
	r.WriteBurst(PATABLE, []byte{0x00, 0xC0})
}

func (r *Radio) Frequency() uint32 {
	return registersToFrequency(r.ReadBurst(FREQ2, 3))
}

func registersToFrequency(freq []byte) uint32 {
	f := uint32(freq[0])<<16 + uint32(freq[1])<<8 + uint32(freq[2])
	return uint32(uint64(f) * FXOSC >> 16)
}

func (r *Radio) SetFrequency(freq uint32) {
	r.WriteBurst(FREQ2, frequencyToRegisters(freq))
}

func frequencyToRegisters(freq uint32) []byte {
	f := (uint64(freq)<<16 + FXOSC/2) / FXOSC
	return []byte{byte(f >> 16), byte(f >> 8), byte(f)}
}

func (r *Radio) ReadIF() uint32 {
	f := r.ReadRegister(FSCTRL1)
	return uint32(uint64(f) * FXOSC >> 10)
}

func (r *Radio) ReadChannelParams() (chanbw uint32, drate uint32) {
	m4 := r.ReadRegister(MDMCFG4)
	chanbw_E := (m4 >> MDMCFG4_CHANBW_E_SHIFT) & 0x3
	chanbw_M := (m4 >> MDMCFG4_CHANBW_M_SHIFT) & 0x3
	drate_E := (m4 >> MDMCFG4_DRATE_E_SHIFT) & 0xF
	drate_M := r.ReadRegister(MDMCFG3)
	chanbw = uint32(FXOSC / ((4 + uint64(chanbw_M)) << (chanbw_E + 3)))
	drate = uint32(((256 + uint64(drate_M)) << drate_E * FXOSC) >> 28)
	return
}

func (r *Radio) ReadModemConfig() (fec bool, minPreamble byte, chanspc uint32) {
	m1 := r.ReadRegister(MDMCFG1)
	fec = m1&MDMCFG1_FEC_EN != 0
	minPreamble = numPreamble[(m1&MDMCFG1_NUM_PREAMBLE_MASK)>>4]
	chanspc_E := m1 & MDMCFG1_CHANSPC_E_MASK
	chanspc_M := r.ReadRegister(MDMCFG0)
	chanspc = uint32(((256 + uint64(chanspc_M)) << chanspc_E * FXOSC) >> 18)
	return
}

func (r *Radio) ReadRSSI() int {
	const rssi_offset = 74 // see data sheet section 17.3
	rssi := r.ReadRegister(RSSI)
	d := int(rssi)
	if d >= 128 {
		d -= 256
	}
	return d/2 - rssi_offset
}

func (r *Radio) ReadPaTable() []byte {
	buf := make([]byte, 9)
	buf[0] = READ_MODE | BURST_MODE | PATABLE
	r.err = r.device.Transfer(buf)
	return buf[1:]
}

// Per section 20 of data sheet, read NUM_RXBYTES
// repeatedly until same value is returned twice.
func (r *Radio) ReadNumRxBytes() byte {
	last := byte(0)
	read := false
	for r.Error() == nil {
		n := r.ReadRegister(RXBYTES)
		if n&RXFIFO_OVERFLOW != 0 {
			r.err = RxFifoOverflow
		}
		n &= NUM_RXBYTES_MASK
		if read && n == last {
			return n
		}
		last = n
		read = true
	}
	return 0
}

func (r *Radio) ReadNumTxBytes() byte {
	n := r.ReadRegister(TXBYTES)
	if n&TXFIFO_UNDERFLOW != 0 {
		r.err = TxFifoUnderflow
	}
	return n & NUM_TXBYTES_MASK
}

func (r *Radio) State() string {
	return StateName(r.ReadState())
}

func (r *Radio) ReadState() byte {
	status := r.Strobe(SNOP)
	return (status >> STATE_SHIFT) & STATE_MASK
}

func StateName(state byte) string {
	return stateName[state]
}

func (r *Radio) ReadMarcState() byte {
	return r.ReadRegister(MARCSTATE) & MARCSTATE_MASK
}

func MarcStateName(state byte) string {
	return marcState[state]
}

var (
	stateName = []string{
		"IDLE",
		"RX",
		"TX",
		"FSTXON",
		"CALIBRATE",
		"SETTLING",
		"RXFIFO_OVERFLOW",
		"TXFIFO_UNDERFLOW",
	}
	marcState = []string{
		"SLEEP",
		"IDLE",
		"XOFF",
		"VCOON_MC",
		"REGON_MC",
		"MANCAL",
		"VCOON",
		"REGON",
		"STARTCAL",
		"BWBOOST",
		"FS_LOCK",
		"IFADCON",
		"ENDCAL",
		"RX",
		"RX_END",
		"RX_RST",
		"TXRX_SWITCH",
		"RXFIFO_OVERFLOW",
		"FSTXON",
		"TX",
		"TX_END",
		"RXTX_SWITCH",
		"TXFIFO_UNDERFLOW",
	}
	strobeString = []string{
		"SRES",
		"SFSTXON",
		"SXOFF",
		"SCAL",
		"SRX",
		"STX",
		"SIDLE",
		"SAFC",
		"SWOR",
		"SPWD",
		"SFRX",
		"SFTX",
		"SWORRST",
		"SNOP",
	}
	numPreamble = []uint8{2, 3, 4, 6, 8, 12, 16, 24}
)
