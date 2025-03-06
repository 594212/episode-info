#!/usr/bin/env bash
set -x

tsc -w &

sass --watch scss/main.scss:dist/style.css &

fuser -k 6969/tcp
python -m http.server 6969 &
wait
