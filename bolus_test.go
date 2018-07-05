package medtronic

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestEncodeBolus(t *testing.T) {
	cases := []struct {
		family Family
		amount Insulin
		actual Insulin
	}{
		{22, 1000, 1000},
		{22, 2550, 2500},
		{23, 575, 575},
		{23, 2575, 2550},
		{23, 11250, 11200},
		{54, 1075, 1050},
		{54, 5075, 5050},
		{54, 10125, 10100},
		{54, 10175, 10100},
	}
	log.SetOutput(ioutil.Discard)
	for _, c := range cases {
		name := fmt.Sprintf("%d_%d", c.family, c.amount)
		t.Run(name, func(t *testing.T) {
			r, err := encodeBolus(c.amount, c.family)
			if err != nil {
				t.Errorf("encodeBolus(%d, %d) raised error (%v)", c.amount, c.family, err)
			}
			a := Insulin(r) * milliUnitsPerStroke(c.family)
			if a != c.actual {
				t.Errorf("encodeBolus(%v, %d) == %d, want %d", c.amount, c.family, a, c.actual)
			}
		})
	}
}
