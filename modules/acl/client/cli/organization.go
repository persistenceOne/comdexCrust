package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/acl/internal/types"
)

// GetOrganizationCmd : returns a query account address of organization
func GetOrganizationCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "[organizationID]",
		Short: "Query address based on organizationID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext()

			strOrganizationID := args[0]
			organizationID, err := types.GetOrganizationIDFromString(strOrganizationID)
			if err != nil {
				return nil
			}
			res, _, err := ctx.QueryStore(types.GetOrganizationKey(organizationID), types.ModuleName)
			if err != nil {
				return err
			}

			if res == nil {
				return cTypes.ErrUnknownAddress("No account with organization id " + strOrganizationID +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			output, err := cdc.MarshalJSONIndent(cTypes.AccAddress(res), "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}
