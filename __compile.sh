#!/bin/bash
set -e

echo "> Cleaning build directory..."
rm -rf ./.build/bin

echo ""
echo "> Compiling Go scripts..."
go run ./src/compile-scripts/main.go

echo ""
echo "Done!"
