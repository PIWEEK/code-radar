#!/bin/bash

(cd back && go run cmd/code-radar/main.go ${1:-https://github.com/PIWEEK/code-radar.git})

