package cc1100

func (dev *Device) InitRF() error {
	err := dev.WriteEach([]byte{
		// Carrier sense: high if RSSI level is above threshold
		IOCFG2, 0x0E,
		IOCFG1, 0x00,
		// Assert when sync word has been sent/received
		IOCFG0, 0x06,

		// Sync word
		SYNC1, 0xFF,
		SYNC0, 0x00,

		// Packet length
		PKTLEN, 0xFF,

		// Always accept sync word
		// Do not append status
		// No address check
		PKTCTRL1, 2 << PKTCTRL1_PQT_SHIFT,

		// No whitening mode
		// Normal format
		// Disable CRC calculation and check
		// Fixed packet length mode
		PKTCTRL0, PKTCTRL0_LENGTH_CONFIG_VARIABLE,

		// Channel number
		CHANNR, 0x00,

		// Intermediate frequency
		// 0x06 * 26 MHz / 2^10 == 152 kHz
		FSCTRL1, 0x06,

		// Frequency offset
		FSCTRL0, 0x00,

		// 24-bit base frequency
		// 0x2340FC * 26 MHz / 2^16 == 916.6 MHz (916599975 Hz)
		FREQ2, 0x23,
		FREQ1, 0x40,
		FREQ0, 0xFC,

		// CHANBW_E = 2, CHANBW_M = 1, DRATE_E = 9
		// Channel BW = 26 MHz / (8 * (4 + CHANBW_M) * 2^CHANBW_E) == 162.5 kHz
		MDMCFG4, ((2 << MDMCFG4_CHANBW_E_SHIFT) |
			(1 << MDMCFG4_CHANBW_M_SHIFT) |
			(9 << MDMCFG4_DRATE_E_SHIFT)),

		// DRATE_M = 74 (0x4A)
		// Data rate = (256 + DRATE_M) * 2^DRATE_E * 26 MHz / 2^28 == 16365 Baud
		MDMCFG3, 0x4A,

		MDMCFG2, (MDMCFG2_DEM_DCFILT_ON |
			MDMCFG2_MOD_FORMAT_ASK_OOK |
			MDMCFG2_SYNC_MODE_30_32_THRES),

		// CHANSPC_E = 1
		MDMCFG1, (MDMCFG1_FEC_DIS |
			MDMCFG1_NUM_PREAMBLE_16 |
			(1 << MDMCFG1_CHANSPC_E_SHIFT)),

		// CHANSPC_M = 248 (0xF8)
		// Channel spacing = (256 + CHANSPC_M) * 2^CHANSPC_E * 26 MHz / 2^18 == 99975 Hz
		MDMCFG0, 0xF8,

		MCSM2, MCSM2_RX_TIME_END_OF_PACKET,

		MCSM1, (MCSM1_CCA_MODE_RSSI_BELOW_UNLESS_RECEIVING |
			MCSM1_RXOFF_MODE_IDLE |
			MCSM1_TXOFF_MODE_IDLE),

		MCSM0, (MCSM0_FS_AUTOCAL_FROM_IDLE |
			MCSM0_MAGIC_3 |
			MCSM0_CLOSE_IN_RX_0DB),

		FOCCFG, (FOCCFG_FOC_PRE_K_3K |
			FOCCFG_FOC_POST_K_PRE_K_OVER_2 |
			FOCCFG_FOC_LIMIT_BW_OVER_2),

		BSCFG, (BSCFG_BS_PRE_K_2K |
			BSCFG_BS_PRE_KP_3KP |
			BSCFG_BS_POST_KI_PRE_KI_OVER_2 |
			BSCFG_BS_LIMIT_0),

		AGCCTRL2, (AGCCTRL2_MAX_DVGA_GAIN_ALL |
			AGCCTRL2_MAX_LNA_GAIN_0 |
			AGCCTRL2_MAGN_TARGET_38dB),

		AGCCTRL1, (AGCCTRL1_AGC_LNA_PRIORITY_0 |
			AGCCTRL1_CARRIER_SENSE_REL_THR_DISABLE |
			AGCCTRL1_CARRIER_SENSE_ABS_THR_0DB),

		AGCCTRL0, (AGCCTRL0_HYST_LEVEL_MEDIUM |
			AGCCTRL0_WAIT_TIME_16 |
			AGCCTRL0_AGC_FREEZE_NORMAL |
			AGCCTRL0_FILTER_LENGTH_32),

		FREND1, ((1 << FREND1_LNA_CURRENT_SHIFT) |
			(1 << FREND1_LNA2MIX_CURRENT_SHIFT) |
			(1 << FREND1_LODIV_BUF_CURRENT_RX_SHIFT) |
			(2 << FREND1_MIX_CURRENT_SHIFT)),

		// Use PA_TABLE 1 for transmitting '1' in ASK
		// (PA_TABLE 0 is always used for '0')
		FREND0, ((1 << FREND0_LODIV_BUF_CURRENT_TX_SHIFT) |
			(1 << FREND0_PA_POWER_SHIFT)),

		FSCAL3, (3 << 6) | (2 << 4) | 0x09,
		FSCAL2, (1 << 5) | 0x0A, // VCO high
		FSCAL1, 0x00,
		FSCAL0, 0x1F,

		TEST2, TEST2_RX_LOW_DATA_RATE_MAGIC,
		TEST1, TEST1_RX_LOW_DATA_RATE_MAGIC,
		TEST0, (2 << 2) | 1, // disable VCO selection calibration
	})
	if err != nil {
		return err
	}

	// Power amplifier output settings
	// (see section 24 of the datasheet)
	err = dev.spiDev.Write([]byte{
		BURST_MODE | PATABLE,
		0x00,
		0xC0,
	})
	if err != nil {
		return err
	}

	return nil
}

