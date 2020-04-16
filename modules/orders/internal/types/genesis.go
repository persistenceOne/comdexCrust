package types

import (
	"github.com/persistenceOne/comdexCrust/types"
)

type GenesisState struct {
	Orders []types.Order `json:"orders"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGensis(data GenesisState) error {
	return nil
}
