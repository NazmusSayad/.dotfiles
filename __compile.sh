#!/bin/bash

set -e
rm -rf ./.build/bin
go run ./src/compile-scripts/main.go
