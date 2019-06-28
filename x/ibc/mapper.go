package ibc

import (
	"fmt"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
)

// Mapper : IBC Mapper
type Mapper struct {
	key       sdk.StoreKey
	cdc       *wire.Codec
	codespace sdk.CodespaceType
}

// NewMapper : The Mapper should not take a CoinKeeper. Rather have the CoinKeeper
// take an Mapper.
func NewMapper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Mapper {
	// XXX: How are these codecs supposed to work?
	return Mapper{
		key:       key,
		cdc:       cdc,
		codespace: codespace,
	}
}

// PostIBCPacket :
// XXX: This is not the public API. This will change in MVP2 and will henceforth
// only be invoked from another module directly and not through a user
// transaction.
// TODO: Handle invalid IBC packets and return errors.
func (ibcm Mapper) PostIBCPacket(ctx sdk.Context, packet IBCTransferMsg) sdk.Error {
	// write everything into the state
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, packet.DestChain)
	bz, err := ibcm.cdc.MarshalBinary(packet)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(packet.DestChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(index + 1)
	if err != nil {
		panic(err)
	}
	store.Set(EgressLengthKey(packet.DestChain), bz)
	
	return nil
}

// ReceiveIBCPacket :
// XXX: In the future every module is able to register it's own handler for
// handling it's own IBC packets. The "ibc" handler will only route the packets
// to the appropriate callbacks.
// XXX: For now this handles all interactions with the CoinKeeper.
// XXX: This needs to do some authentication checking.
func (ibcm Mapper) ReceiveIBCPacket(ctx sdk.Context, packet IBCPacket) sdk.Error {
	return nil
}

// --------------------------
// Functions for accessing the underlying KVStore.

func marshalBinaryPanic(cdc *wire.Codec, value interface{}) []byte {
	res, err := cdc.MarshalBinary(value)
	if err != nil {
		panic(err)
	}
	return res
}

func unmarshalBinaryPanic(cdc *wire.Codec, bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinary(bz, ptr)
	if err != nil {
		panic(err)
	}
}

// GetIngressSequence :
// TODO add description
func (ibcm Mapper) GetIngressSequence(ctx sdk.Context, srcChain string) int64 {
	store := ctx.KVStore(ibcm.key)
	key := IngressSequenceKey(srcChain)
	
	bz := store.Get(key)
	if bz == nil {
		zero := marshalBinaryPanic(ibcm.cdc, int64(0))
		store.Set(key, zero)
		return 0
	}
	
	var res int64
	unmarshalBinaryPanic(ibcm.cdc, bz, &res)
	return res
}

// SetIngressSequence :
// TODO add description
func (ibcm Mapper) SetIngressSequence(ctx sdk.Context, srcChain string, sequence int64) {
	store := ctx.KVStore(ibcm.key)
	key := IngressSequenceKey(srcChain)
	
	bz := marshalBinaryPanic(ibcm.cdc, sequence)
	store.Set(key, bz)
}

// Retrieves the index of the currently stored outgoing IBC packets.
func (ibcm Mapper) getEgressLength(store sdk.KVStore, destChain string) int64 {
	bz := store.Get(EgressLengthKey(destChain))
	if bz == nil {
		zero := marshalBinaryPanic(ibcm.cdc, int64(0))
		store.Set(EgressLengthKey(destChain), zero)
		return 0
	}
	var res int64
	unmarshalBinaryPanic(ibcm.cdc, bz, &res)
	return res
}

// EgressKey : Stores an outgoing IBC packet under "egress/chain_id/index".
func EgressKey(destChain string, index int64) []byte {
	return []byte(fmt.Sprintf("egress/%s/%d", destChain, index))
}

// EgressLengthKey : Stores the number of outgoing IBC packets under "egress/index".
func EgressLengthKey(destChain string) []byte {
	return []byte(fmt.Sprintf("egress/%s", destChain))
}

// IngressSequenceKey : Stores the sequence number of incoming IBC packet under "ingress/index".
func IngressSequenceKey(srcChain string) []byte {
	return []byte(fmt.Sprintf("ingress/%s", srcChain))
}

// *****comdex

