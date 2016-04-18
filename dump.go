package cc1100

import (
	"fmt"
	"log"

	"github.com/ecc1/spi"
)

func DumpRF(dev *spi.Device) {
	awaitIdleState(dev)
	freq, err := ReadFrequency(dev)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Frequency: %d\n", freq)
	fmt.Printf("Channel: %d\n", readReg(dev, CHANNR))
	showFreqSynthControl(dev)
	showModemConfig(dev)
}

func readReg(dev *spi.Device, addr byte) byte {
	v, err := ReadRegister(dev, addr)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func awaitIdleState(dev *spi.Device) {
	for {
		m := readReg(dev, MARCSTATE)
		fmt.Printf("MARCSTATE: %s\n", marcState[m&MARCSTATE_MASK])
		switch m & MARCSTATE_MASK {
		case MARCSTATE_SLEEP, MARCSTATE_IDLE:
			return
		}
	}
}

func showFreqSynthControl(dev *spi.Device) {
	f, err := ReadIF(dev)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Intermediate frequency: %d Hz\n", f)
	fmt.Printf("Frequency offset: %d Hz\n", readReg(dev, FSCTRL0))
}

func showModemConfig(dev *spi.Device) {
	chanbw, drate, err := ReadChannelParams(dev)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Channel bandwidth: %d Hz\n", chanbw)
	fmt.Printf("Data rate: %d Baud\n", drate)

	m2 := readReg(dev, MDMCFG2)
	showBoolCondition("DC blocking filter", m2&MDMCFG2_DEM_DCFILT_OFF == 0)
	showBoolCondition("Manchester encoding", m2&(1<<3) != 0)
	fmt.Printf("Modulation format: %s\n", modFormat[(m2&MDMCFG2_MOD_FORMAT_MASK)>>4])
	fmt.Printf("Sync mode: %s\n", syncMode[m2&MDMCFG2_SYNC_MODE_MASK])

	m1 := readReg(dev, MDMCFG1)
	showBoolCondition("Forward Error Correction", m1&MDMCFG1_FEC_EN != 0)
	fmt.Printf("Min preamble bytes: %d\n", numPreamble[(m1&MDMCFG1_NUM_PREAMBLE_MASK)>>4])

	chanspc_E := m1 & MDMCFG1_CHANSPC_E_MASK
	chanspc_M := readReg(dev, MDMCFG0)
	chanspc := (((256+uint64(chanspc_M))<<chanspc_E)*FXOSC) >> 18
	fmt.Printf("Channel spacing: %d Hz\n", chanspc)
}

func showBoolCondition(name string, cond bool) {
	if cond {
		fmt.Printf("%s: enabled\n", name)
	} else {
		fmt.Printf("%s: disabled\n", name)
	}
}

var (
	marcState = []string{
		"SLEEP", "IDLE", "XOFF", "VCOON_MC",
		"REGON_MC", "MANCAL", "VCOON", "REGON",
		"STARTCAL", "BWBOOST", "FS_LOCK", "IFADCON",
		"ENDCAL", "RX", "RX_END", "RX_RST",
		"TXRX_SWITCH", "RX_OVERFLOW", "FSTXON", "TX",
		"TX_END", "RXTX_SWITCH", "TX_UNDERFLOW",
	}

	modFormat = []string{
		"2-FSK", "GFSK", "-", "OOK", "-", "-", "-", "MSK",
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
)
