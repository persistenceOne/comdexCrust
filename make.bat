cd %GOPATH%\src\github.com\commitHub\commitBlockchain

go mod verify

go build -o %GOBIN%\maincli.exe main\cmd\maincli\main.go 
go build -o %GOBIN%\maind.exe main\cmd\maind\main.go