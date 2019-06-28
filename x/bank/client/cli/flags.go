package cli

import (
	flag "github.com/spf13/pflag"
)

// noLint
const (
	FlagTo                 = "to"
	FlagAmount             = "amount"
	FlagDocumentHash       = "documentHash"
	FlagAssetType          = "assetType"
	FlagAssetPrice         = "assetPrice"
	FlagAssetQuantity      = "assetQuantity"
	FlagQuantityUnit       = "quantityUnit"
	FlagTransactionID      = "transactionID"
	FlagTransactionAmount  = "transactionAmount"
	FlagPegHash            = "pegHash"
	FlagBuyerAddress       = "buyerAddress"
	FlagSellerAddress      = "sellerAddress"
	FlagFiatProofHash      = "fiatProofHash"
	FlagAWBProofHash       = "awbProofHash"
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
	FlagRedeemAsset        = "redeemAsset"
	FlagRedeemFiat         = "redeemFiat"
	FlagReleaseAsset       = "releaseAsset"
	FlagModerated          = "moderated"
)

var (
	fsTo                 = flag.NewFlagSet("", flag.ContinueOnError)
	fsAmount             = flag.NewFlagSet("", flag.ContinueOnError)
	fsDocumentHash       = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetType          = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetPrice         = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetQuantity      = flag.NewFlagSet("", flag.ContinueOnError)
	fsQuantityUnit       = flag.NewFlagSet("", flag.ContinueOnError)
	fsTransactionID      = flag.NewFlagSet("", flag.ContinueOnError)
	fsTransactionAmount  = flag.NewFlagSet("", flag.ContinueOnError)
	fsPegHash            = flag.NewFlagSet("", flag.ContinueOnError)
	fsBuyerAddress       = flag.NewFlagSet("", flag.ContinueOnError)
	fsSellerAddress      = flag.NewFlagSet("", flag.ContinueOnError)
	fsFiatProofHash      = flag.NewFlagSet("", flag.ContinueOnError)
	fsAWBProofHash       = flag.NewFlagSet("", flag.ContinueOnError)
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
	fsRedeemAsset        = flag.NewFlagSet("", flag.ContinueOnError)
	fsRedeemFiat         = flag.NewFlagSet("", flag.ContinueOnError)
	fsReleaseAsset       = flag.NewFlagSet("", flag.ContinueOnError)
	fsModerated          = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsTo.String(FlagTo, "", "Address to send coins")
	fsAmount.String(FlagAmount, "", "Amount of coins to send")
	fsDocumentHash.String(FlagDocumentHash, "", "Hash of the asset doccuments of the asset")
	fsAssetType.String(FlagAssetType, "", "Type of the asset")
	fsAssetPrice.String(FlagAssetPrice, "", "Price of the asset")
	fsAssetQuantity.String(FlagAssetQuantity, "", "Quantity of the assent in integer")
	fsQuantityUnit.String(FlagQuantityUnit, "", "The unit of the qunatity")
	fsTransactionID.String(FlagTransactionID, "", "Fiat deposit transaction ID")
	fsTransactionAmount.String(FlagTransactionAmount, "", "Fiat deposit transaction amount")
	fsPegHash.String(FlagPegHash, "", "Peg Hash to be transferred ")
	fsBuyerAddress.String(FlagBuyerAddress, "", "Buyer's Address ")
	fsSellerAddress.String(FlagSellerAddress, "", "Seller's Address")
	fsFiatProofHash.String(FlagFiatProofHash, "", "Proof of fiat hash")
	fsAWBProofHash.String(FlagAWBProofHash, "", "Proof of awb hash")
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
	fsRedeemAsset.String(FlagRedeemAsset, "", "Redeem assets")
	fsRedeemFiat.String(FlagRedeemFiat, "", "Redeem fiats")
	fsReleaseAsset.String(FlagReleaseAsset, "", "Release assets")
	fsModerated.Bool(FlagModerated, false, "private")
}
