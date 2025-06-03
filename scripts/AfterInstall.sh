#!/bin/bash
set -e

echo "Detecting OS type..."

OS_NAME=$(grep -E '^NAME=' /etc/os-release | cut -d= -f2 | tr -d '"')

if [[ "$OS_NAME" == "Amazon Linux" ]]; then
    USERNAME="ec2-user"
elif [[ "$OS_NAME" == "Ubuntu" ]]; then
    USERNAME="ubuntu"
else
    USERNAME=$(whoami)  # Fallback to the current user
fi

echo "Detected OS: $OS_NAME"
echo "Using user: $USERNAME"


echo "Moving application files..."
if [[ -f "/opt/Build/gobuild" ]]; then
    sudo mv /opt/Build/gobuild /opt/test/
else
    echo "Error: /opt/Build/gobuild not found!" >&2
    exit 1
fi

echo "Setting permissions..."
sudo chmod +x /opt/test/gobuild
sudo chown -R $USERNAME:$USERNAME /opt/test/gobuild

echo "Starting service..."
sudo systemctl restart test || sudo systemctl start test

echo "Service started successfully."
exit 0
