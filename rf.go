package cc1100

import (
	"github.com/ecc1/spi"
)

func InitRF(dev *spi.Device) error {
	err := Write(dev, []byte{
		IOCFG2, 0x0E,
		IOCFG1, GDO1_DS,
		IOCFG0, 0x06,

		// sync word
		SYNC1, 0xFF,
		SYNC0, 0x00,

		// packet length (default)
		PKTLEN, 0xFF,

		// always accept sync word
		// do not append status
		// no address check

		PKTCTRL1, 0x00,

		// no whitening mode
		// normal format
		// disable CRC calculation and check
		// fixed packet length mode
		PKTCTRL0, 0x00,

		// channel number
		CHANNR, 0x00,

		// intermediate frequency
		// 0x06 * 26 MHz / 2^10 == 152 kHz
		FSCTRL1, 0x06,

		// frequency offset
		FSCTRL0, 0x00,

		// 24-bit base frequency
		// 0x2340FC * 26 MHz / 2^16 == 916.6 MHz (916599975 Hz)
		FREQ2, 0x23,
		FREQ1, 0x40,
		FREQ0, 0xFC,

		// CHANBW_E = 1, CHANBW_M = 1, DRATE_E = 9
		// channel BW = 26 MHz / (8 * (4 + CHANBW_M) * 2^CHANBW_E) == 325 kHz
		MDMCFG4, ((1 << MDMCFG4_CHANBW_E_SHIFT) |
			(1 << MDMCFG4_CHANBW_M_SHIFT) |
			(9 << MDMCFG4_DRATE_E_SHIFT)),

		// DRATE_M = 74 (0x4A)
		// data rate = (256 + DRATE_M) * 2^DRATE_E * 26 MHz / 2^28 == 16365 Baud
		MDMCFG3, 0x4A,

		MDMCFG2, (MDMCFG2_DEM_DCFILT_ON |
			MDMCFG2_MOD_FORMAT_ASK_OOK |
			MDMCFG2_SYNC_MODE_30_32),

		// CHANSPC_E = 1
		MDMCFG1, (MDMCFG1_FEC_DIS |
			MDMCFG1_NUM_PREAMBLE_16 |
			(1 << MDMCFG1_CHANSPC_E_SHIFT)),

		// CHANSPC_M = 248 (0xF8)
		// channel spacing
		// (256 + CHANSPC_M) * 2^CHANSPC_E * 26 MHz / 2^18 == 99975 Hz
		MDMCFG0, 0xF8,

		MCSM2, MCSM2_RX_TIME_END_OF_PACKET, // (default)

		MCSM1, (MCSM1_CCA_MODE_RSSI_BELOW_UNLESS_RECEIVING |
			MCSM1_RXOFF_MODE_IDLE |
			MCSM1_TXOFF_MODE_IDLE), // (default)

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
			AGCCTRL2_MAGN_TARGET_33dB), // (default)

		AGCCTRL1, (AGCCTRL1_AGC_LNA_PRIORITY_1 |
			AGCCTRL1_CARRIER_SENSE_REL_THR_DISABLE |
			AGCCTRL1_CARRIER_SENSE_ABS_THR_0DB), // (default)

		AGCCTRL0, (AGCCTRL0_HYST_LEVEL_MEDIUM |
			AGCCTRL0_WAIT_TIME_16 |
			AGCCTRL0_AGC_FREEZE_NORMAL |
			AGCCTRL0_FILTER_LENGTH_16), // (default)

		FREND1, ((1 << FREND1_LNA_CURRENT_SHIFT) |
			(1 << FREND1_LNA2MIX_CURRENT_SHIFT) |
			(1 << FREND1_LODIV_BUF_CURRENT_RX_SHIFT) |
			(2 << FREND1_MIX_CURRENT_SHIFT)), // (default)

		// use PA_TABLE 1 for transmitting '1' in ASK
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
	// (see Table 72 on page 207 of the datasheet)
	err = dev.Transfer([]byte{
		WRITE_BURST | PATABLE,
		0x00,
		0xC0, // 10dBm, 36mA
	})
	if err != nil {
		return err
	}

	return nil
}

func ReadFrequency(dev *spi.Device) (uint32, error) {
	freq2, err := ReadRegister(dev, FREQ2)
	if err != nil {
		return 0, err
	}
	freq1, err := ReadRegister(dev, FREQ1)
	if err != nil {
		return 0, err
	}
	freq0, err := ReadRegister(dev, FREQ0)
	if err != nil {
		return 0, err
	}
	f := uint32(freq2)<<16 + uint32(freq1)<<8 + uint32(freq0)
	return uint32(uint64(f) * FXOSC >> 16), nil
}

func ReadIF(dev *spi.Device) (uint32, error) {
	f, err := ReadRegister(dev, FSCTRL1)
	if err != nil {
		return 0, err
	}
	return uint32(uint64(f) * FXOSC >> 10), nil
}

func ReadChannelParams(dev *spi.Device) (chanbw uint32, drate uint32, err error) {
	var m4 byte
	m4, err = ReadRegister(dev, MDMCFG4)
	if err != nil {
		return
	}
	chanbw_E := (m4 >> MDMCFG4_CHANBW_E_SHIFT) & 0x3
	chanbw_M := (m4 >> MDMCFG4_CHANBW_M_SHIFT) & 0x3
	drate_E := (m4 >> MDMCFG4_DRATE_E_SHIFT) & 0xF

	var drate_M byte
	drate_M, err = ReadRegister(dev, MDMCFG3)
	if err != nil {
		return
	}

	chanbw = uint32(FXOSC / ((4 + uint64(chanbw_M)) << (chanbw_E + 3)))
	drate = uint32((((256+uint64(drate_M))<<drate_E)*FXOSC) >> 28)
	return
}
