package types

import (
	"github.com/persistenceOne/persistenceSDK/types"
)

type GenesisState struct {
	FiatPegs []types.FiatPeg `json:"fiatPegs"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
