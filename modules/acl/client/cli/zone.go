package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/acl/internal/types"
)

// GetZoneCmd : returns a query account address of zone
func GetZoneCmd(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "[zoneID]",
		Short: "Query address based on zone id  ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext()

			strZoneID := args[0]
			zoneID, err := types.GetZoneIDFromString(strZoneID)
			if err != nil {
				return nil
			}
			res, _, err := ctx.QueryStore(types.GetZoneKey(zoneID), types.ModuleName)
			if err != nil {
				return err
			}

			if res == nil {
				return cTypes.ErrUnknownAddress("No account with zone id " + strZoneID +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			output, err := cdc.MarshalJSON(cTypes.AccAddress(res))
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
}
