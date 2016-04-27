package cc1100

import (
	"fmt"
	"log"
)

func (dev *Device) DumpRF() {
	freq, err := dev.ReadFrequency()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Frequency: %d\n", freq)
	fmt.Printf("Channel: %d\n", dev.readReg(CHANNR))
	dev.showFreqSynthControl()
	dev.showModemConfig()
}

func (dev *Device) readReg(addr byte) byte {
	v, err := dev.ReadRegister(addr)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func (dev *Device) showFreqSynthControl() {
	f, err := dev.ReadIF()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Intermediate frequency: %d Hz\n", f)
	fmt.Printf("Frequency offset: %d Hz\n", dev.readReg(FSCTRL0))
}

func (dev *Device) showModemConfig() {
	chanbw, drate, err := dev.ReadChannelParams()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Channel bandwidth: %d Hz\n", chanbw)
	fmt.Printf("Data rate: %d Baud\n", drate)

	m2 := dev.readReg(MDMCFG2)
	showBoolCondition("DC blocking filter", m2&MDMCFG2_DEM_DCFILT_OFF == 0)
	showBoolCondition("Manchester encoding", m2&(1<<3) != 0)
	fmt.Printf("Modulation format: %s\n", modFormat[(m2&MDMCFG2_MOD_FORMAT_MASK)>>4])
	fmt.Printf("Sync mode: %s\n", syncMode[m2&MDMCFG2_SYNC_MODE_MASK])

	fec, minPreamble, chanspc, err := dev.ReadModemConfig()
	if err != nil {
		log.Fatal(err)
	}
	showBoolCondition("Forward Error Correction", fec)
	fmt.Printf("Min preamble bytes: %d\n", minPreamble)
	fmt.Printf("Channel spacing: %d Hz\n", chanspc)
}

func showBoolCondition(name string, cond bool) {
	if cond {
		fmt.Printf("%s: enabled\n", name)
	} else {
		fmt.Printf("%s: disabled\n", name)
	}
}

func strobeName(strobe byte) string {
	return strobeString[strobe-SRES]
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

	modFormat = []string{
		"2-FSK",
		"GFSK",
		"-",
		"OOK",
		"-",
		"-",
		"-",
		"MSK",
	}

	syncMode = []string{
		"No preamble/sync",
		"15/16 sync word bits detected",
		"16/16 sync word bits detected",
		"30/32 sync word bits detected",
		"No preamble/sync, carrier-sense above threshold",
		"15/16 + carrier-sense above threshold",
		"16/16 + carrier-sense above threshold",
		"30/32 + carrier-sense above threshold",
	}

	numPreamble = []uint8{2, 3, 4, 6, 8, 12, 16, 24}

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
)
