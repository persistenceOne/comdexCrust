package keeper_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/exported"
	bankTypes "github.com/commitHub/commitBlockchain/modules/bank/internal/types"
	"github.com/commitHub/commitBlockchain/simApp"
	"github.com/commitHub/commitBlockchain/types"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmtime "github.com/tendermint/tendermint/types/time"
)

type prerequisites func()

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

func TestBaseKeeper_DelegateCoins(t *testing.T) {
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

	type args struct {
		ctx           cTypes.Context
		delegatorAddr cTypes.AccAddress
		moduleAccAddr cTypes.AccAddress
		amt           cTypes.Coins
	}
	arg1 := args{ctx, addr3, addrModule, delCoins}
	arg2 := args{ctx, addr2, addr4, delCoins}
	arg3 := args{ctx, addr2, addrModule, delCoins}
	arg4 := args{ctx, addr1, addrModule, delCoins}

	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Account does not exist.",
			arg1,
			cTypes.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", arg1.delegatorAddr))},

		{"Module account is absent.",
			arg2,
			cTypes.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", arg2.moduleAccAddr))},

		{"Require the ability for a non-vesting account to delegate.",
			arg3,
			nil},

		{"Require the ability for a vesting account to delegate.",
			arg4,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.BankKeeper.DelegateCoins(tt.args.ctx, tt.args.delegatorAddr, tt.args.moduleAccAddr, tt.args.amt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseKeeper.DelegateCoins() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseKeeper_UndelegateCoins(t *testing.T) {
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

	app.BankKeeper.DelegateCoins(ctx, addr2, addrModule, delCoins)

	app.BankKeeper.DelegateCoins(ctx, addr1, addrModule, delCoins)

	type args struct {
		ctx           cTypes.Context
		moduleAccAddr cTypes.AccAddress
		delegatorAddr cTypes.AccAddress
		amt           cTypes.Coins
	}
	arg1 := args{ctx, addrModule, addr2, delCoins}
	arg2 := args{ctx, addrModule, addr1, delCoins}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Require the ability for a non-vesting account to undelegate.",
			arg1,
			nil},
		{"Require the ability for a vesting account to undelegate.",
			arg2,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.BankKeeper.UndelegateCoins(tt.args.ctx, tt.args.moduleAccAddr, tt.args.delegatorAddr, tt.args.amt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseKeeper.UndelegateCoins() = %v, want %v", got, tt.want)
			}
		})
	}
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

	type args struct {
		ctx     cTypes.Context
		inputs  []bankTypes.Input
		outputs []bankTypes.Output
	}
	arg := args{ctx, inputs, outputs}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Adding and subtracting coins.",
			arg,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.BankKeeper.InputOutputCoins(tt.args.ctx, tt.args.inputs, tt.args.outputs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.InputOutputCoins() = %v, want %v", got, tt.want)
			}
		})
	}
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

	type args struct {
		ctx      cTypes.Context
		fromAddr cTypes.AccAddress
		toAddr   cTypes.AccAddress
		amt      cTypes.Coins
	}
	arg1 := args{ctx, addr1, addr2, sendCoin2}
	arg2 := args{ctx, addr1, addr2, sendCoin1}

	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Insufficient coins.",
			arg1,
			cTypes.ErrInsufficientCoins(fmt.Sprintf("insufficient account funds; %s < %s", ak.GetAccount(ctx, arg1.fromAddr).SpendableCoins(ctx.BlockHeader().Time), arg1.amt))},

		{"Sending coins.",
			arg2,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.BankKeeper.SendCoins(tt.args.ctx, tt.args.fromAddr, tt.args.toAddr, tt.args.amt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.SendCoins() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_GetSendEnabled(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	type args struct {
		ctx cTypes.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Send Enabled",
			args{ctx}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.BankKeeper.GetSendEnabled(tt.args.ctx); got != tt.want {
				t.Errorf("BaseSendKeeper.GetSendEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_SetSendEnabled(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	type args struct {
		ctx     cTypes.Context
		enabled bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Set Send Enabled",
			args{ctx, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.BankKeeper.SetSendEnabled(tt.args.ctx, tt.args.enabled)
		})
	}
}

func TestBaseSendKeeper_DefineZones(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	genesisAccount := getGenesisAccount(ctx, app.AccountKeeper)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, addr1))
	zone := cTypes.AccAddress([]byte("zone"))

	defineZone1 := bankTypes.NewDefineZone(addr1, zone, acl.DefaultZoneID)
	defineZone2 := bankTypes.NewDefineZone(genesisAccount.GetAddress(), zone, acl.DefaultZoneID)
	defineZone3 := bankTypes.NewDefineZone(genesisAccount.GetAddress(), addr1, acl.DefaultZoneID)

	type args struct {
		ctx        cTypes.Context
		defineZone bankTypes.DefineZone
	}
	arg1 := args{ctx, defineZone1}
	arg2 := args{ctx, defineZone2}
	arg3 := args{ctx, defineZone3}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"From account is not Genesis.",
			arg1,
			cTypes.ErrInternal(fmt.Sprintf("Account %v is not the genesis account. Zones can only be defined by the genesis account.", arg1.defineZone.From.String()))},

		{"Add Zone",
			arg2,
			nil},

		{"Zone ID already exists",
			arg3,
			cTypes.NewError(cTypes.CodespaceType("acl"), 102, "zone with this given id already exist")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.BankKeeper.DefineZones(tt.args.ctx, tt.args.defineZone); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.DefineZones() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_DefineOrganizations(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	organization := cTypes.AccAddress([]byte("organization"))
	addr1 := cTypes.AccAddress([]byte("addr1"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, addr1))

	defineOrganization1 := bankTypes.NewDefineOrganization(addr1, organization, acl.DefaultOrganizationID, acl.DefaultZoneID)
	defineOrganization2 := bankTypes.NewDefineOrganization(zone, organization, acl.DefaultOrganizationID, acl.DefaultZoneID)
	defineOrganization3 := bankTypes.NewDefineOrganization(zone, addr1, acl.DefaultOrganizationID, acl.DefaultZoneID)

	type args struct {
		ctx                cTypes.Context
		defineOrganization bankTypes.DefineOrganization
	}
	arg1 := args{ctx, defineOrganization1}
	arg2 := args{ctx, defineOrganization2}
	arg3 := args{ctx, defineOrganization3}

	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"From account is not zone.",
			arg1,
			cTypes.ErrInternal(fmt.Sprintf("Account %v is not the zone account. Organizations can only be defined by the zone account.", arg1.defineOrganization.From.String()))},

		{"Add Organization",
			arg2,
			nil},

		{"Organization Id already exists.",
			arg3,
			cTypes.NewError("acl", 102, "organization with given id already exist")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.BankKeeper.DefineOrganizations(tt.args.ctx, tt.args.defineOrganization); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.DefineOrganizations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_DefineACLs(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)
	zoneAccount := app.AccountKeeper.NewAccountWithAddress(ctx, zone)
	app.AccountKeeper.SetAccount(ctx, zoneAccount)

	organization := cTypes.AccAddress([]byte("organization"))
	aclOrganization := acl.Organization{
		Address: organization,
		ZoneID:  acl.DefaultZoneID,
	}
	app.ACLKeeper.SetOrganization(ctx, acl.DefaultOrganizationID, aclOrganization)
	organizationAccount := app.AccountKeeper.NewAccountWithAddress(ctx, organization)
	app.AccountKeeper.SetAccount(ctx, organizationAccount)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, addr1))

	trader := cTypes.AccAddress([]byte("trader"))
	aclAccount := acl.BaseACLAccount{
		Address:        trader,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{IssueAsset: true},
	}

	defineACL1 := bankTypes.NewDefineACL(addr1, trader, &aclAccount)
	defineACL2 := bankTypes.NewDefineACL(zone, trader, &aclAccount)

	type args struct {
		ctx       cTypes.Context
		defineACL bankTypes.DefineACL
	}
	arg1 := args{ctx, defineACL1}
	arg2 := args{ctx, defineACL2}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"From account does not have permission",
			arg1,
			cTypes.ErrInternal(fmt.Sprintf("Account %v does not have access to define acl for account %v.", arg1.defineACL.From.String(), arg1.defineACL.To.String()))},

		{"Define ACL",
			arg2,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.BankKeeper.DefineACLs(tt.args.ctx, tt.args.defineACL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.DefineACLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_IssueAssetsToWallets(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	addr1 := cTypes.AccAddress([]byte("addr1"))

	trader := cTypes.AccAddress([]byte("trader"))
	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
	}

	issueAsset1 := bankTypes.NewIssueAsset(trader, trader, &assetPeg)

	aclAccount := acl.BaseACLAccount{
		Address:        trader,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{IssueAsset: false},
	}

	issueAsset2 := bankTypes.NewIssueAsset(trader, trader, &assetPeg)
	issueAsset3 := bankTypes.NewIssueAsset(trader, trader, &assetPeg)
	issueAsset4 := bankTypes.NewIssueAsset(addr1, trader, &assetPeg)
	issueAsset5 := bankTypes.NewIssueAsset(zone, trader, &assetPeg)

	type args struct {
		ctx        cTypes.Context
		issueAsset bankTypes.IssueAsset
	}
	arg1 := args{ctx, issueAsset1}
	arg2 := args{ctx, issueAsset2}
	arg3 := args{ctx, issueAsset3}
	arg4 := args{ctx, issueAsset4}
	arg5 := args{ctx, issueAsset5}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"ACL account not defined.",
			arg1,
			func() {},
			cTypes.NewError("acl", 102, "acl for this account not defined")},

		{"Moderated is false and IssueAsset of ACL set to false.",
			arg2,
			func() {
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal(fmt.Sprintf("Assets cant be issued to account %v.", arg2.issueAsset.ToAddress.String()))},

		{"Moderated is false and IssueAsset of ACL set to true. Asset is issued.",
			arg3,
			func() {
				aclAccount.SetACL(acl.ACL{IssueAsset: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			nil},

		{"Moderated is true, from address is not zone address.",
			arg4,
			func() {
				assetPeg.SetModerated(true)
			},
			cTypes.ErrInternal("Unauthorised transaction.")},

		{"Moderated is true, from address is not zone address.",
			arg5,
			func() {},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.BankKeeper.IssueAssetsToWallets(tt.args.ctx, tt.args.issueAsset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.IssueAssetsToWallets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_IssueFiatsToWallets(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	addr1 := cTypes.AccAddress([]byte("addr1"))

	trader := cTypes.AccAddress([]byte("trader"))
	aclAccount := acl.BaseACLAccount{
		Address:        trader,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{IssueFiat: false},
	}

	fiatPeg := types.BaseFiatPeg{
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	issueFiat1 := bankTypes.NewIssueFiat(addr1, trader, &fiatPeg)
	issueFiat2 := bankTypes.NewIssueFiat(zone, trader, &fiatPeg)
	issueFiat3 := bankTypes.NewIssueFiat(zone, trader, &fiatPeg)

	type args struct {
		ctx       cTypes.Context
		issueFiat bankTypes.IssueFiat
	}
	arg1 := args{ctx, issueFiat1}
	arg2 := args{ctx, issueFiat2}
	arg3 := args{ctx, issueFiat3}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"ACL Account not set",
			arg1,
			func() {},
			cTypes.ErrInternal("To account acl not defined.")},
		{"Issue Fiat of ACL set to false",
			arg2,
			func() {
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal(fmt.Sprintf("Fiats can't be issued to account %v.", arg2.issueFiat.ToAddress.String()))},
		{"Fiat issued.",
			arg3,
			func() {
				aclAccount.SetACL(acl.ACL{IssueFiat: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.BankKeeper.IssueFiatsToWallets(tt.args.ctx, tt.args.issueFiat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.IssueFiatsToWallets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_ReleaseLockedAssets2(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	ak := app.AccountKeeper

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	trader := cTypes.AccAddress([]byte("trader"))
	traderAccount := ak.NewAccountWithAddress(ctx, trader)
	ak.SetAccount(ctx, traderAccount)
	aclAccount := acl.BaseACLAccount{
		Address:        trader,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{ReleaseAsset: false},
	}

	releaseAsset := bankTypes.NewReleaseAsset(zone, trader, []byte(""))
	err := app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	require.Error(t, err, "To account acl not defined.")

	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	traderAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	ak.SetAccount(ctx, traderAccount)

	releaseAsset = bankTypes.NewReleaseAsset(zone, trader, ak.GetAccount(ctx, trader).GetAssetPegWallet()[0].GetPegHash())
	err = app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	require.Error(t, err, "Assets cannot be released for account "+trader.String()+". Access Denied.")

	aclAccount.SetACL(acl.ACL{ReleaseAsset: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	releaseAsset = bankTypes.NewReleaseAsset(zone, trader, []byte("PegHash"))
	err = app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	require.Error(t, err, "Asset peg not found.")

	releaseAsset = bankTypes.NewReleaseAsset(zone, trader, ak.GetAccount(ctx, trader).GetAssetPegWallet()[0].GetPegHash())
	err = app.BankKeeper.ReleaseLockedAssets(ctx, releaseAsset)
	require.NoError(t, err)

}

func TestBaseSendKeeper_SendAssetsToWallets2(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{SendAsset: false},
	}
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	assetPegHash := assetPeg.GetPegHash()
	sellerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, seller)
	sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.AccountKeeper.SetAccount(ctx, sellerAccount)

	sendAsset := bankTypes.NewSendAsset(buyer, buyer, assetPegHash)
	err := app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "Unauthorised transaction.")

	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.Error(t, err, "Unauthorised transaction.")

	aclAccount.SetACL(acl.ACL{SendAsset: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

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

	assetPeg.SetLocked(false)
	sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.AccountKeeper.SetAccount(ctx, sellerAccount)

	sendAsset = bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	err = app.BankKeeper.SendAssetsToWallets(ctx, sendAsset)
	require.NoError(t, err)
}

func TestBaseSendKeeper_SendFiatsToWallets2(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	seller := cTypes.AccAddress([]byte("seller"))
	assetPegHash := []byte("30")

	buyer := cTypes.AccAddress([]byte("buyer"))
	aclAccount := acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{SendFiat: false},
	}
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	fiatPeg := types.BaseFiatPeg{
		PegHash:           []byte("30"),
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	buyerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, buyer)
	buyerAccount.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.AccountKeeper.SetAccount(ctx, buyerAccount)
	fiatPegHash := fiatPeg.GetPegHash()

	sendFiat := bankTypes.NewSendFiat(addr1, seller, fiatPegHash, 500)
	err := app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.Error(t, err, "Unauthorised transaction.")

	sendFiat = bankTypes.NewSendFiat(buyer, seller, fiatPegHash, 500)
	err = app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)
	require.Error(t, err, "Unauthorised transaction.")

	aclAccount.SetACL(acl.ACL{SendFiat: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

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

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	seller := cTypes.AccAddress([]byte("seller"))
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{BuyerExecuteOrder: false, SellerExecuteOrder: true},
	}
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        true,
	}
	assetPegHash := assetPeg.GetPegHash()

	buyer := cTypes.AccAddress([]byte("buyer"))
	aclAccount.SetAddress(buyer)
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, 10000)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	buyerExecuteOrder := bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ := app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "Asset token not found!")

	order := app.OrderKeeper.NewOrder(buyer, seller, assetPegHash)
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(addr1, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "To account acl not defined.")

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account "+buyer.String()+". Access Denied.")

	aclAccount.SetACL(acl.ACL{BuyerExecuteOrder: true, SellerExecuteOrder: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.NoError(t, err)

	//Unmoderated
	assetPeg = types.BaseAssetPeg{
		PegHash:       []byte("31"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
		Locked:        true,
	}
	assetPegHash = assetPeg.GetPegHash()
	negotiation = getNegotiation(ctx, buyer, seller, assetPegHash, 500, 10000)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
	order = app.OrderKeeper.NewOrder(buyer, seller, assetPegHash)
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	order2 := app.OrderKeeper.NewOrder(addr1, seller, assetPegHash)
	order2.SetAssetPegWallet(order.GetAssetPegWallet())
	app.OrderKeeper.SetOrder(ctx, order2)
	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(buyer, addr1, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "To account acl not defined.")

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account %v."+buyer.String()+" Access Denied.")

	aclAccount.SetACL(acl.ACL{BuyerExecuteOrder: false, SellerExecuteOrder: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account %v."+buyer.String()+" Access Denied.")

	aclAccount.SetACL(acl.ACL{BuyerExecuteOrder: true, SellerExecuteOrder: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.NoError(t, err)

}

func TestBaseSendKeeper_SellerExecuteOrder(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	seller := cTypes.AccAddress([]byte("seller"))
	sellerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, seller)
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{IssueAsset: true, IssueFiat: true, SendAsset: true, SendFiat: true, BuyerExecuteOrder: true, SellerExecuteOrder: false, ChangeBuyerBid: true, ChangeSellerBid: true, ConfirmBuyerBid: true, ConfirmSellerBid: true, Negotiation: true, ReleaseAsset: true, RedeemFiat: true, RedeemAsset: true},
	}
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	assetPeg := types.BaseAssetPeg{
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        false,
	}
	assetPegHash := assetPeg.PegHash
	sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.AccountKeeper.SetAccount(ctx, sellerAccount)

	buyer := cTypes.AccAddress([]byte("buyer"))
	aclAccount.SetAddress(buyer)
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	fiatPeg := types.BaseFiatPeg{
		PegHash:           []byte("30"),
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	fiatPegHash := fiatPeg.GetPegHash()
	buyerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, buyer)
	buyerAccount.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.AccountKeeper.SetAccount(ctx, buyerAccount)

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
	aclAccount.SetACL(acl.ACL{SellerExecuteOrder: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	sendFiat := bankTypes.NewSendFiat(buyer, seller, fiatPegHash, negotiation.GetBid())
	app.BankKeeper.SendFiatsToWallets(ctx, sendFiat)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	//Unmoderated
	seller2 := cTypes.AccAddress([]byte("seller2"))
	assetPegHash = []byte("30")
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
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{SellerExecuteOrder: false},
	}
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(addr1, buyer, seller2, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account "+addr1.String()+". Access Denied.")

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller2, buyer, seller2, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Trade cannot be executed for account "+seller2.String()+". Access Denied.")

	aclAccount.SetACL(acl.ACL{SellerExecuteOrder: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller2, buyer, seller2, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

}

func TestBaseSendKeeper_RedeemAsset(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	buyer := cTypes.AccAddress([]byte("buyer"))
	seller := cTypes.AccAddress([]byte("seller"))
	sellerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, seller)
	aclAccount := acl.BaseACLAccount{
		Address:        seller,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{RedeemAsset: false},
	}
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	assetPeg := types.BaseAssetPeg{
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        false,
	}

	assetPegHash := assetPeg.PegHash

	redeemAsset := bankTypes.NewRedeemAsset(buyer, seller, assetPegHash)
	err := app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.Error(t, err, "To account acl not defined.")

	redeemAsset = bankTypes.NewRedeemAsset(zone, seller, assetPegHash)
	err = app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.Error(t, err, "Assets can't be redeemed from account "+seller.String()+".")

	aclAccount.SetACL(acl.ACL{RedeemAsset: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	redeemAsset = bankTypes.NewRedeemAsset(zone, seller, assetPegHash)
	err = app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.Error(t, err, "No Assets Found!")

	sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.AccountKeeper.SetAccount(ctx, sellerAccount)

	redeemAsset = bankTypes.NewRedeemAsset(zone, seller, []byte("PegHash"))
	err = app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.Error(t, err, "No Assets With Given PegHash Found!")

	redeemAsset = bankTypes.NewRedeemAsset(zone, seller, assetPegHash)
	err = app.BankKeeper.RedeemAssetsFromWallets(ctx, redeemAsset)
	require.NoError(t, err)
}

func TestBaseSendKeeper_RedeemFiat(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	addr1 := cTypes.AccAddress([]byte("addr1"))
	buyer := cTypes.AccAddress([]byte("buyer"))
	buyerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, buyer)
	aclAccount := acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{RedeemFiat: false},
	}
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	fiatPeg := types.BaseFiatPeg{
		PegHash:           []byte("30"),
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}

	redeemFiat := bankTypes.NewRedeemFiat(addr1, zone, 500)
	err := app.BankKeeper.RedeemFiatsFromWallets(ctx, redeemFiat)
	require.Error(t, err, "To account acl not defined.")

	redeemFiat = bankTypes.NewRedeemFiat(buyer, zone, 500)
	err = app.BankKeeper.RedeemFiatsFromWallets(ctx, redeemFiat)
	require.Error(t, err, "Fiats can't be redeemed from account "+buyer.String()+".")

	aclAccount.SetACL(acl.ACL{RedeemFiat: true})
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	buyerAccount.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.AccountKeeper.SetAccount(ctx, buyerAccount)

	redeemFiat = bankTypes.NewRedeemFiat(buyer, zone, 50000)
	err = app.BankKeeper.RedeemFiatsFromWallets(ctx, redeemFiat)
	require.Error(t, err, "Redeemed amount higher than the account balance")

	redeemFiat = bankTypes.NewRedeemFiat(buyer, zone, 500)
	err = app.BankKeeper.RedeemFiatsFromWallets(ctx, redeemFiat)
	require.NoError(t, err)
}

func TestBaseSendKeeper_PrivateExchangeOrderTokens(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress([]byte("buyer"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, buyer))
	seller := cTypes.AccAddress([]byte("seller"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, seller))
	aclAccount := acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{BuyerExecuteOrder: true, SellerExecuteOrder: true},
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
	require.NoError(t, err)

	assetPegHash = []byte("30")
	assetPeg.SetPegHash(assetPegHash)
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	//Above test reverses the order
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

func TestBaseSendKeeper_ExchangeOrderTokens(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zone := cTypes.AccAddress([]byte("zone"))
	app.ACLKeeper.SetZoneAddress(ctx, acl.DefaultZoneID, zone)

	buyer := cTypes.AccAddress([]byte("buyer"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, buyer))
	seller := cTypes.AccAddress([]byte("seller"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, seller))
	aclAccount := acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{BuyerExecuteOrder: true, SellerExecuteOrder: true},
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
		Moderated:     true,
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
	app.OrderKeeper.SetOrder(ctx, order)

	buyerExecuteOrder := bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ := app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.Error(t, err, "negotiation not found.")

	sellerExecuteOrder := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "negotiation not found.")

	negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()-1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Fiat tokens not found!")

	order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	assetPeg.SetPegHash([]byte("31"))
	order.SetAssetPegWallet(types.AddAssetPegToWallet(&assetPeg, order.GetAssetPegWallet()))
	app.OrderKeeper.SetOrder(ctx, order)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.Error(t, err, "Asset token not found!")

	negotiation = getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()+1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	assetPegHash = []byte("30")
	assetPeg.SetPegHash(assetPegHash)
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	//Above test reverses the order
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	negotiation.SetBuyerBlockHeight(ctx.BlockHeight() - 1)
	negotiation.SetTime(ctx.BlockHeight() - 1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	//Above test reverses the order
	order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
	order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
	app.OrderKeeper.SetOrder(ctx, order)

	negotiation.SetBuyerBlockHeight(ctx.BlockHeight() + 1)
	negotiation.SetTime(ctx.BlockHeight() + 1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	sellerExecuteOrder = bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	err, _ = app.BankKeeper.SellerExecuteTradeOrder(ctx, sellerExecuteOrder)
	require.NoError(t, err)

	buyerExecuteOrder = bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	err, _ = app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	require.NoError(t, err)
}
