package cli

import (
	"fmt"
	
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/negotiation"
	"github.com/spf13/cobra"
)

// GetNegotiationCmd : command to get negotiation details
func GetNegotiationCmd(storeName string, cdc *wire.Codec, decoder sdk.NegotiationDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "negotiation [negotiationID]",
		Short: "Query negotiation details",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext()
			
			nego := args[0]
			
			negotiationID, err := sdk.GetNegotiationIDHex(nego)
			if err != nil {
				return err
			}
			
			res, err := cliCtx.QueryStore(negotiation.StoreKey([]byte(negotiationID)), storeName)
			if err != nil {
				return err
			}
			
			// Check if account was found
			if res == nil {
				return sdk.ErrUnknownAddress("No negotiation with negotiationID " + nego +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}
			
			// decode the value
			negotiationResponse, err := decoder(res)
			if err != nil {
				return err
			}
			
			// print out whole account
			output, err := wire.MarshalJSONIndent(cdc, negotiationResponse)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().AddFlagSet(fsNegotiationID)
	return cmd
}

// GetNegotiationDecoder : get negotiation decoder
func GetNegotiationDecoder(cdc *wire.Codec) sdk.NegotiationDecoder {
	return func(negotiationBytes []byte) (negotiation sdk.Negotiation, err error) {
		err = cdc.UnmarshalBinaryBare(negotiationBytes, &negotiation)
		if err != nil {
			panic(err)
		}
		return negotiation, err
	}
}
