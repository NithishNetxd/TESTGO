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
if [[ -f "/opt/Build/codepipeline" ]]; then
    sudo mv /opt/Build/codepipeline /opt/code/
else
    echo "Error: /opt/Build/codepipeline not found!" >&2
    exit 1
fi

echo "Setting permissions..."
sudo chmod +x /opt/code/codepipeline
sudo chown -R $USERNAME:$USERNAME /opt/code/codepipeline

echo "Starting service..."
sudo systemctl restart codepipeline || sudo systemctl start codepipeline

echo "Service started successfully."
exit 0
