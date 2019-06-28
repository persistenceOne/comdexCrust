package cli

import (
	"os"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/utils"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	
	"github.com/comdex-blockchain/wire"
	
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	"github.com/comdex-blockchain/x/ibc"
)

// IBCRedeemFiatCmd : IBC redeem fiat command
func IBCRedeemFiatCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeemFiat",
		Short: "Redeems fiat from the redeemer to the issuer",
		
		RunE: func(cmd *cobra.Command, args []string) error {
			
			txCtx := context2.NewTxContextFromCLI().
				WithCodec(cdc)
			
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			
			issuerAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			
			dest := viper.GetString(flagTo)
			
			redeemerAddress, err := sdk.AccAddressFromBech32(dest)
			if err != nil {
				return err
			}
			
			sourceChain := viper.GetString(client.FlagChainID)
			destinationChain := viper.GetString(flagChain)
			
			amount := viper.GetInt64(flagAmount)
			
			pegHashStr := viper.GetString(flagPegHash)
			pegHashHex, err := sdk.GetFiatPegHashHex(pegHashStr)
			
			fiatPeg := sdk.BaseFiatPeg{
				PegHash:        pegHashHex,
				RedeemedAmount: amount,
			}
			
			msg := ibc.BuildRedeemFiatMsg(issuerAddress, redeemerAddress, amount, sdk.FiatPegWallet{fiatPeg}, sourceChain, destinationChain)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	
	cmd.Flags().String(flagTo, "", "Address to issue fiat to")
	cmd.Flags().String(flagChain, "", "Destination chain to send coins")
	cmd.Flags().String(flagPegHash, "", "Peg Hash")
	cmd.Flags().String(flagAmount, "", "amount")
	return cmd
}
