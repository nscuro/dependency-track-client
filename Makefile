build:
	go build -v ./...
.PHONY: build

test:
	go test -v ./...
.PHONY: test

docker:
	docker build -t nscuro/dependency-track-client:latest -f build/docker/Dockerfile .
.PHONY: docker
