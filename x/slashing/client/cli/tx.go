package cli

import (
	"os"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	authctx "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/slashing"
	
	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail",
		Args:  cobra.NoArgs,
		Short: "unjail validator previously jailed for downtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			
			valAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			
			msg := slashing.NewMsgUnjail(sdk.ValAddress(valAddr))
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	
	return cmd
}
