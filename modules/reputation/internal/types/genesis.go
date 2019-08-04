package types

type GenesisState struct {
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGensis(data GenesisState) error {
	return nil
}
