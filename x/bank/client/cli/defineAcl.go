package cli

import (
	"os"
	"strconv"
	
	"github.com/comdex-blockchain/client/context"
	"github.com/comdex-blockchain/client/utils"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	authcmd "github.com/comdex-blockchain/x/auth/client/cli"
	context2 "github.com/comdex-blockchain/x/auth/client/context"
	"github.com/comdex-blockchain/x/bank"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const storeName = "acl"

// DefineACLCmd : assign Acl properties to accout from cli
func DefineACLCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defineACL",
		Short: "Assign Acl properties to address",
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := context2.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().WithAccountDecoder(authcmd.GetAccountDecoder(cdc)).WithCodec(cdc).WithLogger(os.Stdout)
			
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
			
			strOrganizationID := viper.GetString(FlagOrganizationID)
			organizationID, err := sdk.GetOrganizationIDFromString(strOrganizationID)
			if err != nil {
				return nil
			}
			
			strZoneID := viper.GetString(FlagZoneID)
			zoneID, err := sdk.GetZoneIDFromString(strZoneID)
			if err != nil {
				return nil
			}
			aclRequest := BuildACL()
			aclAccount := &sdk.BaseACLAccount{
				Address:        to,
				ZoneID:         zoneID,
				OrganizationID: organizationID,
				ACL:            aclRequest,
			}
			
			msg := bank.BuildMsgDefineACL(from, to, aclAccount)
			
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
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

// BuildACL : build and return the sdk.ACL object from string request
func BuildACL() sdk.ACL {
	var Request sdk.ACL
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
