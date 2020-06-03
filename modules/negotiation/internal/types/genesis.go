package types

type GenesisState struct {
	Negotiations []Negotiation `json:"negotiations"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
