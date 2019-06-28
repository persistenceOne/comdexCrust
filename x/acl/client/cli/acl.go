package cli

import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/acl"
)

// GetACLAccountCmd : returns a query account that will display the
func GetACLAccountCmd(storeName string, cdc *wire.Codec, decoder sdk.ACLAccountDecoder) *cobra.Command {
	return &cobra.Command{
		Use:   "acl [address]",
		Short: "Query acl account ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			// find the key to look up the account
			addr := args[0]
			
			key, err := sdk.AccAddressFromBech32(addr)
			if err != nil {
				return err
			}
			
			// perform query
			ctx := context.NewCLIContext()
			
			res, err := ctx.QueryStore(acl.AccountStoreKey(key), storeName)
			if err != nil {
				return err
			}
			
			// Check if account was found
			if res == nil {
				return sdk.ErrUnknownAddress("No account with address " + addr +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}
			
			// decode the value
			aclAccount, err := decoder(res)
			if err != nil {
				return err
			}
			// print out whole account
			output, err := wire.MarshalJSONIndent(cdc, aclAccount)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
}

// GetACLAccountDecoder : decode and return the ACLAccount
func GetACLAccountDecoder(cdc *wire.Codec) sdk.ACLAccountDecoder {
	return func(aclBytes []byte) (acl sdk.ACLAccount, err error) {
		// acct := new(auth.BaseAccount)
		err = cdc.UnmarshalBinaryBare(aclBytes, &acl)
		if err != nil {
			panic(err)
		}
		return acl, err
	}
}
