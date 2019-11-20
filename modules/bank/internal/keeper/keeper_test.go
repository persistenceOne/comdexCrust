package keeper_test

import (
	"testing"
	"time"

	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/exported"
	bankTypes "github.com/commitHub/commitBlockchain/modules/bank/internal/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/commitHub/commitBlockchain/simApp"
	"github.com/commitHub/commitBlockchain/types"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func getGenesisAccount(ctx cTypes.Context, ak auth.AccountKeeper) exported.Account {
	allAccounts := ak.GetAllAccounts(ctx)
	var genesisAccount exported.Account
	for _, account := range allAccounts {
		if account.GetAccountNumber() == 0 {
			genesisAccount = account
		}
	}
	return genesisAccount
}

func getNegotiation(ctx cTypes.Context, buyer, seller cTypes.AccAddress, assetPegHash types.PegHash, bid, time int64) types.Negotiation {
	negotiation := types.NewNegotiation(buyer, seller, assetPegHash)
	negotiation.SetBid(bid)
	negotiation.SetTime(time)
	negotiation.SetBuyerSignature(types.Signature([]byte(negotiation.GetNegotiationID().Bytes())))
	negotiation.SetSellerSignature(types.Signature([]byte(negotiation.GetNegotiationID().Bytes())))
	negotiation.SetBuyerBlockHeight(ctx.BlockHeight())
	negotiation.SetSellerBlockHeight(ctx.BlockHeight())
	negotiation.SetBuyerContractHash("BuyerContractHash")
	negotiation.SetSellerContractHash("SellerContractHash")
	return negotiation
}

func TestBaseSendKeeper_DelegateCoins(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	origCoins := cTypes.NewCoins(cTypes.NewInt64Coin("stake", 100))
	delCoins := cTypes.NewCoins(cTypes.NewInt64Coin("stake", 50))

	addr1 := cTypes.AccAddress([]byte("addr1"))
	addr2 := cTypes.AccAddress([]byte("addr2"))
	addr3 := cTypes.AccAddress([]byte("addr3"))
	addr4 := cTypes.AccAddress([]byte("addr4"))
	addrModule := cTypes.AccAddress([]byte("moduleAcc"))

	bacc := auth.NewBaseAccountWithAddress(addr1)
	bacc.SetCoins(origCoins)
	macc := ak.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	vacc := auth.NewContinuousVestingAccount(&bacc, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := ak.NewAccountWithAddress(ctx, addr2)
	ak.SetAccount(ctx, vacc)
	ak.SetAccount(ctx, acc)
	ak.SetAccount(ctx, macc)
	app.BankKeeper.SetCoins(ctx, addr2, origCoins)

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))

	// Negative Cases
	err := app.BankKeeper.DelegateCoins(ctx, addr3, addrModule, delCoins)
	require.Error(t, err, "account "+addr3.String()+" does not exist")
	err = app.BankKeeper.DelegateCoins(ctx, addr2, addr4, delCoins)
	require.Error(t, err, "module account "+addr4.String()+" does not exist")

	// require the ability for a non-vesting account to delegate
	err = app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins)
	acc = ak.GetAccount(ctx, addr2)
	macc = ak.GetAccount(ctx, addrModule)
	require.NoError(t, err)
	require.Equal(t, origCoins.Sub(delCoins), acc.GetCoins())
	require.Equal(t, delCoins, macc.GetCoins())

	// require the ability for a vesting account to delegate
	err = app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins)
	vacc = ak.GetAccount(ctx, addr1).(*auth.ContinuousVestingAccount)
	require.NoError(t, err)
	require.Equal(t, delCoins, vacc.GetCoins())
}

