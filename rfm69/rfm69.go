package rfm69

// http://www.hoperf.com/upload/rf/RFM69HCW-V1.1.pdf

const (
	FXOSC        = 32000000
	SpiWriteMode = 1 << 7
)

// Common Configuration Registers
const (
	RegFifo       = 0x00 // FIFO read/write access
	RegOpMode     = 0x01 // Operating modes of the transceiver
	RegDataModul  = 0x02 // Data operation mode and Modulation settings
	RegBitrateMsb = 0x03 // Bit Rate setting, Most Significant Bits
	RegBitrateLsb = 0x04 // Bit Rate setting, Least Significant Bits
	RegFdevMsb    = 0x05 // Frequency Deviation setting, Most Significant Bits
	RegFdevLsb    = 0x06 // Frequency Deviation setting, Least Significant Bits
	RegFrfMsb     = 0x07 // RF Carrier Frequency, Most Significant Bits
	RegFrfMid     = 0x08 // RF Carrier Frequency, Intermediate Bits
	RegFrfLsb     = 0x09 // RF Carrier Frequency, Least Significant Bits
	RegOsc1       = 0x0A // RF Oscillators Settings
	RegAfcCtrl    = 0x0B // AFC control in low modulation index situations
	RegListen1    = 0x0D // Listen Mode settings
	RegListen2    = 0x0E // Listen Mode Idle duration
	RegListen3    = 0x0F // Listen Mode Rx duration
	RegVersion    = 0x10
)

// Transmitter Registers
const (
	RegPaLevel = 0x11 // PA selection and Output Power control
	RegPaRamp  = 0x12 // Control of the PA ramp time in FSK mode
	RegOcp     = 0x13 // Over Current Protection control
)

// Receiver Registers
const (
	RegLna        = 0x18 // LNA settings
	RegRxBw       = 0x19 // Channel Filter BW Control
	RegAfcBw      = 0x1A // Channel Filter BW control during the AFC routine
	RegOokPeak    = 0x1B // OOK demodulator selection and control in peak mode
	RegOokAvg     = 0x1C // Average threshold control of the OOK demodulator
	RegOokFix     = 0x1D // Fixed threshold control of the OOK demodulator
	RegAfcFei     = 0x1E // AFC and FEI control and status
	RegAfcMsb     = 0x1F // MSB of the frequency correction of the AFC
	RegAfcLsb     = 0x20 // LSB of the frequency correction of the AFC
	RegFeiMsb     = 0x21 // MSB of the calculated frequency error
	RegFeiLsb     = 0x22 // LSB of the calculated frequency error
	RegRssiConfig = 0x23 // RSSI-related settings
	RegRssiValue  = 0x24 // RSSI value in dBm
)

// IRQ and Pin Mapping Registers
const (
	RegDioMapping1 = 0x25 // Mapping of pins DIO0 to DIO3
	RegDioMapping2 = 0x26 // Mapping of pins DIO4 and DIO5, ClkOut frequency
	RegIrqFlags1   = 0x27 // Status register: PLL Lock state, Timeout, RSSI > Threshold...
	RegIrqFlags2   = 0x28 // Status register: FIFO handling flags...
	RegRssiThresh  = 0x29 // RSSI Threshold control
	RegRxTimeout1  = 0x2A // Timeout duration between Rx request and RSSI detection
	RegRxTimeout2  = 0x2B // Timeout duration between RSSI detection and PayloadReady
)

// Packet Engine Registers
const (
	RegPreambleMsb   = 0x2C // Preamble length, MSB
	RegPreambleLsb   = 0x2D // Preamble length, LSB
	RegSyncConfig    = 0x2E // Sync Word Recognition control
	RegSyncValue1    = 0x2F // Sync Word bytes, 1 through 8
	RegSyncValue2    = 0x30
	RegSyncValue3    = 0x31
	RegSyncValue4    = 0x32
	RegSyncValue5    = 0x33
	RegSyncValue6    = 0x34
	RegSyncValue7    = 0x35
	RegSyncValue8    = 0x36
	RegPacketConfig1 = 0x37 // Packet mode settings
	RegPayloadLength = 0x38 // Payload length setting
	RegNodeAdrs      = 0x39 // Node address
	RegBroadcastAdrs = 0x3A // Broadcast address
	RegAutoModes     = 0x3B // Auto modes settings
	RegFifoThresh    = 0x3C // Fifo threshold, Tx start condition
	RegPacketConfig2 = 0x3D // Packet mode settings
	RegAesKey1       = 0x3E // 16 bytes of the cypher key
	RegAesKey2       = 0x3F
	RegAesKey3       = 0x40
	RegAesKey4       = 0x41
	RegAesKey5       = 0x42
	RegAesKey6       = 0x43
	RegAesKey7       = 0x44
	RegAesKey8       = 0x45
	RegAesKey9       = 0x46
	RegAesKey10      = 0x47
	RegAesKey11      = 0x48
	RegAesKey12      = 0x49
	RegAesKey13      = 0x4A
	RegAesKey14      = 0x4B
	RegAesKey15      = 0x4C
	RegAesKey16      = 0x4D
)

