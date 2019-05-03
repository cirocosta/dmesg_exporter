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
			-o /usr/bin/dmesg_exporter \
			-ldflags "-X main.version=$(cat ./VERSION) -extldflags \"-static\""


FROM base AS tests

	RUN set -x && \
		go test -v ./...


FROM alpine

	COPY \
		--from=builder \
		/usr/bin/dmesg_exporter \
		/usr/bin/dmesg_exporter

	ENTRYPOINT [ "/usr/bin/dmesg_exporter" ]
