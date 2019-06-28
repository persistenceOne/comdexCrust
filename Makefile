PACKAGES=$(shell go list ./... | grep -v '/vendor/')

all: get_tools get_vendor_deps build

get_tools:
	go get github.com/golang/dep/cmd/dep

build:
	go build -o ${GOBIN}/assetcli asset/cmd/assetcli/main.go && go build -o ${GOBIN}/assetd asset/cmd/assetd/main.go
	go build -o ${GOBIN}/maincli main/cmd/maincli/main.go && go build -o ${GOBIN}/maind main/cmd/maind/main.go
	go build -o ${GOBIN}/fiatcli fiat/cmd/fiatcli/main.go && go build -o ${GOBIN}/fiatd fiat/cmd/fiatd/main.go
	go build -o ${GOBIN}/blockExplorer blockExplorer/main/main.go

get_vendor_deps:
	@rm -rf vendor/
	@dep ensure -v

test:
	@go test $(PACKAGES)

benchmark:
	@go test -bench=. $(PACKAGES)

.PHONY: all build test benchmark
