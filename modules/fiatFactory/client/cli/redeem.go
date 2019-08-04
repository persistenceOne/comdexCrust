package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"

	fiatFactoryTypes "github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
	"github.com/commitHub/commitBlockchain/types"
)

func RedeemFiatCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeem",
		Short: "Redeem fiat with the given details and redeem to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ownerAddressStr := viper.GetString(FlagOwnerAddress)

			ownerAddress, err := cTypes.AccAddressFromBech32(ownerAddressStr)
			if err != nil {
				return nil
			}

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetFiatPegHashHex(pegHashStr)

			amount := viper.GetInt64(FlagAmount)

			fiatPeg := types.BaseFiatPeg{
				PegHash:        pegHashHex,
				RedeemedAmount: amount,
			}

			msg := fiatFactoryTypes.BuildRedeemFiatMsg(cliCtx.GetFromAddress(), ownerAddress, amount, types.FiatPegWallet{fiatPeg})

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsOwnerAddress)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsAmount)

	return cmd
}
