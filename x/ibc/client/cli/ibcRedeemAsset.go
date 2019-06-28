package cli

import (
	"os"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/client/utils"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	
	"github.com/comdex-blockchain/wire"
	
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	"github.com/comdex-blockchain/x/ibc"
)

const (
	flagPegHash = "pegHash"
)

// IBCRedeemAssetCmd : IBC redeem asset command
func IBCRedeemAssetCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeemAsset",
		Short: "Redeems asset from the redeemer to the issuer",
		
		RunE: func(cmd *cobra.Command, args []string) error {
			
			txCtx := context2.NewTxContextFromCLI().
				WithCodec(cdc)
			
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			
			issuerAddress, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			
			dest := viper.GetString(flagTo)
			
			redeemerAddress, err := sdk.AccAddressFromBech32(dest)
			if err != nil {
				return err
			}
			
			pegHash, err := sdk.GetAssetPegHashHex(viper.GetString(flagPegHash))
			if err != nil {
				return err
			}
			
			sourceChain := viper.GetString(client.FlagChainID)
			destinationChain := viper.GetString(flagChain)
			
			msg := ibc.BuildRedeemAssetMsg(issuerAddress, redeemerAddress, pegHash, sourceChain, destinationChain)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}
	
	cmd.Flags().String(flagTo, "", "Address to issue asset to")
	cmd.Flags().String(flagPegHash, "", "Peg Hash")
	// cmd.Flags().String(flagChainID, "", "ID of Destination chain to send coins")
	cmd.Flags().String(flagChain, "", "Destination chain to send coins")
	return cmd
}
