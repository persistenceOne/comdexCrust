package orders

import (
	"github.com/commitHub/commitBlockchain/modules/orders/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/orders/internal/types"
)

const (
	ModuleName   = types.ModuleName
	QuerierRoute = types.QuerierRoute
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	
	DefaultCodeSpace = types.DefaultCodeSpace
)

var (
	RegisterCodec       = types.RegisterCodec
	ModuleCdc           = types.ModuleCdc
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGensis
	
	NewQuerier = keeper.NewQuerier
	
	GetNegotiationKey = types.GetOrderKey
	OrdersKey         = types.OrdersKey
	NewKeeper         = keeper.NewKeeper
	
	ErrUnauthorized = types.ErrUnauthorized
)

type (
	GenesisState = types.GenesisState
	Keeper = keeper.Keeper
	Order = types.Order
	BaseOrder = types.BaseOrder
	
	ACLKeeper = types.ACLKeeper
	AccountKeeper = types.AccountKeeper
)
