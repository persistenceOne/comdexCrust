package acl

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/modules/acl/internal/types"
)

func InitGenesis(ctx cTypes.Context, keeper Keeper, data GenesisState) (err error) {
	
	_ = keeper.SetZoneAddress(ctx, GetZoneKey(DefaultZoneID), data.ZoneID)
	
	_ = keeper.SetOrganization(ctx, GetOrganizationKey(DefaultOrganizationID), data.Organization)
	
	for _, acl := range data.Accounts {
		_ = keeper.SetACLAccount(ctx, acl)
	}
	
	return nil
	
}

func ExportGenesisState(ctx cTypes.Context, keeper Keeper) GenesisState {
	zoneID, _ := keeper.GetZoneAddress(ctx, GetZoneKey(DefaultZoneID))
	organizationID, _ := keeper.GetOrganization(ctx, DefaultOrganizationID)
	
	return GenesisState{
		Accounts:     keeper.GetACLAccounts(ctx),
		ZoneID:       zoneID,
		Organization: organizationID,
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
