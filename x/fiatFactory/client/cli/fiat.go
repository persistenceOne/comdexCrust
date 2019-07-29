package cli

import (
	"fmt"

	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/x/fiatFactory"

	"github.com/commitHub/commitBlockchain/client/context"
	wire "github.com/commitHub/commitBlockchain/wire"
	"github.com/spf13/cobra"
)

//GetFiatCmd : command to get aeest details
func GetFiatCmd(storeName string, cdc *wire.Codec, decoder sdk.FiatPegDecoder) *cobra.Command {
	return &cobra.Command{
		Use:   "fiat [pegHash]",
		Short: "Query fiat transaction details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// find the key to look up the account
			pegHash := args[0]

			// perform query
			ctx := context.NewCLIContext()
			pegHashHex, err := sdk.GetFiatPegHashHex(pegHash)
			if err != nil {
				return err
			}
			res, err := ctx.QueryStore(fiatFactory.FiatPegHashStoreKey(pegHashHex), storeName)
			if err != nil {
				return err
			}

			// Check if account was found
			if res == nil {
				return sdk.ErrUnknownAddress("No fiat with pegHash " + pegHash +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			// decode the value
			fiatPeg, err := decoder(res)
			if err != nil {
				return err
			}

			// print out whole account
			output, err := wire.MarshalJSONIndent(cdc, fiatPeg)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
}

//GetFiatPegDecoder : get fiat peg decoder
func GetFiatPegDecoder(cdc *wire.Codec) sdk.FiatPegDecoder {
	return func(fiatBytes []byte) (fiat sdk.FiatPeg, err error) {
		err = cdc.UnmarshalBinaryBare(fiatBytes, &fiat)
		if err != nil {
			panic(err)
		}
		return fiat, err
	}
}
