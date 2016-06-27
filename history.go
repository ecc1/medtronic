package medtronic

const (
	CurrentPage CommandCode = 0x9D
	History     CommandCode = 0x80
)

func (pump *Pump) CurrentPage() int {
	data := pump.Execute(CurrentPage)
	if pump.Error() != nil {
		return 0
	}
	if len(data) < 5 || data[0] != 4 {
		pump.BadResponse(CurrentPage, data)
		return 0
	}
	page := fourByteInt(data[1:5])
	if page < 0 || page > 36 {
		page = 36
	}
	return page
}

func (pump *Pump) History(page int) []byte {
	return pump.Download(History, page)
}
