package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/assetFactory"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//RedeemAssetCmd : redeem asset tx and sign it with the give key
func RedeemAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem",
		Short: "Redeem assets with the given details and issues to the redeemed address",
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

			toStr := viper.GetString(FlagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			ownerAddressStr := viper.GetString(FlagOwnerAddress)

			ownerAddress, err := sdk.AccAddressFromBech32(ownerAddressStr)
			if err != nil {
				return err
			}

			pegHash, err := sdk.GetAssetPegHashHex(viper.GetString(FlagPegHash))
			if err != nil {
				return err
			}

			msg := assetFactory.BuildRedeemAssetMsg(from, ownerAddress, to, pegHash)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsOwnerAddress)
	return cmd
}
