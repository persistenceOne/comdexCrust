package cli

import (
	"fmt"
	
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/acl/internal/types"
)

// GetACLAccountCmd : returns a query account that will display the
func GetACLAccountCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "[address]",
		Short: "Query acl account ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			// find the key to look up the account
			addr := args[0]
			
			key, err := cTypes.AccAddressFromBech32(addr)
			if err != nil {
				return err
			}
			
			// perform query
			ctx := context.NewCLIContext()
			
			res, _, err := ctx.QueryStore(types.GetACLAccountKey(key), types.ModuleName)
			if err != nil {
				return err
			}
			
			// Check if account was found
			if res == nil {
				return cTypes.ErrUnknownAddress("No account with address" + addr +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}
			
			// decode the value
			var account types.ACLAccount
			err = cdc.UnmarshalBinaryLengthPrefixed(res, &account)
			if err != nil {
				return err
			}
			// print out whole account
			output, err := cdc.MarshalJSONIndent(account, "", " ")
			if err != nil {
				return err
			}
			
			fmt.Println(string(output))
			return nil
		},
	}
}
