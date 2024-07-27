#!/bin/bash

sudo cp /var/log/nginx/access.log /home/linuxuser/
sudo chown linuxuser /home/linuxuser/access.log

/home/linuxuser/goAnalyzeNGINX

sudo cp -r /home/linuxuser/countries.json /var/www/website/projects/map

rm /home/linuxuser/countries.json
rm /home/linuxuser/access.log
