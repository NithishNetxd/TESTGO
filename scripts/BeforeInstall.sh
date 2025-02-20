#!/bin/bash

cd /opt/code/
sudo service codepipe stop
sudo mv codepipeline BKP/codepipeline_$(date +%F)
