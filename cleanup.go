package medtronic

import (
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func (pump *Pump) closeWhenSignaled() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, unix.SIGTERM)
	<-ch
	pump.Close()
	os.Exit(0)
}
