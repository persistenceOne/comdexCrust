package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const BondDenom = "ucommit"

// CapabilityKey represent the Cosmos SDK keys for object-capability
// generation in the IBC protocol as defined in https://github.com/cosmos/ics/tree/master/spec/ics-005-port-allocation#data-structures
type CapabilityKey sdk.StoreKey
