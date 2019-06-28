package cli

import (
	"os"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/bank/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// BuyerExecuteOrderCmd : executes the exchange escrow transaction from the buyers side
func BuyerExecuteOrderCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buyerExecuteOrder",
		Short: "executes the exchange escrow transaction from the buyers side",
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
			
			buyerAddressString := viper.GetString(FlagBuyerAddress)
			buyerAddress, err := sdk.AccAddressFromBech32(buyerAddressString)
			if err != nil {
				return err
			}
			
			sellerAddressString := viper.GetString(FlagSellerAddress)
			sellerAddress, err := sdk.AccAddressFromBech32(sellerAddressString)
			if err != nil {
				return err
			}
			
			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)
			if err != nil {
				return err
			}
			
			fiatProofHashStr := viper.GetString(FlagFiatProofHash)
			
			msg := client.BuildBuyerExecuteOrderMsg(from, buyerAddress, sellerAddress, pegHashHex, fiatProofHashStr)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsBuyerAddress)
	cmd.Flags().AddFlagSet(fsSellerAddress)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsFiatProofHash)
	return cmd
}
