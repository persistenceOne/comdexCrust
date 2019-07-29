package cli

import (
	"fmt"

	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/context"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/fiatFactory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetFiatCmd :
func GetFiatCmd(storeName string, cdc *wire.Codec, decoder sdk.FiatPegDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fiat [peghash] [nodeURI]",
		Short: "Query fiat from main chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			pegHash := args[0]
			trustNode := viper.GetBool(client.FlagTrustNode)
			nodeURI := viper.GetString(client.FlagNode)

			ctx := context.NewCLIContext()
			ctx = ctx.WithNodeURI(nodeURI)
			ctx = ctx.WithTrustNode(trustNode)

			pegHashHex, err := sdk.GetFiatPegHashHex(pegHash)
			if err != nil {
				return nil
			}

			res, err := ctx.QueryStore(fiatFactory.FiatPegHashStoreKey(pegHashHex), storeName)
			if err != nil {
				return nil
			}

			fiatPeg, err := decoder(res)
			if err != nil {
				return nil
			}

			output, err := wire.MarshalJSONIndent(cdc, fiatPeg)
			if err != nil {
				return nil
			}

			fmt.Println(string(output))
			return nil
		},
	}
	cmd.Flags().String(flagPegHash, "", "PegHash to query")
	return cmd

}
