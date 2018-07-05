package medtronic

const (
	// MaxGlucosePages is the maximum number of glucose history pages.
	MaxGlucosePages = 32
)

// GlucosePage downloads the given glucose page.
func (pump *Pump) GlucosePage(page int) []byte {
	return pump.Download(glucosePage, page)
}

// ISIGPage downloads the given ISIG page.
func (pump *Pump) ISIGPage(page int) []byte {
	return pump.Download(isigPage, page)
}

// VcntrPage downloads the given vcntr page.
func (pump *Pump) VcntrPage(page int) []byte {
	return pump.Download(vcntrPage, page)
}

// CalibrationFactor returns the CGM calibration factor.
func (pump *Pump) CalibrationFactor() int {
	data := pump.Execute(calibrationFactor)
	if pump.Error() != nil {
		return 0
	}
	if len(data) < 3 || data[0] != 2 {
		pump.BadResponse(calibrationFactor, data)
		return 0
	}
	return int(twoByteUint(data[1:3]))

}

// CGMCurrentGlucosePage returns the current CGM glucose page number.
func (pump *Pump) CGMCurrentGlucosePage() int {
	data := pump.Execute(cgmPageCount)
	if pump.Error() != nil {
		return 0
	}
	if len(data) < 13 || data[0] != 12 {
		pump.BadResponse(cgmPageCount, data)
		return 0
	}
	return int(fourByteUint(data[1:5]))
}

// CGMWriteTimestamp writes a new sensor timestamp to the CGM history.
func (pump *Pump) CGMWriteTimestamp() {
	pump.Execute(cgmWriteTimestamp)
}
