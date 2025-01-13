#!/bin/bash

cd /opt/stfc/
sudo pkill stfc
sudo mv stfc BKP/stfc_$(date +%F)
