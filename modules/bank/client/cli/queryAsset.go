package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/assetFactory"
	"github.com/commitHub/commitBlockchain/types"
)

func GetAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pegHash [pegHash]",
		Short: "Query asset from main chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			pegHash := args[0]
			nodeURI := viper.GetString(client.FlagNode)
			trustnode := viper.GetBool(client.FlagTrustNode)

			ctx := context.NewCLIContext()
			ctx = ctx.WithNodeURI(nodeURI)
			ctx = ctx.WithTrustNode(trustnode)

			pegHashHex, err := types.GetAssetPegHashHex(pegHash)
			if err != nil {
				return err
			}

			res, _, err := ctx.QueryStore(assetFactory.AssetPegHashStoreKey(pegHashHex), assetFactory.ModuleName)
			if err != nil {
				return err
			}

			if res == nil {
				return cTypes.ErrUnknownAddress("No asset with pegHash " + pegHash +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			var assetPeg types.AssetPeg
			err = ctx.Codec.UnmarshalBinaryLengthPrefixed(res, &assetPeg)
			if err != nil {
				return err
			}

			output, err := ctx.Codec.MarshalJSONIndent(assetPeg, "", " ")
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().String(FlagPegHash, "", "pegHash to query")
	return cmd
}
