package cli

import (
	"fmt"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/assetFactory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetAssetCmd :
func GetAssetCmd(storeName string, cdc *wire.Codec, decoder sdk.AssetPegDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asset [pegHash] [nodeURI]",
		Short: "Query asset from main chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			
			pegHash := args[0]
			nodeURI := viper.GetString(client.FlagNode)
			trustnode := viper.GetBool(client.FlagTrustNode)
			
			ctx := context.NewCLIContext()
			ctx = ctx.WithNodeURI(nodeURI)
			ctx = ctx.WithTrustNode(trustnode)
			
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHash)
			if err != nil {
				return err
			}
			
			res, err := ctx.QueryStore(assetFactory.AssetPegHashStoreKey(pegHashHex), storeName)
			if err != nil {
				return err
			}
			
			if res == nil {
				return sdk.ErrUnknownAddress("No asset with pegHash " + pegHash +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}
			
			assetPeg, err := decoder(res)
			if err != nil {
				return err
			}
			
			output, err := wire.MarshalJSONIndent(cdc, assetPeg)
			if err != nil {
				return err
			}
			
			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().String(flagPegHash, "", "pegHash to query")
	return cmd
}
