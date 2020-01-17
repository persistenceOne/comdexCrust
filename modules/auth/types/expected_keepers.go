package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceSDK/modules/supply/exported"
)

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx cTypes.Context, senderAddr cTypes.AccAddress, recipientModule string, amt cTypes.Coins) cTypes.Error
	GetModuleAccount(ctx cTypes.Context, moduleName string) exported.ModuleAccountI
	GetModuleAddress(moduleName string) cTypes.AccAddress
}
