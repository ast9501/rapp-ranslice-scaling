#!/bin/bash
curl -i -u 'sample1:sample1' -X POST -d @sample-fault-event.json --header "Content-Type: application/json" https://192.168.101.136:30417/eventListener/v5 -k
