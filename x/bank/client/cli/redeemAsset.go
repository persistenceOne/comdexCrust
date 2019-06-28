package cli

import (
	"os"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/acl"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/bank/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagPegHash = "pegHash"
)

// RedeemAssetCmd : create a init redeemasset tx and sign it with the give key
func RedeemAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeemAsset",
		Short: "Redeem asset from redeemerAddress to the given issuerAddress",
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
			
			redeemerAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			
			issuerAddress, err := sdk.AccAddressFromBech32(viper.GetString(flagTo))
			if err != nil {
				return err
			}
			
			pegHash, err := sdk.GetAssetPegHashHex(viper.GetString(flagPegHash))
			if err != nil {
				return err
			}
			
			res, err := cliCtx.QueryStore(acl.AccountStoreKey(issuerAddress), "acl")
			if err != nil {
				return err
			}
			
			// the query will return empty if there is no data for this account
			if len(res) == 0 {
				return sdk.ErrInternal("Unauthorized transaction")
			}
			
			// decode the value
			decoder := acl.GetACLAccountDecoder(cdc)
			account, err := decoder(res)
			if err != nil {
				return err
			}
			if !account.GetACL().RedeemAsset {
				return sdk.ErrInternal("Unauthorized transaction")
			}
			
			msg := client.BuildRedeemAssetMsg(issuerAddress, redeemerAddress, pegHash)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	
	return cmd
}
