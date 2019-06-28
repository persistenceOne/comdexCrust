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

const (
	flagRedeemAmount = "redeemAmount"
)

// RedeemFiatCmd : create a redeem fiat tx and sign it with the given key
func RedeemFiatCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeemFiat",
		Short: "Redeem fiat with the given details and returns some/total amount to the issuer address",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			txCtx := context2.NewTxContextFromCLI().
				WithCodec(cdc)
			
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			
			toStr := viper.GetString(flagTo)
			
			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}
			
			amount := viper.GetInt64(FlagAmount)
			
			msg := client.BuildRedeemFiatMsg(from, to, amount)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsAmount)
	return cmd
}
