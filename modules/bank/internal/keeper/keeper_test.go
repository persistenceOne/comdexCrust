package keeper_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/persistenceOne/persistenceSDK/modules/acl"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/auth/exported"
	bankTypes "github.com/persistenceOne/persistenceSDK/modules/bank/internal/types"
	"github.com/persistenceOne/persistenceSDK/simApp"
	"github.com/persistenceOne/persistenceSDK/types"
	cTypes "github.com/cosmos/cosmos-sdk/types"
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

func TestBaseSendKeeper_ReleaseLockedAssets(t *testing.T) {
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
	releaseAsset1 := bankTypes.NewReleaseAsset(zone, trader, assetPeg.GetPegHash())
	releaseAsset2 := bankTypes.NewReleaseAsset(zone, trader, assetPeg.GetPegHash())
	releaseAsset3 := bankTypes.NewReleaseAsset(zone, trader, []byte("PegHash"))
	releaseAsset4 := bankTypes.NewReleaseAsset(zone, trader, assetPeg.GetPegHash())

	type args struct {
		ctx          cTypes.Context
		releaseAsset bankTypes.ReleaseAsset
	}
	arg1 := args{ctx, releaseAsset1}
	arg2 := args{ctx, releaseAsset2}
	arg3 := args{ctx, releaseAsset3}
	arg4 := args{ctx, releaseAsset4}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"ACL not defined.",
			arg1,
			func() {},
			cTypes.ErrInternal("To account acl not defined.")},
		{"Relase Asset of ACL set to false.",
			arg2,
			func() {
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
				traderAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				ak.SetAccount(ctx, traderAccount)
			},
			cTypes.ErrInternal(fmt.Sprintf("Assets cannot be released for account %v. Access Denied.", arg2.releaseAsset.OwnerAddress.String()))},
		{"Invalid pegHash is given.",
			arg3,
			func() {
				aclAccount.SetACL(acl.ACL{ReleaseAsset: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal("Asset peg not found.")},
		{"Locked status of asset is set to true.",
			arg4,
			func() {},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.BankKeeper.ReleaseLockedAssets(tt.args.ctx, tt.args.releaseAsset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.ReleaseLockedAssets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_SendAssetsToWallets(t *testing.T) {
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

	negotiationID := types.NegotiationID(append(append(buyer.Bytes(), seller.Bytes()...), assetPegHash.Bytes()...))
	var sellerSignature types.Signature

	sendAsset1 := bankTypes.NewSendAsset(buyer, buyer, assetPegHash)
	sendAsset2 := bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	sendAsset3 := bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	sendAsset4 := bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	sendAsset5 := bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	sendAsset6 := bankTypes.NewSendAsset(seller, buyer, []byte("WrongPegHash"))
	sendAsset7 := bankTypes.NewSendAsset(seller, buyer, assetPegHash)
	sendAsset8 := bankTypes.NewSendAsset(seller, buyer, assetPegHash)

	type args struct {
		ctx       cTypes.Context
		sendAsset bankTypes.SendAsset
	}
	arg1 := args{ctx, sendAsset1}
	arg2 := args{ctx, sendAsset2}
	arg3 := args{ctx, sendAsset3}
	arg4 := args{ctx, sendAsset4}
	arg5 := args{ctx, sendAsset5}
	arg6 := args{ctx, sendAsset6}
	arg7 := args{ctx, sendAsset7}
	arg8 := args{ctx, sendAsset8}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"ACL not defined.",
			arg1,
			func() {},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"SendAsset of ACL set to false.",
			arg2,
			func() {},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"Negotiation not found.",
			arg3,
			func() {
				aclAccount.SetACL(acl.ACL{SendAsset: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.NewError("negotiation", 600, "negotiation not found.")},
		{"Negotiation time expired.",
			arg4,
			func() {
				negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()-1)
				sellerSignature = negotiation.GetSellerSignature()
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			cTypes.ErrInvalidSequence("Negotiation time expired.")},
		{"Signatures not present",
			arg5,
			func() {
				negotiation, _ := app.NegotiationKeeper.GetNegotiation(ctx, negotiationID)
				negotiation.SetTime(ctx.BlockHeight() + 1)
				negotiation.SetSellerSignature(nil)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			cTypes.ErrInternal("Signatures are not present")},
		{"Invalid Asset PegHash",
			arg6,
			func() {
				negotiation, _ := app.NegotiationKeeper.GetNegotiation(ctx, negotiationID)
				negotiation.SetSellerSignature(sellerSignature)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
				app.NegotiationKeeper.SetNegotiation(ctx, getNegotiation(ctx, buyer, seller, []byte("WrongPegHash"), 500, ctx.BlockHeight()+1))

			},
			cTypes.ErrInsufficientCoins("Asset not found.")},
		{"Asset is locked.",
			arg7,
			func() {},
			cTypes.ErrInsufficientCoins("Asset locked.")},
		{"Sending asset to order.",
			arg8,
			func() {
				assetPeg.SetLocked(false)
				sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.AccountKeeper.SetAccount(ctx, sellerAccount)
			},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.BankKeeper.SendAssetsToWallets(tt.args.ctx, tt.args.sendAsset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.SendAssetsToWallets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_SendFiatsToWallets(t *testing.T) {
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
	var buyerSignature types.Signature
	var negotiationID types.NegotiationID

	sendFiat1 := bankTypes.NewSendFiat(addr1, seller, assetPegHash, 500)
	sendFiat2 := bankTypes.NewSendFiat(buyer, seller, assetPegHash, 500)
	sendFiat3 := bankTypes.NewSendFiat(buyer, seller, assetPegHash, 500)
	sendFiat4 := bankTypes.NewSendFiat(buyer, seller, assetPegHash, 500)
	sendFiat5 := bankTypes.NewSendFiat(buyer, seller, assetPegHash, 500)
	sendFiat6 := bankTypes.NewSendFiat(buyer, seller, assetPegHash, 5000000)
	sendFiat7 := bankTypes.NewSendFiat(buyer, seller, assetPegHash, 500)

	type args struct {
		ctx      cTypes.Context
		sendFiat bankTypes.SendFiat
	}
	arg1 := args{ctx, sendFiat1}
	arg2 := args{ctx, sendFiat2}
	arg3 := args{ctx, sendFiat3}
	arg4 := args{ctx, sendFiat4}
	arg5 := args{ctx, sendFiat5}
	arg6 := args{ctx, sendFiat6}
	arg7 := args{ctx, sendFiat7}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"ACL not defined.",
			arg1,
			func() {},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"Send Fiat of ACL set to false.",
			arg2,
			func() {},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"Negotiation is not defined.",
			arg3,
			func() {
				aclAccount.SetACL(acl.ACL{SendFiat: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.NewError("negotiation", 600, "negotiation not found.")},
		{"Negotiation time expired.",
			arg4,
			func() {
				negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()-1)
				negotiationID = negotiation.GetNegotiationID()
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			cTypes.ErrInvalidSequence("Negotiation time expired.")},
		{"Signatures are not present.",
			arg5, func() {
				negotiation, _ := app.NegotiationKeeper.GetNegotiation(ctx, negotiationID)
				negotiation.SetTime(ctx.BlockHeight() + 1)
				buyerSignature = negotiation.GetBuyerSignature()
				negotiation.SetBuyerSignature(nil)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			cTypes.ErrInternal("Signatures are not present")},
		{"Sending amount more than  present.",
			arg6,
			func() {
				negotiation, _ := app.NegotiationKeeper.GetNegotiation(ctx, negotiationID)
				negotiation.SetBuyerSignature(buyerSignature)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			cTypes.ErrInsufficientCoins(fmt.Sprintf("Insufficient funds"))},
		{"Sending fiat to order.",
			arg7,
			func() {},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.BankKeeper.SendFiatsToWallets(tt.args.ctx, tt.args.sendFiat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.SendFiatsToWallets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_RedeemAssetsFromWallets(t *testing.T) {
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

	assetPegHash := assetPeg.GetPegHash()

	redeemAsset1 := bankTypes.NewRedeemAsset(buyer, seller, assetPegHash)
	redeemAsset2 := bankTypes.NewRedeemAsset(zone, seller, assetPegHash)
	redeemAsset3 := bankTypes.NewRedeemAsset(zone, seller, assetPegHash)
	redeemAsset4 := bankTypes.NewRedeemAsset(zone, seller, []byte("PegHash"))
	redeemAsset5 := bankTypes.NewRedeemAsset(zone, seller, assetPegHash)

	type args struct {
		ctx         cTypes.Context
		redeemAsset bankTypes.RedeemAsset
	}
	arg1 := args{ctx, redeemAsset1}
	arg2 := args{ctx, redeemAsset2}
	arg3 := args{ctx, redeemAsset3}
	arg4 := args{ctx, redeemAsset4}
	arg5 := args{ctx, redeemAsset5}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"From address is not zone.",
			arg1,
			func() {},
			cTypes.ErrInternal("Unauthorised transaction.")},
		{"Redeem Asset of ACL set to false.",
			arg2,
			func() {},
			cTypes.ErrInternal(fmt.Sprintf("Assets can't be redeemed from account %v.", arg2.redeemAsset.RedeemerAddress.String()))},
		{"No asset is present.",
			arg3,
			func() {
				aclAccount.SetACL(acl.ACL{RedeemAsset: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal("No Assets Found!")},
		{"Peg hash not present.",
			arg4,
			func() {
				sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.AccountKeeper.SetAccount(ctx, sellerAccount)
			},
			cTypes.ErrInternal("No Assets With Given PegHash Found!")},
		{"Asset is redeemed.",
			arg5,
			func() {},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.BankKeeper.RedeemAssetsFromWallets(tt.args.ctx, tt.args.redeemAsset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.RedeemAssetsFromWallets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_RedeemFiatsFromWallets(t *testing.T) {
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

	redeemFiat1 := bankTypes.NewRedeemFiat(addr1, zone, 500)
	redeemFiat2 := bankTypes.NewRedeemFiat(buyer, zone, 500)
	redeemFiat3 := bankTypes.NewRedeemFiat(buyer, zone, 50000)
	redeemFiat4 := bankTypes.NewRedeemFiat(buyer, zone, 500)

	type args struct {
		ctx        cTypes.Context
		redeemFiat bankTypes.RedeemFiat
	}
	arg1 := args{ctx, redeemFiat1}
	arg2 := args{ctx, redeemFiat2}
	arg3 := args{ctx, redeemFiat3}
	arg4 := args{ctx, redeemFiat4}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"ACL of account is not defined.",
			arg1,
			func() {},
			cTypes.ErrInternal("To account acl not defined.")},
		{"Send Fiat of ACL set to false.",
			arg2,
			func() {},
			cTypes.ErrInternal(fmt.Sprintf("Fiats can't be redeemed from account %v.", arg2.redeemFiat.RedeemerAddress.String()))},
		{"Redeeming amount is higher than the net account balance.",
			arg3,
			func() {
				aclAccount.SetACL(acl.ACL{RedeemFiat: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
				buyerAccount.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
				app.AccountKeeper.SetAccount(ctx, buyerAccount)
			},
			cTypes.ErrInsufficientCoins(fmt.Sprintf("Redeemed amount higher than the account balance"))},
		{"Redeeming amount from wallet.",
			arg4,
			func() {},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.BankKeeper.RedeemFiatsFromWallets(tt.args.ctx, tt.args.redeemFiat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.RedeemFiatsFromWallets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseSendKeeper_BuyerExecuteTradeOrder(t *testing.T) {

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

	buyerExecuteOrder1 := bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	buyerExecuteOrder2 := bankTypes.NewBuyerExecuteOrder(addr1, buyer, seller, assetPegHash, "fiatProofHash")
	buyerExecuteOrder3 := bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
	buyerExecuteOrder4 := bankTypes.NewBuyerExecuteOrder(zone, seller, buyer, assetPegHash, "fiatProofHash")
	buyerExecuteOrder5 := bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")

	//Unmoderated
	assetPeg2 := types.BaseAssetPeg{
		PegHash:       []byte("31"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
		Locked:        true,
	}
	assetPegHash2 := assetPeg2.GetPegHash()
	negotiation2 := getNegotiation(ctx, buyer, seller, assetPegHash2, 500, ctx.BlockHeight()+1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation2)

	buyerExecuteOrder6 := bankTypes.NewBuyerExecuteOrder(buyer, addr1, seller, assetPegHash2, "fiatProofHash")
	buyerExecuteOrder7 := bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash2, "fiatProofHash")
	buyerExecuteOrder8 := bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash2, "fiatProofHash")
	buyerExecuteOrder9 := bankTypes.NewBuyerExecuteOrder(seller, seller, buyer, assetPegHash2, "fiatProofHash")
	buyerExecuteOrder10 := bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash2, "fiatProofHash")

	var emptyFiatPegWallets []types.FiatPegWallet
	// onSuccessEmptyFiatPegWallets := append(emptyFiatPegWallets, types.FiatPegWallet{})

	type args struct {
		ctx               cTypes.Context
		buyerExecuteOrder bankTypes.BuyerExecuteOrder
	}
	arg1 := args{ctx, buyerExecuteOrder1}
	arg2 := args{ctx, buyerExecuteOrder2}
	arg3 := args{ctx, buyerExecuteOrder3}
	arg4 := args{ctx, buyerExecuteOrder4}
	arg5 := args{ctx, buyerExecuteOrder5}
	arg6 := args{ctx, buyerExecuteOrder6}
	arg7 := args{ctx, buyerExecuteOrder7}
	arg8 := args{ctx, buyerExecuteOrder8}
	arg9 := args{ctx, buyerExecuteOrder9}
	arg10 := args{ctx, buyerExecuteOrder10}
	tests := []struct {
		name  string
		args  args
		pre   prerequisites
		want  cTypes.Error
		want1 []types.FiatPegWallet
	}{
		{"Moderated Asset: Asset not present in the order.",
			arg1,
			func() {},
			cTypes.ErrInsufficientCoins("Asset token not found!"),
			emptyFiatPegWallets},
		{"Moderated Asset: From account is not zone.",
			arg2,
			func() {
				order := app.OrderKeeper.NewOrder(buyer, seller, assetPegHash)
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.OrderKeeper.SetOrder(ctx, order)
			},
			cTypes.ErrInternal("Unauthorised transaction."),
			emptyFiatPegWallets},
		{"Moderated Asset: BuyerExecuteOrder of buyer's ACL account set to false.",
			arg3,
			func() {},
			cTypes.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", arg3.buyerExecuteOrder.BuyerAddress.String())),
			emptyFiatPegWallets},
		{"Moderated Asset: Order fails due to negotiation not found",
			arg4,
			func() {
				aclAccount2 := aclAccount
				aclAccount2.SetAddress(seller)
				aclAccount2.SetACL(acl.ACL{BuyerExecuteOrder: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount2)
				order := app.OrderKeeper.NewOrder(seller, buyer, assetPegHash)
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.OrderKeeper.SetOrder(ctx, order)
			},
			cTypes.NewError("negotiation", 600, "negotiation not found."),
			emptyFiatPegWallets},
		{"Moderated Asset: ACL is defined to true and from address is zone, executes the order.",
			arg5,
			func() {
				aclAccount.SetACL(acl.ACL{BuyerExecuteOrder: true, SellerExecuteOrder: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			nil,
			append(emptyFiatPegWallets, nil),
		},
		{"Unmoderated Asset: Buyer account is different and acl not defined for it.",
			arg6,
			func() {
				order2 := app.OrderKeeper.NewOrder(addr1, seller, assetPegHash2)
				order2.SetAssetPegWallet(types.AssetPegWallet{assetPeg2})
				app.OrderKeeper.SetOrder(ctx, order2)
			},
			cTypes.NewError("acl", 102, "acl for this account not defined"),
			emptyFiatPegWallets},
		{"Unmoderated Asset: Zone is the from address",
			arg7,
			func() {
				order2 := app.OrderKeeper.NewOrder(buyer, seller, assetPegHash2)
				order2.SetAssetPegWallet(types.AssetPegWallet{assetPeg2})
				app.OrderKeeper.SetOrder(ctx, order2)
			},
			cTypes.ErrUnauthorized(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", arg7.buyerExecuteOrder.MediatorAddress.String())),
			emptyFiatPegWallets},
		{"Unmoderated Asset: BuyerExecuteOrder of buyer's ACL account set to false.",
			arg8,
			func() {
				aclAccount.SetACL(acl.ACL{BuyerExecuteOrder: false, SellerExecuteOrder: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", arg8.buyerExecuteOrder.BuyerAddress.String())),
			emptyFiatPegWallets},
		{"Unoderated Asset: Order fails due to negotiation not found",
			arg9,
			func() {
				aclAccount2 := aclAccount
				aclAccount2.SetAddress(seller)
				aclAccount2.SetACL(acl.ACL{BuyerExecuteOrder: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount2)
				order := app.OrderKeeper.NewOrder(seller, buyer, assetPegHash2)
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg2})
				app.OrderKeeper.SetOrder(ctx, order)
			},
			cTypes.NewError("negotiation", 600, "negotiation not found."),
			emptyFiatPegWallets},
		{"Unmoderated Asset: ACL is defined to true and from address is buyer, executes the order.",
			arg10,
			func() {
				aclAccount.SetACL(acl.ACL{BuyerExecuteOrder: true, SellerExecuteOrder: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			nil,
			append(emptyFiatPegWallets, nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			got, got1 := app.BankKeeper.BuyerExecuteTradeOrder(tt.args.ctx, tt.args.buyerExecuteOrder)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.BuyerExecuteTradeOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("BaseSendKeeper.BuyerExecuteTradeOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBaseSendKeeper_SellerExecuteTradeOrder(t *testing.T) {
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
		PegHash:       []byte("30"),
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     true,
		Locked:        false,
	}
	assetPegHash := assetPeg.GetPegHash()
	app.AccountKeeper.SetAccount(ctx, sellerAccount)

	buyer := cTypes.AccAddress([]byte("buyer"))
	aclAccount.SetAddress(buyer)
	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	fiatPeg := types.BaseFiatPeg{
		PegHash:           []byte("30"),
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	buyerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, buyer)
	app.AccountKeeper.SetAccount(ctx, buyerAccount)

	negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, 10000)
	negotiationID := negotiation.GetNegotiationID()
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	var emptyAssetPegWallets []types.AssetPegWallet

	sellerExecuteOrder1 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder2 := bankTypes.NewSellerExecuteOrder(addr1, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder3 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder4 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")

	//Unmoderated
	seller2 := cTypes.AccAddress([]byte("seller2"))
	assetPegHash2 := []byte("31")
	assetPeg2 := types.BaseAssetPeg{
		PegHash:       assetPegHash2,
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
		Locked:        false,
	}
	negotiation2 := getNegotiation(ctx, buyer, seller2, assetPegHash2, 500, ctx.BlockHeight()+1)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation2)

	aclAccount2 := acl.BaseACLAccount{
		Address:        seller2,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{SellerExecuteOrder: false},
	}

	sellerExecuteOrder5 := bankTypes.NewSellerExecuteOrder(seller, buyer, seller2, assetPegHash2, "awbProofHash")
	sellerExecuteOrder6 := bankTypes.NewSellerExecuteOrder(addr1, buyer, seller2, assetPegHash2, "awbProofHash")
	sellerExecuteOrder7 := bankTypes.NewSellerExecuteOrder(seller2, buyer, seller2, assetPegHash2, "awbProofHash")
	sellerExecuteOrder8 := bankTypes.NewSellerExecuteOrder(seller2, buyer, seller2, assetPegHash2, "awbProofHash")

	type args struct {
		ctx                cTypes.Context
		sellerExecuteOrder bankTypes.SellerExecuteOrder
	}
	arg1 := args{ctx, sellerExecuteOrder1}
	arg2 := args{ctx, sellerExecuteOrder2}
	arg3 := args{ctx, sellerExecuteOrder3}
	arg4 := args{ctx, sellerExecuteOrder4}
	arg5 := args{ctx, sellerExecuteOrder5}
	arg6 := args{ctx, sellerExecuteOrder6}
	arg7 := args{ctx, sellerExecuteOrder7}
	arg8 := args{ctx, sellerExecuteOrder8}
	tests := []struct {
		name  string
		args  args
		pre   prerequisites
		want  cTypes.Error
		want1 []types.AssetPegWallet
	}{
		{"Moderated Asset: Asset not present in the order.",
			arg1,
			func() {},
			cTypes.ErrInsufficientCoins("Asset token not found!"),
			emptyAssetPegWallets},
		{"Moderated Asset: From address is not zone.",
			arg2,
			func() {
				order := app.OrderKeeper.NewOrder(buyer, seller, assetPegHash)
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.OrderKeeper.SetOrder(ctx, order)
			},
			cTypes.ErrInternal("Unauthorised transaction."),
			emptyAssetPegWallets},
		{"Moderated Asset: SellerExecuteOrder of buyer's ACL account set to false.",
			arg3,
			func() {},
			cTypes.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", arg3.sellerExecuteOrder.SellerAddress.String())),
			emptyAssetPegWallets},
		{"Moderated Asset: ACL is defined to true and from address is zone, executes the order.",
			arg4,
			func() {
				aclAccount.SetAddress(seller)
				aclAccount.SetACL(acl.ACL{SellerExecuteOrder: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

				order := app.OrderKeeper.GetOrder(ctx, negotiationID)
				order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
				app.OrderKeeper.SetOrder(ctx, order)
			},
			nil,
			append(emptyAssetPegWallets, types.AssetPegWallet{assetPeg}),
		},
		{"Unmoderated Asset: Seller account is different and acl not defined for it.",
			arg5,
			func() {
				order2 := app.OrderKeeper.NewOrder(buyer, seller2, assetPegHash2)
				order2.SetAssetPegWallet(types.AssetPegWallet{assetPeg2})
				app.OrderKeeper.SetOrder(ctx, order2)
			},
			cTypes.NewError("acl", 102, "acl for this account not defined"),
			emptyAssetPegWallets},
		{"Unmoderated Asset: Unknown account is the from address.",
			arg6,
			func() {
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount2)
			},
			cTypes.ErrUnauthorized(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", arg6.sellerExecuteOrder.MediatorAddress.String())),
			emptyAssetPegWallets},
		{"Unmoderated Asset: SellerExecuteOrder of seller's ACL account set to false.",
			arg7,
			func() {},
			cTypes.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", arg7.sellerExecuteOrder.SellerAddress.String())),
			emptyAssetPegWallets},
		{"Unmoderated Asset: ACL is defined to true, executes the order.",
			arg8,
			func() {
				aclAccount2.SetACL(acl.ACL{SellerExecuteOrder: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount2)
			},
			nil,
			append(emptyAssetPegWallets, types.AssetPegWallet{assetPeg2})},
	}
	for _, tt := range tests {
		tt.pre()
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := app.BankKeeper.SellerExecuteTradeOrder(tt.args.ctx, tt.args.sellerExecuteOrder)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.SellerExecuteTradeOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("BaseSendKeeper.SellerExecuteTradeOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
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

	orderID := order.GetNegotiationID()

	fiatPegHash := []byte("fiatPegHash")
	fiatPeg := types.BaseFiatPeg{
		PegHash:           fiatPegHash,
		TransactionID:     "DEFAULT",
		TransactionAmount: 1000,
	}
	app.OrderKeeper.SetOrder(ctx, order)

	sellerExecuteOrder1 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder2 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder3 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder4 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder5 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder6 := bankTypes.NewSellerExecuteOrder(zone, buyer, seller, assetPegHash, "awbProofHash")

	type args struct {
		ctx                cTypes.Context
		sellerExecuteOrder bankTypes.SellerExecuteOrder
	}
	arg1 := args{ctx, sellerExecuteOrder1}
	arg2 := args{ctx, sellerExecuteOrder2}
	arg3 := args{ctx, sellerExecuteOrder3}
	arg4 := args{ctx, sellerExecuteOrder4}
	arg5 := args{ctx, sellerExecuteOrder5}
	arg6 := args{ctx, sellerExecuteOrder6}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"Negotiation not present.",
			arg1,
			func() {},
			cTypes.NewError("negotiation", 600, "negotiation not found.")},
		{"Fiat peg not present.",
			arg2,
			func() {
				negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()-1)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			cTypes.ErrInsufficientCoins("Fiat tokens not found!")},
		{"Asset peg not found.",
			arg3,
			func() {
				order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
				app.OrderKeeper.SetOrder(ctx, order)
				assetPeg.SetPegHash([]byte("31"))
				order.SetAssetPegWallet(types.AddAssetPegToWallet(&assetPeg, order.GetAssetPegWallet()))
				app.OrderKeeper.SetOrder(ctx, order)
			},
			cTypes.ErrInsufficientCoins("Asset token not found!")},
		{"Order is reversed as more than one asset in order.",
			arg4,
			func() {
				negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()+1)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			nil},
		{"Order is reversed as negotiation time expired.",
			arg5,
			func() {
				assetPegHash = []byte("30")
				assetPeg.SetPegHash(assetPegHash)
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.OrderKeeper.SetOrder(ctx, order)

				//Above test reverses the order
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
				app.OrderKeeper.SetOrder(ctx, order)

				negotiation, _ := app.NegotiationKeeper.GetNegotiation(ctx, orderID)
				negotiation.SetBuyerBlockHeight(ctx.BlockHeight() - 1)
				negotiation.SetTime(ctx.BlockHeight() - 1)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			nil},
		{"Order is executed successfully.",
			arg6,
			func() {
				//Above test reverses the order
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
				app.OrderKeeper.SetOrder(ctx, order)

				negotiation, _ := app.NegotiationKeeper.GetNegotiation(ctx, orderID)
				negotiation.SetBuyerBlockHeight(ctx.BlockHeight() + 1)
				negotiation.SetTime(ctx.BlockHeight() + 1)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

				buyerExecuteOrder := bankTypes.NewBuyerExecuteOrder(zone, buyer, seller, assetPegHash, "fiatProofHash")
				app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
			},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			got, _ := app.BankKeeper.SellerExecuteTradeOrder(tt.args.ctx, tt.args.sellerExecuteOrder)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.SellerExecuteTradeOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
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

	// buyerExecuteOrder := bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash, "fiatProofHash")
	// err, _ := app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
	// require.Error(t, err, "negotiation not found.")

	sellerExecuteOrder1 := bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder2 := bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder3 := bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "awbProofHash")
	sellerExecuteOrder4 := bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "")
	sellerExecuteOrder5 := bankTypes.NewSellerExecuteOrder(seller, buyer, seller, assetPegHash, "awbProofHash")

	type args struct {
		ctx                cTypes.Context
		sellerExecuteOrder bankTypes.SellerExecuteOrder
	}
	arg1 := args{ctx, sellerExecuteOrder1}
	arg2 := args{ctx, sellerExecuteOrder2}
	arg3 := args{ctx, sellerExecuteOrder3}
	arg4 := args{ctx, sellerExecuteOrder4}
	arg5 := args{ctx, sellerExecuteOrder5}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"Negotiation not present.",
			arg1,
			func() {},
			cTypes.NewError("negotiation", 600, "negotiation not found.")},
		{"Asset peg not found.",
			arg2,
			func() {
				assetPeg.SetPegHash([]byte("31"))
				order.SetAssetPegWallet(types.AddAssetPegToWallet(&assetPeg, order.GetAssetPegWallet()))
				app.OrderKeeper.SetOrder(ctx, order)

				negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()-1)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			cTypes.ErrInsufficientCoins("Asset token not found!")},
		{"Order is reversed as more than one asset in order.",
			arg3,
			func() {
				negotiation := getNegotiation(ctx, buyer, seller, assetPegHash, 500, ctx.BlockHeight()+1)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			nil},
		{"Order is reversed as no awbProofHash.",
			arg4,
			func() {
				//Above test reverses the order
				assetPegHash = []byte("30")
				assetPeg.SetPegHash(assetPegHash)
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.OrderKeeper.SetOrder(ctx, order)
			},
			nil},
		{"Order is executed successfully.",
			arg5,
			func() {
				//Above test reverses the order
				order.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				order.SetFiatPegWallet(types.FiatPegWallet{fiatPeg})
				app.OrderKeeper.SetOrder(ctx, order)
				buyerExecuteOrder := bankTypes.NewBuyerExecuteOrder(buyer, buyer, seller, assetPegHash, "fiatProofHash")
				app.BankKeeper.BuyerExecuteTradeOrder(ctx, buyerExecuteOrder)
			},
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			got, _ := app.BankKeeper.SellerExecuteTradeOrder(tt.args.ctx, tt.args.sellerExecuteOrder)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseSendKeeper.SellerExecuteTradeOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
