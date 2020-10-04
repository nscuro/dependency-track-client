DIST_DIR=./dist
LDFLAGS="-s -w -X github.com/nscuro/dependency-track-client/internal/version.Version=$(shell git describe --tags)"
GCFLAGS="all=-trimpath=$(shell pwd)"
ASMFLAGS="all=-trimpath=$(shell pwd)"

build:
	go build -v ./...
.PHONY: build

test:
	go test -v ./...
.PHONY: test

bench:
	go test -bench=. ./...
.PHONY: bench

install:
	GO111MODULE=on CGO_ENABLED=0 go install -ldflags=${LDFLAGS} -gcflags=${GCFLAGS} -asmflags=${ASMFLAGS} -v ./cmd/dtrack
.PHONY: install

docker:
	docker build -t nscuro/dependency-track-client:latest -f ./build/docker/Dockerfile .
.PHONY: docker

pre-dist:
	mkdir -p ${DIST_DIR}
.PHONY: 

windows:
	GOOS=windows GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 \
	go build -ldflags=${LDFLAGS} -gcflags=${GCFLAGS} -asmflags=${ASMFLAGS} \
		-o ${DIST_DIR}/dtrack-windows-amd64.exe ./cmd/dtrack
.PHONY: windows

darwin:
	GOOS=darwin GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 \
	go build -ldflags=${LDFLAGS} -gcflags=${GCFLAGS} -asmflags=${ASMFLAGS} \
		-o ${DIST_DIR}/dtrack-darwin-amd64 ./cmd/dtrack
.PHONY: darwin

linux:
	GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 \
	go build -ldflags=${LDFLAGS} -gcflags=${GCFLAGS} -asmflags=${ASMFLAGS} \
		-o ${DIST_DIR}/dtrack-linux-amd64 ./cmd/dtrack
.PHONY: linux

clean:
	rm -rf ${DIST_DIR}; go clean ./...
.PHONY: clean

all: clean build test pre-dist windows darwin linux
.PHONY: all