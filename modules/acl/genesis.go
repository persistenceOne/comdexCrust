package acl

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceSDK/modules/acl/internal/types"
)

func InitGenesis(ctx cTypes.Context, keeper Keeper, data GenesisState) (err error) {

	_ = keeper.SetZoneAddress(ctx, GetZoneKey(DefaultZoneID), DefaultZoneID)

	_ = keeper.SetOrganization(ctx, GetOrganizationKey(DefaultOrganizationID), Organization{})

	for _, acl := range data.Accounts {
		_ = keeper.SetACLAccount(ctx, &acl)
	}

	return nil

}

func ExportGenesisState(ctx cTypes.Context, keeper Keeper) GenesisState {
	zones := keeper.GetZones(ctx)
	organizations := keeper.GetOrganizations(ctx)

	return GenesisState{
		Accounts:     keeper.GetACLAccounts(ctx),
		ZoneID:       zones,
		Organization: organizations,
	}
}

func DefaultACLAccount(zoneID types.ZoneID, organizationID types.OrganizationID, address cTypes.AccAddress) types.ACLAccount {
	return &types.BaseACLAccount{
		Address:        address,
		ZoneID:         zoneID,
		OrganizationID: organizationID,
		ACL: types.ACL{
			IssueAsset:         true,
			IssueFiat:          true,
			SendAsset:          true,
			SendFiat:           true,
			BuyerExecuteOrder:  true,
			SellerExecuteOrder: true,
			ChangeBuyerBid:     true,
			ChangeSellerBid:    true,
			ConfirmBuyerBid:    true,
			ConfirmSellerBid:   true,
			Negotiation:        true,
			RedeemFiat:         true,
			RedeemAsset:        true,
			ReleaseAsset:       true,
		},
	}
}
