package cli

import (
	"fmt"
	
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/reputation"
	"github.com/spf13/cobra"
)

// GetReputationCmd : command to get order details
func GetReputationCmd(storeName string, cdc *wire.Codec, decoder sdk.ReputationDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reputation [address]",
		Short: "Query reputation details",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext()
			addrStr := args[0]
			
			addr, err := sdk.AccAddressFromBech32(addrStr)
			if err != nil {
				return err
			}
			res, err := cliCtx.QueryStore(reputation.AccountStoreKey(addr), storeName)
			if err != nil {
				return err
			}
			
			// Check if account was found
			if res == nil {
				return sdk.ErrUnknownAddress("No reputation with address " + addr.String() +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}
			
			// decode the value
			ratings, err := decoder(res)
			if err != nil {
				return err
			}
			
			// print out whole account
			output, err := wire.MarshalJSONIndent(cdc, ratings)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	
	return cmd
}

// GetReputationDecoder : get order decoder
func GetReputationDecoder(cdc *wire.Codec) sdk.ReputationDecoder {
	return func(reputationBytes []byte) (reputation sdk.AccountReputation, err error) {
		err = cdc.UnmarshalBinaryBare(reputationBytes, &reputation)
		if err != nil {
			panic(err)
		}
		return reputation, err
	}
}
