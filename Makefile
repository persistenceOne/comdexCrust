PACKAGES := $(shell go list ./... | grep -v '/simulation')
VERSION := $(shell git branch | grep \* | cut -d ' ' -f2)
COMMIT := $(shell git rev-parse --short HEAD)
GOSUM := $(shell which gosum)

export GO111MODULE=on


BUILD_TAGS := -s  -w \
	-X github.com/persistenceOne/comdexCrust/version.Version=${VERSION} \
	-X github.com/persistenceOne/comdexCrust/version.Commit=${COMMIT}

ifneq (${GOSUM},)
	ifneq (${wildcard go.sum},)
		BUILD_TAGS += -X github.com/persistenceOne/comdexCrust/version.VendorHash=$(shell ${GOSUM} go.sum)
	endif
endif

BUILD_FLAGS += -ldflags "${BUILD_TAGS}"

all: install

build: go.sum
ifeq (${OS},Windows_NT)
	go build -mod=readonly ${BUILD_FLAGS} -o bin/maind.exe main/cmd/maind/
	go build -mod=readonly ${BUILD_FLAGS} -o bin/maincli.exe main/cmd/maincli/
else
	go build -mod=readonly ${BUILD_FLAGS} -o bin/maind main/cmd/maind/
	go build -mod=readonly ${BUILD_FLAGS} -o bin/maincli main/cmd/maincli/
endif

install: go.sum
	go install -mod=readonly ${BUILD_FLAGS} ./main/cmd/maind
	go install -mod=readonly ${BUILD_FLAGS} ./main/cmd/maincli

go.sum:
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

.PHONY: all build install  go.sum
