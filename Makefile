NAME=euvsvirus-backend
PKG=github.com/euvsvirus-banan/backend
VERSION=$(shell cat VERSION)
REVISION=$(shell git rev-parse HEAD)

PROTOBUF_VERSION=3.11.4
GRPC_VERSION=v1.29.1
PROTOC_VERSION=v1.3.1

GOLANGCI_LINT_VERSION=v1.23.6
GOLANGCI_LINT_OPTS=--out-format=line-number \
				   --print-issued-lines=false

BUILDINFO_PKG="github.com/euvsvirus-banan/backend/internal/version"
LDFLAGS = -ldflags "\
		-X '$(BUILDINFO_PKG).Project=$(NAME)' \
		-X '$(BUILDINFO_PKG).Version=$(VERSION)' \
		-X '$(BUILDINFO_PKG).GitRevision=$(REVISION)' \
		-X '$(BUILDINFO_PKG).BuildDate=$(shell date -u)' \
		-X '$(BUILDINFO_PKG).GoVersion=$(shell go version)' \
		"


.PHONY: build-protobuf
build-protobuf:
	docker build \
		--build-arg USER=$(shell id -un) \
		--build-arg USERID=$(shell id -u) \
		--build-arg GROUP=$(shell id -gn) \
		--build-arg GROUPID=$(shell id -g) \
		--build-arg PROTOBUF_VERSION=$(PROTOBUF_VERSION) \
		--build-arg GRPC_VERSION=$(GRPC_VERSION) \
		--build-arg PROTOC_VERSION=$(PROTOC_VERSION) \
		-f Dockerfile.protobuf \
		-t protobuf-${NAME} \
		.

.PHONY: protobuf
protobuf:
	docker run \
		-v $(PWD):/go/src/$(PKG) \
		-w /go/src/$(PKG) \
		protobuf-${NAME} \
			protoc \
				--proto_path=. \
				--gofast_out=plugins=grpc:. \
				users/rpc/userspb/service.proto
	docker run \
		-v $(PWD):/go/src/$(PKG) \
		-w /go/src/$(PKG) \
		protobuf-${NAME} \
			protoc \
				--proto_path=. \
				--gofast_out=plugins=grpc:. \
				requests/rpc/requestspb/service.proto


.PHONY: docker-build
docker-build:
	docker build \
		--build-arg PKG=$(PKG) \
		-f Dockerfile \
		-t $(PKG):dev \
		.

.PHONY: build
build:
	go build -o /tmp/backend $(LDFLAGS) $(PKG)/cmd

.PHONY: docker-run
docker-run: docker-build
	docker run \
		--name $(NAME) \
		-p 65010:65010 \
		-v $(PWD)/data:/euvsvirus-backend \
		--rm \
		-it \
		$(PKG):dev \
			--addr 0.0.0.0:65010

.PHONY: linter
linter:
	docker run \
		-v $(PWD):/go/src/$(PKG) \
		-w /go/src/$(PKG) \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
			golangci-lint run \
				--out-format=line-number \
				--print-issued-lines=false \
				--timeout 3m
