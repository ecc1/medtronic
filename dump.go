package cc1101

import (
	"log"
)

func (r *Radio) DumpRF() {
	log.Printf("State: %s\n", r.State())

	freq, err := r.Frequency()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Frequency: %d\n", freq)
	log.Printf("Channel: %d\n", r.readReg(CHANNR))
	r.showFreqSynthControl()
	r.showModemConfig()
	pa, err := r.ReadPaTable()
	if err != nil {
		log.Fatal(err)
	}
	n := r.readReg(FREND0) & FREND0_PA_POWER_MASK
	log.Printf("PATABLE: % X using 0..%d\n", pa, n)
}

func (r *Radio) readReg(addr byte) byte {
	v, err := r.ReadRegister(addr)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func (r *Radio) showFreqSynthControl() {
	f, err := r.ReadIF()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Intermediate frequency: %d Hz\n", f)
	log.Printf("Frequency offset: %d Hz\n", r.readReg(FSCTRL0))
}

func (r *Radio) showModemConfig() {
	chanbw, drate, err := r.ReadChannelParams()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Channel bandwidth: %d Hz\n", chanbw)
	log.Printf("Data rate: %d Baud\n", drate)

	m2 := r.readReg(MDMCFG2)
	showBoolCondition("DC blocking filter", m2&MDMCFG2_DEM_DCFILT_OFF == 0)
	showBoolCondition("Manchester encoding", m2&(1<<3) != 0)
	log.Printf("Modulation format: %s\n", modFormat[(m2&MDMCFG2_MOD_FORMAT_MASK)>>4])
	log.Printf("Sync mode: %s\n", syncMode[m2&MDMCFG2_SYNC_MODE_MASK])

	fec, minPreamble, chanspc, err := r.ReadModemConfig()
	if err != nil {
		log.Fatal(err)
	}
	showBoolCondition("Forward Error Correction", fec)
	log.Printf("Min preamble bytes: %d\n", minPreamble)
	log.Printf("Channel spacing: %d Hz\n", chanspc)
}

func showBoolCondition(name string, cond bool) {
	if cond {
		log.Printf("%s: enabled\n", name)
	} else {
		log.Printf("%s: disabled\n", name)
	}
}

func strobeName(strobe byte) string {
	return strobeString[strobe-SRES]
}

var (
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
