#!/bin/ash

dep ensure
GOOS=darwin GOARCH=amd64 packr build -o binaries/amd64/${version}/darwin/diferencia_darwin_amd64
GOOS=windows GOARCH=amd64 packr build -o binaries/amd64/${version}/windows/diferencia_windows_amd64.exe
GOOS=linux GOARCH=amd64 packr build -o binaries/amd64/${version}/linux/diferencia_linux_amd64