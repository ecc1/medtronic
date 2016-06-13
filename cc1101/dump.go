package cc1101

import (
	"log"
)

func (r *Radio) DumpRF() {
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	log.Printf("State: %s", r.State())
	log.Printf("Frequency: %d", r.Frequency())
	log.Printf("Channel: %d", r.ReadRegister(CHANNR))
	r.showFreqSynthControl()
	r.showModemConfig()
	pa := r.ReadPaTable()
	n := r.ReadRegister(FREND0) & FREND0_PA_POWER_MASK
	log.Printf("PATABLE: % X using 0..%d", pa, n)
}

func (r *Radio) showFreqSynthControl() {
	log.Printf("Intermediate frequency: %d Hz", r.ReadIF())
	log.Printf("Frequency offset: %d Hz", r.ReadRegister(FSCTRL0))
}

func (r *Radio) showModemConfig() {
	chanbw, drate := r.ReadChannelParams()
	log.Printf("Channel bandwidth: %d Hz", chanbw)
	log.Printf("Data rate: %d Baud", drate)

	m2 := r.ReadRegister(MDMCFG2)
	showBoolCondition("DC blocking filter", m2&MDMCFG2_DEM_DCFILT_OFF == 0)
	showBoolCondition("Manchester encoding", m2&(1<<3) != 0)
	log.Printf("Modulation format: %s", modFormat[(m2&MDMCFG2_MOD_FORMAT_MASK)>>4])
	log.Printf("Sync mode: %s", syncMode[m2&MDMCFG2_SYNC_MODE_MASK])

	fec, minPreamble, chanspc := r.ReadModemConfig()
	showBoolCondition("Forward Error Correction", fec)
	log.Printf("Min preamble bytes: %d", minPreamble)
	log.Printf("Channel spacing: %d Hz", chanspc)
}

func showBoolCondition(name string, cond bool) {
	if cond {
		log.Printf("%s: enabled", name)
	} else {
		log.Printf("%s: disabled", name)
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
)
