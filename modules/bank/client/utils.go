package client

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/types"

	bankTypes "github.com/commitHub/commitBlockchain/modules/bank/internal/types"
)

func BuildMsg(from cTypes.AccAddress, to cTypes.AccAddress, coins cTypes.Coins) cTypes.Msg {
	msg := bankTypes.NewMsgSend(from, to, coins)
	return msg
}

func BuildIssueAssetMsg(from cTypes.AccAddress, to cTypes.AccAddress, assetPeg types.AssetPeg) cTypes.Msg {
	issueAsset := bankTypes.NewIssueAsset(from, to, assetPeg)
	msg := bankTypes.NewMsgBankIssueAssets([]bankTypes.IssueAsset{issueAsset})
	return msg
}

func BuildRedeemAssetMsg(issuerAddress cTypes.AccAddress, redeemerAddress cTypes.AccAddress, pegHash types.PegHash) cTypes.Msg {
	redeemAsset := bankTypes.NewRedeemAsset(issuerAddress, redeemerAddress, pegHash)
	msg := bankTypes.NewMsgBankRedeemAssets([]bankTypes.RedeemAsset{redeemAsset})
	return msg
}

func BuildIssueFiatMsg(from cTypes.AccAddress, to cTypes.AccAddress, fiatPeg types.FiatPeg) cTypes.Msg {

	issueFiat := bankTypes.NewIssueFiat(from, to, fiatPeg)
	msg := bankTypes.NewMsgBankIssueFiats([]bankTypes.IssueFiat{issueFiat})
	return msg
}

func BuildRedeemFiatMsg(from cTypes.AccAddress, to cTypes.AccAddress, amount int64) cTypes.Msg {

	redeemFiat := bankTypes.NewRedeemFiat(from, to, amount)
	msg := bankTypes.NewMsgBankRedeemFiats([]bankTypes.RedeemFiat{redeemFiat})
	return msg
}

func BuildSendAssetMsg(from cTypes.AccAddress, to cTypes.AccAddress, pegHash types.PegHash) cTypes.Msg {

	sendAsset := bankTypes.NewSendAsset(from, to, pegHash)
	msg := bankTypes.NewMsgBankSendAssets([]bankTypes.SendAsset{sendAsset})
	return msg
}

func BuildSendFiatMsg(from cTypes.AccAddress, to cTypes.AccAddress, pegHash types.PegHash, amount int64) cTypes.Msg {

	sendFiat := bankTypes.NewSendFiat(from, to, pegHash, amount)
	msg := bankTypes.NewMsgBankSendFiats([]bankTypes.SendFiat{sendFiat})
	return msg
}

func BuildBuyerExecuteOrderMsg(from cTypes.AccAddress, buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash, fiatProofHash string) cTypes.Msg {

	buyerExecuteOrder := bankTypes.NewBuyerExecuteOrder(from, buyerAddress, sellerAddress, pegHash, fiatProofHash)
	msg := bankTypes.NewMsgBankBuyerExecuteOrders([]bankTypes.BuyerExecuteOrder{buyerExecuteOrder})
	return msg
}

func BuildSellerExecuteOrderMsg(from cTypes.AccAddress, buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash, awbProofHash string) cTypes.Msg {

	sellerExecuteOrder := bankTypes.NewSellerExecuteOrder(from, buyerAddress, sellerAddress, pegHash, awbProofHash)
	msg := bankTypes.NewMsgBankSellerExecuteOrders([]bankTypes.SellerExecuteOrder{sellerExecuteOrder})
	return msg
}

func BuildReleaseAssetMsg(from cTypes.AccAddress, to cTypes.AccAddress, pegHash types.PegHash) cTypes.Msg {

	releaseAsset := bankTypes.NewReleaseAsset(from, to, pegHash)
	msg := bankTypes.NewMsgBankReleaseAssets([]bankTypes.ReleaseAsset{releaseAsset})
	return msg
}