func TestUndelegateCoins(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper
	now := tmtime.Now()
	endTime := now.Add(24 * time.Hour)

	origCoins := cTypes.NewCoins(cTypes.NewInt64Coin("stake", 100))
	delCoins := cTypes.NewCoins(cTypes.NewInt64Coin("stake", 50))

	addr1 := cTypes.AccAddress([]byte("addr1"))
	addr2 := cTypes.AccAddress([]byte("addr2"))
	addrModule := cTypes.AccAddress([]byte("moduleAcc"))

	bacc := auth.NewBaseAccountWithAddress(addr1)
	bacc.SetCoins(origCoins)
	macc := ak.NewAccountWithAddress(ctx, addrModule) // we don't need to define an actual module account bc we just need the address for testing
	vacc := auth.NewContinuousVestingAccount(&bacc, ctx.BlockHeader().Time.Unix(), endTime.Unix())
	acc := ak.NewAccountWithAddress(ctx, addr2)
	ak.SetAccount(ctx, vacc)
	ak.SetAccount(ctx, acc)
	ak.SetAccount(ctx, macc)
	app.BankKeeper.SetCoins(ctx, addr2, origCoins)

	ctx = ctx.WithBlockTime(now.Add(12 * time.Hour))

	// require the ability for a non-vesting account to delegate
	err := app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins)
	require.NoError(t, err)

	acc = ak.GetAccount(ctx, addr2)
	macc = ak.GetAccount(ctx, addrModule)
	require.Equal(t, origCoins.Sub(delCoins), acc.GetCoins())
	require.Equal(t, delCoins, macc.GetCoins())

	// require the ability for a non-vesting account to undelegate
	err = app.BankKeeper.UndelegateCoins(ctx, addrModule, addr2, delCoins)
	require.NoError(t, err)

	acc = ak.GetAccount(ctx, addr2)
	macc = ak.GetAccount(ctx, addrModule)
	require.Equal(t, origCoins, acc.GetCoins())
	require.True(t, macc.GetCoins().Empty())

	// require the ability for a vesting account to delegate
	err = app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins)
	require.NoError(t, err)

	vacc = ak.GetAccount(ctx, addr1).(*auth.ContinuousVestingAccount)
	macc = ak.GetAccount(ctx, addrModule)
	require.Equal(t, origCoins.Sub(delCoins), vacc.GetCoins())
	require.Equal(t, delCoins, macc.GetCoins())

	// require the ability for a vesting account to undelegate
	err = app.BankKeeper.UndelegateCoins(ctx, addrModule, addr1, delCoins)
	require.NoError(t, err)

	vacc = ak.GetAccount(ctx, addr1).(*auth.ContinuousVestingAccount)
	macc = ak.GetAccount(ctx, addrModule)
	require.Equal(t, origCoins, vacc.GetCoins())
	require.True(t, macc.GetCoins().Empty())
}

func TestBaseSendKeeper_InputOutputCoins(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	addr2 := cTypes.AccAddress([]byte("addr2"))

	coins1 := cTypes.NewCoins(cTypes.NewInt64Coin("foo", 100))
	coins2 := cTypes.NewCoins(cTypes.NewInt64Coin("foo", 100))
	inputCoins := cTypes.NewCoins(cTypes.NewInt64Coin("foo", 10))
	outputCoins := cTypes.NewCoins(cTypes.NewInt64Coin("foo", 10))

	app.BankKeeper.SetCoins(ctx, addr1, coins1)
	app.BankKeeper.SetCoins(ctx, addr2, coins2)

	inputs := []bankTypes.Input{bankTypes.NewInput(addr1, inputCoins)}
	outputs := []bankTypes.Output{bankTypes.NewOutput(addr2, outputCoins)}

	err := app.BankKeeper.InputOutputCoins(ctx, inputs, outputs)
	require.NoError(t, err)
}

func TestBaseSendKeeper_SendCoins(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	addr1 := cTypes.AccAddress([]byte("addr1"))
	addr2 := cTypes.AccAddress([]byte("addr2"))

	coins1 := cTypes.NewCoins(cTypes.NewInt64Coin("foo", 100))
	coins2 := cTypes.NewCoins(cTypes.NewInt64Coin("foo", 100))
	sendCoin1 := cTypes.NewCoins(cTypes.NewInt64Coin("foo", 50))
	sendCoin2 := cTypes.NewCoins(cTypes.NewInt64Coin("foo", 200))

	app.BankKeeper.SetCoins(ctx, addr1, coins1)
	app.BankKeeper.SetCoins(ctx, addr2, coins2)
	err := app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoin2)
	require.Error(t, err, "insufficient account funds; 100foo < 200foo")

	err = app.BankKeeper.SendCoins(ctx, addr1, addr2, sendCoin1)

	require.NoError(t, err)
	require.Equal(t, cTypes.NewCoins(cTypes.NewInt64Coin("foo", 50)), ak.GetAccount(ctx, addr1).GetCoins())
	require.Equal(t, cTypes.NewCoins(cTypes.NewInt64Coin("foo", 150)), ak.GetAccount(ctx, addr2).GetCoins())

}

