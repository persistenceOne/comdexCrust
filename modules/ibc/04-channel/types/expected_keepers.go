package types

import (
	clientexported "github.com/commitHub/commitBlockchain/modules/ibc/02-client/exported"
	connection "github.com/commitHub/commitBlockchain/modules/ibc/03-connection"
	commitment "github.com/commitHub/commitBlockchain/modules/ibc/23-commitment"
	commitTypes "github.com/commitHub/commitBlockchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ClientKeeper expected account IBC client keeper
type ClientKeeper interface {
	GetConsensusState(ctx sdk.Context, clientID string) (clientexported.ConsensusState, bool)
}

// ConnectionKeeper expected account IBC connection keeper
type ConnectionKeeper interface {
	GetConnection(ctx sdk.Context, connectionID string) (connection.ConnectionEnd, bool)
	VerifyMembership(
		ctx sdk.Context, connection connection.ConnectionEnd, height uint64,
		proof commitment.ProofI, path string, value []byte,
	) bool
	VerifyNonMembership(
		ctx sdk.Context, connection connection.ConnectionEnd, height uint64,
		proof commitment.ProofI, path string,
	) bool
}

// PortKeeper expected account IBC port keeper
type PortKeeper interface {
	Authenticate(key commitTypes.CapabilityKey, portID string) bool
}
