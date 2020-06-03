package types

type GenesisState struct {
}

func DefaultGenesisState() {}

func ValidateGenesis(data GenesisState) error {
	return nil
}
