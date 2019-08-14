package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Accounts     []ACLAccount        `json:"accounts"`
	ZoneID       []cTypes.AccAddress `json:"zone_id"`
	Organization []Organization      `json:"organization"`
}

func NewGenesisState(accounts []ACLAccount, zoneID cTypes.AccAddress, organization Organization) GenesisState {
	return GenesisState{
		Accounts:     accounts,
		ZoneID:       []cTypes.AccAddress{zoneID},
		Organization: []Organization{organization},
	}
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts:     nil,
		ZoneID:       nil,
		Organization: nil,
	}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
