package utils

import (
	"fmt"
	"os"
	"time"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/rest"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/auth"
	authctx "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/common"
)

// SendTx implements a auxiliary handler that facilitates sending a series of
// messages in a signed transaction given a TxContext and a QueryContext. It
// ensures that the account exists, has a proper number and sequence set. In
// addition, it builds and signs a transaction with the supplied messages.
// Finally, it broadcasts the signed transaction to a node.
func SendTx(txCtx authctx.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg) error {
	txCtx, err := prepareTxContext(txCtx, cliCtx)
	if err != nil {
		return err
	}
	autogas := cliCtx.DryRun || (cliCtx.Gas == 0)
	if autogas {
		txCtx, err = EnrichCtxWithGas(txCtx, cliCtx, cliCtx.FromAddressName, msgs)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "estimated gas = %v\n", txCtx.Gas)
	}
	if cliCtx.DryRun {
		return nil
	}
	
	// passphrase, err := keys.GetPassphrase(cliCtx.FromAddressName)
	// if err != nil {
	// 	return err
	// }
	
	// build and sign the transaction
	txBytes, err := txCtx.BuildAndSign(cliCtx.FromAddressName, "1234567890", msgs)
	if err != nil {
		return err
	}
	// broadcast to a Tendermint node
	return cliCtx.EnsureBroadcastTx(txBytes)
}

// SendTxWithResponse : send tx with response
func SendTxWithResponse(txCtx authctx.TxContext, cliCtx context.CLIContext, msgs []sdk.Msg, passphrase string) ([]byte, error) {
	txCtx, err := prepareTxContext(txCtx, cliCtx)
	if err != nil {
		return nil, err
	}
	autogas := cliCtx.DryRun || (cliCtx.Gas == 0)
	if autogas {
		txCtx, err = EnrichCtxWithGas(txCtx, cliCtx, cliCtx.FromAddressName, msgs)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stdout, "estimated gas = %v\n", txCtx.Gas)
	}
	if cliCtx.DryRun {
		return nil, nil
	}
	txBytes, err := txCtx.BuildAndSign(cliCtx.FromAddressName, passphrase, msgs)
	if err != nil {
		return nil, err
	}
	return cliCtx.EnsureBroadcastTxWithResponse(txBytes)
}

// -------------------------------------------------------------------------------------------//
// SendTxSWithResponse : send tx with response
func SendTxSWithResponse(txSCtx []authctx.TxContext, cliCtx []context.CLIContext, msgs []sdk.Msg, passphrase []string) ([]byte, error) {
	stdTxs := auth.StdTx{}
	for i, txCtx := range txSCtx {
		txCtx, err := prepareTxContext(txCtx, cliCtx[i])
		if err != nil {
			return nil, err
		}
		txSCtx[i] = txCtx
		
		autogas := cliCtx[i].DryRun || (cliCtx[i].Gas == 0)
		if autogas {
			txCtx, err = EnrichCtxWithGas(txCtx, cliCtx[i], cliCtx[i].FromAddressName, []sdk.Msg{msgs[i]})
			if err != nil {
				return nil, err
			}
			fmt.Fprintf(os.Stdout, "estimated gas = %v\n", txCtx.Gas)
		}
		if cliCtx[i].DryRun {
			return nil, nil
		}
		
		var count = int64(0)
		for j := 0; j < i; j++ {
			if txCtx.AccountNumber == txSCtx[j].AccountNumber {
				count++
			}
		}
		txCtx.Sequence = txCtx.Sequence + count
		
		txBytes, err := txCtx.BuildAndSign(cliCtx[i].FromAddressName, passphrase[i], msgs)
		if err != nil {
			return nil, err
		}
		stdtx := auth.StdTx{}
		err = cliCtx[i].Codec.UnmarshalBinary(txBytes, &stdtx)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			stdTxs.Msgs = stdtx.Msgs
			stdTxs.Memo = stdtx.Memo
			stdTxs.Fee = stdtx.Fee
			stdTxs.Signatures = stdtx.Signatures
		} else {
			stdTxs.Signatures = append(stdTxs.Signatures, stdtx.Signatures...)
		}
	}
	txBytes, err := cliCtx[0].Codec.MarshalBinary(stdTxs)
	if err != nil {
		return nil, err
	}
	return cliCtx[0].EnsureBroadcastTxWithResponse(txBytes)
}

