package constants

import (
	"flag"
	"time"
	
	"github.com/gorilla/websocket"
)

const IP = "127.0.0.1:2259"
const BlockchainIP = "tcp://0.0.0.0:36657"
const BlockHeightURL = "http://localhost:36657/block?height="

var Addr = flag.String("addr", IP, "http service address")
var Upgrader = websocket.Upgrader{} // use default options

const Subscriber = "Test-Client"
const EventBlock = "tm.event='NewBlock'"
const EventTx = "tm.event='Tx'"

const Timeout = 10000 * time.Millisecond
