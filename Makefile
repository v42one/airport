PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = $(shell cat .version)
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)

GOBUILD = CGO_ENABLED=0 STATIC=0 go build -ldflags "-extldflags -static -s -w -X $(PKG)/pkg/version.Version=$(VERSION)+sha.$(COMMIT_SHA)"
GOBIN ?= ./bin

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

APP ?= clash-proxy-operator
WORKSPACE ?= ./cmd/$(APP)

HUB ?= docker.io/morlay
TAG ?= master
NAMESPACE ?= clash-proxy

up:
	WATCH_NAMESPACE=$(NAMESPACE) go run $(WORKSPACE)

build:
	$(GOBUILD) -o $(GOBIN)/$(APP)-$(GOOS)-$(GOARCH) $(WORKSPACE)

PLATFORMS = amd64 arm64

BUILDER ?= docker

buildx:
	for arch in $(PLATFORMS); do \
  		$(MAKE) build GOOS=linux GOARCH=$${arch}; \
  	done

dockerx:
	docker buildx build \
		--push \
		--build-arg=BUILDER=$(BUILDER) \
		--build-arg=APP=$(APP) \
		$(foreach arch,$(PLATFORMS),--platform=linux/$(arch)) \
		--tag $(HUB)/$(APP):$(TAG) \
		-f $(WORKSPACE)/Dockerfile .

dockerx.dev: buildx
	$(MAKE) dockerx BUILDER=local


WORKING_DIR = ./deploy

include $(WORKING_DIR)/components/Makefile

apply.%:
	$(MAKE) apply COMPONENT=$*