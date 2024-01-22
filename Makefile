SHELL = /bin/sh
.PHONY: build


build:
	mkdir -p dist/bin
	go build -o dist/bin/hotrod ./cmd/hotrod

docker-build:
	docker build -t signadot/hotrod:latest .

generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/route/route.proto
