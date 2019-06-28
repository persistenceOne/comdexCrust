package cli

import (
	"github.com/comdex-blockchain/client/context"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/negotiation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"os"
	
	keybase "github.com/comdex-blockchain/client/keys"
	"github.com/comdex-blockchain/client/utils"
	"github.com/comdex-blockchain/crypto/keys"
)

// ConfirmSellerBidCmd : Change or create a negotiation bid
func ConfirmSellerBidCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirmSellerBid",
		Short: "Confirm  a seller negotiation bid",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			var kb keys.Keybase
			
			kb, err := keybase.GetKeyBase()
			if err != nil {
				return err
			}
			
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
			
			toStr := viper.GetString(FlagTo)
			
			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}
			
			bidInt64 := viper.GetInt64(FlagBid)
			timeInt64 := viper.GetInt64(FlagTime)
			pegHashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := sdk.GetAssetPegHashHex(pegHashStr)
			negotiationID := sdk.NegotiationID(append(append(to.Bytes(), from.Bytes()...), pegHashHex.Bytes()...))
			
			// passphrase, err := keybase.GetPassphrase(cliCtx.FromAddressName)
			// if err != nil {
			// 	return err
			// }
			
			SignBytes := negotiation.NewSignNegotiationBody(to, from, pegHashHex, bidInt64, timeInt64)
			signature, _, err := kb.Sign(cliCtx.FromAddressName, "1234567890", SignBytes.GetSignBytes())
			if err != nil {
				return err
			}
			proposedNegotiation := &sdk.BaseNegotiation{
				NegotiationID:   negotiationID,
				BuyerAddress:    to,
				SellerAddress:   from,
				PegHash:         pegHashHex,
				Bid:             bidInt64,
				Time:            timeInt64,
				SellerSignature: signature,
				BuyerSignature:  nil,
			}
			
			msg := negotiation.BuildMsgConfirmSellerBid(proposedNegotiation)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
			
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsBid)
	cmd.Flags().AddFlagSet(fsTime)
	return cmd
}
