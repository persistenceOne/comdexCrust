package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	
	"github.com/comdex-blockchain/wire"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	
	"os"
	
	"github.com/comdex-blockchain/client/utils"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	"github.com/comdex-blockchain/x/ibc"
)

const (
	flagTransactionID     = "transactionID"
	flagTransactionAmount = "transactionAmount"
)

// IBCIssueFiatCmd : IBC issue fiat command
func IBCIssueFiatCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issueFiat",
		Short: "Initializes fiat with the given details and issues to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			txCtx := context2.NewTxContextFromCLI().WithCodec(cdc)
			
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			
			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			
			dest := viper.GetString(flagTo)
			
			to, err := sdk.AccAddressFromBech32(dest)
			if err != nil {
				return err
			}
			
			transactionIDStr := viper.GetString(flagTransactionID)
			transactionAmountInt64 := viper.GetInt64(flagTransactionAmount)
			sourceChain := viper.GetString(client.FlagChainID)
			destinationChain := viper.GetString(flagChain)
			
			fiatPeg := sdk.BaseFiatPeg{
				TransactionID:     transactionIDStr,
				TransactionAmount: transactionAmountInt64,
			}
			
			msg := ibc.BuildIssueFiatMsg(from, to, &fiatPeg, sourceChain, destinationChain)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	
	cmd.Flags().String(flagTo, "", "Address to issue asset to")
	cmd.Flags().String(flagChain, "", "Destination chain to send coins")
	// cmd.Flags().String(flagChainID, "", "ID of Destination chain to send coins")
	cmd.Flags().String(flagTransactionID, "", "Fiat deposit transaction ID")
	cmd.Flags().String(flagTransactionAmount, "", "Fiat deposit transaction amount")
	return cmd
}
