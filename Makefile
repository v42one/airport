PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = $(shell cat .version)
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)

GOBUILD = CGO_ENABLED=0 STATIC=0 go build -ldflags "-extldflags -static -s -w -X $(PKG)/pkg/version.Version=$(VERSION)+sha.$(COMMIT_SHA)"
GOBIN ?= ./bin

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

APP ?= clash-proxy-provider
WORKSPACE ?= ./cmd/$(APP)

HUB ?= docker.io/morlay
TAG ?= master
NAMESPACE ?= clash-proxy
KUBECONFIG = ${HOME}/.kube/config--hw-sg.yaml

up: tidy fmt
	WATCH_NAMESPACE=$(NAMESPACE) KUBECONFIG=$(KUBECONFIG) go run $(WORKSPACE)

build: tidy
	$(GOBUILD) -o $(GOBIN)/$(APP)-$(GOOS)-$(GOARCH) $(WORKSPACE)

fmt:
	goimports -w -l .

tidy:
	go mod tidy

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

apply.shadowsocks:
	cd ./deploy && cuem k show -o _output/ ./shadowsocks
	cd ./deploy && cuem k apply ./shadowsocks

apply.clash-proxy-provider:
	cd ./deploy && cuem k show -o _output/ ./clash-proxy-provider
	cd ./deploy && cuem k apply ./clash-proxy-provider