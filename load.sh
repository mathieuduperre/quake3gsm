#!/bin/bash
curl -H "Content-Type: application/json" --data @global.json http://172.16.7.1:24601/v1/sites
./add-site.sh
./add-plugin.sh
./add-service.sh
