package rest

import (
	"io/ioutil"
	"net/http"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	"github.com/comdex-blockchain/crypto/keys"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	authctx "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/ibc"
	
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application

type transferBody struct {
	Amount        sdk.Coins `json:"amount"`
	From          string    `json:"from" valid:"required~Enter the FromName"`
	Password      string    `json:"password" valid:"required~Enter the Password"`
	SourceChainID string    `json:"sourceChainID" valid:"required~Enter the SourceChainID, matches(^[a-zA-Z]+(-[A-Za-z]+)*$)~Invalid SourceChainID"`
	AccountNumber int64     `json:"accountNumber"`
	Sequence      int64     `json:"sequence"`
	Gas           int64     `json:"gas"`
	GasAdjustment string    `json:"gasAdjustment"`
}

// TransferRequestHandlerFn - http request handler to transfer coins to a address
// on a different chain via IBC
// nolint: gocyclo
func TransferRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		destChainID := vars["destchain"]
		bech32addr := vars["address"]
		
		to, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		var m transferBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		err = cdc.UnmarshalJSON(body, &m)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		_, err = govalidator.ValidateStruct(m)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		info, err := kb.Get(m.From)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		
		// build message
		packet := ibc.NewIBCPacket(sdk.AccAddress(info.GetPubKey().Address()), to, m.Amount, m.SourceChainID, destChainID)
		msg := ibc.IBCTransferMsg{IBCPacket: packet}
		
		txCtx := authctx.TxContext{
			Codec:         cdc,
			ChainID:       m.SourceChainID,
			AccountNumber: m.AccountNumber,
			Sequence:      m.Sequence,
			Gas:           m.Gas,
		}
		
		adjustment, ok := utils.ParseFloat64OrReturnBadRequest(w, m.GasAdjustment, client.DefaultGasAdjustment)
		if !ok {
			return
		}
		cliCtx = cliCtx.WithGasAdjustment(adjustment)
		
		if utils.HasDryRunArg(r) || m.Gas == 0 {
			newCtx, err := utils.EnrichCtxWithGas(txCtx, cliCtx, m.From, []sdk.Msg{msg})
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			if utils.HasDryRunArg(r) {
				utils.WriteSimulationResponse(w, txCtx.Gas)
				return
			}
			txCtx = newCtx
		}
		
		txBytes, err := txCtx.BuildAndSign(m.From, m.Password, []sdk.Msg{msg})
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}
		
		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		output, err := cdc.MarshalJSON(res)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		w.Write(utils.ResponseBytesToJSON(output))
	}
}
