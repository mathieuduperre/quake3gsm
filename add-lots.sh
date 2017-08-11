#!/bin/bash
CAP=$1
CAP=${CAP:-1000}
curl -d @hello-plugin.json http://127.0.0.1:24601/v1/plugins
for i in $(seq 1 $CAP);
do
  sed -i .bck 's/helloSvc[[:digit:]]*/helloSvc'"${i}"'/g' hello-service.json
  curl http://127.0.0.1:24601/v1/services -d @hello-service.json
  sleep 0.1
done
