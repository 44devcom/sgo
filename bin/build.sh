#!/bin/bash

GOOS=linux   GOARCH=amd64 go build -o ../dist/linux-amd64/sgo ../src/start.go
GOOS=darwin  GOARCH=amd64 go build -o ../dist/darwin-amd64/sgo ../src/start.go
GOOS=darwin  GOARCH=arm64 go build -o ../dist/darwin-arm64/sgo ../src/start.go
GOOS=windows GOARCH=amd64 go build -o ../dist/windows-amd64/sgo.exe ../src/start.go