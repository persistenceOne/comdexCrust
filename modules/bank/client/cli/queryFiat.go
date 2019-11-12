package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/fiatFactory"
	"github.com/commitHub/commitBlockchain/types"
)

func GetFiatCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "[peghash] [nodeURI]",
		Short: "Query fiat from main chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			pegHash := args[0]
			trustNode := viper.GetBool(client.FlagTrustNode)
			nodeURI := viper.GetString(client.FlagNode)

			ctx := context.NewCLIContext()
			ctx = ctx.WithNodeURI(nodeURI)
			ctx = ctx.WithTrustNode(trustNode)

			pegHashHex, err := types.GetFiatPegHashHex(pegHash)
			if err != nil {
				return nil
			}

			res, _, err := ctx.QueryStore(fiatFactory.FiatPegHashStoreKey(pegHashHex), fiatFactory.ModuleName)
			if err != nil {
				return nil
			}

			var fiatPeg types.FiatPeg
			err = ctx.Codec.UnmarshalBinaryLengthPrefixed(res, &fiatPeg)
			if err != nil {
				return err
			}

			output, err := ctx.Codec.MarshalJSONIndent(fiatPeg, "", " ")
			if err != nil {
				return nil
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().String(FlagPegHash, "", "PegHash to query")
	return cmd

}
