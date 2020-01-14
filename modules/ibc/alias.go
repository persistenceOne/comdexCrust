package ibc

// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/commitHub/commitBlockchain/modules/ibc/keeper
// ALIASGEN: github.com/commitHub/commitBlockchain/modules/ibc/types

import (
	"github.com/commitHub/commitBlockchain/modules/ibc/keeper"
	"github.com/commitHub/commitBlockchain/modules/ibc/types"
)

const (
	ModuleName       = types.ModuleName
	StoreKey         = types.StoreKey
	QuerierRoute     = types.QuerierRoute
	RouterKey        = types.RouterKey
	DefaultCodespace = types.DefaultCodespace
)

var (
	// functions aliases
	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)

type (
	Keeper = keeper.Keeper
)