func TestBaseSendKeeper_GetEnabled(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	enabled := app.BankKeeper.GetSendEnabled(ctx)
	require.True(t, enabled)
}

func TestBaseSendKeeper_SetEnabled(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	app.BankKeeper.SetSendEnabled(ctx, true)
	enabled := app.BankKeeper.GetSendEnabled(ctx)
	require.True(t, enabled)
}

func TestBaseSendKeeper_DefineZones(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, addr1))
	zone := cTypes.AccAddress([]byte("zone"))

	defineZone := bankTypes.NewDefineZone(addr1, zone, acl.DefaultZoneID)
	err := app.BankKeeper.DefineZones(ctx, defineZone)
	require.Error(t, err, "Account "+addr1.String()+" is not the genesis account. Zones can only be defined by the genesis account.")

	defineZone = bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	err = app.BankKeeper.DefineZones(ctx, defineZone)
	require.NoError(t, err)

	defineZone = bankTypes.NewDefineZone(genesisAccount.GetAddress(), addr1, acl.DefaultZoneID)
	err = app.BankKeeper.DefineZones(ctx, defineZone)
	require.Error(t, err, "zone with this given id already exist")
}

func TestBaseSendKeeper_DefineOrganizations(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	addr1 := cTypes.AccAddress([]byte("addr1"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, addr1))

	defineOrganization := bankTypes.NewDefineOrganization(addr1, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	err := app.BankKeeper.DefineOrganizations(ctx, defineOrganization)
	require.Error(t, err, "Account"+addr1.String()+" is not the zone account. Organizations can only be defined by the zone account.")

	defineOrganization = bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	err = app.BankKeeper.DefineOrganizations(ctx, defineOrganization)
	require.NoError(t, err)

	defineOrganization = bankTypes.NewDefineOrganization(zone, addr1, acl.DefaultOrganizationID, defineZone.ZoneID)
	err = app.BankKeeper.DefineOrganizations(ctx, defineOrganization)
	require.Error(t, err, "organization with this given id already exist")
}

func TestBaseSendKeeper_DefineACLs(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, addr1))

	trader := cTypes.AccAddress([]byte("trader"))
	aclAccount := acl.BaseACLAccount{
		Address:        trader,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}

	defineACL := bankTypes.NewDefineACL(addr1, trader, &aclAccount)
	err := app.BankKeeper.DefineACLs(ctx, defineACL)
	require.Error(t, err, "Account "+addr1.String()+" does not have access to define acl for account "+trader.String()+".")

	defineACL = bankTypes.NewDefineACL(organization, trader, &aclAccount)
	err = app.BankKeeper.DefineACLs(ctx, defineACL)
	require.NoError(t, err)
}

