package medtronic

const (
	HistoryPageCount CommandCode = 0x9D
	HistoryPage      CommandCode = 0x80
)

func (pump *Pump) HistoryPageCount() int {
	data := pump.Execute(HistoryPageCount)
	if pump.Error() != nil {
		return 0
	}
	if len(data) < 5 || data[0] != 4 {
		pump.BadResponse(HistoryPageCount, data)
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
