package negotiation

import (
	"github.com/commitHub/commitBlockchain/modules/negotiation/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/negotiation/internal/types"
)

const (
	StoreKey     = types.StoreKey
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
	
	CodeInvalidSignature = types.CodeInvalidSignature
	DefaultCodeSpace     = types.DefaultCodeSpace
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc
	
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	
	NewNegotiation         = types.NewNegotiation
	NewSignNegotiationBody = types.NewSignNegotiationBody
	
	EventTypeChangeNegotiationBid  = types.EventTypeChangeNegotiationBid
	EventTypeConfirmNegotiationBid = types.EventTypeConfirmNegotiationBid
	
	AttributeKeyNegotiationID = types.AttributeKeyNegotiationID
	AttributeKeyBuyerAddress  = types.AttributeKeyBuyerAddress
	AttributeKeySellerAddress = types.AttributeKeySellerAddress
	AttributeKeyPegHash       = types.AttributeKeyPegHash
	
	ErrCodeVerifySignature = types.ErrVerifySignature
	ErrCodeInvalidBid      = types.ErrInvalidBid
	
	NewQuerier = keeper.NewQuerier
	NewKeeper  = keeper.NewKeeper
	
	GetNegotiationKey = types.GetNegotiationKey
	
	ErrUnauthorized = types.ErrUnauthorized
	
	BuildMsgChangeBuyerBid   = types.BuildMsgChangeBuyerBid
	BuildMsgChangeSellerBid  = types.BuildMsgChangeSellerBid
	BuildMsgConfirmBuyerBid  = types.BuildMsgConfirmBuyerBid
	BuildMsgConfirmSellerBid = types.BuildMsgConfirmSellerBid
	
	GetNegotiationIDFromString = types.GetNegotiationIDFromString
)

type (
	GenesisState = types.GenesisState
	Keeper       = keeper.Keeper
	
	NegotiationID = types.NegotiationID
	Negotiation   = types.Negotiation
	
	BaseNegotiation = types.BaseNegotiation
	
	MsgChangeBuyerBids   = types.MsgChangeBuyerBids
	MsgChangeSellerBids  = types.MsgChangeSellerBids
	MsgConfirmBuyerBids  = types.MsgConfirmBuyerBids
	MsgConfirmSellerBids = types.MsgConfirmSellerBids
	
	ChangeBid  = types.ChangeBid
	ConfirmBid = types.ConfirmBid
	
	Signature = types.Signature
)
