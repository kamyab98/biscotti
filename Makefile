COMMIT := $(shell git rev-parse HEAD)
VERSION ?= $(shell git describe --tags ${COMMIT})
IMG ?= docker.yektanet.tech/projects/platform/biscotti:${VERSION}

CGO_ENABLED=0
GOOS=linux
GOARCH=amd64

mod:
	go mod tidy

build:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go build -o ./bin/biscotti ./cmd/biscotti/main.go

docker-build:
	docker build -t ${IMG} .

docker-push:
	docker push ${IMG}

run:
	CGO_ENABLED=${CGO_ENABLED} GOOS=${GOOS} GOARCH=${GOARCH} go run ./cmd/biscotti/main.go
