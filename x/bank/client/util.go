package client

import (
	sdk "github.com/commitHub/commitBlockchain/types"
	bank "github.com/commitHub/commitBlockchain/x/bank"
)

//BuildMsg : build the sendTx msg
func BuildMsg(from sdk.AccAddress, to sdk.AccAddress, coins sdk.Coins) sdk.Msg {
	input := bank.NewInput(from, coins)
	output := bank.NewOutput(to, coins)
	msg := bank.NewMsgSend([]bank.Input{input}, []bank.Output{output})
	return msg
}

//BuildIssueAssetMsg : build the issueAssetTx
func BuildIssueAssetMsg(from sdk.AccAddress, to sdk.AccAddress, assetPeg sdk.AssetPeg) sdk.Msg {
	issueAsset := bank.NewIssueAsset(from, to, assetPeg)
	msg := bank.NewMsgBankIssueAssets([]bank.IssueAsset{issueAsset})
	return msg
}

//BuildRedeemAssetMsg : build the RedeemAsset msg
func BuildRedeemAssetMsg(issuerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, pegHash sdk.PegHash) sdk.Msg {
	redeemAsset := bank.NewRedeemAsset(issuerAddress, redeemerAddress, pegHash)
	msg := bank.NewMsgBankRedeemAssets([]bank.RedeemAsset{redeemAsset})
	return msg
}

//BuildIssueFiatMsg : butild the issueFiatTx
func BuildIssueFiatMsg(from sdk.AccAddress, to sdk.AccAddress, fiatPeg sdk.FiatPeg) sdk.Msg {

	issueFiat := bank.NewIssueFiat(from, to, fiatPeg)
	msg := bank.NewMsgBankIssueFiats([]bank.IssueFiat{issueFiat})
	return msg
}

//BuildRedeemFiatMsg : build the redeemFiatTx
func BuildRedeemFiatMsg(from sdk.AccAddress, to sdk.AccAddress, amount int64) sdk.Msg {

	redeemFiat := bank.NewRedeemFiat(from, to, amount)
	msg := bank.NewMsgBankRedeemFiats([]bank.RedeemFiat{redeemFiat})
	return msg
}

//BuildSendAssetMsg : build the sendAssetTx
func BuildSendAssetMsg(from sdk.AccAddress, to sdk.AccAddress, pegHash sdk.PegHash) sdk.Msg {

	sendAsset := bank.NewSendAsset(from, to, pegHash)
	msg := bank.NewMsgBankSendAssets([]bank.SendAsset{sendAsset})
	return msg
}

//BuildSendFiatMsg : build the SendFiatMsg
func BuildSendFiatMsg(from sdk.AccAddress, to sdk.AccAddress, pegHash sdk.PegHash, amount int64) sdk.Msg {

	sendFiat := bank.NewSendFiat(from, to, pegHash, amount)
	msg := bank.NewMsgBankSendFiats([]bank.SendFiat{sendFiat})
	return msg
}

//BuildBuyerExecuteOrderMsg : build the BuyerExecuteOrderMsg
func BuildBuyerExecuteOrderMsg(from sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, fiatProofHash string) sdk.Msg {

	buyerExecuteOrder := bank.NewBuyerExecuteOrder(from, buyerAddress, sellerAddress, pegHash, fiatProofHash)
	msg := bank.NewMsgBankBuyerExecuteOrders([]bank.BuyerExecuteOrder{buyerExecuteOrder})
	return msg
}

//BuildSellerExecuteOrderMsg : butild the SellerExecuteOrderMsg
func BuildSellerExecuteOrderMsg(from sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, awbProofHash string) sdk.Msg {

	sellerExecuteOrder := bank.NewSellerExecuteOrder(from, buyerAddress, sellerAddress, pegHash, awbProofHash)
	msg := bank.NewMsgBankSellerExecuteOrders([]bank.SellerExecuteOrder{sellerExecuteOrder})
	return msg
}

//BuildReleaseAssetMsg : butild the releaseAssetMessage
func BuildReleaseAssetMsg(from sdk.AccAddress, to sdk.AccAddress, pegHash sdk.PegHash) sdk.Msg {

	releaseAsset := bank.NewReleaseAsset(from, to, pegHash)
	msg := bank.NewMsgBankReleaseAssets([]bank.ReleaseAsset{releaseAsset})
	return msg
}
