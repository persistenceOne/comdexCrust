package types

import (
	"github.com/commitHub/commitBlockchain/types"
)

type GenesisState struct {
	Negotiations []types.Negotiation `json:"negotiations"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
