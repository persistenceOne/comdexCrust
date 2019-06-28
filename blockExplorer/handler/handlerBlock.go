package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	
	"github.com/comdex-blockchain/blockExplorer/constants"
	"github.com/comdex-blockchain/blockExplorer/dataTypes"
	
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/rpc/client"
)

func HandleBlock(w http.ResponseWriter, r *http.Request) {
	c, err := constants.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	w.WriteHeader(http.StatusOK)
	blockChannel := make(chan interface{})
	client := client.NewHTTP(constants.BlockchainIP, "/websocket")
	err = client.Start()
	if err != nil {
		fmt.Println("Start Error: ", err)
		return
	}
	defer client.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), constants.Timeout)
	defer cancel()
	query := query.MustParse(constants.EventBlock)
	err = client.Subscribe(ctx, constants.Subscriber, query, blockChannel)
	if err != nil {
		fmt.Println(err)
	}
SendBlock:
	for {
		select {
		case interfaceBlockData := <-blockChannel:
			var blockData dataTypes.RawBlock
			rawBlockData, err := json.Marshal(interfaceBlockData)
			if err != nil {
				log.Println("Json Marshal Error: ", err)
				break
			}
			err = json.Unmarshal(rawBlockData, &blockData)
			if err != nil {
				log.Println("Json UnMarshal Error:", err)
				break
			}
			fmt.Println("Sending Message")
			err = c.WriteJSON(blockData)
			if err != nil {
				log.Println("Writer Error:", err)
				break SendBlock
			}
		}
	}
}
