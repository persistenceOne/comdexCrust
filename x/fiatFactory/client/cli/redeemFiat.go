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

//RedeemFiatCmd : create a init fiat tx and sign it with the given key
func RedeemFiatCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem",
		Short: "Redeem fiat with the given details and redeem to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := context2.NewTxContextFromCLI().
				WithCodec(cdc)

			cliCtx := context.NewCLIContext().
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc)).
				WithCodec(cdc).WithLogger(os.Stdout)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			ownerAddressStr := viper.GetString(FlagOwnerAddress)

			ownerAddress, err := sdk.AccAddressFromBech32(ownerAddressStr)
			if err != nil {
				return nil
			}

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := sdk.GetFiatPegHashHex(pegHashStr)

			amount := viper.GetInt64(FlagAmount)

			fiatPeg := sdk.BaseFiatPeg{
				PegHash:        pegHashHex,
				RedeemedAmount: amount,
			}

			msg := fiatFactory.BuildRedeemFiatMsg(from, ownerAddress, amount, sdk.FiatPegWallet{fiatPeg})

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsOwnerAddress)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsAmount)
	return cmd
}