func TestBaseSendKeeper_IssueAssetsToWallets(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	trader := cTypes.AccAddress([]byte("trader"))
	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
	}

	issueAsset := bankTypes.NewIssueAsset(trader, trader, &assetPeg)
	err := app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	require.Error(t, err, "acl for this account not defined")

	aclAccount := acl.BaseACLAccount{
		Address:        trader,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: false},
	}

	defineACL := bankTypes.NewDefineACL(zone, trader, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	issueAsset = bankTypes.NewIssueAsset(trader, trader, &assetPeg)
	err = app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	require.Error(t, err, "Assets cant be issued to account "+trader.String()+".")

	aclAccount.SetACL(acl.ACL{IssueAsset: true})
	defineACL = bankTypes.NewDefineACL(zone, trader, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	issueAsset = bankTypes.NewIssueAsset(trader, trader, &assetPeg)
	err = app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	require.NoError(t, err)

	assetPeg.SetModerated(true)

	issueAsset = bankTypes.NewIssueAsset(addr1, trader, &assetPeg)
	err = app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	require.Error(t, err, "Unauthorised transaction.")

	issueAsset = bankTypes.NewIssueAsset(zone, trader, &assetPeg)
	err = app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	require.NoError(t, err)
}

func TestBaseSendKeeper_IssueFiatsToWallets(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	trader := cTypes.AccAddress([]byte("trader"))
	aclAccount := acl.BaseACLAccount{
		Address:        trader,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueFiat: false},
	}

	fiatPeg := types.BaseFiatPeg{
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	issueFiat := bankTypes.NewIssueFiat(addr1, trader, &fiatPeg)
	err := app.BankKeeper.IssueFiatsToWallets(ctx, issueFiat)
	require.Error(t, err, "Unauthorised transaction.")

	defineACL := bankTypes.NewDefineACL(zone, trader, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	issueFiat = bankTypes.NewIssueFiat(zone, trader, &fiatPeg)
	err = app.BankKeeper.IssueFiatsToWallets(ctx, issueFiat)
	require.Error(t, err, "Fiats can't be issued to account "+trader.String()+".")

	aclAccount.SetACL(acl.ACL{IssueFiat: true})
	defineACL = bankTypes.NewDefineACL(zone, trader, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	issueFiat = bankTypes.NewIssueFiat(zone, trader, &fiatPeg)
	err = app.BankKeeper.IssueFiatsToWallets(ctx, issueFiat)
	require.NoError(t, err)
}

func TestBaseSendKeeper_ReleasedLockedAssets(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	trader := cTypes.AccAddress([]byte("trader"))
	ak.SetAccount(ctx, ak.NewAccountWithAddress(ctx, trader))
	aclAccount := acl.BaseACLAccount{
		Address:        trader,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, ReleaseAsset: false},
	}

	releaseAsset := bankTypes.NewReleaseAsset(zone, trader, []byte(""))
	err := app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	require.Error(t, err, "To account acl not defined.")

	defineACL := bankTypes.NewDefineACL(zone, trader, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	issueAsset := bankTypes.NewIssueAsset(zone, trader, &assetPeg)
	app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)

	releaseAsset = bankTypes.NewReleaseAsset(zone, trader, ak.GetAccount(ctx, trader).GetAssetPegWallet()[0].GetPegHash())
	err = app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	require.Error(t, err, "Assets cannot be released for account "+trader.String()+". Access Denied.")

	aclAccount.SetACL(acl.ACL{IssueAsset: true, ReleaseAsset: true})
	defineACL = bankTypes.NewDefineACL(zone, trader, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	releaseAsset = bankTypes.NewReleaseAsset(zone, trader, []byte("PegHash"))
	err = app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	require.Error(t, err, "Asset peg not found.")

	releaseAsset = bankTypes.NewReleaseAsset(zone, trader, ak.GetAccount(ctx, trader).GetAssetPegWallet()[0].GetPegHash())
	err = app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	require.NoError(t, err)

}

func TestBaseSendKeeper_SendAssetsToWallets(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, SendAsset: false, ReleaseAsset: true},
	}

	defineACL := bankTypes.NewDefineACL(zone, seller, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)
	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	issueAsset := bankTypes.NewIssueAsset(zone, seller, &assetPeg)
	app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	assetPegHash := ak.GetAccount(ctx, seller).GetAssetPegWallet()[0].GetPegHash()

	sendAsset := bankTypes.NewSendAsset(buyer, buyer, assetPegHash)
	err := app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "Unauthorised transaction.")

	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "Unauthorised transaction.")

	aclAccount.SetACL(acl.ACL{IssueAsset: true, SendAsset: true, ReleaseAsset: true})
	defineACL = bankTypes.NewDefineACL(zone, seller, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "negotiation not found.")

	negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()-1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "Negotiation time expired.")

	negotiation.SetTime(ctx.BlockHeight() + 1)
	sellerSignature := negotiation.GetSellerSignature()
	negotiation.SetSellerSignature(nil)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "Signatures are not present")

	negotiation.SetSellerSignature(sellerSignature)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	app.NegotiationKeeper.SetNegotiation(ctx, getNegotiation(ctx, buyer, seller, []byte("WrongPegHash"), 500, ctx.BlockHeight()+1))
	sendAsset = bankTypes.NewSendAsset(seller, buyer, []byte("WrongPegHash"))
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "Asset not found.")

	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "Asset locked.")

	releaseAsset := bankTypes.NewReleaseAsset(zone, seller, assetPegHash)
	app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)

	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.NoError(t, err)
}

