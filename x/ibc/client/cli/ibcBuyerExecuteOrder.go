package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	wire "github.com/commitHub/commitBlockchain/wire"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/bank/client/cli"
	"github.com/commitHub/commitBlockchain/x/ibc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// IBCBuyerExecuteOrder :
func IBCBuyerExecuteOrder(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buyerExecuteOrder",
		Short: "ibc buyer execute order",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := context2.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			buyerAddressString := viper.GetString(cli.FlagBuyerAddress)
			buyerAddress, err := sdk.AccAddressFromBech32(buyerAddressString)
			if err != nil {
				return err
			}

			sellerAddressString := viper.GetString(cli.FlagSellerAddress)
			sellerAddress, err := sdk.AccAddressFromBech32(sellerAddressString)
			if err != nil {
				return err
			}
			pegHashStr := viper.GetString(flagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)

			fiatPeg := sdk.BaseFiatPeg{
				PegHash: pegHashHex,
			}
			if err != nil {
				return err
			}

			fiatProofHashStr := viper.GetString(cli.FlagFiatProofHash)
			destinationChainID := viper.GetString(flagChain)
			sourceChainID := viper.GetString(client.FlagChainID)

			msg := ibc.BuildBuyerExecuteOrder(from, buyerAddress, sellerAddress, pegHashHex, fiatProofHashStr, sdk.FiatPegWallet{fiatPeg}, sourceChainID, destinationChainID)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagChain, "", "destination chainID")
	cmd.Flags().String(cli.FlagFiatProofHash, "", "fiat proof hash")
	cmd.Flags().String(cli.FlagBuyerAddress, "", "buyerAddress")
	cmd.Flags().String(cli.FlagSellerAddress, "", "sellerAddress")
	cmd.Flags().String(flagPegHash, "", "pegHash")

	return cmd
}
