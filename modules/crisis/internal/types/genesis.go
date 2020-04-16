package types

import (
	"fmt"

	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/comdexCrust/types"
)

// GenesisState - crisis genesis state
type GenesisState struct {
	ConstantFee cTypes.Coin `json:"constant_fee" yaml:"constant_fee"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(constantFee cTypes.Coin) GenesisState {
	return GenesisState{
		ConstantFee: constantFee,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() GenesisState {
	return GenesisState{
		ConstantFee: cTypes.NewCoin(types.BondDenom, cTypes.NewInt(10000000)),
	}
}

// ValidateGenesis - validate crisis genesis data
func ValidateGenesis(data GenesisState) error {
	if !data.ConstantFee.IsPositive() {
		return fmt.Errorf("constant fee must be positive: %s", data.ConstantFee)
	}
	return nil
}
