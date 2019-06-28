cd %GOPATH%\src\github.com\comdex-blockchain

go get github.com/golang/dep/cmd/dep
go get ./...

dep init -v
dep ensure -v

go build -o %GOBIN%\assetcli.exe asset\cmd\assetcli\main.go 
go build -o %GOBIN%\assetd.exe asset\cmd\assetd\main.go

go build -o %GOBIN%\maincli.exe main\cmd\maincli\main.go 
go build -o %GOBIN%\maind.exe main\cmd\maind\main.go

go build -o %GOBIN%\fiatcli.exe fiat\cmd\fiatcli\main.go 
go build -o %GOBIN%\fiatd.exe fiat\cmd\fiatd\main.go  

go build -o %GOBIN%\blockExplorer.exe blockExplorer\main\main.go  

PAUSE