// Temperature Sensor Registers
const (
	RegTemp1 = 0x4E // Temperature Sensor control
	RegTemp2 = 0x4F // Temperature readout
)

// Test Registers
const (
	RegTest     = 0x50 // Internal test registers
	RegTestLna  = 0x58 // Sensitivity boost
	RegTestPa1  = 0x5A // High Power PA settings
	RegTestPa2  = 0x5C // High Power PA settings
	RegTestDagc = 0x6F // Fading Margin Improvement
	RegTestAfc  = 0x71 // AFC offset for low modulation index AFC
)

type RfConfiguration struct {
	// Omit RegFifo to avoid reading or writing it with this struct.
	RegOpMode        byte // Operating modes of the transceiver
	RegDataModul     byte // Data operation mode and Modulation settings
	RegBitrateMsb    byte // Bit Rate setting, Most Significant Bits
	RegBitrateLsb    byte // Bit Rate setting, Least Significant Bits
	RegFdevMsb       byte // Frequency Deviation setting, Most Significant Bits
	RegFdevLsb       byte // Frequency Deviation setting, Least Significant Bits
	RegFrfMsb        byte // RF Carrier Frequency, Most Significant Bits
	RegFrfMid        byte // RF Carrier Frequency, Intermediate Bits
	RegFrfLsb        byte // RF Carrier Frequency, Least Significant Bits
	RegOsc1          byte // RF Oscillators Settings
	RegAfcCtrl       byte // AFC control in low modulation index situations
	reserved0C       byte
	RegListen1       byte // Listen Mode settings
	RegListen2       byte // Listen Mode Idle duration
	RegListen3       byte // Listen Mode Rx duration
	RegVersion       byte // PA selection and Output Power control
	RegPaLevel       byte // PA selection and Output Power control
	RegPaRamp        byte // Control of the PA ramp time in FSK mode
	RegOcp           byte // Over Current Protection control
	reserved14       byte
	reserved15       byte
	reserved16       byte
	reserved17       byte
	RegLna           byte // LNA settings
	RegRxBw          byte // Channel Filter BW Control
	RegAfcBw         byte // Channel Filter BW control during the AFC routine
	RegOokPeak       byte // OOK demodulator selection and control in peak mode
	RegOokAvg        byte // Average threshold control of the OOK demodulator
	RegOokFix        byte // Fixed threshold control of the OOK demodulator
	RegAfcFei        byte // AFC and FEI control and status
	RegAfcMsb        byte // MSB of the frequency correction of the AFC
	RegAfcLsb        byte // LSB of the frequency correction of the AFC
	RegFeiMsb        byte // MSB of the calculated frequency error
	RegFeiLsb        byte // LSB of the calculated frequency error
	RegRssiConfig    byte // RSSI-related settings
	RegRssiValue     byte // RSSI value in dBm
	RegDioMapping1   byte // Mapping of pins DIO0 to DIO3
	RegDioMapping2   byte // Mapping of pins DIO4 and DIO5, ClkOut frequency
	RegIrqFlags1     byte // Status register: PLL Lock state, Timeout, RSSI > Threshold...
	RegIrqFlags2     byte // Status register: FIFO handling flags...
	RegRssiThresh    byte // RSSI Threshold control
	RegRxTimeout1    byte // Timeout duration between Rx request and RSSI detection
	RegRxTimeout2    byte // Timeout duration between RSSI detection and PayloadReady
	RegPreambleMsb   byte // Preamble length, MSB
	RegPreambleLsb   byte // Preamble length, LSB
	RegSyncConfig    byte // Sync Word Recognition control
	RegSyncValue1    byte // Sync Word bytes, 1 through 8
	RegSyncValue2    byte
	RegSyncValue3    byte
	RegSyncValue4    byte
	RegSyncValue5    byte
	RegSyncValue6    byte
	RegSyncValue7    byte
	RegSyncValue8    byte // Packet mode settings
	RegPacketConfig1 byte // Packet mode settings
	RegPayloadLength byte // Payload length setting
	RegNodeAdrs      byte // Node address
	RegBroadcastAdrs byte // Broadcast address
	RegAutoModes     byte // Auto modes settings
	RegFifoThresh    byte // Fifo threshold, Tx start condition
	RegPacketConfig2 byte // Packet mode settings
	RegAesKey1       byte // 16 bytes of the cypher key
	RegAesKey2       byte
	RegAesKey3       byte
	RegAesKey4       byte
	RegAesKey5       byte
	RegAesKey6       byte
	RegAesKey7       byte
	RegAesKey8       byte
	RegAesKey9       byte
	RegAesKey10      byte
	RegAesKey11      byte
	RegAesKey12      byte
	RegAesKey13      byte
	RegAesKey14      byte
	RegAesKey15      byte
	RegAesKey16      byte // Temperature Sensor control
	RegTemp1         byte // Temperature Sensor control
	RegTemp2         byte // Temperature readout Omit test
	// Omit test registers to avoid undefined behavior
	// when reading or writing this struct.
}

