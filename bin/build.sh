#!/bin/bash

set -e

cd "$(dirname "$0")/.."

echo "Running tests..."
go test ./...

export CGO_ENABLED=0

echo "Building..."
GOOS=linux   GOARCH=amd64 go build -o dist/linux-amd64/sgo src/start.go
# PIE required on Android/Termux (rejects ET_EXEC / e_type: 2)
GOOS=linux GOARCH=arm64 go build -buildmode=pie -o dist/linux-arm64/sgo src/start.go
GOOS=darwin  GOARCH=amd64 go build -o dist/darwin-amd64/sgo src/start.go
GOOS=darwin  GOARCH=arm64 go build -o dist/darwin-arm64/sgo src/start.go
GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/sgo.exe src/start.go

echo "Done."
