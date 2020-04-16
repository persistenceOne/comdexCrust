package orders

import (
	"github.com/persistenceOne/comdexCrust/modules/orders/internal/keeper"
	"github.com/persistenceOne/comdexCrust/modules/orders/internal/types"
	types2 "github.com/persistenceOne/comdexCrust/types"
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
	Keeper       = keeper.Keeper
	Order        = types2.Order
	BaseOrder    = types2.BaseOrder

	ACLKeeper     = types.ACLKeeper
	AccountKeeper = types.AccountKeeper
)
