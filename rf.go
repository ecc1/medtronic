package cc1100

import (
	"errors"
)

var (
	RxFifoOverflow  = errors.New("RXFIFO overflow")
	TxFifoUnderflow = errors.New("TXFIFO underflow")
)

func (r *Radio) InitRF() error {
	err := r.WriteEach([]byte{
		// High impedance (3-state)
		IOCFG2, 0x2E,
		IOCFG1, 0x2E,
		// Assert when sync word has been sent/received
		IOCFG0, 0x06,

		// 4 bytes in RX FIFO, 61 bytes in TX FIFO
		FIFOTHR, 0x00,

		SYNC1, 0xFF,
		SYNC0, 0x00,

		PKTLEN, 0xFF,
		PKTCTRL1, 2 << PKTCTRL1_PQT_SHIFT,
		PKTCTRL0, PKTCTRL0_LENGTH_CONFIG_INFINITE,

		CHANNR, 0x00,

		// Intermediate frequency
		// 0x06 * 24 MHz / 2^10 == 140625 Hz
		FSCTRL1, 0x06,

		FSCTRL0, 0x00,

		// CHANBW_E = 2, CHANBW_M = 1, DRATE_E = 9
		// Channel BW = 24 MHz / (8 * (4 + CHANBW_M) * 2^CHANBW_E) == 150 kHz
		MDMCFG4, 2<<MDMCFG4_CHANBW_E_SHIFT | 1<<MDMCFG4_CHANBW_M_SHIFT | 9<<MDMCFG4_DRATE_E_SHIFT,

		// DRATE_M = 102 (0x66)
		// Data rate = (256 + DRATE_M) * 2^DRATE_E * 24 MHz / 2^28 == 16388 Baud
		MDMCFG3, 0x66,

		MDMCFG2, MDMCFG2_DEM_DCFILT_ON | MDMCFG2_MOD_FORMAT_ASK_OOK | MDMCFG2_SYNC_MODE_30_32_THRES,

		// CHANSPC_E = 2
		MDMCFG1, MDMCFG1_FEC_DIS | MDMCFG1_NUM_PREAMBLE_16 | 2<<MDMCFG1_CHANSPC_E_SHIFT,

		// CHANSPC_M = 26 (0x1A)
		// Channel spacing = (256 + CHANSPC_M) * 2^CHANSPC_E * 24 MHz / 2^18 == 103271 Hz
		MDMCFG0, 0x1A,

		MCSM2, MCSM2_RX_TIME_END_OF_PACKET,
		MCSM1, MCSM1_CCA_MODE_RSSI_BELOW_UNLESS_RECEIVING | MCSM1_RXOFF_MODE_IDLE | MCSM1_TXOFF_MODE_IDLE,
		MCSM0, MCSM0_FS_AUTOCAL_FROM_IDLE,
		FOCCFG, FOCCFG_FOC_PRE_K_3K | FOCCFG_FOC_POST_K_PRE_K_OVER_2 | FOCCFG_FOC_LIMIT_BW_OVER_2,
		BSCFG, BSCFG_BS_PRE_K_2K | BSCFG_BS_PRE_KP_3KP | BSCFG_BS_POST_KI_PRE_KI_OVER_2 | BSCFG_BS_LIMIT_0,
		AGCCTRL2, AGCCTRL2_MAX_DVGA_GAIN_ALL | AGCCTRL2_MAX_LNA_GAIN_0 | AGCCTRL2_MAGN_TARGET_38dB,
		AGCCTRL1, AGCCTRL1_AGC_LNA_PRIORITY_0 | AGCCTRL1_CARRIER_SENSE_REL_THR_DISABLE | AGCCTRL1_CARRIER_SENSE_ABS_THR_0DB,
		AGCCTRL0, AGCCTRL0_HYST_LEVEL_MEDIUM | AGCCTRL0_WAIT_TIME_16 | AGCCTRL0_AGC_FREEZE_NORMAL | AGCCTRL0_FILTER_LENGTH_32,
		FREND1, 1<<FREND1_LNA_CURRENT_SHIFT | 1<<FREND1_LNA2MIX_CURRENT_SHIFT | 1<<FREND1_LODIV_BUF_CURRENT_RX_SHIFT | 2<<FREND1_MIX_CURRENT_SHIFT,

		// Use PA_TABLE 1 for transmitting '1' in ASK
		// (PA_TABLE 0 is always used for '0')
		FREND0, 1<<FREND0_LODIV_BUF_CURRENT_TX_SHIFT | 1<<FREND0_PA_POWER_SHIFT,

		FSCAL3, 3<<6 | 2<<4 | 0x09,
		FSCAL2, 1<<5 | 0x0A, // VCO high
		FSCAL1, 0x00,
		FSCAL0, 0x1F,

		TEST2, TEST2_RX_LOW_DATA_RATE_MAGIC,
		TEST1, TEST1_RX_LOW_DATA_RATE_MAGIC,
		TEST0, 2<<2 | 1, // disable VCO selection calibration
	})
	if err != nil {
		return err
	}

	// Power amplifier output settings (see section 24 of the data sheet)
	err = r.WriteBurst(PATABLE, []byte{0x00, 0xC0})
	if err != nil {
		return err
	}

	return nil
}

