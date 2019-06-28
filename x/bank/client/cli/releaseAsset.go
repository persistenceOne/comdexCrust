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

// ReleaseAssetCmd : create a release asset tx and sign it with the give key
func ReleaseAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "releaseAsset",
		Short: "Releases an locked asset peg.",
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
			
			toStr := viper.GetString(flagTo)
			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}
			
			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)
			if err != nil {
				return err
			}
			
			msg := client.BuildReleaseAssetMsg(from, to, pegHashHex)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	return cmd
}
