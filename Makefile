PACKAGES=$(shell go list ./... | grep -v '/vendor/')
VERSION := $(shell echo $(shell git describe --always) | sed 's/^v//')
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TAGS := netgo
BUILD_TAGS := $(strip ${BUILD_TAGS})

LD_FLAGS := -s -w \
	-X github.com/commitHub/commitBlockchain/version.Version=${VERSION} \
	-X github.com/commitHub/commitBlockchain/version.GitCommit=${COMMIT} \
	-X github.com/commitHub/commitBlockchain/version.BuildTags=${BUILD_TAGS}

BUILD_FLAGS := -tags "${BUILD_TAGS}" -ldflags "${LD_FLAGS}"

all: get_tools get_vendor_deps build

get_tools:
	go get github.com/golang/dep/cmd/dep

build:
	go build ${BUILD_FLAGS} -o ${GOBIN}/assetcli asset/cmd/assetcli/main.go && go build ${BUILD_FLAGS} -o ${GOBIN}/assetd asset/cmd/assetd/main.go
	go build ${BUILD_FLAGS} -o ${GOBIN}/maincli main/cmd/maincli/main.go && go build ${BUILD_FLAGS} -o ${GOBIN}/maind main/cmd/maind/main.go
	go build ${BUILD_FLAGS} -o ${GOBIN}/fiatcli fiat/cmd/fiatcli/main.go && go build ${BUILD_FLAGS} -o ${GOBIN}/fiatd fiat/cmd/fiatd/main.go
	go build ${BUILD_FLAGS} -o ${GOBIN}/blockExplorer blockExplorer/main/main.go

get_vendor_deps:
	@rm -rf vendor/
	@dep ensure -v

test:
	@go test $(PACKAGES)

benchmark:
	@go test -bench=. $(PACKAGES)

.PHONY: all build test benchmark
