package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/auth/client/utils"
	negotiationTypes "github.com/persistenceOne/persistenceSDK/modules/negotiation/internal/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

func ConfirmBuyerBidCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-buyer-bid",
		Short: "Confirm negotiation bid",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)
			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			bid := viper.GetInt64(FlagBid)
			time := viper.GetInt64(FlagTime)
			hashStr := viper.GetString(FlagPegHash)
			pegHashHex, err := types.GetAssetPegHashHex(hashStr)
			buyerContractHash := viper.GetString(FlagBuyerContractHash)
			negotiationID := types.NegotiationID(append(append(cliCtx.GetFromAddress().Bytes(), to.Bytes()...), pegHashHex...))

			kb, err := keys.NewKeyBaseFromHomeFlag()
			if err != nil {
				return err
			}

			passphrase, err := keys.GetPassphrase(cliCtx.GetFromName())
			if err != nil {
				return err
			}

			SignBytes := negotiationTypes.NewSignNegotiationBody(cliCtx.GetFromAddress(), to, pegHashHex, bid, time)
			signature, _, err := kb.Sign(cliCtx.GetFromName(), passphrase, SignBytes.GetSignBytes())
			if err != nil {
				return err
			}

			proposedNegotiation := &types.BaseNegotiation{
				NegotiationID:     negotiationID,
				BuyerAddress:      cliCtx.GetFromAddress(),
				SellerAddress:     to,
				PegHash:           pegHashHex,
				Bid:               bid,
				Time:              time,
				BuyerContractHash: buyerContractHash,
				BuyerSignature:    signature,
				SellerSignature:   nil,
			}

			msg := negotiationTypes.BuildMsgConfirmBuyerBid(proposedNegotiation)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsBid)
	cmd.Flags().AddFlagSet(fsTime)
	cmd.Flags().AddFlagSet(fsBuyerContractHash)
	return cmd
}
