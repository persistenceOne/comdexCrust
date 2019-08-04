PACKAGES := $(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
GOSUM := $(shell which gosum)

export GO111MODULE=on

BUILD_TAGS := netgo
BUILD_TAGS := $(strip ${BUILD_TAGS})

BUILD_FLAGS := -tags "${BUILD_TAGS}"

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