func (r *Radio) Frequency() (uint32, error) {
	freq, err := r.ReadBurst(FREQ2, 3)
	if err != nil {
		return 0, err
	}
	f := uint32(freq[0])<<16 + uint32(freq[1])<<8 + uint32(freq[2])
	return uint32(uint64(f) * FXOSC >> 16), nil
}

func (r *Radio) SetFrequency(freq uint32) error {
	f := (uint64(freq)<<16 + FXOSC/2) / FXOSC
	return r.WriteBurst(FREQ2, []byte{
		byte(f >> 16),
		byte(f >> 8),
		byte(f),
	})
}

func (r *Radio) ReadIF() (uint32, error) {
	f, err := r.ReadRegister(FSCTRL1)
	if err != nil {
		return 0, err
	}
	return uint32(uint64(f) * FXOSC >> 10), nil
}

func (r *Radio) ReadChannelParams() (chanbw uint32, drate uint32, err error) {
	var m4 byte
	m4, err = r.ReadRegister(MDMCFG4)
	if err != nil {
		return
	}
	chanbw_E := (m4 >> MDMCFG4_CHANBW_E_SHIFT) & 0x3
	chanbw_M := (m4 >> MDMCFG4_CHANBW_M_SHIFT) & 0x3
	drate_E := (m4 >> MDMCFG4_DRATE_E_SHIFT) & 0xF

	var drate_M byte
	drate_M, err = r.ReadRegister(MDMCFG3)
	if err != nil {
		return
	}

	chanbw = uint32(FXOSC / ((4 + uint64(chanbw_M)) << (chanbw_E + 3)))
	drate = uint32(((256 + uint64(drate_M)) << drate_E * FXOSC) >> 28)
	return
}

func (r *Radio) ReadModemConfig() (fec bool, minPreamble byte, chanspc uint32, err error) {
	var m1 byte
	m1, err = r.ReadRegister(MDMCFG1)
	if err != nil {
		return
	}
	fec = m1&MDMCFG1_FEC_EN != 0
	minPreamble = numPreamble[(m1&MDMCFG1_NUM_PREAMBLE_MASK)>>4]
	chanspc_E := m1 & MDMCFG1_CHANSPC_E_MASK
	var chanspc_M byte
	chanspc_M, err = r.ReadRegister(MDMCFG0)
	if err != nil {
		return
	}
	chanspc = uint32(((256 + uint64(chanspc_M)) << chanspc_E * FXOSC) >> 18)
	return
}

func (r *Radio) ReadRSSI() (int, error) {
	const rssi_offset = 74 // see data sheet section 17.3
	rssi, err := r.ReadRegister(RSSI)
	if err != nil {
		return 0, err
	}
	d := int(rssi)
	if d >= 128 {
		d -= 256
	}
	return d/2 - rssi_offset, nil
}

func (r *Radio) ReadPaTable() ([]byte, error) {
	buf := make([]byte, 9)
	buf[0] = READ_MODE | BURST_MODE | PATABLE
	err := r.device.Transfer(buf)
	if err != nil {
		return nil, err
	}
	return buf[1:], nil
}

// Per section 20 of data sheet, read NUM_RXBYTES
// repeatedly until same value is returned twice.
func (r *Radio) ReadNumRxBytes() (byte, error) {
	last := byte(0)
	read := false
	for {
		n, err := r.ReadRegister(RXBYTES)
		if err != nil {
			return 0, err
		}
		if n&RXFIFO_OVERFLOW != 0 {
			err = RxFifoOverflow
		}
		n &= NUM_RXBYTES_MASK
		if read && n == last {
			return n, err
		}
		last = n
		read = true
	}
}

func (r *Radio) ReadNumTxBytes() (byte, error) {
	n, err := r.ReadRegister(TXBYTES)
	if err != nil {
		return 0, err
	}
	if n&TXFIFO_UNDERFLOW != 0 {
		err = TxFifoUnderflow
	}
	return n & NUM_TXBYTES_MASK, err
}

func (r *Radio) State() string {
	s, err := r.ReadState()
	if err != nil {
		return err.Error()
	}
	return StateName(s)
}

func (r *Radio) ReadState() (byte, error) {
	status, err := r.Strobe(SNOP)
	if err != nil {
		return 0, err
	}
	return (status >> STATE_SHIFT) & STATE_MASK, nil
}

func StateName(state byte) string {
	return stateName[state]
}

func (r *Radio) ReadMarcState() (byte, error) {
	state, err := r.ReadRegister(MARCSTATE)
	if err != nil {
		return 0, err
	}
	return state & MARCSTATE_MASK, nil
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
)
