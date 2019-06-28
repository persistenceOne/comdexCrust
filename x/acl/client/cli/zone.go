package cli

import (
	"fmt"
	
	"github.com/comdex-blockchain/x/acl"
	
	"github.com/spf13/cobra"
	
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// GetZoneCmd : returns a query account address of zone
func GetZoneCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "zone [zoneID]",
		Short: "Query address based on zone id  ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCLIContext()
			
			strZoneID := args[0]
			zoneID, err := sdk.GetZoneIDFromString(strZoneID)
			if err != nil {
				return nil
			}
			res, err := ctx.QueryStore(acl.ZoneStoreKey(zoneID), storeName)
			if err != nil {
				return err
			}
			
			if res == nil {
				return sdk.ErrUnknownAddress("No account with zone id " + strZoneID +
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