func TestBaseSendKeeper_SendFiatsToWallets(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	seller := cTypes.AccAddress([]byte("seller"))
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}
	defineACL := bankTypes.NewDefineACL(zone, seller, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)
	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	issueAsset := bankTypes.NewIssueAsset(zone, seller, &assetPeg)
	app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	assetPegHash := ak.GetAccount(ctx, seller).GetAssetPegWallet()[0].GetPegHash()
	releaseAsset := bankTypes.NewReleaseAsset(zone, seller, assetPegHash)
	app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)

	buyer := cTypes.AccAddress([]byte("buyer"))
	aclAccount = acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendFiat: false, ReleaseAsset: true},
	}
	defineACL = bankTypes.NewDefineACL(zone, buyer, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	fiatPeg := types.BaseFiatPeg{
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	issueFiat := bankTypes.NewIssueFiat(zone, buyer, &fiatPeg)
	app.BankKeeper.IssueFiatsToWallets(ctx, issueFiat)
	fiatPegHash := ak.GetAccount(ctx, buyer).GetFiatPegWallet()[0].GetPegHash()

	sendFiat := bankTypes.NewSendFiat(addr1, seller, fiatPegHash, 500)
	err := app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.Error(t, err, "Unauthorised transaction.")

	sendFiat = bankTypes.NewSendFiat(buyer, seller, fiatPegHash, 500)
	err = app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.Error(t, err, "Unauthorised transaction.")

	aclAccount.SetACL(acl.ACL{IssueAsset: true, IssueFiat: true, SendFiat: true, ReleaseAsset: true})
	app.BankKeeper.DefineACLs(ctx, defineACL)

	sendFiat = bankTypes.NewSendFiat(buyer, seller, fiatPegHash, 500)
	err = app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.Error(t, err, "negotiation not found.")

	negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()-1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sendFiat = bankTypes.NewSendFiat(buyer, seller, fiatPegHash, 500)
	err = app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.Error(t, err, "Negotiation time expired.")

	negotiation.SetTime(ctx.BlockHeight() + 1)
	buyerSignature := negotiation.GetBuyerSignature()
	negotiation.SetBuyerSignature(nil)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sendFiat = bankTypes.NewSendFiat(buyer, seller, fiatPegHash, 500)
	err = app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.Error(t, err, "Signatures are not present")

	negotiation.SetBuyerSignature(buyerSignature)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sendFiat = bankTypes.NewSendFiat(buyer, seller, fiatPegHash, 5000000)
	err = app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.Error(t, err, "Insufficient funds")

	sendFiat = bankTypes.NewSendFiat(buyer, seller, fiatPegHash, negotiation.GetBid())
	err = app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.NoError(t, err)
}

func TestBaseSendKeeper_BuyerExecuteOrder(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	seller := cTypes.AccAddress([]byte("seller"))
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}
	defineACL := bankTypes.NewDefineACL(zone, seller, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)
	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	issueAsset := bankTypes.NewIssueAsset(zone, seller, &assetPeg)
	app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	assetPegHash := ak.GetAccount(ctx, seller).GetAssetPegWallet()[0].GetPegHash()
	releaseAsset := bankTypes.NewReleaseAsset(zone, seller, assetPegHash)
	app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)

	buyer := cTypes.AccAddress([]byte("buyer"))
	aclAccount = acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: false, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}
	defineACL = bankTypes.NewDefineACL(zone, buyer, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	fiatPeg := types.BaseFiatPeg{
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	issueFiat := bankTypes.NewIssueFiat(zone, buyer, &fiatPeg)
	app.BankKeeper.IssueFiatsToWallets(ctx, issueFiat)

	negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, 10000)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	buyerExecuteOrder := bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ := app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "Asset token not found!")

	sendAsset := bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(addr1, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "To account acl not defined.")

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account "+buyer.String()+". Access Denied.")

	aclAccount.SetACL(acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true})
	defineACL = bankTypes.NewDefineACL(zone, buyer, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.NoError(t, err)

	//Unmoderated
	assetPeg = types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
		Locked:        true,
	}
	issueAsset = bankTypes.NewIssueAsset(seller, seller, &assetPeg)
	app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	assetPegWallet := ak.GetAccount(ctx, seller).GetAssetPegWallet()
	assetPegHash = ak.GetAccount(ctx, seller).GetAssetPegWallet()[1].GetPegHash()
	releaseAsset = bankTypes.NewReleaseAsset(zone, seller, assetPegHash)
	app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	negotiation = getNegotiation(ctx, buyer, seller, assetPegHash, 500, 10000)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)

	order := app.OrderKeeper.NewOrder(addr1, seller, assetPegHash)
	order.SetAssetPegWallet(assetPegWallet)
	app.OrderKeeper.SetOrder(ctx, order)
	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(buyer, addr1, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "To account acl not defined.")

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account %v."+buyer.String()+" Access Denied.")

	aclAccount.SetACL(acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: false, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true})
	defineACL = bankTypes.NewDefineACL(zone, buyer, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account %v."+buyer.String()+" Access Denied.")

	aclAccount.SetACL(acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true})
	defineACL = bankTypes.NewDefineACL(zone, buyer, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.NoError(t, err)

}

