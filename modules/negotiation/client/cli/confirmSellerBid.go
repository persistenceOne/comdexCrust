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

func ConfirmSellerBidCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "confirm-seller-bid",
		Short: "Confirm  a seller negotiation bid",
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
			sellerContractHash := viper.GetString(FlagSellerContractHash)
			negotiationID := types.NegotiationID(append(append(to.Bytes(), cliCtx.GetFromAddress().Bytes()...), pegHashHex...))

			kb, err := keys.NewKeyBaseFromHomeFlag()
			if err != nil {
				return err
			}

			passphrase, err := keys.GetPassphrase(cliCtx.GetFromName())
			if err != nil {
				return err
			}

			SignBytes := negotiationTypes.NewSignNegotiationBody(to, cliCtx.GetFromAddress(), pegHashHex, bid, time)
			signature, _, err := kb.Sign(cliCtx.GetFromName(), passphrase, SignBytes.GetSignBytes())
			if err != nil {
				return err
			}

			proposedNegotiation := &types.BaseNegotiation{
				NegotiationID:      negotiationID,
				BuyerAddress:       to,
				SellerAddress:      cliCtx.GetFromAddress(),
				PegHash:            pegHashHex,
				Bid:                bid,
				Time:               time,
				SellerContractHash: sellerContractHash,
				SellerSignature:    signature,
				BuyerSignature:     nil,
			}

			msg := negotiationTypes.BuildMsgConfirmSellerBid(proposedNegotiation)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}

	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	cmd.Flags().AddFlagSet(fsBid)
	cmd.Flags().AddFlagSet(fsTime)
	cmd.Flags().AddFlagSet(fsSellerContractHash)
	return cmd
}
