#!/bin/sh -e

sleep_seconds=270

if [ ! "$NIGHTSCOUT_SITE" ]; then
    echo NIGHTSCOUT_SITE environement variable is not set
    exit 1
fi

glucose_file=$(tempfile)
trap "rm -f -- '$glucose_file'" EXIT

download_glucose() {
    curl -s "$NIGHTSCOUT_SITE/api/v1/entries.json?count=15" | jq . > $glucose_file
}

while true; do
    download_glucose && bgproxy -f $glucose_file || true
    sleep $sleep_seconds
done