func TestBaseSendKeeper_SellerExecuteOrder(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	seller := cTypes.AccAddress([]byte("seller"))
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: false, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}
	defineACL := bankTypes.NewDefineACL(zone, seller, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)
	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	issueAsset := bankTypes.NewIssueAsset(zone, seller, &assetPeg)
	app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	assetPegHash := ak.GetAccount(ctx, seller).GetAssetPegWallet()[0].GetPegHash()
	releaseAsset := bankTypes.NewReleaseAsset(zone, seller, assetPegHash)
	app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)

	buyer := cTypes.AccAddress([]byte("buyer"))
	aclAccount = acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: false, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}
	defineACL = bankTypes.NewDefineACL(zone, buyer, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	fiatPeg := types.BaseFiatPeg{
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	issueFiat := bankTypes.NewIssueFiat(zone, buyer, &fiatPeg)
	app.BankKeeper.IssueFiatsToWallets(ctx, issueFiat)
	fiatPegHash := ak.GetAccount(ctx, buyer).GetFiatPegWallet()[0].GetPegHash()

	negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, 10000)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sellerExecuteOrder := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ := app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Asset token not found!")

	sendAsset := bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(addr1, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "To account acl not defined.")

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account "+seller.String()+". Access Denied.")

	aclAccount.SetAddress(seller)
	aclAccount.SetACL(acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true})
	defineACL = bankTypes.NewDefineACL(zone, seller, &aclAccount)
	err = app.BankKeeper.DefineACLs(ctx, defineACL)

	sendFiat := bankTypes.NewSendFiat(buyer, seller, fiatPegHash, negotiation.GetBid())
	app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	//Unmoderated
	seller2 := cTypes.AccAddress([]byte("seller2"))
	assetPegHash = []byte("DEFAULT")
	assetPeg = types.BaseAssetPeg{
		PegHash:       assetPegHash,
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
		Locked:        false,
	}
	negotiation = getNegotiation(ctx, buyer, seller2, assetPegHash, 500, ctx.BlockHeight()+1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	order := app.OrderKeeper.NewOrder(buyer, seller2, assetPegHash)
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller, buyer, seller2, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "To account acl not defined.")

	aclAccount = acl.BaseACLAccount{
		Address:        seller2,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: false, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}
	defineACL = bankTypes.NewDefineACL(zone, seller2, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(addr1, buyer, seller2, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account "+addr1.String()+". Access Denied.")

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller2, buyer, seller2, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account "+seller2.String()+". Access Denied.")

	aclAccount.SetACL(acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true})
	defineACL = bankTypes.NewDefineACL(zone, seller2, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller2, buyer, seller2, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

}

func TestBaseSendKeeper_RedeemAsset(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueAsset: true, ReleaseAsset: true, RedeemAsset: false},
	}
	defineACL := bankTypes.NewDefineACL(zone, seller, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)
	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}

	assetPegHash := []byte("PegHash")

	redeemAsset := bankTypes.NewRedeemAsset(buyer, seller, assetPegHash)
	err := app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.Error(t, err, "To account acl not defined.")

	redeemAsset = bankTypes.NewRedeemAsset(zone, seller, assetPegHash)
	err = app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.Error(t, err, "Assets can't be redeemed from account "+seller.String()+".")

	aclAccount.SetACL(acl.ACL{IssueAsset: true, ReleaseAsset: true, RedeemAsset: true})
	app.BankKeeper.DefineACLs(ctx, defineACL)

	redeemAsset = bankTypes.NewRedeemAsset(zone, seller, assetPegHash)
	err = app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.Error(t, err, "No Assets Found!")

	issueAsset := bankTypes.NewIssueAsset(zone, seller, &assetPeg)
	app.BankKeeper.IssueAssetsToWallets(ctx, issueAsset)
	assetPegHash = ak.GetAccount(ctx, seller).GetAssetPegWallet()[0].GetPegHash()
	releaseAsset := bankTypes.NewReleaseAsset(zone, seller, assetPegHash)
	app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)

	redeemAsset = bankTypes.NewRedeemAsset(zone, seller, []byte("PegHash"))
	err = app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.Error(t, err, "No Assets With Given PegHash Found!")

	redeemAsset = bankTypes.NewRedeemAsset(zone, seller, assetPegHash)
	err = app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.NoError(t, err)
}

