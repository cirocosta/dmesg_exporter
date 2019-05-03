build:
	go build -i -v

fmt:
	go fmt ./...

test:
	go test ./...

image:
	docker build -t dmesg_exporter .

test-image:
	docker build -t dmesg_exporter --target tests .
