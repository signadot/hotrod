#!/bin/bash

# Embed web assets into the binary.
go run github.com/mjibson/esc -pkg frontend -o services/frontend/gen_assets.go  -prefix services/frontend/web_assets services/frontend/web_assets

go build -o hotrod