// Configuration register values after reset,
// according to data sheet section 6.
var ResetRfConfiguration = RfConfiguration{
	RegOpMode:        0x04,
	RegDataModul:     0x00,
	RegBitrateMsb:    0x1A,
	RegBitrateLsb:    0x0B,
	RegFdevMsb:       0x00,
	RegFdevLsb:       0x52,
	RegFrfMsb:        0xE4,
	RegFrfMid:        0xC0,
	RegFrfLsb:        0x00,
	RegOsc1:          0x41,
	RegAfcCtrl:       0x40, // unused bits; disagrees with data sheet
	reserved0C:       0x02,
	RegListen1:       0x92,
	RegListen2:       0xF5,
	RegListen3:       0x20,
	RegVersion:       0x24,
	RegPaLevel:       0x9F,
	RegPaRamp:        0x09,
	RegOcp:           0x1A,
	reserved14:       0x40,
	reserved15:       0xB0,
	reserved16:       0x7B,
	reserved17:       0x9B,
	RegLna:           0x08,
	RegRxBw:          0x86,
	RegAfcBw:         0x8A,
	RegOokPeak:       0x40,
	RegOokAvg:        0x80,
	RegOokFix:        0x06,
	RegAfcFei:        0x10,
	RegAfcMsb:        0x00,
	RegAfcLsb:        0x00,
	RegFeiMsb:        0x00,
	RegFeiLsb:        0x00,
	RegRssiConfig:    0x02,
	RegRssiValue:     0xFF,
	RegDioMapping1:   0x00,
	RegDioMapping2:   0x05,
	RegIrqFlags1:     0x80,
	RegIrqFlags2:     0x00,
	RegRssiThresh:    0xFF,
	RegRxTimeout1:    0x00,
	RegRxTimeout2:    0x00,
	RegPreambleMsb:   0x00,
	RegPreambleLsb:   0x03,
	RegSyncConfig:    0x98,
	RegSyncValue1:    0x00,
	RegSyncValue2:    0x00,
	RegSyncValue3:    0x00,
	RegSyncValue4:    0x00,
	RegSyncValue5:    0x00,
	RegSyncValue6:    0x00,
	RegSyncValue7:    0x00,
	RegSyncValue8:    0x00,
	RegPacketConfig1: 0x10,
	RegPayloadLength: 0x40,
	RegNodeAdrs:      0x00,
	RegBroadcastAdrs: 0x00,
	RegAutoModes:     0x00,
	RegFifoThresh:    0x0F,
	RegPacketConfig2: 0x02,
	RegAesKey1:       0x00,
	RegAesKey2:       0x00,
	RegAesKey3:       0x00,
	RegAesKey4:       0x00,
	RegAesKey5:       0x00,
	RegAesKey6:       0x00,
	RegAesKey7:       0x00,
	RegAesKey8:       0x00,
	RegAesKey9:       0x00,
	RegAesKey10:      0x00,
	RegAesKey11:      0x00,
	RegAesKey12:      0x00,
	RegAesKey13:      0x00,
	RegAesKey14:      0x00,
	RegAesKey15:      0x00,
	RegAesKey16:      0x00,
	RegTemp1:         0x01,
	RegTemp2:         0x00,
}

