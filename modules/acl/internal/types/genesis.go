package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Accounts     []BaseACLAccount    `json:"accounts"`
	ZoneID       []cTypes.AccAddress `json:"zoneID"`
	Organization []Organization      `json:"organization"`
}

func NewGenesisState(accounts []BaseACLAccount, zoneID cTypes.AccAddress, organizations []Organization) GenesisState {
	return GenesisState{
		Accounts:     accounts,
		ZoneID:       []cTypes.AccAddress{zoneID},
		Organization: organizations,
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
