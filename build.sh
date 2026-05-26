#!/bin/bash
set -e

echo "Building Speedtest Application..."

mkdir -p build

# Build Linux 386
echo "Building Linux 386..."
CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -v -o build/speedtest-linux-386 .
chmod +x build/speedtest-linux-386

# Build Linux AMD64
echo "Building Linux AMD64..."
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v -o build/speedtest-linux-amd64 .
chmod +x build/speedtest-linux-amd64

# Build Windows 386
echo "Building Windows 386..."
CGO_ENABLED=1 GOOS=windows GOARCH=386 go build -v -o build/speedtest-windows-386.exe .

# Build Windows AMD64
echo "Building Windows AMD64..."
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -v -o build/speedtest-windows-amd64.exe .

echo -e "\nBuild completed successfully!"
ls -la build/