// Default (recommended) values, according to data sheet section 6.
var DefaultRfConfiguration = RfConfiguration{
	RegOpMode:        0x04,
	RegDataModul:     0x00,
	RegBitrateMsb:    0x1A,
	RegBitrateLsb:    0x0B,
	RegFdevMsb:       0x00,
	RegFdevLsb:       0x52,
	RegFrfMsb:        0xE4,
	RegFrfMid:        0xC0,
	RegFrfLsb:        0x00,
	RegOsc1:          0x41,
	RegAfcCtrl:       0x00,
	reserved0C:       0x02,
	RegListen1:       0x92,
	RegListen2:       0xF5,
	RegListen3:       0x20,
	RegVersion:       0x24,
	RegPaLevel:       0x9F,
	RegPaRamp:        0x09,
	RegOcp:           0x1A,
	reserved14:       0x40,
	reserved15:       0xB0,
	reserved16:       0x7B,
	reserved17:       0x9B,
	RegLna:           0x88,
	RegRxBw:          0x55,
	RegAfcBw:         0x8B,
	RegOokPeak:       0x40,
	RegOokAvg:        0x80,
	RegOokFix:        0x06,
	RegAfcFei:        0x10,
	RegAfcMsb:        0x00,
	RegAfcLsb:        0x00,
	RegFeiMsb:        0x00,
	RegFeiLsb:        0x00,
	RegRssiConfig:    0x02,
	RegRssiValue:     0xFF,
	RegDioMapping1:   0x00,
	RegDioMapping2:   0x07,
	RegIrqFlags1:     0x80,
	RegIrqFlags2:     0x00,
	RegRssiThresh:    0xE4,
	RegRxTimeout1:    0x00,
	RegRxTimeout2:    0x00,
	RegPreambleMsb:   0x00,
	RegPreambleLsb:   0x03,
	RegSyncConfig:    0x98,
	RegSyncValue1:    0x01,
	RegSyncValue2:    0x01,
	RegSyncValue3:    0x01,
	RegSyncValue4:    0x01,
	RegSyncValue5:    0x01,
	RegSyncValue6:    0x01,
	RegSyncValue7:    0x01,
	RegSyncValue8:    0x01,
	RegPacketConfig1: 0x10,
	RegPayloadLength: 0x40,
	RegNodeAdrs:      0x00,
	RegBroadcastAdrs: 0x00,
	RegAutoModes:     0x00,
	RegFifoThresh:    0x8F,
	RegPacketConfig2: 0x02,
	RegAesKey1:       0x00,
	RegAesKey2:       0x00,
	RegAesKey3:       0x00,
	RegAesKey4:       0x00,
	RegAesKey5:       0x00,
	RegAesKey6:       0x00,
	RegAesKey7:       0x00,
	RegAesKey8:       0x00,
	RegAesKey9:       0x00,
	RegAesKey10:      0x00,
	RegAesKey11:      0x00,
	RegAesKey12:      0x00,
	RegAesKey13:      0x00,
	RegAesKey14:      0x00,
	RegAesKey15:      0x00,
	RegAesKey16:      0x00,
	RegTemp1:         0x01,
	RegTemp2:         0x00,
}

// RegOpMode
const (
	SequencerOff    = 1 << 7
	ListenOn        = 1 << 6
	ListenAbort     = 1 << 5
	ModeShift       = 2
	ModeMask        = 7 << 2
	SleepMode       = 0 << 2
	StandbyMode     = 1 << 2
	FreqSynthMode   = 2 << 2
	TransmitterMode = 3 << 2
	ReceiverMode    = 4 << 2
)

// RegDataModul
const (
	PacketMode                   = 0 << 5
	ContinuousModeWithBitSync    = 2 << 5
	ContinuousModeWithoutBitSync = 3 << 5
	ModulationTypeMask           = 3 << 3
	ModulationTypeFSK            = 0 << 3
	ModulationTypeOOK            = 1 << 3
	ModulationShapingShift       = 0
)