// KafkaConsumerMsgs : msgs to consume 5 second delay
func KafkaConsumerMsgs(cliCtx context.CLIContext, cdc *wire.Codec, kafkaState rest.KafkaState) {
	
	quit := make(chan bool)
	
	txSCtx := []authctx.TxContext{}
	passwords := []string{}
	cliSCtx := []context.CLIContext{}
	listMsgs := []sdk.Msg{}
	ticketIDs := []rest.Ticket{}
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				kafkaMsg := rest.KafkaTopicConsumer("Topic", kafkaState.Consumers, cdc)
				if kafkaMsg.Msg != nil {
					txCtx, cliCtxFromMsg, password := rest.KafkaTxAndKafkaCliFromKafkaMsg(kafkaMsg, cliCtx)
					txSCtx = append(txSCtx, txCtx)
					cliSCtx = append(cliSCtx, cliCtxFromMsg)
					passwords = append(passwords, password)
					listMsgs = append(listMsgs, kafkaMsg.Msg)
					ticketIDs = append(ticketIDs, kafkaMsg.TicketID)
				}
			}
		}
	}()
	
	time.Sleep(rest.SleepTimer)
	quit <- true
	if len(listMsgs) == 0 {
		return
	}
	
	HandleMultipleMsgOutput(txSCtx, cliSCtx, listMsgs, passwords, ticketIDs, kafkaState)
	
}

// HandleMultipleMsgOutput : handles output type
func HandleMultipleMsgOutput(txSCtx []authctx.TxContext, cliCtx []context.CLIContext, msgs []sdk.Msg, passphrase []string, ticketIDs []rest.Ticket, kafkaState rest.KafkaState) {
	output, err := SendTxSWithResponse(txSCtx, cliCtx, msgs, passphrase)
	if err != nil {
		for _, ticketID := range ticketIDs {
			rest.AddResponseToDB(ticketID, []byte(err.Error()), kafkaState.KafkaDB, cliCtx[0].Codec)
		}
		return
	}
	
	for i, ticketID := range ticketIDs {
		rest.AddResponseToDB(ticketID, ResponseBytesToJSON(output), kafkaState.KafkaDB, cliCtx[i].Codec)
	}
	return
}

// SimulateMsgs simulates the transaction and returns the gas estimate and the adjusted value.
func SimulateMsgs(txCtx authctx.TxContext, cliCtx context.CLIContext, name string, msgs []sdk.Msg, gas int64) (estimated, adjusted int64, err error) {
	txBytes, err := txCtx.WithGas(gas).BuildWithPubKey(name, msgs)
	if err != nil {
		return
	}
	estimated, adjusted, err = CalculateGas(cliCtx.Query, cliCtx.Codec, txBytes, cliCtx.GasAdjustment)
	return
}

// EnrichCtxWithGas calculates the gas estimate that would be consumed by the
// transaction and set the transaction's respective value accordingly.
func EnrichCtxWithGas(txCtx authctx.TxContext, cliCtx context.CLIContext, name string, msgs []sdk.Msg) (authctx.TxContext, error) {
	_, adjusted, err := SimulateMsgs(txCtx, cliCtx, name, msgs, 0)
	if err != nil {
		return txCtx, err
	}
	return txCtx.WithGas(adjusted), nil
}

// CalculateGas simulates the execution of a transaction and returns
// both the estimate obtained by the query and the adjusted amount.
func CalculateGas(queryFunc func(string, common.HexBytes) ([]byte, error), cdc *amino.Codec, txBytes []byte, adjustment float64) (estimate, adjusted int64, err error) {
	// run a simulation (via /app/simulate query) to
	// estimate gas and update TxContext accordingly
	rawRes, err := queryFunc("/app/simulate", txBytes)
	if err != nil {
		return
	}
	estimate, err = parseQueryResponse(cdc, rawRes)
	if err != nil {
		return
	}
	adjusted = adjustGasEstimate(estimate, adjustment)
	return
}

func adjustGasEstimate(estimate int64, adjustment float64) int64 {
	return int64(adjustment * float64(estimate))
}

func parseQueryResponse(cdc *amino.Codec, rawRes []byte) (int64, error) {
	var simulationResult sdk.Result
	if err := cdc.UnmarshalBinary(rawRes, &simulationResult); err != nil {
		return 0, err
	}
	return simulationResult.GasUsed, nil
}

func prepareTxContext(txCtx authctx.TxContext, cliCtx context.CLIContext) (authctx.TxContext, error) {
	if err := cliCtx.EnsureAccountExists(); err != nil {
		return txCtx, err
	}
	
	from, err := cliCtx.GetFromAddress()
	if err != nil {
		return txCtx, err
	}
	
	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txCtx.AccountNumber == 0 {
		accNum, err := cliCtx.GetAccountNumber(from)
		if err != nil {
			return txCtx, err
		}
		txCtx = txCtx.WithAccountNumber(accNum)
	}
	
	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if txCtx.Sequence == 0 {
		accSeq, err := cliCtx.GetAccountSequence(from)
		if err != nil {
			return txCtx, err
		}
		txCtx = txCtx.WithSequence(accSeq)
	}
	return txCtx, nil
}
