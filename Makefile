SHELL = /bin/sh
.PHONY: build


build:
	mkdir -p dist/bin
	go build -o dist/bin/hotrod ./cmd/hotrod

docker-build:
	docker build -t signadot/hotrod:latest .
