#!/bin/env bash

(cd front && npm install && npm run build)

cp -R front/public/* back/cmd/code-radar-embed/

(cd back && go build -o ../target/code-radar cmd/code-radar-embed/main.go)
