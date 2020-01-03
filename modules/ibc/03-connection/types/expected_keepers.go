package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	client "github.com/commitHub/commitBlockchain/modules/ibc/02-client"
	clientexported "github.com/commitHub/commitBlockchain/modules/ibc/02-client/exported"
	commitment "github.com/commitHub/commitBlockchain/modules/ibc/23-commitment"
)

// ClientKeeper expected account IBC client keeper
type ClientKeeper interface {
	GetConsensusState(ctx sdk.Context, clientID string) (clientexported.ConsensusState, bool)
	GetClientState(ctx sdk.Context, clientID string) (client.State, bool)
	VerifyMembership(
		ctx sdk.Context, clientID string, height uint64,
		proof commitment.ProofI, path commitment.PathI, value []byte,
	) bool
	VerifyNonMembership(
		ctx sdk.Context, clientID string, height uint64,
		proof commitment.ProofI, path commitment.PathI,
	) bool
}
