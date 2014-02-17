#!/bin/bash

source ~/src/golang-crosscompile/crosscompile.bash

rm -rf build/
mkdir build/

go-linux-386 build -o build/dm-linux-386
go-linux-amd64 build -o build/dm-linux-amd64

go-darwin-386 build -o build/dm-mac-386
go-darwin-amd64 build -o build/dm-mac-amd64

go-windows-386 build -o build/dm-win-386
go-windows-amd64 build -o build/dm-win-amd64
