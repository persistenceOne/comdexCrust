package cli

import (
	"os"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/fiatFactory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// IssueFiatCmd : create a init fiat tx and sign it with the given key
func IssueFiatCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Initializes fiat with the given details and issues to the given address",
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
			
			toStr := viper.GetString(FlagTo)
			
			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}
			
			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := sdk.GetFiatPegHashHex(pegHashStr)
			transactionIDStr := viper.GetString(FlagTransactionID)
			transactionAmountInt64 := viper.GetInt64(FlagTransactionAmount)
			
			fiatPeg := sdk.BaseFiatPeg{
				PegHash:           pegHashHex,
				TransactionID:     transactionIDStr,
				TransactionAmount: transactionAmountInt64,
			}
			fiatPegI := sdk.ToFiatPeg(fiatPeg)
			
			msg := fiatFactory.BuildIssueFiatMsg(from, to, fiatPegI)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsTransactionID)
	cmd.Flags().AddFlagSet(fsTransactionAmount)
	return cmd
}