// PostIBCMsgIssueAssetsPacket : post issue asset packet to asset chain
func (ibcm Mapper) PostIBCMsgIssueAssetsPacket(ctx sdk.Context, msg MsgIssueAssets) sdk.Error {
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, msg.IssueAssets[0].DestinationChain)
	bz, err := ibcm.cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(msg.IssueAssets[0].DestinationChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	store.Set(EgressLengthKey(msg.IssueAssets[0].DestinationChain), bz)
	return nil
}

// PostIBCMsgRedeemAssetsPacket :
func (ibcm Mapper) PostIBCMsgRedeemAssetsPacket(ctx sdk.Context, msg MsgRedeemAssets) sdk.Error {
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, msg.RedeemAssets[0].DestinationChain)
	bz, err := ibcm.cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(msg.RedeemAssets[0].DestinationChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	store.Set(EgressLengthKey(msg.RedeemAssets[0].DestinationChain), bz)
	return nil
}

// PostIBCMsgSendAssetsPacket :
func (ibcm Mapper) PostIBCMsgSendAssetsPacket(ctx sdk.Context, msg MsgSendAssets) sdk.Error {
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, msg.SendAssets[0].DestinationChain)
	bz, err := ibcm.cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(msg.SendAssets[0].DestinationChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	store.Set(EgressLengthKey(msg.SendAssets[0].DestinationChain), bz)
	return nil
}

// PostIBCTransferMsg : post ibc transafer msg
func (ibcm Mapper) PostIBCTransferMsg(ctx sdk.Context, packet IBCTransferMsg) sdk.Error {
	// write everything into the state
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, packet.DestChain)
	bz, err := ibcm.cdc.MarshalBinary(packet)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(packet.DestChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	store.Set(EgressLengthKey(packet.DestChain), bz)
	
	return nil
}

// PostIBCMsgIssueFiatsPacket : post issue fiat packet to fiat chain
func (ibcm Mapper) PostIBCMsgIssueFiatsPacket(ctx sdk.Context, msg MsgIssueFiats) sdk.Error {
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, msg.IssueFiats[0].DestinationChain)
	bz, err := ibcm.cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(msg.IssueFiats[0].DestinationChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	store.Set(EgressLengthKey(msg.IssueFiats[0].DestinationChain), bz)
	return nil
}

// PostIBCMsgRedeemFiatsPacket :
func (ibcm Mapper) PostIBCMsgRedeemFiatsPacket(ctx sdk.Context, msg MsgRedeemFiats) sdk.Error {
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, msg.RedeemFiats[0].DestinationChain)
	bz, err := ibcm.cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(msg.RedeemFiats[0].DestinationChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	store.Set(EgressLengthKey(msg.RedeemFiats[0].DestinationChain), bz)
	return nil
}

// PostIBCMsgSendFiatsPacket :
func (ibcm Mapper) PostIBCMsgSendFiatsPacket(ctx sdk.Context, msg MsgSendFiats) sdk.Error {
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, msg.SendFiats[0].DestinationChain)
	bz, err := ibcm.cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(msg.SendFiats[0].DestinationChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	store.Set(EgressLengthKey(msg.SendFiats[0].DestinationChain), bz)
	return nil
}

// PostIBCMsgBuyerExecuteOrdersPacket :
func (ibcm Mapper) PostIBCMsgBuyerExecuteOrdersPacket(ctx sdk.Context, msg MsgBuyerExecuteOrders) sdk.Error {
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, msg.BuyerExecuteOrders[0].DestinationChain)
	
	bz, err := ibcm.cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(msg.BuyerExecuteOrders[0].DestinationChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressLengthKey(msg.BuyerExecuteOrders[0].DestinationChain), bz)
	return nil
}

// PostIBCMsgSellerExecuteOrdersPacket :
func (ibcm Mapper) PostIBCMsgSellerExecuteOrdersPacket(ctx sdk.Context, msg MsgSellerExecuteOrders) sdk.Error {
	store := ctx.KVStore(ibcm.key)
	index := ibcm.getEgressLength(store, msg.SellerExecuteOrders[0].DestinationChain)
	
	bz, err := ibcm.cdc.MarshalBinary(msg)
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressKey(msg.SellerExecuteOrders[0].DestinationChain, index), bz)
	bz, err = ibcm.cdc.MarshalBinary(int64(index + 1))
	if err != nil {
		panic(err)
	}
	
	store.Set(EgressLengthKey(msg.SellerExecuteOrders[0].DestinationChain), bz)
	return nil
}

// #####comdex
