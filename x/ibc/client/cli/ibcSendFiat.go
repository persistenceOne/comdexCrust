package cli

import (
	"os"

	"github.com/commitHub/commitBlockchain/client"
	"github.com/commitHub/commitBlockchain/client/context"
	"github.com/commitHub/commitBlockchain/client/utils"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	authcmd "github.com/commitHub/commitBlockchain/x/auth/client/cli"
	context2 "github.com/commitHub/commitBlockchain/x/auth/client/context"
	"github.com/commitHub/commitBlockchain/x/ibc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//IBCSendFiatCmd : create a send fiat tx and sign it with the give key
func IBCSendFiatCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendFiat",
		Short: "Sends an fiat peg to an order transaction with a given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := context2.NewTxContextFromCLI().
				WithCodec(cdc)

			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			transactionAmountInt64 := viper.GetInt64(flagTransactionAmount)
			sourceChain := viper.GetString(client.FlagChainID)
			destinationChain := viper.GetString(flagChain)

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			amount := viper.GetInt64(flagAmount)

			pegHashStr := viper.GetString(flagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)

			fiatPeg := sdk.BaseFiatPeg{
				PegHash:           pegHashHex,
				TransactionAmount: transactionAmountInt64,
			}

			msg := ibc.BuildSendFiatMsg(from, to, pegHashHex, amount, sdk.FiatPegWallet{fiatPeg}, sourceChain, destinationChain)

			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagTo, "", "Address to issue fiat to")
	cmd.Flags().String(flagAmount, "", "Amount of coins to send")
	cmd.Flags().String(flagPegHash, "", "Peg Hash")
	cmd.Flags().String(flagChain, "", "destination chainID")
	return cmd
}
