package cli

import (
	"fmt"
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

// IssueAssetCmd : create a init asset tx and sign it with the given key
func IssueAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issueAsset",
		Short: "Initializes asset with the given details and issues to the given address",
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
			moderated := viper.GetBool(FlagModerated)
			var to sdk.AccAddress
			toStr := viper.GetString(flagTo)
			if !moderated && toStr == "" {
				return sdk.ErrInternal(fmt.Sprintf("must provide toAddress."))
			}
			if toStr == "" {
				to = from
			} else {
				to, err = sdk.AccAddressFromBech32(viper.GetString(flagTo))
				if err != nil {
					return err
				}
				if moderated {
					if to.String() != from.String() {
						return sdk.ErrInternal(fmt.Sprintf("Wrong toAddress."))
					}
				}
			}
			
			assetPeg := &sdk.BaseAssetPeg{
				AssetQuantity: viper.GetInt64(FlagAssetQuantity),
				AssetType:     viper.GetString(FlagAssetType),
				AssetPrice:    viper.GetInt64(FlagAssetPrice),
				DocumentHash:  viper.GetString(FlagDocumentHash),
				QuantityUnit:  viper.GetString(FlagQuantityUnit),
				Moderated:     moderated,
			}
			
			msg := client.BuildIssueAssetMsg(from, to, assetPeg)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsDocumentHash)
	cmd.Flags().AddFlagSet(fsAssetType)
	cmd.Flags().AddFlagSet(fsAssetPrice)
	cmd.Flags().AddFlagSet(fsAssetQuantity)
	cmd.Flags().AddFlagSet(fsQuantityUnit)
	cmd.Flags().AddFlagSet(fsModerated)
	return cmd
}
