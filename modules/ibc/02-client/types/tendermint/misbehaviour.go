package tendermint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/commitHub/commitBlockchain/modules/ibc/02-client/exported"
)

// CheckMisbehaviour checks if the evidence provided is a misbehaviour
func CheckMisbehaviour(evidence exported.Evidence) sdk.Error {
	// TODO: check evidence
	return nil
}