func TestBaseSendKeeper_RedeemFiat(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	genesisAccount := getGenesisAccount(ctx, ak)
	zone := cTypes.AccAddress([]byte("zone"))
	defineZone := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	app.BankKeeper.DefineZones(ctx, defineZone)

	organization := cTypes.AccAddress([]byte("organization"))
	defineOrganization := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, defineZone.ZoneID)
	app.BankKeeper.DefineOrganizations(ctx, defineOrganization)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	aclAccount := acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         defineOrganization.ZoneID,
		OrganizationID: defineOrganization.OrganizationID,
		ACL:            acl.ACL{IssueFiat: true, SendFiat: true, RedeemFiat: false},
	}
	defineACL := bankTypes.NewDefineACL(zone, buyer, &aclAccount)
	app.BankKeeper.DefineACLs(ctx, defineACL)

	fiatPeg := types.BaseFiatPeg{
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}

	redeemFiat := bankTypes.NewRedeemFiat(seller, zone, 500)
	err := app.BankKeeper.RedeemFiatsFromWallets(ctx, redeemFiat)
	require.Error(t, err, "To account acl not defined.")

	redeemFiat = bankTypes.NewRedeemFiat(buyer, zone, 500)
	err = app.BankKeeper.RedeemFiatsFromWallets(ctx, redeemFiat)
	require.Error(t, err, "Fiats can't be redeemed from account "+buyer.String()+".")

	aclAccount.SetACL(acl.ACL{IssueFiat: true, SendFiat: true, RedeemFiat: true})
	app.BankKeeper.DefineACLs(ctx, defineACL)

	issueFiat := bankTypes.NewIssueFiat(zone, buyer, &fiatPeg)
	app.BankKeeper.IssueFiatsToWallets(ctx, issueFiat)

	redeemFiat = bankTypes.NewRedeemFiat(buyer, zone, 50000)
	err = app.BankKeeper.RedeemFiatsFromWallets(ctx, redeemFiat)
	require.Error(t, err, "Redeemed amount higher than the account balance")

	redeemFiat = bankTypes.NewRedeemFiat(buyer, zone, 500)
	err = app.BankKeeper.RedeemFiatsFromWallets(ctx, redeemFiat)
	require.NoError(t, err)
}

func TestBaseSendKeeper_PrivateExecuteOrder(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, buyer))
	seller := cTypes.AccAddress([]byte("seller"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, seller))
	aclAccount := acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
	aclAccount.SetAddress(seller)
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	assetPegHash := []byte("30")
	assetPeg := types.BaseAssetPeg{
		PegHash:       assetPegHash,
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
		Locked:        false,
	}

	order := app.OrderKeeper.NewOrder(buyer, seller, assetPegHash)
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})

	fiatPegHash := []byte("fiatPegHash")
	fiatPeg := types.BaseFiatPeg{
		PegHash:           fiatPegHash,
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})

	app.OrderKeeper.SetOrder(ctx, order)

	buyerExecuteOrder := bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ := app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "negotiation not found.")

	sellerExecuteOrder := bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "negotiation not found.")

	assetPeg.SetPegHash([]byte("31"))
	order.SetAssetPegWallet(types.AddAssetPegToWallet(&assetPeg, order.GetAssetPegWallet()))
	app.OrderKeeper.SetOrder(ctx, order)

	negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()-1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Asset token not found!")

	negotiation = getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()+1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)

	assetPegHash = []byte("30")
	assetPeg.SetPegHash(assetPegHash)
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.NoError(t, err)
}
