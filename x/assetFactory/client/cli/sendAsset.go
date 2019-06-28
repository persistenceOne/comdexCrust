package cli

import (
	"os"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/assetFactory"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// SendAssetCmd : Send an asset to order with a buyer
func SendAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendAsset",
		Short: "Send an asset to order with a buyer",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			txCtx := context2.NewTxContextFromCLI().WithCodec(cdc)
			
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
			
			ownerAddressStr := viper.GetString(FlagOwnerAddress)
			
			ownerAddress, err := sdk.AccAddressFromBech32(ownerAddressStr)
			if err != nil {
				return nil
			}
			
			toStr := viper.GetString(FlagTo)
			
			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}
			
			pegHashStr := viper.GetString(FlagPegHash)
			pegHash, err := sdk.GetAssetPegHashHex(pegHashStr)
			
			msg := assetFactory.BuildSendAssetMsg(from, ownerAddress, to, pegHash)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsOwnerAddress)
	return cmd
}
