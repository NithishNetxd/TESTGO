#!/bin/bash
set -e

cd /opt/code/
echo "Stopping existing service..."
sudo systemctl stop codepipe || sudo service codepipe stop

echo "Backing up old binary..."
sudo mv codepipeline BKP/codepipeline_$(date +%F)

echo "Service stopped and backup created."
exit 0
