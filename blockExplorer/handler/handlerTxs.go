package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	
	"github.com/comdex-blockchain/blockExplorer/constants"
	"github.com/comdex-blockchain/blockExplorer/dataTypes"
	
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/rpc/client"
)

func HandleTxs(w http.ResponseWriter, r *http.Request) {
	c, err := constants.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	w.WriteHeader(http.StatusOK)
	txChannel := make(chan interface{})
	client := client.NewHTTP(constants.BlockchainIP, "/websocket")
	err = client.Start()
	if err != nil {
		fmt.Println("Start Error: ", err)
		return
	}
	defer client.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), constants.Timeout)
	defer cancel()
	query := query.MustParse(constants.EventTx)
	err = client.Subscribe(ctx, constants.Subscriber, query, txChannel)
	if err != nil {
		fmt.Println(err)
	}
SendTxHash:
	for {
		select {
		case txData := <-txChannel:
			var rawTxData dataTypes.RawTx
			var rawComdexBlock dataTypes.ComdexBlock
			txRawData, err := json.Marshal(txData)
			if err != nil {
				log.Println("Json Marshal Error: ", err)
				break
			}
			err = json.Unmarshal(txRawData, &rawTxData)
			if err != nil {
				log.Println("Json UnMarshal Error:", err)
				break
			}
			blockHeight := rawTxData.Height
			blockResponse, err := http.Get(constants.BlockHeightURL + strconv.FormatInt(blockHeight, 10))
			if err != nil {
				log.Println("Block Get Response Error:", err)
				break
			}
			err = json.NewDecoder(blockResponse.Body).Decode(&rawComdexBlock)
			if err != nil {
				log.Println("Block Decoder Error:", err)
				break
			}
			txHash := rawComdexBlock.Result.Block.Header.Data_hash
			comdexTx := dataTypes.ComdexTx{
				Hash:   txHash,
				Height: blockHeight,
			}
			err = c.WriteJSON(comdexTx)
			if err != nil {
				log.Println("Writer Error:", err)
				break SendTxHash
			}
		}
	}
}
