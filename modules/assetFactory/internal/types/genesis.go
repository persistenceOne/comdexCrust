package types

import (
	"github.com/commitHub/commitBlockchain/types"
)

type GenesisState struct {
	AssetPegs []types.AssetPeg `json:"assetPegs"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
