package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/commitHub/commitBlockchain/types"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"
	"github.com/commitHub/commitBlockchain/modules/bank/client"
)

func BuyerExecuteOrderCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buyerExecuteOrder",
		Short: "executes the exchange escrow transaction from the buyers side",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			buyerAddressString := viper.GetString(FlagBuyerAddress)
			buyerAddress, err := cTypes.AccAddressFromBech32(buyerAddressString)
			if err != nil {
				return err
			}

			sellerAddressString := viper.GetString(FlagSellerAddress)
			sellerAddress, err := cTypes.AccAddressFromBech32(sellerAddressString)
			if err != nil {
				return err
			}

			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(pegHashStr)
			if err != nil {
				return err
			}

			fiatProofHashStr := viper.GetString(FlagFiatProofHash)

			msg := client.BuildBuyerExecuteOrderMsg(cliCtx.GetFromAddress(), buyerAddress, sellerAddress, pegHashHex, fiatProofHashStr)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsBuyerAddress)
	cmd.Flags().AddFlagSet(fsSellerAddress)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsFiatProofHash)
	return cmd
}
