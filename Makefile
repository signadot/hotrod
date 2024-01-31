GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)


RELEASE_TAG ?= $(shell git describe)
RELEASE_OSES ?= linux
RELEASE_ARCHES ?= amd64 arm64

DOCKER ?= docker


SHELL = /bin/bash
.PHONY: build


build:
	mkdir -p dist/$(GOOS)/$(GOARCH)/bin
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o dist/$(GOOS)/$(GOARCH)/bin/hotrod ./cmd/hotrod

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

tag-release:
	(cd k8s/base && kustomize edit set signadot/hotrod:$(RELEASE_TAG))
	git commit -m tag-release-$(RELEASE_TAG) k8s/base
	git tag -a -m release-$(RELEASE_TAG) $(RELEASE_TAG)
	git push origin $(RELEASE_TAG)

release: build-release release-images.txt tag-release

	for os in $(RELEASE_OSES); do \
 		for arch in $(RELEASE_ARCHES); do \
			GOOS=$$os GOARCH=$$arch $(MAKE) push-docker; \
		done; \
	done;
	$(DOCKER) manifest create --amend signadot/hotrod:$(RELEASE_TAG) \
		$(shell cat dist/release-images.txt)
	$(DOCKER) manifest push signadot/hotrod:$(RELEASE_TAG)

generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		services/route/route.proto
