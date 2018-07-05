package medtronic

const (
	// MaxHistoryPages is the maximum number of pump history pages.
	MaxHistoryPages = 36

	// Max512HistoryPages is the maximum number of pump history pages for model x12 pumps.
	Max512HistoryPages = 32
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
			return Max512HistoryPages
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
	if page > MaxHistoryPages {
		page = MaxHistoryPages
	}
	return int(page)
}
