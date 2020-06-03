package types

type GenesisState struct {
	Orders []Order `json:"orders"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGensis(data GenesisState) error {
	return nil
}
