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

func SellerExecuteOrderCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sellerExecuteOrder",
		Short: "executes the exchange escrow transaction from the sellers side",
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

			awbProofHashStr := viper.GetString(FlagAWBProofHash)

			msg := client.BuildSellerExecuteOrderMsg(cliCtx.GetFromAddress(), buyerAddress,
				sellerAddress, pegHashHex, awbProofHashStr)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsBuyerAddress)
	cmd.Flags().AddFlagSet(fsSellerAddress)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsAWBProofHash)
	return cmd
}
