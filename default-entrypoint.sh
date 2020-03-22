#!/bin/sh
cd /build
go mod download
go build -buildmode=c-shared -o /dist/${BINARY_NAME}.so -v .