#!/bin/bash
mkdir -p bin
go build -o bin/goserve main.go
echo "Build complete: bin/goserve"
