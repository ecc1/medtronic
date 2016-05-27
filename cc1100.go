package cc1100

// Register definitions for Texas Instruments CC1101.

const (
	// Crystal frequency in Hz.
	FXOSC = 24000000

	// SPI transaction header bits for read/write and burst/single access.
	READ_MODE  = 1 << 7
	BURST_MODE = 1 << 6
)

// Configuration registers (read/write).
const (
	IOCFG2   = 0x00 // GDO2 output pin configuration
	IOCFG1   = 0x01 // GDO1 output pin configuration
	IOCFG0   = 0x02 // GDO0 output pin configuration
	FIFOTHR  = 0x03 // RX FIFO and TX FIFO thresholds
	SYNC1    = 0x04 // Sync word, high byte
	SYNC0    = 0x05 // Sync word, low byte
	PKTLEN   = 0x06 // Packet length
	PKTCTRL1 = 0x07 // Packet automation control
	PKTCTRL0 = 0x08 // Packet automation control
	ADDR     = 0x09 // Device address
	CHANNR   = 0x0A // Channel number
	FSCTRL1  = 0x0B // Frequency synthesizer control
	FSCTRL0  = 0x0C // Frequency synthesizer control
	FREQ2    = 0x0D // Frequency control word, high byte
	FREQ1    = 0x0E // Frequency control word, middle byte
	FREQ0    = 0x0F // Frequency control word, low byte
	MDMCFG4  = 0x10 // Modem configuration
	MDMCFG3  = 0x11 // Modem configuration
	MDMCFG2  = 0x12 // Modem configuration
	MDMCFG1  = 0x13 // Modem configuration
	MDMCFG0  = 0x14 // Modem configuration
	DEVIATN  = 0x15 // Modem deviation setting
	MCSM2    = 0x16 // Main Radio Control State Machine configuration
	MCSM1    = 0x17 // Main Radio Control State Machine configuration
	MCSM0    = 0x18 // Main Radio Control State Machine configuration
	FOCCFG   = 0x19 // Frequency Offset Compensation configuration
	BSCFG    = 0x1A // Bit Synchronization configuration
	AGCCTRL2 = 0x1B // AGC control
	AGCCTRL1 = 0x1C // AGC control
	AGCCTRL0 = 0x1D // AGC control
	WOREVT1  = 0x1E // High byte Event 0 timeout
	WOREVT0  = 0x1F // Low byte Event 0 timeout
	WORCTRL  = 0x20 // Wake On Radio control
	FREND1   = 0x21 // Front end RX configuration
	FREND0   = 0x22 // Front end TX configuration
	FSCAL3   = 0x23 // Frequency synthesizer calibration
	FSCAL2   = 0x24 // Frequency synthesizer calibration
	FSCAL1   = 0x25 // Frequency synthesizer calibration
	FSCAL0   = 0x26 // Frequency synthesizer calibration
	RCCTRL1  = 0x27 // RC oscillator configuration
	RCCTRL0  = 0x28 // RC oscillator configuration
	FSTEST   = 0x29 // Frequency synthesizer calibration control
	PTEST    = 0x2A // Production test
	AGCTEST  = 0x2B // AGC test
	TEST2    = 0x2C // Various test settings
	TEST1    = 0x2D // Various test settings
	TEST0    = 0x2E // Various test settings
)

