#!/bin/bash
set -e

cd /opt/code/
sudo chmod +x codepipeline

echo "Starting service..."
sudo systemctl restart codepipeline || sudo service codepipeline start

echo "Service started successfully."
exit 0
