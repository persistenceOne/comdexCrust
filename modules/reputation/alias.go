package reputation

import (
	"github.com/commitHub/commitBlockchain/modules/reputation/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/reputation/internal/types"
)

const (
	ModuleName   = types.ModuleName
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
)

var (
	RegisterCodec       = types.RegisterCodec
	ModuleCdc           = types.ModuleCdc
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGensis
	
	EventTypeSetBuyerRatingToFeedback  = types.EventTypeSetBuyerRatingToFeedback
	EventTypeSetSellerRatingToFeedback = types.EventTypeSetSellerRatingToFeedback
	
	AttributeKeyPegHash = types.AttributeKeyPegHash
	AttributeKeyRating  = types.AttributeKeyRating
	AttributeKeyFrom    = types.AttributeKeyFrom
	AttributeKeyTo      = types.AttributeKeyTo
	
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)

type (
	Keeper             = keeper.Keeper
	GenesisState       = types.GenesisState
	MsgBuyerFeedbacks  = types.MsgBuyerFeedbacks
	MsgSellerFeedbacks = types.MsgSellerFeedbacks
	
	AccountReputation     = types.AccountReputation
	BaseAccountReputation = types.BaseAccountReputation
)