// Command strobes (write-only).
const (
	// Reset chip.
	SRES = 0x30

	// Enable and calibrate frequency synthesizer
	// (if MCSM0.FS_AUTOCAL=1). If in RX/TX (with CCA):
	// Go to a wait state where only the synthesizer is running
	// (for quick RX / TX turnaround).
	SFSTXON = 0x31

	// Turn off crystal oscillator.
	SXOFF = 0x32

	// Calibrate frequency synthesizer and turn it off.
	// SCAL can be strobed from IDLE mode without setting
	//  manual calibration mode (MCSM0.FS_AUTOCAL=0)
	SCAL = 0x33

	// Enable RX. Perform calibration first if coming from IDLE
	// and MCSM0.FS_AUTOCAL=1.
	SRX = 0x34

	// In IDLE state: Enable TX. Perform calibration first if
	// MCSM0.FS_AUTOCAL=1. If in RX state and CCA is enabled:
	// Only go to TX if channel is clear.
	STX = 0x35

	// Exit RX / TX, turn off frequency synthesizer and exit
	// Wake-On-Radio mode if applicable.
	SIDLE = 0x36

	// Perform AFC adjustment of the frequency synthesizer.
	SAFC = 0x37

	// Start automatic RX polling sequence (Wake-on-Radio)
	// if WORCTRL.RC_PD=0.
	SWOR = 0x38

	// Enter power down mode when CSn goes high.
	SPWD = 0x39

	// Flush the RX FIFO buffer.
	// Only issue SFRX in IDLE or RXFIFO_OVERFLOW states.
	SFRX = 0x3A

	// Flush the TX FIFO buffer.
	// Only issue SFTX in IDLE or TXFIFO_UNDERFLOW states.
	SFTX = 0x3B

	// Reset real time clock to Event1 value.
	SWORRST = 0x3C

	// No operation. May be used to get access to the chip status byte.
	SNOP = 0x3D
)

// Status registers (read-only).
// Since these must be read with the burst access bit set,
// it is included in the address for simplicity.
const (
	PARTNUM        = 0x70 // Part number for CC1101
	VERSION        = 0x71 // Current version number
	FREQEST        = 0x72 // Frequency Offset Estimate
	LQI            = 0x73 // Demodulator estimate for Link Quality
	RSSI           = 0x74 // Received signal strength indication
	MARCSTATE      = 0x75 // Control state machine state
	WORTIME1       = 0x76 // High byte of WOR timer
	WORTIME0       = 0x77 // Low byte of WOR timer
	PKTSTATUS      = 0x78 // Current GDOx status and packet status
	VCO_VC_DAC     = 0x79 // Current setting from PLL calibration module
	TXBYTES        = 0x7A // Underflow and number of bytes in the TX FIFO
	RXBYTES        = 0x7B // Overflow and number of bytes in the RX FIFO
	RCCTRL1_STATUS = 0x7C // Last RC oscillator calibration result
	RCCTRL0_STATUS = 0x7D // Last RC oscillator calibration result
)

const (
	// PA Table
	PATABLE = 0x3E

	// FIFOs
	TXFIFO = 0x3F
	RXFIFO = 0x3F
)

const (
	STATE_IDLE = iota
	STATE_RX
	STATE_TX
	STATE_FSTXON
	STATE_CALIBRATE
	STATE_SETTLING
	STATE_RXFIFO_OVERFLOW
	STATE_TXFIFO_UNDERFLOW

	// status bits 6:4
	STATE_MASK  = 0x7
	STATE_SHIFT = 4

	CHIP_RDY = 0x80
)

