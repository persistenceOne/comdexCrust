package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/acl"
	"github.com/persistenceOne/comdexCrust/modules/auth"
	"github.com/persistenceOne/comdexCrust/modules/auth/client/utils"
	bankTypes "github.com/persistenceOne/comdexCrust/modules/bank/internal/types"
)

func DefineACLCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defineACL",
		Short: "Assign Acl properties to address",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)

			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			strOrganizationID := viper.GetString(FlagOrganizationID)
			organizationID, err := acl.GetOrganizationIDFromString(strOrganizationID)
			if err != nil {
				return nil
			}

			strZoneID := viper.GetString(FlagZoneID)
			zoneID, err := acl.GetZoneIDFromString(strZoneID)
			if err != nil {
				return nil
			}
			aclRequest := BuildACL()
			aclAccount := &acl.BaseACLAccount{
				Address:        to,
				ZoneID:         zoneID,
				OrganizationID: organizationID,
				ACL:            aclRequest,
			}

			msg := bankTypes.BuildMsgDefineACL(cliCtx.GetFromAddress(), to, aclAccount)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsOrganizationID)
	cmd.Flags().AddFlagSet(fsZoneID)
	cmd.Flags().AddFlagSet(fsIssueAsset)
	cmd.Flags().AddFlagSet(fsIssueFiat)
	cmd.Flags().AddFlagSet(fsSendAsset)
	cmd.Flags().AddFlagSet(fsSendFiat)
	cmd.Flags().AddFlagSet(fsBuyerExecuteOrder)
	cmd.Flags().AddFlagSet(fsSellerExecuteOrder)
	cmd.Flags().AddFlagSet(fsChangeBuyerBid)
	cmd.Flags().AddFlagSet(fsChangeSellerBid)
	cmd.Flags().AddFlagSet(fsConfirmBuyerBid)
	cmd.Flags().AddFlagSet(fsConfirmSellerBid)
	cmd.Flags().AddFlagSet(fsNegotiation)
	cmd.Flags().AddFlagSet(fsRedeemFiat)
	cmd.Flags().AddFlagSet(fsRedeemAsset)
	cmd.Flags().AddFlagSet(fsReleaseAsset)
	return cmd
}

func BuildACL() acl.ACL {
	var Request acl.ACL
	data, err := strconv.ParseBool(viper.GetString(FlagIssueAsset))
	if err == nil {
		Request.IssueAsset = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagIssueFiat))
	if err == nil {
		Request.IssueFiat = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagSendAsset))
	if err == nil {
		Request.SendAsset = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagSendFiat))
	if err == nil {
		Request.SendFiat = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagBuyerExecuteOrder))
	if err == nil {
		Request.BuyerExecuteOrder = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagSellerExecuteOrder))
	if err == nil {
		Request.SellerExecuteOrder = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagChangeBuyerBid))
	if err == nil {
		Request.ChangeBuyerBid = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagChangeSellerBid))
	if err == nil {
		Request.ChangeSellerBid = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagConfirmBuyerBid))
	if err == nil {
		Request.ConfirmBuyerBid = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagConfirmSellerBid))
	if err == nil {
		Request.ConfirmSellerBid = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagNegotiation))
	if err == nil {
		Request.Negotiation = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagRedeemFiat))
	if err == nil {
		Request.RedeemFiat = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagRedeemAsset))
	if err == nil {
		Request.RedeemAsset = data
	}
	data, err = strconv.ParseBool(viper.GetString(FlagReleaseAsset))
	if err == nil {
		Request.ReleaseAsset = data
	}
	return Request
}
