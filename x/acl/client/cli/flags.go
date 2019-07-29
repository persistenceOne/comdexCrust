package cli

import (
	flag "github.com/spf13/pflag"
)

//noLint
const (
	FlagTo                 = "to"
	FlagOrganizationID     = "organizationID"
	FlagZoneID             = "zoneID"
	FlagIssueAsset         = "issueAsset"
	FlagIssueFiat          = "issueFiat"
	FlagSendAsset          = "sendAsset"
	FlagSendFiat           = "sendFiat"
	FlagBuyerExecuteOrder  = "buyerExecuteOrder"
	FlagSellerExecuteOrder = "sellerExecuteOrder"
	FlagChangeBuyerBid     = "changeBuyerBid"
	FlagChangeSellerBid    = "changeSellerBid"
	FlagConfirmBuyerBid    = "confirmBuyerBid"
	FlagConfirmSellerBid   = "confirmSellerBid"
	FlagNegotiation        = "negotiation"
	FlagRedeemFiat         = "redeemFiat"
	FlagRedeemAsset        = "redeemAsset"
	FlagReleaseAsset       = "releaseAsset"
)

//onlint
var (
	fsTo                 = flag.NewFlagSet("", flag.ContinueOnError)
	fsOrganizationID     = flag.NewFlagSet("", flag.ContinueOnError)
	fsZoneID             = flag.NewFlagSet("", flag.ContinueOnError)
	fsIssueAsset         = flag.NewFlagSet("", flag.ContinueOnError)
	fsIssueFiat          = flag.NewFlagSet("", flag.ContinueOnError)
	fsSendAsset          = flag.NewFlagSet("", flag.ContinueOnError)
	fsSendFiat           = flag.NewFlagSet("", flag.ContinueOnError)
	fsBuyerExecuteOrder  = flag.NewFlagSet("", flag.ContinueOnError)
	fsSellerExecuteOrder = flag.NewFlagSet("", flag.ContinueOnError)
	fsChangeBuyerBid     = flag.NewFlagSet("", flag.ContinueOnError)
	fsChangeSellerBid    = flag.NewFlagSet("", flag.ContinueOnError)
	fsConfirmBuyerBid    = flag.NewFlagSet("", flag.ContinueOnError)
	fsConfirmSellerBid   = flag.NewFlagSet("", flag.ContinueOnError)
	fsNegotiation        = flag.NewFlagSet("", flag.ContinueOnError)
	fsRedeemFiat         = flag.NewFlagSet("", flag.ContinueOnError)
	fsRedeemAsset        = flag.NewFlagSet("", flag.ContinueOnError)
	fsReleaseAsset       = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsTo.String(FlagTo, "", "To account address")
	fsOrganizationID.String(FlagOrganizationID, "", "Organization Id")
	fsZoneID.String(FlagZoneID, "", "Zone id")
	fsIssueAsset.String(FlagIssueAsset, "", "Issue asset")
	fsIssueFiat.String(FlagIssueFiat, "", "Issue fiat")
	fsSendAsset.String(FlagSendAsset, "", "send Asset")
	fsSendFiat.String(FlagSendFiat, "", "Send fiat")
	fsBuyerExecuteOrder.String(FlagBuyerExecuteOrder, "", "buyer execute order")
	fsSellerExecuteOrder.String(FlagSellerExecuteOrder, "", "seller execute order")
	fsChangeBuyerBid.String(FlagChangeBuyerBid, "", "change buyer Id")
	fsChangeSellerBid.String(FlagChangeSellerBid, "", "change seller Id")
	fsConfirmBuyerBid.String(FlagConfirmBuyerBid, "", "Conform buyer Id")
	fsConfirmSellerBid.String(FlagConfirmSellerBid, "", "Conform seller Id")
	fsNegotiation.String(FlagNegotiation, "", "Negotiation")
	fsRedeemFiat.String(FlagRedeemFiat, "", "Redeem fiat")
	fsRedeemAsset.String(FlagRedeemAsset, "", "Redeem assets")
	fsReleaseAsset.String(FlagReleaseAsset, "", "Release assets")
}
