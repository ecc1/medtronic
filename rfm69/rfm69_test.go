package rfm69

import (
	"testing"
	"unsafe"
)

func TestRfConfiguration(t *testing.T) {
	have := int(unsafe.Sizeof(RfConfiguration{}))
	want := RegTemp2 - RegOpMode + 1
	if have != want {
		t.Errorf("Sizeof(RfConfiguration) == %d, want %d", have, want)
	}
}
