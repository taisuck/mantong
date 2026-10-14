#!/bin/sh
# Linux용 바이너리를 output/mantong 으로 빌드한다.
set -eu
cd "$(dirname "$0")/.."
mkdir -p output
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o output/mantong ./src
echo "built: output/mantong"
