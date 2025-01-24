#!/bin/bash

# Check if the system is Ubuntu
if [ -f /etc/os-release ]; then
    . /etc/os-release
    if [[ "$ID" != "ubuntu" ]]; then
        echo "This script only works on Ubuntu systems."
        exit 1
    fi
else
    echo "Cannot determine OS. This script only works on Ubuntu systems."
    exit 1
fi

# Define installation path
INSTALL_PATH="/usr/local/bin"

if ! command -v figlet &> /dev/null; then
    echo "figlet not found, installing..."
    sudo apt-get update
    sudo apt-get install -y figlet
fi

# Build the project
echo "Building the project..."
go build -o vy main.go

# Move the binary to /usr/local/bin
echo "Installing vy Command line globally..."
sudo cp vy $INSTALL_PATH/

echo "vy has been installed. You can now run it from anywhere."
