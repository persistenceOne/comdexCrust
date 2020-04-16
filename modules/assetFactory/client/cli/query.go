package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/comdexCrust/modules/assetFactory/internal/keeper"
	assetFactoryTypes "github.com/persistenceOne/comdexCrust/modules/assetFactory/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

func QueryAssetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pegHash [pegHash]",
		Short: "Query asset peg",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			pegHash := args[0]

			ctx := context.NewCLIContext()
			pegHashHex, err := types.GetAssetPegHashHex(pegHash)
			if err != nil {
				return err
			}

			res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", assetFactoryTypes.QuerierRoute, keeper.QueryAsset, pegHashHex), nil)
			if err != nil {
				return err
			}

			if res == nil {
				return cTypes.ErrUnknownAddress("No asset with pegHash " + pegHash +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}
			fmt.Println(string(res))
			return nil
		},
	}
}
