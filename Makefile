GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

RELEASE_TAG ?= $(shell git describe)
RELEASE_OSES ?= linux
RELEASE_ARCHES ?= amd64 arm64

DOCKER ?= docker


SHELL = /bin/bash
.PHONY: build

build: build-frontend-app
	mkdir -p dist/$(GOOS)/$(GOARCH)/bin
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o dist/$(GOOS)/$(GOARCH)/bin/hotrod ./cmd/hotrod

build-frontend-app:
	cd services/frontend/react_app && ./scripts/build.sh

dev-build-docker: build
	$(DOCKER) build -t signadot/hotrod:latest \
		--platform $(GOOS)/$(GOARCH) \
		.

build-docker: build
	$(DOCKER) build -t signadot/hotrod:$(RELEASE_TAG)-$(GOOS)-$(GOARCH) \
		--platform $(GOOS)/$(GOARCH) \
		--provenance false \
		.

push-docker: build-docker
	$(DOCKER) push signadot/hotrod:$(RELEASE_TAG)-$(GOOS)-$(GOARCH)


build-release:
	for os in $(RELEASE_OSES); do \
		for arch in $(RELEASE_ARCHES); do \
			GOOS=$$os GOARCH=$$arch $(MAKE) build-docker; \
		done; \
	done;

release-images.txt:
	mkdir -p dist
	rm -f dist/release-images.txt
	for os in $(RELEASE_OSES); do \
 		for arch in $(RELEASE_ARCHES); do \
			echo signadot/hotrod:${RELEASE_TAG}-$$os-$$arch >> dist/release-images.txt; \
		done; \
	done;

release-image: build-release release-images.txt
	for os in $(RELEASE_OSES); do \
 		for arch in $(RELEASE_ARCHES); do \
			GOOS=$$os GOARCH=$$arch $(MAKE) push-docker; \
		done; \
	done;
	$(DOCKER) manifest create --amend signadot/hotrod:$(RELEASE_TAG) \
		$(shell cat dist/release-images.txt)
	$(DOCKER) manifest push signadot/hotrod:$(RELEASE_TAG)

tag-release:
	./tag-release.sh $(RELEASE_TAG)

release: release-image tag-release

generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/route/route.proto

