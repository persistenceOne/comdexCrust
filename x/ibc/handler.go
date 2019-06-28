package ibc

import (
	"reflect"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/x/acl"
	"github.com/comdex-blockchain/x/assetFactory"
	"github.com/comdex-blockchain/x/bank"
	"github.com/comdex-blockchain/x/fiatFactory"
	"github.com/comdex-blockchain/x/negotiation"
	"github.com/comdex-blockchain/x/order"
	"github.com/comdex-blockchain/x/reputation"
)

// NewHandler : handles IBC bank msgs
func NewHandler(ibcm Mapper, ck bank.Keeper, ak acl.Keeper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, reputationKeeper reputation.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case IBCTransferMsg:
			return handleIBCTransferMsg(ctx, ibcm, ck, msg)
		case IBCReceiveMsg:
			return handleIBCReceiveMsg(ctx, ibcm, ck, msg)
		case MsgIssueAssets:
			return handleIBCHubIssueAssetMsg(ctx, ibcm, ck, msg, ak)
		case MsgRedeemAssets:
			return handleIBCHubRedeemAssetMsg(ctx, ibcm, ck, msg, ak)
		case MsgIssueFiats:
			return handleIBCHubIssueFiatMsg(ctx, ibcm, ck, msg, ak)
		case MsgRedeemFiats:
			return handleIBCHubRedeemFiatMsg(ctx, ibcm, ck, msg, ak)
		case MsgSendAssets:
			return handleIBCHubSendAssetMsg(ctx, ibcm, ck, msg, ak, negotiationKeeper, orderKeeper, reputationKeeper)
		case MsgSendFiats:
			return handleIBCHubSendFiatMsg(ctx, ibcm, ck, msg, ak, negotiationKeeper, orderKeeper, reputationKeeper)
		case MsgBuyerExecuteOrders:
			return handleIBCHubBuyerExecuteOrdersMsg(ctx, ibcm, ck, msg, ak, negotiationKeeper, orderKeeper, reputationKeeper)
		case MsgSellerExecuteOrders:
			return handleIBCHubSellerExecuteOrdersMsg(ctx, ibcm, ck, msg, ak, negotiationKeeper, orderKeeper, reputationKeeper)
		
		default:
			errMsg := "Unrecognized IBC Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// NewAssetHandler : handles IBC asset msgs
func NewAssetHandler(ibcm Mapper, ck bank.Keeper, ak assetFactory.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case IBCTransferMsg:
			return handleIBCTransferMsg(ctx, ibcm, ck, msg)
		case IBCReceiveMsg:
			return handleIBCReceiveMsg(ctx, ibcm, ck, msg)
		case MsgRelayIssueAssets:
			return handleIBCRelayIssueAssetMsg(ctx, ibcm, ak, msg)
		case MsgRelayRedeemAssets:
			return handleIBCRelayRedeemAssetMsg(ctx, ibcm, ak, msg)
		case MsgRelaySendAssets:
			return handleIBCRelaySendAssetMsg(ctx, ibcm, ak, msg)
		case MsgRelaySellerExecuteOrders:
			return handlerIBCRelaySellerExecuteOrdersMsg(ctx, ibcm, ak, msg)
		default:
			errMsg := "Unrecognized IBC Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// NewFiatHandler : handles IBC fiat msgs
func NewFiatHandler(ibcm Mapper, ck bank.Keeper, fk fiatFactory.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case IBCTransferMsg:
			return handleIBCTransferMsg(ctx, ibcm, ck, msg)
		case IBCReceiveMsg:
			return handleIBCReceiveMsg(ctx, ibcm, ck, msg)
		case MsgRelayIssueFiats:
			return handleIBCRelayIssueFiatMsg(ctx, ibcm, fk, msg)
		case MsgRelayRedeemFiats:
			return handleIBCRelayRedeemFiatMsg(ctx, ibcm, fk, msg)
		case MsgRelaySendFiats:
			return handleIBCRelaySendFiatMsg(ctx, ibcm, fk, msg)
		case MsgRelayBuyerExecuteOrders:
			return handleIBCRelayBuyerExecuteOrderMsg(ctx, ibcm, fk, msg)
		default:
			errMsg := "Unrecognized IBC Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// IBCTransferMsg deducts coins from the account and creates an egress IBC packet.
func handleIBCTransferMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg IBCTransferMsg) sdk.Result {
	packet := msg
	
	_, _, err := ck.SubtractCoins(ctx, packet.SrcAddr, packet.Coins)
	if err != nil {
		return err.Result()
	}
	
	err = ibcm.PostIBCPacket(ctx, packet)
	if err != nil {
		return err.Result()
	}
	
	return sdk.Result{}
}

// IBCReceiveMsg adds coins to the destination address and creates an ingress IBC packet.
func handleIBCReceiveMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg IBCReceiveMsg) sdk.Result {
	packet := msg.IBCPacket
	
	seq := ibcm.GetIngressSequence(ctx, packet.SrcChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	
	_, _, err := ck.AddCoins(ctx, packet.DestAddr, packet.Coins)
	if err != nil {
		return err.Result()
	}
	
	ibcm.SetIngressSequence(ctx, packet.SrcChain, seq+1)
	
	return sdk.Result{}
}

// *****comdex
func handleIBCHubIssueAssetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg MsgIssueAssets, ak acl.Keeper) sdk.Result {
	tags, err, issuedAssetPegs := ck.IssueAssetsToWallets(ctx, (toHubMsgIssueAssets(msg)).IssueAssets, ak)
	if err != nil {
		return err.Result()
	}
	
	for i := range msg.IssueAssets {
		msg.IssueAssets[i].AssetPeg.SetPegHash(issuedAssetPegs[i].GetPegHash())
	}
	
	err = ibcm.PostIBCMsgIssueAssetsPacket(ctx, msg)
	if err != nil {
		return err.Result()
	}
	
	return sdk.Result{
		Tags: tags,
	}
}

func handleIBCRelayIssueAssetMsg(ctx sdk.Context, ibcm Mapper, ak assetFactory.Keeper, msg MsgRelayIssueAssets) sdk.Result {
	
	seq := ibcm.GetIngressSequence(ctx, msg.IssueAssets[0].SourceChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	
	var typesIssueAsset []assetFactory.IssueAsset
	for _, issueAsset := range msg.IssueAssets {
		typesIssueAsset = append(typesIssueAsset, assetFactory.NewIssueAsset(issueAsset.IssuerAddress, issueAsset.ToAddress, issueAsset.AssetPeg))
	}
	_, err := ak.IssueAssetsToWallets(ctx, typesIssueAsset)
	if err != nil {
		return err.Result()
	}
	
	ibcm.SetIngressSequence(ctx, msg.IssueAssets[0].SourceChain, seq+1)
	
	return sdk.Result{}
}

func handleIBCHubRedeemAssetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg MsgRedeemAssets, ak acl.Keeper) sdk.Result {
	tags, err, _ := ck.RedeemAssetsFromWallets(ctx, (toHubMsgRedeemAssets(msg)).RedeemAssets, ak)
	if err != nil {
		return err.Result()
	}
	
	err = ibcm.PostIBCMsgRedeemAssetsPacket(ctx, msg)
	if err != nil {
		return err.Result()
	}
	
	return sdk.Result{
		Tags: tags,
	}
}

func handleIBCRelayRedeemAssetMsg(ctx sdk.Context, ibcm Mapper, ak assetFactory.Keeper, msg MsgRelayRedeemAssets) sdk.Result {
	
	seq := ibcm.GetIngressSequence(ctx, msg.RedeemAssets[0].SourceChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	
	var typesRedeemAsset []assetFactory.RedeemAsset
	for _, redeemAsset := range msg.RedeemAssets {
		typesRedeemAsset = append(typesRedeemAsset, assetFactory.NewRedeemAsset(msg.Relayer, redeemAsset.RedeemerAddress, redeemAsset.IssuerAddress, redeemAsset.PegHash))
	}
	_, err := ak.RedeemAssetsToWallets(ctx, typesRedeemAsset)
	if err != nil {
		return err.Result()
	}
	
	ibcm.SetIngressSequence(ctx, msg.RedeemAssets[0].SourceChain, seq+1)
	
	return sdk.Result{}
}

func handleIBCHubSendAssetMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg MsgSendAssets, ak acl.Keeper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err, sentAssetPegs := ck.SendAssetsToWallets(ctx, orderKeeper, negotiationKeeper, (toHubMsgSendAssets(msg)).SendAssets, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	
	for i := range msg.SendAssets {
		msg.SendAssets[i].PegHash = sentAssetPegs[i].GetPegHash()
	}
	
	err = ibcm.PostIBCMsgSendAssetsPacket(ctx, msg)
	if err != nil {
		return err.Result()
	}
	
	return sdk.Result{
		Tags: tags,
	}
}

func handleIBCRelaySendAssetMsg(ctx sdk.Context, ibcm Mapper, ak assetFactory.Keeper, msg MsgRelaySendAssets) sdk.Result {
	
	seq := ibcm.GetIngressSequence(ctx, msg.SendAssets[0].SourceChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	
	var typesSendAsset []assetFactory.SendAsset
	for _, SendAsset := range msg.SendAssets {
		typesSendAsset = append(typesSendAsset, assetFactory.NewSendAsset(msg.Relayer, SendAsset.FromAddress, SendAsset.ToAddress, SendAsset.PegHash))
	}
	_, err := ak.SendAssetsToOrders(ctx, typesSendAsset)
	if err != nil {
		return err.Result()
	}
	
	ibcm.SetIngressSequence(ctx, msg.SendAssets[0].SourceChain, seq+1)
	
	return sdk.Result{}
}

func handleIBCHubIssueFiatMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg MsgIssueFiats, ak acl.Keeper) sdk.Result {
	tags, err, issuedFiatPegs := ck.IssueFiatsToWallets(ctx, (toHubMsgIssueFiats(msg)).IssueFiats, ak)
	if err != nil {
		return err.Result()
	}
	
	for i := range msg.IssueFiats {
		msg.IssueFiats[i].FiatPeg.SetPegHash(issuedFiatPegs[i].GetPegHash())
	}
	
	err = ibcm.PostIBCMsgIssueFiatsPacket(ctx, msg)
	if err != nil {
		return err.Result()
	}
	
	return sdk.Result{
		Tags: tags,
	}
}

func handleIBCRelayIssueFiatMsg(ctx sdk.Context, ibcm Mapper, ak fiatFactory.Keeper, msg MsgRelayIssueFiats) sdk.Result {
	
	seq := ibcm.GetIngressSequence(ctx, msg.IssueFiats[0].SourceChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	
	var typesIssueFiat []fiatFactory.IssueFiat
	for _, issueFiat := range msg.IssueFiats {
		typesIssueFiat = append(typesIssueFiat, fiatFactory.NewIssueFiat(issueFiat.IssuerAddress, issueFiat.ToAddress, issueFiat.FiatPeg))
	}
	_, err := ak.IssueFiatsToWallets(ctx, typesIssueFiat)
	if err != nil {
		return err.Result()
	}
	
	ibcm.SetIngressSequence(ctx, msg.IssueFiats[0].SourceChain, seq+1)
	
	return sdk.Result{}
}

func handleIBCHubRedeemFiatMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg MsgRedeemFiats, ak acl.Keeper) sdk.Result {
	tags, err, redeemerFiatPegWallets := ck.RedeemFiatsFromWallets(ctx, (toHubMsgRedeemFiats(msg)).RedeemFiats, ak)
	if err != nil {
		return err.Result()
	}
	
	for i := range msg.RedeemFiats {
		msg.RedeemFiats[i].FiatPegWallet = redeemerFiatPegWallets[i]
	}
	
	err = ibcm.PostIBCMsgRedeemFiatsPacket(ctx, msg)
	if err != nil {
		return err.Result()
	}
	
	return sdk.Result{
		Tags: tags,
	}
}

func handleIBCRelayRedeemFiatMsg(ctx sdk.Context, ibcm Mapper, ak fiatFactory.Keeper, msg MsgRelayRedeemFiats) sdk.Result {
	
	seq := ibcm.GetIngressSequence(ctx, msg.RedeemFiats[0].SourceChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	
	var typesRedeemFiat []fiatFactory.RedeemFiat
	for _, redeemFiat := range msg.RedeemFiats {
		typesRedeemFiat = append(typesRedeemFiat, fiatFactory.NewRedeemFiat(msg.Relayer, redeemFiat.RedeemerAddress, redeemFiat.Amount, redeemFiat.FiatPegWallet))
	}
	_, err := ak.RedeemFiatsFromWallets(ctx, typesRedeemFiat)
	if err != nil {
		return err.Result()
	}
	
	ibcm.SetIngressSequence(ctx, msg.RedeemFiats[0].SourceChain, seq+1)
	
	return sdk.Result{}
}

func handleIBCHubSendFiatMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg MsgSendFiats, ak acl.Keeper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err, sentFiatPegWallets := ck.SendFiatsToWallets(ctx, orderKeeper, negotiationKeeper, (toHubMsgSendFiats(msg)).SendFiats, ak, reputationKeeper)
	if err != nil {
		return err.Result()
	}
	
	for i := range msg.SendFiats {
		msg.SendFiats[i].FiatPegWallet = sentFiatPegWallets[i]
	}
	
	err = ibcm.PostIBCMsgSendFiatsPacket(ctx, msg)
	if err != nil {
		return err.Result()
	}
	
	return sdk.Result{
		Tags: tags,
	}
}

func handleIBCRelaySendFiatMsg(ctx sdk.Context, ibcm Mapper, ak fiatFactory.Keeper, msg MsgRelaySendFiats) sdk.Result {
	
	seq := ibcm.GetIngressSequence(ctx, msg.SendFiats[0].SourceChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	
	var typesSendFiat []fiatFactory.SendFiat
	for _, SendFiat := range msg.SendFiats {
		typesSendFiat = append(typesSendFiat, fiatFactory.NewSendFiat(msg.Relayer, SendFiat.FromAddress, SendFiat.ToAddress, SendFiat.PegHash, SendFiat.FiatPegWallet))
	}
	_, err := ak.SendFiatsToOrders(ctx, typesSendFiat)
	if err != nil {
		return err.Result()
	}
	ibcm.SetIngressSequence(ctx, msg.SendFiats[0].SourceChain, seq+1)
	
	return sdk.Result{}
}

func handleIBCHubBuyerExecuteOrdersMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg MsgBuyerExecuteOrders, ak acl.Keeper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err, fiatPegWallets := ck.BuyerExecuteTradeOrders(ctx, negotiationKeeper, orderKeeper, (toHubMsgBuyerExecuteOrdermsg(msg)).BuyerExecuteOrders, ak, reputationKeeper) // add toHub
	if err != nil {
		return err.Result()
	}
	for i := range msg.BuyerExecuteOrders {
		msg.BuyerExecuteOrders[i].FiatPegWallet = fiatPegWallets[i]
	}
	err = ibcm.PostIBCMsgBuyerExecuteOrdersPacket(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handleIBCRelayBuyerExecuteOrderMsg(ctx sdk.Context, ibcm Mapper, fk fiatFactory.Keeper, msg MsgRelayBuyerExecuteOrders) sdk.Result {
	
	seq := ibcm.GetIngressSequence(ctx, msg.BuyerExecuteOrders[0].SourceChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	var typeBuyerExecuteOrders []fiatFactory.SendFiat
	for _, in := range msg.BuyerExecuteOrders {
		typeBuyerExecuteOrders = append(typeBuyerExecuteOrders, fiatFactory.NewSendFiat(msg.Relayer, in.BuyerAddress, in.SellerAddress, in.PegHash, in.FiatPegWallet))
	}
	_, err := fk.ExecuteFiatOrders(ctx, typeBuyerExecuteOrders)
	if err != nil {
		return err.Result()
	}
	ibcm.SetIngressSequence(ctx, msg.BuyerExecuteOrders[0].SourceChain, seq+1)
	
	return sdk.Result{}
}

func handleIBCHubSellerExecuteOrdersMsg(ctx sdk.Context, ibcm Mapper, ck bank.Keeper, msg MsgSellerExecuteOrders, ak acl.Keeper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, reputationKeeper reputation.Keeper) sdk.Result {
	tags, err, _ := ck.SellerExecuteTradeOrders(ctx, negotiationKeeper, orderKeeper, (toHubMsgSellerExecuteOrdermsg(msg)).SellerExecuteOrders, ak, reputationKeeper) // add toHub
	if err != nil {
		return err.Result()
	}
	err = ibcm.PostIBCMsgSellerExecuteOrdersPacket(ctx, msg)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{
		Tags: tags,
	}
}

func handlerIBCRelaySellerExecuteOrdersMsg(ctx sdk.Context, ibcm Mapper, ak assetFactory.Keeper, msg MsgRelaySellerExecuteOrders) sdk.Result {
	seq := ibcm.GetIngressSequence(ctx, msg.SellerExecuteOrders[0].SourceChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}
	var typeSellerExecuteOrders []assetFactory.SendAsset
	for _, in := range msg.SellerExecuteOrders {
		typeSellerExecuteOrders = append(typeSellerExecuteOrders, assetFactory.NewSendAsset(msg.Relayer, in.SellerAddress, in.BuyerAddress, in.PegHash))
	}
	_, err := ak.ExecuteAssetOrders(ctx, typeSellerExecuteOrders)
	if err != nil {
		return err.Result()
	}
	ibcm.SetIngressSequence(ctx, msg.SellerExecuteOrders[0].SourceChain, seq+1)
	
	return sdk.Result{}
}

// #####comdex
