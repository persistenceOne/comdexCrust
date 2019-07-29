package cli

import (
	"fmt"

	sdk "github.com/commitHub/commitBlockchain/types"

	"github.com/commitHub/commitBlockchain/client/context"
	wire "github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/assetFactory"
	"github.com/spf13/cobra"
)

//GetAssetCmd : command to get aeest details
func GetAssetCmd(storeName string, cdc *wire.Codec, decoder sdk.AssetPegDecoder) *cobra.Command {
	return &cobra.Command{
		Use:   "asset [pegHash]",
		Short: "Query asset details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// find the key to look up the account
			pegHash := args[0]

			// perform query
			ctx := context.NewCLIContext()

			pegHashHex, err := sdk.GetAssetPegHashHex(pegHash)
			if err != nil {
				return err
			}

			//fetch details from store
			res, err := ctx.QueryStore(assetFactory.AssetPegHashStoreKey(pegHashHex), storeName)
			if err != nil {
				return err
			}

			// Check if account was found
			if res == nil {
				return sdk.ErrUnknownAddress("No asset with pegHash " + pegHash +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			// decode the value
			assetPeg, err := decoder(res)
			if err != nil {
				return err
			}

			// print out whole account
			output, err := wire.MarshalJSONIndent(cdc, assetPeg)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
}

//GetAssetPegDecoder : get asset peg decoder
func GetAssetPegDecoder(cdc *wire.Codec) sdk.AssetPegDecoder {
	return func(assetBytes []byte) (asset sdk.AssetPeg, err error) {
		err = cdc.UnmarshalBinaryBare(assetBytes, &asset)
		if err != nil {
			panic(err)
		}
		return asset, err
	}
}
