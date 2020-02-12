package types

import (
	"github.com/persistenceOne/persistenceSDK/types"
)

type GenesisState struct {
	Reputations []types.BaseAccountReputation
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGensis(data GenesisState) error {
	return nil
}