package medtronic

const (
	lastHistoryPage Command = 0x9D
	historyPage     Command = 0x80
)

// LastHistoryPage returns the pump's last history page number.
func (pump *Pump) LastHistoryPage() int {
	data := pump.Execute(lastHistoryPage)
	if pump.Error() != nil {
		return 0
	}
	if len(data) < 5 || data[0] != 4 {
		pump.BadResponse(lastHistoryPage, data)
		return 0
	}
	page := fourByteUint(data[1:5])
	if page > 35 {
		page = 35
	}
	return int(page)
}

// HistoryPage downloads the given history page.
func (pump *Pump) HistoryPage(page int) []byte {
	return pump.Download(historyPage, page)
}
