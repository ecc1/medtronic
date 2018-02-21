package packet

// Medtronic packet types.
const (
	MySentry = 0xA2
	Meter    = 0xA5 // Bayer Contour glucometer
	RFRemote = 0xA6 // MMT-503 remote
	Pump     = 0xA7
	Sensor   = 0xA8
)
