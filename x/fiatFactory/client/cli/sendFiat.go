package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/fiatFactory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//SendFiatCmd : Send an fiat to order with a buyer
func SendFiatCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendFiat",
		Short: "Send an fiat to order with a buyer",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := context2.NewTxContextFromCLI().WithCodec(cdc)

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			ownerAddressStr := viper.GetString(FlagOwnerAddress)

			ownerAddress, err := sdk.AccAddressFromBech32(ownerAddressStr)
			if err != nil {
				return nil
			}

			toStr := viper.GetString(FlagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			assetPegHashStr := viper.GetString(FlagAssetPegHash)
			assetPegHash, err := sdk.GetAssetPegHashHex(assetPegHashStr)

			fiatPegHashStr := viper.GetString(FlagFiatPegHash)
			fiatPegHash, err := sdk.GetFiatPegHashHex(fiatPegHashStr)

			fiatPeg := sdk.BaseFiatPeg{
				PegHash:           fiatPegHash,
				TransactionAmount: viper.GetInt64(FlagAmount),
			}

			msg := fiatFactory.BuildSendFiatMsg(from, ownerAddress, to, assetPegHash, sdk.FiatPegWallet{fiatPeg})

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsAssetPegHash)
	cmd.Flags().AddFlagSet(fsFiatPegHash)
	cmd.Flags().AddFlagSet(fsAmount)
	cmd.Flags().AddFlagSet(fsOwnerAddress)
	return cmd
}
