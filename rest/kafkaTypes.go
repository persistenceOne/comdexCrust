package rest

import (
	"strconv"
	
	"github.com/Shopify/sarama"
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// Ticket : is a type that implements string
type Ticket string

// KafkaMsg : is a store that can be stored in kafka queues
type KafkaMsg struct {
	Msg      sdk.Msg     `json:"msg"`
	TicketID Ticket      `json:"ticketID"`
	KafkaTx  KafkaTxCtx  `json:"kafkaTxCtx"`
	KafkaCli KafkaCliCtx `json:"kafkaCliCtx"`
}

// NewKafkaMsgFromRest : makes a msg to send to kafka queue
func NewKafkaMsgFromRest(msg sdk.Msg, ticketID Ticket, txCtx context2.TxContext, cliCtx context.CLIContext, passPhrase string) KafkaMsg {
	kafkaTx := KafkaTxCtx{
		PassPhrase:    passPhrase,
		AccountNumber: txCtx.AccountNumber,
		Sequence:      txCtx.Sequence,
		Gas:           txCtx.Gas,
		ChainID:       txCtx.ChainID,
		Memo:          txCtx.Memo,
		Fee:           txCtx.Fee,
	}
	kafkaCli := KafkaCliCtx{
		Height:          cliCtx.Height,
		Gas:             cliCtx.Gas,
		GasAdjustment:   strconv.FormatFloat(cliCtx.GasAdjustment, 'E', -1, 64),
		NodeURI:         cliCtx.NodeURI,
		FromAddressName: cliCtx.FromAddressName,
		AccountStore:    cliCtx.AccountStore,
		TrustNode:       cliCtx.TrustNode,
		UseLedger:       cliCtx.UseLedger,
		Async:           cliCtx.Async,
		JSON:            cliCtx.JSON,
		PrintResponse:   cliCtx.PrintResponse,
		DryRun:          cliCtx.DryRun,
	}
	
	return KafkaMsg{
		Msg:      msg,
		TicketID: ticketID,
		KafkaTx:  kafkaTx,
		KafkaCli: kafkaCli,
	}
	
}

// KafkaTxAndKafkaCliFromKafkaMsg : sets the txctx and clictx again to consume
func KafkaTxAndKafkaCliFromKafkaMsg(msg KafkaMsg, cliCtx context.CLIContext) (context2.TxContext, context.CLIContext, string) {
	
	txCtx := context2.TxContext{
		Codec:         cliCtx.Codec,
		Gas:           msg.KafkaTx.Gas,
		ChainID:       msg.KafkaTx.ChainID,
		AccountNumber: msg.KafkaTx.AccountNumber,
		Sequence:      msg.KafkaTx.Sequence,
	}
	
	cliCtx.Height = msg.KafkaCli.Height
	cliCtx.Gas = msg.KafkaCli.Gas
	gasAdjustment, err := strconv.ParseFloat(msg.KafkaCli.GasAdjustment, 64)
	if err != nil {
		panic(err)
	}
	cliCtx.GasAdjustment = gasAdjustment
	cliCtx.NodeURI = msg.KafkaCli.NodeURI
	cliCtx.FromAddressName = msg.KafkaCli.FromAddressName
	cliCtx.AccountStore = msg.KafkaCli.AccountStore
	cliCtx.TrustNode = msg.KafkaCli.TrustNode
	cliCtx.UseLedger = msg.KafkaCli.UseLedger
	cliCtx.Async = msg.KafkaCli.Async
	cliCtx.JSON = msg.KafkaCli.JSON
	cliCtx.PrintResponse = msg.KafkaCli.PrintResponse
	cliCtx.DryRun = msg.KafkaCli.DryRun
	
	return txCtx, cliCtx, msg.KafkaTx.PassPhrase
}

// KafkaTxCtx : auth.tx without codec
type KafkaTxCtx struct {
	PassPhrase    string
	AccountNumber int64
	Sequence      int64
	Gas           int64
	ChainID       string
	Memo          string
	Fee           string
}

// KafkaCliCtx : client tx without codec
type KafkaCliCtx struct {
	Height          int64
	Gas             int64
	GasAdjustment   string
	NodeURI         string
	FromAddressName string
	AccountStore    string
	TrustNode       bool
	UseLedger       bool
	Async           bool
	JSON            bool
	PrintResponse   bool
	DryRun          bool
}

// TicketIDResponse : is a json structure to send TicketID to user
type TicketIDResponse struct {
	TicketID Ticket `json:"TicketID" valid:"required~TicketID is mandatory,length(20)~RelayerAddress length should be 20" `
}

// KafkaState : is a struct showing the state of kafka
type KafkaState struct {
	KafkaDB   *dbm.GoLevelDB
	Admin     sarama.ClusterAdmin
	Consumer  sarama.Consumer
	Consumers map[string]sarama.PartitionConsumer
	Producer  sarama.SyncProducer
	Topics    []string
}

// NewKafkaState : returns a kafka state
func NewKafkaState(kafkaPorts []string) KafkaState {
	kafkaDB, _ := dbm.NewGoLevelDB("KafkaDB", DefaultCLIHome)
	admin := KafkaAdmin(kafkaPorts)
	producer := NewProducer(kafkaPorts)
	consumer := NewConsumer(kafkaPorts)
	var consumers = make(map[string]sarama.PartitionConsumer)
	
	for _, topic := range Topics {
		partitionConsumer := PartitionConsumers(consumer, topic)
		consumers[topic] = partitionConsumer
	}
	
	return KafkaState{
		KafkaDB:   kafkaDB,
		Admin:     admin,
		Consumer:  consumer,
		Consumers: consumers,
		Producer:  producer,
		Topics:    Topics,
	}
}