func (dev *Device) ReadFrequency() (uint32, error) {
	freq2, err := dev.ReadRegister(FREQ2)
	if err != nil {
		return 0, err
	}
	freq1, err := dev.ReadRegister(FREQ1)
	if err != nil {
		return 0, err
	}
	freq0, err := dev.ReadRegister(FREQ0)
	if err != nil {
		return 0, err
	}
	f := uint32(freq2)<<16 + uint32(freq1)<<8 + uint32(freq0)
	return uint32(uint64(f) * FXOSC >> 16), nil
}

func (dev *Device) WriteFrequency(freq uint32) error {
	f := (uint64(freq)<<16 + FXOSC/2) / FXOSC
	return dev.WriteEach([]byte{
		FREQ2, byte(f >> 16),
		FREQ1, byte(f >> 8),
		FREQ0, byte(f),
	})
}

func (dev *Device) ReadIF() (uint32, error) {
	f, err := dev.ReadRegister(FSCTRL1)
	if err != nil {
		return 0, err
	}
	return uint32(uint64(f) * FXOSC >> 10), nil
}

func (dev *Device) ReadChannelParams() (chanbw uint32, drate uint32, err error) {
	var m4 byte
	m4, err = dev.ReadRegister(MDMCFG4)
	if err != nil {
		return
	}
	chanbw_E := (m4 >> MDMCFG4_CHANBW_E_SHIFT) & 0x3
	chanbw_M := (m4 >> MDMCFG4_CHANBW_M_SHIFT) & 0x3
	drate_E := (m4 >> MDMCFG4_DRATE_E_SHIFT) & 0xF

	var drate_M byte
	drate_M, err = dev.ReadRegister(MDMCFG3)
	if err != nil {
		return
	}

	chanbw = uint32(FXOSC / ((4 + uint64(chanbw_M)) << (chanbw_E + 3)))
	drate = uint32((((256 + uint64(drate_M)) << drate_E) * FXOSC) >> 28)
	return
}

func (dev *Device) ReadModemConfig() (fec bool, minPreamble byte, chanspc uint32, err error) {
	var m1 byte
	m1, err = dev.ReadRegister(MDMCFG1)
	if err != nil {
		return
	}
	fec = m1&MDMCFG1_FEC_EN != 0
	minPreamble = numPreamble[(m1&MDMCFG1_NUM_PREAMBLE_MASK)>>4]
	chanspc_E := m1 & MDMCFG1_CHANSPC_E_MASK
	var chanspc_M byte
	chanspc_M, err = dev.ReadRegister(MDMCFG0)
	if err != nil {
		return
	}
	chanspc = uint32((((256 + uint64(chanspc_M)) << chanspc_E) * FXOSC) >> 18)
	return
}

func (dev *Device) ReadRSSI() (int, error) {
	const rssi_offset = 74 // see data sheet section 17.3
	r, err := dev.ReadRegister(RSSI)
	if err != nil {
		return 0, err
	}
	d := int(r)
	if d >= 128 {
		d -= 256
	}
	return d/2 - rssi_offset, nil
}

func (dev *Device) ReadPaTable() ([]byte, error) {
	buf := make([]byte, 9)
	buf[0] = READ_MODE | BURST_MODE | PATABLE
	err := dev.spiDev.Transfer(buf)
	if err != nil {
		return nil, err
	}
	return buf[1:], nil
}
