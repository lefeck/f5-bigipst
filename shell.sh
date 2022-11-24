#!/usr/bin/env bash

mkdir -p bigipst-bin/{bigipst-linux-amd64,bigipst-windows-amd64,bigipst-darwin-amd64}
# build Linux executable package
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o f5-bigipst main.go
mv f5-bigipst bigipst-bin/bigipst-linux-amd64
cd bigipst-bin
tar -zcf bigipst-linux-amd64.tar.gz  ./bigipst-linux-amd64
rm -rf bigipst-linux-amd64

# build Windows executable package
cd ..
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o f5-bigipst.exe main.go
mv f5-bigipst.exe bigipst-bin/bigipst-windows-amd64
cd bigipst-bin
tar -zcf bigipst-windows-amd64.tar.gz ./bigipst-windows-amd64
rm -rf ./bigipst-windows-amd64

# build Mac OS executable package
cd ..
go build -o f5-bigipsts main.go
mv f5-bigipsts bigipst-bin/bigipst-darwin-amd64
cd bigipst-bin
tar -zcf bigipst-darwin-amd64.tar.gz ./bigipst-darwin-amd64
rm -rf ./bigipst-darwin-amd64