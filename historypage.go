package medtronic

const (
	historyPage      Command = 0x80
	historyPageCount Command = 0x9D
)

// HistoryPage downloads the given history page.
func (pump *Pump) HistoryPage(page int) []byte {
	return pump.Download(historyPage, page)
}

// HistoryPageCount returns the number of pump history pages.
func (pump *Pump) HistoryPageCount() int {
	data := pump.Execute(historyPageCount)
	if pump.Error() != nil {
		e, ok := pump.Error().(InvalidCommandError)
		if ok && e.PumpError == CommandRefused && pump.Family() == 12 {
			pump.SetError(nil)
			return 32
		}
		return 0
	}
	if len(data) < 5 || data[0] != 4 {
		pump.BadResponse(historyPageCount, data)
		return 0
	}
	page := fourByteUint(data[1:5])
	if page == 0 {
		// Pumps can return 0 when first turned on.
		return 1
	}
	if page > 36 {
		page = 36
	}
	return int(page)
}
