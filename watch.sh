#!/usr/bin/env bash
set -x
port=6969
fuser -k $port/tcp

tsc -w &

sass --watch scss/:dist/ &

go run server.go --addr="localhost:$port" &
wait
