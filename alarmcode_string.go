// Code generated by "stringer -type AlarmCode"; DO NOT EDIT.

package medtronic

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[BatteryOutLimitExceeded-3]
	_ = x[NoDelivery-4]
	_ = x[BatteryDepleted-5]
	_ = x[AutoOff-6]
	_ = x[DeviceReset-16]
	_ = x[ReprogramError-61]
	_ = x[EmptyReservoir-62]
}

const (
	_AlarmCode_name_0 = "BatteryOutLimitExceededNoDeliveryBatteryDepletedAutoOff"
	_AlarmCode_name_1 = "DeviceReset"
	_AlarmCode_name_2 = "ReprogramErrorEmptyReservoir"
)

var (
	_AlarmCode_index_0 = [...]uint8{0, 23, 33, 48, 55}
	_AlarmCode_index_2 = [...]uint8{0, 14, 28}
)

func (i AlarmCode) String() string {
	switch {
	case 3 <= i && i <= 6:
		i -= 3
		return _AlarmCode_name_0[_AlarmCode_index_0[i]:_AlarmCode_index_0[i+1]]
	case i == 16:
		return _AlarmCode_name_1
	case 61 <= i && i <= 62:
		i -= 61
		return _AlarmCode_name_2[_AlarmCode_index_2[i]:_AlarmCode_index_2[i+1]]
	default:
		return "AlarmCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