const (
	GDO2_INV      = 1 << 6
	GDO2_CFG_MASK = 0x3f

	GDO1_DS       = 1 << 7
	GDO1_INV      = 1 << 6
	GDO1_CFG_MASK = 0x3f

	GDO0_TEMP_SENSOR_ENABLE = 1 << 7
	GDO0_INV                = 1 << 6
	GDO0_CFG_MASK           = 0x3f

	FIFOTHR_MASK = 0xf
	// FiFOTHR value n corresponds to 4*(n+1) bytes in RX FIFO
	// or 65 - 4*(n+1) bytes in TX FIFO

	PKTCTRL1_PQT_MASK                = 0x7 << 5
	PKTCTRL1_PQT_SHIFT               = 5
	PKTCTRL1_APPEND_STATUS           = 1 << 2
	PKTCTRL1_ADR_CHK_NONE            = 0 << 0
	PKTCTRL1_ADR_CHK_NO_BROADCAST    = 1 << 0
	PKTCTRL1_ADR_CHK_00_BROADCAST    = 2 << 0
	PKTCTRL1_ADR_CHK_00_FF_BROADCAST = 3 << 0

	PKT_APPEND_STATUS_0_RSSI_MASK  = 0xff
	PKT_APPEND_STATUS_0_RSSI_SHIFT = 0
	PKT_APPEND_STATUS_1_CRC_OK     = 1 << 7
	PKT_APPEND_STATUS_1_LQI_MASK   = 0x7f
	PKT_APPEND_STATUS_1_LQI_SHIFT  = 0

	PKTCTRL0_WHITE_DATA             = 1 << 6
	PKTCTRL0_PKT_FORMAT_NORMAL      = 0 << 4
	PKTCTRL0_PKT_FORMAT_RANDOM      = 2 << 4
	PKTCTRL0_CRC_EN                 = 1 << 2
	PKTCTRL0_LENGTH_CONFIG_FIXED    = 0 << 0
	PKTCTRL0_LENGTH_CONFIG_VARIABLE = 1 << 0
	PKTCTRL0_LENGTH_CONFIG_INFINITE = 2 << 0

	MDMCFG4_CHANBW_E_SHIFT = 6
	MDMCFG4_CHANBW_M_SHIFT = 4
	MDMCFG4_DRATE_E_SHIFT  = 0

	MDMCFG3_DRATE_M_SHIFT = 0

	MDMCFG2_DEM_DCFILT_OFF = 1 << 7
	MDMCFG2_DEM_DCFILT_ON  = 0 << 7

	MDMCFG2_MOD_FORMAT_MASK    = 7 << 4
	MDMCFG2_MOD_FORMAT_2_FSK   = 0 << 4
	MDMCFG2_MOD_FORMAT_GFSK    = 1 << 4
	MDMCFG2_MOD_FORMAT_ASK_OOK = 3 << 4
	MDMCFG2_MOD_FORMAT_MSK     = 7 << 4

	MDMCFG2_SYNC_MODE_MASK        = 0x7 << 0
	MDMCFG2_SYNC_MODE_NONE        = 0x0 << 0
	MDMCFG2_SYNC_MODE_15_16       = 0x1 << 0
	MDMCFG2_SYNC_MODE_16_16       = 0x2 << 0
	MDMCFG2_SYNC_MODE_30_32       = 0x3 << 0
	MDMCFG2_SYNC_MODE_NONE_THRES  = 0x4 << 0
	MDMCFG2_SYNC_MODE_15_16_THRES = 0x5 << 0
	MDMCFG2_SYNC_MODE_16_16_THRES = 0x6 << 0
	MDMCFG2_SYNC_MODE_30_32_THRES = 0x7 << 0

	MDMCFG1_FEC_EN  = 1 << 7
	MDMCFG1_FEC_DIS = 0 << 7

	MDMCFG1_NUM_PREAMBLE_MASK = 7 << 4
	MDMCFG1_NUM_PREAMBLE_2    = 0 << 4
	MDMCFG1_NUM_PREAMBLE_3    = 1 << 4
	MDMCFG1_NUM_PREAMBLE_4    = 2 << 4
	MDMCFG1_NUM_PREAMBLE_6    = 3 << 4
	MDMCFG1_NUM_PREAMBLE_8    = 4 << 4
	MDMCFG1_NUM_PREAMBLE_12   = 5 << 4
	MDMCFG1_NUM_PREAMBLE_16   = 6 << 4
	MDMCFG1_NUM_PREAMBLE_24   = 7 << 4

	MDMCFG1_CHANSPC_E_MASK  = 3 << 0
	MDMCFG1_CHANSPC_E_SHIFT = 0

	MDMCFG0_CHANSPC_M_SHIFT = 0

	DEVIATN_DEVIATION_E_SHIFT = 4
	DEVIATN_DEVIATION_M_SHIFT = 0

	MCSM2_RX_TIME_RSSI          = 1 << 4
	MCSM2_RX_TIME_QUAL          = 1 << 3
	MCSM2_RX_TIME_MASK          = 0x7 << 0
	MCSM2_RX_TIME_SHIFT         = 0
	MCSM2_RX_TIME_END_OF_PACKET = 7

	MCSM1_CCA_MODE_ALWAYS                      = 0 << 4
	MCSM1_CCA_MODE_RSSI_BELOW                  = 1 << 4
	MCSM1_CCA_MODE_UNLESS_RECEIVING            = 2 << 4
	MCSM1_CCA_MODE_RSSI_BELOW_UNLESS_RECEIVING = 3 << 4
	MCSM1_RXOFF_MODE_IDLE                      = 0 << 2
	MCSM1_RXOFF_MODE_FSTXON                    = 1 << 2
	MCSM1_RXOFF_MODE_TX                        = 2 << 2
	MCSM1_RXOFF_MODE_RX                        = 3 << 2
	MCSM1_TXOFF_MODE_IDLE                      = 0 << 0
	MCSM1_TXOFF_MODE_FSTXON                    = 1 << 0
	MCSM1_TXOFF_MODE_TX                        = 2 << 0
	MCSM1_TXOFF_MODE_RX                        = 3 << 0

	MCSM0_FS_AUTOCAL_NEVER           = 0 << 4
	MCSM0_FS_AUTOCAL_FROM_IDLE       = 1 << 4
	MCSM0_FS_AUTOCAL_TO_IDLE         = 2 << 4
	MCSM0_FS_AUTOCAL_TO_IDLE_EVERY_4 = 3 << 4
	MCSM0_PO_TIMEOUT_SHIFT           = 2
	MCSM0_PO_TIMEOUT_MASK            = 0x3 << 2
	MCSM0_PIN_CTRL_EN                = 1 << 1
	MCSM0_XOSC_FORCE_ON              = 1 << 0

	FOCCFG_FOC_BS_CS_GATE          = 1 << 5
	FOCCFG_FOC_PRE_K_1K            = 0 << 3
	FOCCFG_FOC_PRE_K_2K            = 1 << 3
	FOCCFG_FOC_PRE_K_3K            = 2 << 3
	FOCCFG_FOC_PRE_K_4K            = 3 << 3
	FOCCFG_FOC_POST_K_PRE_K        = 0 << 2
	FOCCFG_FOC_POST_K_PRE_K_OVER_2 = 1 << 2
	FOCCFG_FOC_LIMIT_0             = 0 << 0
	FOCCFG_FOC_LIMIT_BW_OVER_8     = 1 << 0
	FOCCFG_FOC_LIMIT_BW_OVER_4     = 2 << 0
	FOCCFG_FOC_LIMIT_BW_OVER_2     = 3 << 0

	BSCFG_BS_PRE_K_1K              = 0 << 6
	BSCFG_BS_PRE_K_2K              = 1 << 6
	BSCFG_BS_PRE_K_3K              = 2 << 6
	BSCFG_BS_PRE_K_4K              = 3 << 6
	BSCFG_BS_PRE_KP_1KP            = 0 << 4
	BSCFG_BS_PRE_KP_2KP            = 1 << 4
	BSCFG_BS_PRE_KP_3KP            = 2 << 4
	BSCFG_BS_PRE_KP_4KP            = 3 << 4
	BSCFG_BS_POST_KI_PRE_KI        = 0 << 3
	BSCFG_BS_POST_KI_PRE_KI_OVER_2 = 1 << 3
	BSCFG_BS_POST_KP_PRE_KP        = 0 << 2
	BSCFG_BS_POST_KP_KP            = 1 << 2
	BSCFG_BS_LIMIT_0               = 0 << 0
	BSCFG_BS_LIMIT_3_125           = 1 << 0
	BSCFG_BS_LIMIT_6_25            = 2 << 0
	BSCFG_BS_LIMIT_12_5            = 3 << 0

	AGCCTRL2_MAX_DVGA_GAIN_ALL   = 0 << 6
	AGCCTRL2_MAX_DVGA_GAIN_BUT_1 = 1 << 6
	AGCCTRL2_MAX_DVGA_GAIN_BUT_2 = 2 << 6
	AGCCTRL2_MAX_DVGA_GAIN_BUT_3 = 3 << 6
	AGCCTRL2_MAX_LNA_GAIN_0      = 0 << 3
	AGCCTRL2_MAX_LNA_GAIN_2_6    = 1 << 3
	AGCCTRL2_MAX_LNA_GAIN_6_1    = 2 << 3
	AGCCTRL2_MAX_LNA_GAIN_7_4    = 3 << 3
	AGCCTRL2_MAX_LNA_GAIN_9_2    = 4 << 3
	AGCCTRL2_MAX_LNA_GAIN_11_5   = 5 << 3
	AGCCTRL2_MAX_LNA_GAIN_14_6   = 6 << 3
	AGCCTRL2_MAX_LNA_GAIN_17_1   = 7 << 3
	AGCCTRL2_MAGN_TARGET_24dB    = 0 << 0
	AGCCTRL2_MAGN_TARGET_27dB    = 1 << 0
	AGCCTRL2_MAGN_TARGET_30dB    = 2 << 0
	AGCCTRL2_MAGN_TARGET_33dB    = 3 << 0
	AGCCTRL2_MAGN_TARGET_36dB    = 4 << 0
	AGCCTRL2_MAGN_TARGET_38dB    = 5 << 0
	AGCCTRL2_MAGN_TARGET_40dB    = 6 << 0
	AGCCTRL2_MAGN_TARGET_42dB    = 7 << 0

	AGCCTRL1_AGC_LNA_PRIORITY_0              = 0 << 6
	AGCCTRL1_AGC_LNA_PRIORITY_1              = 1 << 6
	AGCCTRL1_CARRIER_SENSE_REL_THR_DISABLE   = 0 << 4
	AGCCTRL1_CARRIER_SENSE_REL_THR_6DB       = 1 << 4
	AGCCTRL1_CARRIER_SENSE_REL_THR_10DB      = 2 << 4
	AGCCTRL1_CARRIER_SENSE_REL_THR_14DB      = 3 << 4
	AGCCTRL1_CARRIER_SENSE_ABS_THR_DISABLE   = 0x8 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_7DB_BELOW = 0x9 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_6DB_BELOW = 0xa << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_5DB_BELOW = 0xb << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_4DB_BELOW = 0xc << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_3DB_BELOW = 0xd << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_2DB_BELOW = 0xe << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_1DB_BELOW = 0xf << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_0DB       = 0x0 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_1DB_ABOVE = 0x1 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_2DB_ABOVE = 0x2 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_3DB_ABOVE = 0x3 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_4DB_ABOVE = 0x4 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_5DB_ABOVE = 0x5 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_6DB_ABOVE = 0x6 << 0
	AGCCTRL1_CARRIER_SENSE_ABS_THR_7DB_ABOVE = 0x7 << 0

	AGCCTRL0_HYST_LEVEL_NONE          = 0 << 6
	AGCCTRL0_HYST_LEVEL_LOW           = 1 << 6
	AGCCTRL0_HYST_LEVEL_MEDIUM        = 2 << 6
	AGCCTRL0_HYST_LEVEL_HIGH          = 3 << 6
	AGCCTRL0_WAIT_TIME_8              = 0 << 4
	AGCCTRL0_WAIT_TIME_16             = 1 << 4
	AGCCTRL0_WAIT_TIME_24             = 2 << 4
	AGCCTRL0_WAIT_TIME_32             = 3 << 4
	AGCCTRL0_AGC_FREEZE_NORMAL        = 0 << 2
	AGCCTRL0_AGC_FREEZE_SYNC          = 1 << 2
	AGCCTRL0_AGC_FREEZE_MANUAL_ANALOG = 2 << 2
	AGCCTRL0_AGC_FREEZE_MANUAL_BOTH   = 3 << 2
	AGCCTRL0_FILTER_LENGTH_8          = 0 << 0
	AGCCTRL0_FILTER_LENGTH_16         = 1 << 0
	AGCCTRL0_FILTER_LENGTH_32         = 2 << 0
	AGCCTRL0_FILTER_LENGTH_64         = 3 << 0

	FREND1_LNA_CURRENT_SHIFT          = 6
	FREND1_LNA2MIX_CURRENT_SHIFT      = 4
	FREND1_LODIV_BUF_CURRENT_RX_SHIFT = 2
	FREND1_MIX_CURRENT_SHIFT          = 0

	FREND0_LODIV_BUF_CURRENT_TX_MASK  = 0x3 << 4
	FREND0_LODIV_BUF_CURRENT_TX_SHIFT = 4
	FREND0_PA_POWER_MASK              = 0x7
	FREND0_PA_POWER_SHIFT             = 0

	TEST2_NORMAL_MAGIC           = 0x88
	TEST2_RX_LOW_DATA_RATE_MAGIC = 0x81

	TEST1_TX_MAGIC               = 0x31
	TEST1_RX_LOW_DATA_RATE_MAGIC = 0x35

	TEST0_7_2_MASK       = (0xfc)
	TEST0_VCO_SEL_CAL_EN = 1 << 1
	TEST0_0_MASK         = 1

	LQI_CRC_OK       = 1 << 7
	LQI_LQI_EST_MASK = 0x7f

	MARCSTATE_MASK         = 0x1f
	MARCSTATE_SLEEP        = 0x00
	MARCSTATE_IDLE         = 0x01
	MARCSTATE_VCOON_MC     = 0x03
	MARCSTATE_REGON_MC     = 0x04
	MARCSTATE_MANCAL       = 0x05
	MARCSTATE_VCOON        = 0x06
	MARCSTATE_REGON        = 0x07
	MARCSTATE_STARTCAL     = 0x08
	MARCSTATE_BWBOOST      = 0x09
	MARCSTATE_FS_LOCK      = 0x0a
	MARCSTATE_IFADCON      = 0x0b
	MARCSTATE_ENDCAL       = 0x0c
	MARCSTATE_RX           = 0x0d
	MARCSTATE_RX_END       = 0x0e
	MARCSTATE_RX_RST       = 0x0f
	MARCSTATE_TXRX_SWITCH  = 0x10
	MARCSTATE_RX_OVERFLOW  = 0x11
	MARCSTATE_FSTXON       = 0x12
	MARCSTATE_TX           = 0x13
	MARCSTATE_TX_END       = 0x14
	MARCSTATE_RXTX_SWITCH  = 0x15
	MARCSTATE_TX_UNDERFLOW = 0x16

	PKTSTATUS_CRC_OK      = 1 << 7
	PKTSTATUS_CS          = 1 << 6
	PKTSTATUS_PQT_REACHED = 1 << 5
	PKTSTATUS_CCA         = 1 << 4
	PKTSTATUS_SFD         = 1 << 3
	PKTSTATUS_GDO2        = 1 << 2
	PKTSTATUS_GDO0        = 1 << 0

	ENCCCS_MODE_CBC     = 0 << 4
	ENCCCS_MODE_CFB     = 1 << 4
	ENCCCS_MODE_OFB     = 2 << 4
	ENCCCS_MODE_CTR     = 3 << 4
	ENCCCS_MODE_ECB     = 4 << 4
	ENCCCS_MODE_CBC_MAC = 5 << 4
	ENCCCS_RDY          = 1 << 3
	ENCCCS_CMD_ENCRYPT  = 0 << 1
	ENCCCS_CMD_DECRYPT  = 1 << 1
	ENCCCS_CMD_LOAD_KEY = 2 << 1
	ENCCCS_CMD_LOAD_IV  = 3 << 1
	ENCCCS_START        = 1 << 0

	NUM_TXBYTES_MASK = 0x7f
	TXFIFO_UNDERFLOW = 1 << 7

	NUM_RXBYTES_MASK = 0x7f
	RXFIFO_OVERFLOW  = 1 << 7
)
