package types

type GenesisState struct {
	Reputations []AccountReputation
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGensis(data GenesisState) error {
	return nil
}
