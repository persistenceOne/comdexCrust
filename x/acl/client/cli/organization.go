package cli

import (
	"fmt"

	"github.com/commitHub/commitBlockchain/x/acl"

	"github.com/spf13/cobra"

	"github.com/commitHub/commitBlockchain/client/context"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
)

//GetOrganizationCmd : returns a query account address of organization
func GetOrganizationCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "organization [organizationID]",
		Short: "Query address based on organizationID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext()

			strOrganizationID := args[0]
			organizationID, err := sdk.GetOrganizationIDFromString(strOrganizationID)
			if err != nil {
				return nil
			}
			res, err := ctx.QueryStore(acl.OrganizationStoreKey(organizationID), storeName)
			if err != nil {
				return err
			}

			if res == nil {
				return sdk.ErrUnknownAddress("No account with organization id " + strOrganizationID +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			output, err := wire.MarshalJSONIndent(cdc, sdk.AccAddress(res))
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}
