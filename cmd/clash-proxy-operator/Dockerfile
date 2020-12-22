ARG BUILDER=docker
FROM --platform=${BUILDPLATFORM} golang:1.15-buster AS builder-docker

ARG GOPROXY
WORKDIR /go/src
COPY ./ ./

ENV GOBIN=/go/bin
ARG TARGETARCH
RUN GOARCH=${TARGETARCH} make build

FROM busybox as builder-local

COPY ./bin /go/bin

FROM builder-${BUILDER} AS builder

FROM ghcr.io/querycap/distroless/static-debian10:latest

ARG TARGETARCH
COPY --from=builder /go/bin/clash-proxy-operator-linux-${TARGETARCH} /go/bin/clash-proxy-operator

ENV SS_PROXY_CONFIG_MAP_NAME=ss-proxy

ENTRYPOINT ["/go/bin/clash-proxy-operator"]