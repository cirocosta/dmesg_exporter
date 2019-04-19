FROM golang:alpine AS golang


FROM golang AS base

	ENV CGO_ENABLED=0
	RUN apk add --update git

	ADD . /src
	WORKDIR /src

	RUN go mod download


FROM base AS builder

	RUN set -x && \
		go build \
			-tags netgo -v -a \
			-ldflags "-X main.version=$(cat ./VERSION) -extldflags \"-static\"" && \
		mv \
			./dmesg_exporter \
			/usr/bin/dmesg_exporter


FROM base AS tests

	RUN set -x && \
		go test -v ./...


FROM alpine

	COPY \
		--from=builder \
		/usr/bin/dmesg_exporter \
		/usr/bin/dmesg_exporter
