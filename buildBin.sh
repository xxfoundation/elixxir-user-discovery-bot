#!/bin/bash
# My first script

rm release/*
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-w -s" -o release/registration.darwin64 main.go
mv release/registration.darwin64 release/udb.binary
cp release/udb.binary ../localenvironment/binaries
