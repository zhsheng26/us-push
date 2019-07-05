#!/usr/bin/env bash
cd ./server/

go build -o udp_server


cd ././client/

go build -o client
