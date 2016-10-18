package medtronic

const (
	LastHistoryPage Command = 0x9D
	HistoryPage     Command = 0x80
)

func (pump *Pump) LastHistoryPage() int {
	data := pump.Execute(LastHistoryPage)
	if pump.Error() != nil {
		return 0
	}
	if len(data) < 5 || data[0] != 4 {
		pump.BadResponse(LastHistoryPage, data)
		return 0
	}
	page := fourByteInt(data[1:5])
	if page < 0 || page > 36 {
		page = 36
	}
	return page
}

func (pump *Pump) HistoryPage(page int) []byte {
	return pump.Download(HistoryPage, page)
}
