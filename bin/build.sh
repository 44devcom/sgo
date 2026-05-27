#!/bin/bash

set -e

cd "$(dirname "$0")/.."

echo "Running tests..."
go test ./...

export CGO_ENABLED=0

echo "Building..."
# Termux/Bionic: linux/arm64 PIE has TLS align 8; Bionic needs 64. GOOS=android avoids that.
GOOS=android GOARCH=arm64 go build -o dist/android-arm64/sgo src/start.go
GOOS=darwin  GOARCH=amd64 go build -o dist/darwin-amd64/sgo src/start.go
GOOS=darwin  GOARCH=arm64 go build -o dist/darwin-arm64/sgo src/start.go
GOOS=linux   GOARCH=amd64 go build -o dist/linux-amd64/sgo src/start.go
GOOS=linux   GOARCH=arm64 go build -o dist/linux-arm64/sgo src/start.go
GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/sgo.exe src/start.go

echo "Done."
