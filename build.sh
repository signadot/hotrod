#!/bin/bash


# esc -pkg frontend -o ./services/frontend/gen_assets.go  -prefix ./services/frontend/web_assets ./services/frontend/web_assets
go build -o hotrod-linux-amd64
# docker build -t foxish/example-hotrod:latest .
# docker push foxish/example-hotrod:latest