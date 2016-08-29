The `medtronic` package provides a unified API for
communicating with Medtronic insulin pumps using SPI-connected radio
modules.

Documentation: <https://godoc.org/github.com/ecc1/medtronic>

Decoding of messages to and from the pump is derived almost entirely from
[Ben West's pioneering "Decoding Carelink" work.](https://github.com/bewest/decoding-carelink)

The `apps` directory contains a number of command-line applications,
including a "Swiss army knife" application `mdt`
(analogous to the the `openaps use pump ...` commands).

The `spilink` directory contains a backend server that can be used with
[a proof-of-concept mmeowlink driver](https://github.com/ecc1/mmeowlink/tree/spilink)
for [openaps.](https://openaps.org)
