package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/auth/client/utils"
	"github.com/persistenceOne/persistenceSDK/modules/bank/client"
	"github.com/persistenceOne/persistenceSDK/types"
)

func IssueFiatCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issueFiat",
		Short: "Initializes fiat with the given details and issues to the given address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)

			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			transactionIDStr := viper.GetString(FlagTransactionID)
			transactionAmountInt64 := viper.GetInt64(FlagTransactionAmount)

			fiatPeg := &types.BaseFiatPeg{
				TransactionID:     transactionIDStr,
				TransactionAmount: transactionAmountInt64,
			}
			msg := client.BuildIssueFiatMsg(cliCtx.GetFromAddress(), to, fiatPeg)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsTransactionID)
	cmd.Flags().AddFlagSet(fsTransactionAmount)
	return cmd
}
