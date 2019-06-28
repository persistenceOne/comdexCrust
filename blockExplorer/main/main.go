package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	
	"github.com/comdex-blockchain/blockExplorer/constants"
	"github.com/comdex-blockchain/blockExplorer/handler"
)

func checkOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	if origin[0] == "http://localhost:9000" {
		return true
	}
	// u, err := url.Parse(origin[0])
	// if err != nil {
	// 	return false
	// }
	return false
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

// https://tendermint.com/rpc/

func main() {
	flag.Parse()
	log.SetFlags(0)
	fmt.Println(constants.Upgrader.HandshakeTimeout)
	constants.Upgrader.CheckOrigin = checkOrigin
	http.HandleFunc("/block", handler.HandleBlock)
	http.HandleFunc("/tx", handler.HandleTxs)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*constants.Addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
Block Explorer
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>
You need to make WebSocket connection on ws://127.0.0.1:2259/tx (for Tx Hash) and on ws://127.0.0.1:2259/block (for Blocks).
<p>
</body>
</html>
`))
