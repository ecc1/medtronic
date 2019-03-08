#!/bin/sh -e

programs="cgmhistory cgmpage cgmupdate fakemeter historypage listen mdt mmtune pumphistory setbasals sniff"
others="cmd/pumphistory/openaps.jq"
go_ldflags="-s -w"
target_archs="arm 386"
radios="rfm69 spi uart"
dest=binaries
package=$(basename $(pwd))

for arch in $target_archs; do
    echo "Building binaries for $arch architecture"
    for radio in $radios; do
	dir=$dest/$arch/$radio
	mkdir -pv $dir
	gocmd="GOARCH=$arch go build -tags \"$radio\" -ldflags \"$go_ldflags\""
	echo "Build command: $gocmd"
	for prog in $programs; do
	    (cd cmd/$prog && \
	     eval $gocmd && \
	     mv -v $prog ../../$dir/)
	done
	for other in $others; do
	    cp -v $other $dir
	done
	tarball=${package}-${arch}-${radio}.tar.xz
	echo Building $tarball
	tar --create --file $dest/$tarball --xz --verbose --directory $dir .
    done
done

ls -l $dest/*.xz