// RegOsc1
const (
	RcCalStart = 1 << 7
	RcCalDone  = 1 << 6
)

// RegAfcCtrl
const (
	AfcLowBetaOn = 1 << 5
)

// RegListen1
const (
	ListenResolIdleShift = 6
	ListenResolRxShift   = 4
	ListenCriteria       = 1 << 3
	ListenEndShift       = 1
)

// RegPaLevel
// See http://blog.andrehessling.de/2015/02/07/figuring-out-the-power-level-settings-of-hoperfs-rfm69-hwhcw-modules/
const (
	Pa0On            = 1 << 7
	Pa1On            = 1 << 6
	Pa2On            = 1 << 5
	OutputPowerShift = 0
)

// RegOcp
const (
	OcpOn        = 1 << 4
	OcpTrimShift = 0
)

// RegLna
const (
	LnaZin              = 1 << 7
	LnaCurrentGainShift = 3
	LnaGainSelectShift  = 0
)

// RegRxBw
const (
	DccFreqShift  = 5
	DccFreqMask   = 7 << 5
	RxBwMantShift = 3
	RxBwMantMask  = 3 << 3
	RxBwMant16    = 0 << 3
	RxBwMant20    = 1 << 3
	RxBwMant24    = 2 << 3
	RxBwExpShift  = 0
	RxBwExpMask   = 7 << 0
)

// RegOokPeak
const (
	OokThreshTypeShift     = 6
	OokPeakThreshStepShift = 3
	OokPeakThreshDecShift  = 0
)

// RegOokAvg
const (
	OokAverageThreshFiltShift = 6
)

// RegAfcFei
const (
	FeiDone        = 1 << 6
	FeiStart       = 1 << 5
	AfcDone        = 1 << 4
	AfcAutoclearOn = 1 << 3
	AfcAutoOn      = 1 << 2
	AfcClear       = 1 << 1
	AfcStart       = 1 << 0
)

// RegRssiConfig
const (
	RssiDone  = 1 << 1
	RssiStart = 1 << 0
)

// RegDioMapping1
const (
	Dio0MappingShift = 6
	Dio1MappingShift = 4
	Dio2MappingShift = 2
	Dio3MappingShift = 0
)

// RegDioMapping2
const (
	Dio4MappingShift = 6
	Dio5MappingShift = 4
	ClkOutShift      = 0
)

// RegIrqFlags1
const (
	ModeReady        = 1 << 7
	RxReady          = 1 << 6
	TxReady          = 1 << 5
	PllLock          = 1 << 4
	Rssi             = 1 << 3
	Timeout          = 1 << 2
	AutoMode         = 1 << 1
	SyncAddressMatch = 1 << 0
)

// RegIrqFlags2
const (
	FifoFull     = 1 << 7
	FifoNotEmpty = 1 << 6
	FifoLevel    = 1 << 5
	FifoOverrun  = 1 << 4
	PacketSent   = 1 << 3
	PayloadReady = 1 << 2
	CrcOk        = 1 << 1
)

// RegSyncConfig
const (
	SyncOn            = 1 << 7
	FifoFillCondition = 1 << 6
	SyncSizeShift     = 3
	SyncTolShift      = 0
)

// RegPacketConfig1
const (
	FixedLength           = 0 << 7
	VariableLength        = 1 << 7
	DcFreeShift           = 5
	CrcOn                 = 1 << 4
	CrcOff                = 0 << 4
	CrcAutoClearOff       = 1 << 3
	AddressFilteringShift = 1
)

// RegAutoModes
const (
	EnterConditionShift   = 5
	ExitConditionShift    = 2
	IntermediateModeShift = 0
)

// RegFifoThresh
const (
	TxStartFifoNotEmpty = 1 << 7
	TxStartFifoLevel    = 0 << 7
	FifoThresholdShift  = 0
)

// RegPacketConfig2
const (
	InterPacketRxDelayShift = 4
	RestartRx               = 1 << 2
	AutoRxRestartOn         = 1 << 1
	AutoRxRestartOff        = 0 << 1
	AesOn                   = 1 << 0
)

// RegTemp1
const (
	TempMeasStart   = 1 << 3
	TempMeasRunning = 1 << 2
)
