package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/bank/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTransactionID     = "transactionID"
	flagTransactionAmount = "transactionAmount"
)

//IssueFiatCmd : create a init fiat tx and sign it with the given key
func IssueFiatCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issueFiat",
		Short: "Initializes fiat with the given details and issues to the given address",
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

			transactionIDStr := viper.GetString(flagTransactionID)
			transactionAmountInt64 := viper.GetInt64(flagTransactionAmount)

			fiatPeg := &sdk.BaseFiatPeg{
				TransactionID:     transactionIDStr,
				TransactionAmount: transactionAmountInt64,
			}
			msg := client.BuildIssueFiatMsg(from, to, fiatPeg)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsTransactionID)
	cmd.Flags().AddFlagSet(fsTransactionAmount)
	return cmd
}
