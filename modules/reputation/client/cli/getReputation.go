package cli

import (
	"fmt"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/reputation/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/reputation/internal/types"
)

func GetReputationCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "[address]",
		Short: "Query reputation details",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			cliCtx := context.NewCLIContext()
			addrStr := args[0]
			
			addr, err := cTypes.AccAddressFromBech32(addrStr)
			if err != nil {
				return err
			}
			res, _, err := cliCtx.QueryStore(keeper.AccountStoreKey(addr), types.ModuleName)
			if err != nil {
				return err
			}
			
			// Check if account was found
			if res == nil {
				return cTypes.ErrUnknownAddress("No reputation with address " + addr.String() +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}
			
			var reputation types.AccountReputation
			err = cdc.UnmarshalBinaryBare(res, &reputation)
			if err != nil {
				panic(err)
			}
			output, err := cdc.MarshalJSONIndent(reputation, "", " ")
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
	
	return cmd
}
