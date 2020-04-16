package keeper_test

import (
	"reflect"
	"testing"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	aclTypes "github.com/persistenceOne/comdexCrust/modules/acl/internal/types"
	"github.com/persistenceOne/comdexCrust/modules/auth"
	"github.com/persistenceOne/comdexCrust/modules/auth/exported"
	"github.com/persistenceOne/comdexCrust/simApp"
)

type prerequiste func()

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

func TestKeeper_DefineZoneAddress(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	addr2 := cTypes.AccAddress([]byte("addr2"))
	zoneID := []byte("ABCDEF1234")
	type args struct {
		ctx    cTypes.Context
		to     cTypes.AccAddress
		zoneID aclTypes.ZoneID
	}
	arg1 := args{ctx, addr1, zoneID}
	arg2 := args{ctx, addr2, zoneID}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Setting Zone.",
			arg1,
			nil},
		{"Setting Zone again.",
			arg2,
			aclTypes.ErrInvalidID(aclTypes.DefaultCodeSpace, "zone with this given id already exist")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ACLKeeper.DefineZoneAddress(tt.args.ctx, tt.args.to, tt.args.zoneID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.DefineZoneAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_DefineOrganizationAddress(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	addr2 := cTypes.AccAddress([]byte("addr2"))
	zoneID := []byte("ABCDEF1234")
	orgID := []byte("ABC123456")
	type args struct {
		ctx            cTypes.Context
		to             cTypes.AccAddress
		organizationID aclTypes.OrganizationID
		zoneID         aclTypes.ZoneID
	}
	arg1 := args{ctx, addr1, orgID, zoneID}
	arg2 := args{ctx, addr2, orgID, zoneID}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Setting Organization.",
			arg1,
			nil},
		{"Setting Organization again.",
			arg2,
			aclTypes.ErrInvalidID(aclTypes.DefaultCodeSpace, "organization with given id already exist")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ACLKeeper.DefineOrganizationAddress(tt.args.ctx, tt.args.to, tt.args.organizationID, tt.args.zoneID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.DefineOrganizationAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_DefineACLAccount(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	zoneID := []byte("ABCDEF1234")
	orgID := []byte("ABC123456")
	type args struct {
		ctx        cTypes.Context
		toAddress  cTypes.AccAddress
		aclAccount aclTypes.ACLAccount
	}
	arg1 := args{ctx, addr1, &aclTypes.BaseACLAccount{
		Address:        addr1,
		ZoneID:         zoneID,
		OrganizationID: orgID,
	}}
	tests := []struct {
		name string
		args args
		want cTypes.Error
	}{
		{"Setting ACL",
			arg1,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ACLKeeper.DefineACLAccount(tt.args.ctx, tt.args.toAddress, tt.args.aclAccount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.DefineACLAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_CheckZoneAndGetACL(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	zone := cTypes.AccAddress([]byte("zone"))
	zoneID := []byte("ABCDEF1234")
	orgID := []byte("ABC123456")
	aclAccount := aclTypes.BaseACLAccount{
		Address:        addr1,
		ZoneID:         zoneID,
		OrganizationID: orgID,
	}
	type args struct {
		ctx  cTypes.Context
		from cTypes.AccAddress
		to   cTypes.AccAddress
	}
	arg1 := args{ctx, zone, addr1}
	arg2 := args{ctx, addr1, addr1}
	tests := []struct {
		name  string
		args  args
		pre   prerequiste
		want  aclTypes.ACL
		want1 cTypes.Error
	}{
		{"ACL account is not defined.",
			arg1,
			func() {},
			aclTypes.ACL{},
			cTypes.ErrInternal("To account acl not defined.")},
		{"Zone is not set.",
			arg1,
			func() {
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			aclTypes.ACL{},
			cTypes.ErrInternal("To account zone not found.")},
		{"From address is not zone.",
			arg2,
			func() {
				app.ACLKeeper.SetZoneAddress(ctx, zoneID, zone)
			},
			aclTypes.ACL{},
			cTypes.ErrInternal("Unauthorised transaction."),
		},
		{"Get ACL.",
			arg1,
			func() {},
			aclAccount.GetACL(),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			got, got1 := app.ACLKeeper.CheckZoneAndGetACL(tt.args.ctx, tt.args.from, tt.args.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.CheckZoneAndGetACL() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Keeper.CheckZoneAndGetACL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKeeper_CheckValidGenesisAddress(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	genesisAccount := getGenesisAccount(ctx, app.AccountKeeper)
	addr1 := cTypes.AccAddress([]byte("addr1"))
	app.AccountKeeper.SetAccount(ctx, app.AccountKeeper.NewAccountWithAddress(ctx, addr1))
	type args struct {
		ctx     cTypes.Context
		address cTypes.AccAddress
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Checking valid genesis.",
			args{ctx, genesisAccount.GetAddress()},
			true},
		{"Checking invalid genesis.",
			args{ctx, addr1},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ACLKeeper.CheckValidGenesisAddress(tt.args.ctx, tt.args.address); got != tt.want {
				t.Errorf("Keeper.CheckValidGenesisAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_CheckValidOrganizationAddress(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	orgAddr := cTypes.AccAddress([]byte("org"))
	zoneID := []byte("ABCDEF1234")
	orgID := []byte("ABC123456")

	organization := aclTypes.Organization{
		Address: orgAddr,
		ZoneID:  zoneID,
	}

	type args struct {
		ctx            cTypes.Context
		zoneID         aclTypes.ZoneID
		organizationID aclTypes.OrganizationID
		address        cTypes.AccAddress
	}
	arg1 := args{ctx, zoneID, []byte("Invalid Org ID"), orgAddr}
	arg2 := args{ctx, []byte("Invalid Zone ID"), orgID, orgAddr}
	arg3 := args{ctx, zoneID, orgID, orgAddr}
	tests := []struct {
		name string
		args args
		pre  prerequiste
		want bool
	}{
		{"Organization with ID does not exist.",
			arg1,
			func() {},
			false},
		{"Organization with incoorect zone ID.",
			arg2,
			func() {
				app.ACLKeeper.SetOrganization(ctx, orgID, organization)
			},
			false},
		{"Organization is valid.",
			arg3,
			func() {},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.ACLKeeper.CheckValidOrganizationAddress(tt.args.ctx, tt.args.zoneID, tt.args.organizationID, tt.args.address); got != tt.want {
				t.Errorf("Keeper.CheckValidOrganizationAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_CheckValidZoneAddress(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	zoneAddr := cTypes.AccAddress([]byte("zone"))
	addr1 := cTypes.AccAddress([]byte("addr1"))
	zoneID := []byte("ABCDEF1234")

	type args struct {
		ctx     cTypes.Context
		id      aclTypes.ZoneID
		address cTypes.AccAddress
	}
	arg1 := args{ctx, zoneID, zoneAddr}
	arg2 := args{ctx, zoneID, addr1}
	tests := []struct {
		name string
		args args
		pre  prerequiste
		want bool
	}{
		{"Zone is not set.",
			arg1,
			func() {},
			false},
		{"Invalid zone address.",
			arg2,
			func() {
				app.ACLKeeper.SetZoneAddress(ctx, zoneID, zoneAddr)
			},
			false},

		{"Zone is validated.",
			arg1,
			func() {},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.ACLKeeper.CheckValidZoneAddress(tt.args.ctx, tt.args.id, tt.args.address); got != tt.want {
				t.Errorf("Keeper.CheckValidZoneAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetOrganizationsByZoneID(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	orgAddr := cTypes.AccAddress([]byte("org"))
	zoneID := []byte("ABCDEF1234")
	orgID := []byte("ABC123456")

	organization := aclTypes.Organization{
		Address: orgAddr,
		ZoneID:  zoneID,
	}

	app.ACLKeeper.SetOrganization(ctx, orgID, organization)

	var organizationList []aclTypes.Organization

	type args struct {
		ctx cTypes.Context
		id  aclTypes.ZoneID
	}
	tests := []struct {
		name string
		args args
		want []aclTypes.Organization
	}{
		{"Get Organzations.",
			args{ctx, zoneID},
			append(organizationList, organization),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ACLKeeper.GetOrganizationsByZoneID(tt.args.ctx, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetOrganizationsByZoneID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetZones(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	var zones []cTypes.AccAddress

	type args struct {
		ctx cTypes.Context
	}
	tests := []struct {
		name string
		args args
		want []cTypes.AccAddress
	}{
		{"Get Zones.",
			args{ctx},
			append(zones, aclTypes.DefaultZoneID)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ACLKeeper.GetZones(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetZones() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetACLAccounts(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	addr1 := cTypes.AccAddress("addr1")
	zoneID := []byte("ABCDEF1234")
	orgID := []byte("ABC123456")
	aclAccount := aclTypes.BaseACLAccount{
		Address:        addr1,
		ZoneID:         zoneID,
		OrganizationID: orgID,
	}

	app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

	var accounts []aclTypes.BaseACLAccount

	type args struct {
		ctx cTypes.Context
	}
	tests := []struct {
		name string
		args args
		want []aclTypes.BaseACLAccount
	}{
		{"Get Accounts.",
			args{ctx},
			append(accounts, aclAccount)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := app.ACLKeeper.GetACLAccounts(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetACLAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}
