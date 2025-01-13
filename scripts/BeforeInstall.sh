#!/bin/bash

cd /opt/stfc/
sudo pkill stfc
mv stfc BKP/stfc_$(date +%F)
