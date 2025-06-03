#!/bin/bash
set -e

cd /opt/test
echo "Stopping existing service..."
sudo systemctl stop test || sudo service codepipe stop

echo "Backing up old binary..."
sudo mv gobuild BKP/gobuild$(date +%F)

echo "Service stopped and backup created."
exit 0
