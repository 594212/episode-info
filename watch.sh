#!/usr/bin/env bash
set -x
port=6969
fuser -k $port/tcp

tsc -w &

sass --watch scss/main.scss:dist/style.css &

go run server.go --addr="localhost:$port" &
wait
