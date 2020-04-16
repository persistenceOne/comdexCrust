cd %GOPATH%\src\github.com\persistenceOne\comdexCrust

go mod verify

go build -o %GOBIN%\maincli.exe main\cmd\maincli\main.go 
go build -o %GOBIN%\maind.exe main\cmd\maind\main.go