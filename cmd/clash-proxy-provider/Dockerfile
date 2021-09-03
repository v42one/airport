ARG BUILDER=docker
FROM --platform=${BUILDPLATFORM} golang:1.17-buster AS builder-docker

ARG GOPROXY
WORKDIR /go/src
COPY ./ ./

ENV GOBIN=/go/bin
ARG TARGETARCH
RUN GOARCH=${TARGETARCH} make build

FROM busybox as builder-local

COPY ./bin /go/bin

FROM builder-${BUILDER} AS builder

FROM gcr.io/distroless/static-debian10:latest

ARG TARGETARCH
COPY --from=builder /go/bin/clash-proxy-provider-linux-${TARGETARCH} /go/bin/clash-proxy-provider

ENTRYPOINT ["/go/bin/clash-proxy-provider"]