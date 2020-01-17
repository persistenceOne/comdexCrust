cd %GOPATH%\src\github.com\persistenceOne\persistenceSDK

go mod verify

go build -o %GOBIN%\maincli.exe main\cmd\maincli\main.go 
go build -o %GOBIN%\maind.exe main\cmd\maind\main.go