package cc1101

import (
	"testing"
	"unsafe"
)

func TestRfConfiguration(t *testing.T) {
	have := int(unsafe.Sizeof(RfConfiguration{}))
	want := TEST0 - IOCFG2 + 1
	if have != want {
		t.Errorf("Sizeof(RfConfiguration) == %d, want %d", have, want)
	}
}
