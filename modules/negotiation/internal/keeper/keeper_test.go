package keeper_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/persistenceOne/persistenceSDK/modules/acl"
	negotiationTypes "github.com/persistenceOne/persistenceSDK/modules/negotiation/internal/types"
	"github.com/persistenceOne/persistenceSDK/simApp"
	"github.com/persistenceOne/persistenceSDK/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

type prerequisites func()

func TestKeeper_ChangeNegotiationBidWithACL(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress("buyer")
	seller := cTypes.AccAddress("seller")
	addr1 := cTypes.AccAddress("addr1")
	assetPegHash := []byte("30")
	assetPeg := types.BaseAssetPeg{
		PegHash:       assetPegHash,
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
		Locked:        true,
		TakerAddress:  addr1,
	}
	buyerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, buyer)
	sellerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, seller)

	app.AccountKeeper.SetAccount(ctx, buyerAccount)
	app.AccountKeeper.SetAccount(ctx, sellerAccount)

	negotiation := types.NewNegotiation(buyer, seller, assetPegHash)
	negotiation.SetBid(1000)
	negotiation.SetTime(ctx.BlockHeight() + 1)

	aclAccount := acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{Negotiation: false},
	}

	type args struct {
		ctx       cTypes.Context
		changeBid negotiationTypes.ChangeBid
	}
	arg1 := args{ctx, negotiationTypes.NewChangeBid(negotiation)}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"ACL not defined for buyer.",
			arg1,
			func() {},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"ACL not defined for seller.",
			arg1,
			func() {
				aclAccount.SetAddress(buyer)
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"Negotiation in ACL defined set to false.",
			arg1,
			func() {
				aclAccount.SetAddress(seller)
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"Seller does not have given asset peg.",
			arg1,
			func() {
				aclAccount.SetAddress(seller)
				aclAccount.SetACL(acl.ACL{Negotiation: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
				aclAccount.SetAddress(buyer)
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInvalidCoins("Asset peg not found!")},
		{"Buyer address is not taker address.",
			arg1,
			func() {
				aclAccount.SetAddress(seller)
				aclAccount.SetACL(acl.ACL{Negotiation: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
				aclAccount.SetAddress(buyer)
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)

				sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.AccountKeeper.SetAccount(ctx, sellerAccount)
			},
			cTypes.ErrInternal(fmt.Sprintf("Transaction is not permitted with %s", arg1.changeBid.Negotiation.GetBuyerAddress().String()))},
		{"Change bid.",
			arg1,
			func() {
				assetPeg.SetTakerAddress(buyer)
				sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.AccountKeeper.SetAccount(ctx, sellerAccount)
			},
			nil},
		{"Change bid of already signed negotiation.",
			arg1,
			func() {
				negotiation.SetBuyerSignature(types.Signature([]byte(negotiation.GetNegotiationID().Bytes())))
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			negotiationTypes.ErrVerifySignature(negotiationTypes.DefaultCodeSpace, "Already signed. Cannot change negotiation now")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.NegotiationKeeper.ChangeNegotiationBidWithACL(tt.args.ctx, tt.args.changeBid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.ChangeNegotiationBidWithACL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_ConfirmNegotiationBidWithACL(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	addr1 := cTypes.AccAddress("addr1")
	assetPegHash := []byte("30")
	assetPeg := types.BaseAssetPeg{
		PegHash:       assetPegHash,
		DocumentHash:  "DEFAULT",
		AssetType:     "DEFAULT",
		AssetQuantity: 0,
		AssetPrice:    0,
		QuantityUnit:  "DEFAULT",
		Moderated:     false,
		Locked:        true,
		TakerAddress:  addr1,
	}

	kb := keys.NewInMemory()

	_, buyerSeed, _ := kb.CreateMnemonic("buyer", cKeys.English, "1234567890", cKeys.SigningAlgo("secp256k1"))
	buyerInfo, _ := kb.CreateAccount("buyer", buyerSeed, "", "1234567890", uint32(app.AccountKeeper.GetNextAccountNumber(ctx)), 0)
	buyerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, buyerInfo.GetAddress())
	buyerAccount.SetPubKey(buyerInfo.GetPubKey())
	app.AccountKeeper.SetAccount(ctx, buyerAccount)
	buyer := buyerAccount.GetAddress()

	_, sellerSeed, _ := kb.CreateMnemonic("seller", cKeys.English, "1234567890", cKeys.SigningAlgo("secp256k1"))
	sellerInfo, _ := kb.CreateAccount("seller", sellerSeed, "", "1234567890", uint32(app.AccountKeeper.GetNextAccountNumber(ctx)), 0)
	sellerAccount := app.AccountKeeper.NewAccountWithAddress(ctx, sellerInfo.GetAddress())
	sellerAccount.SetPubKey(sellerInfo.GetPubKey())
	app.AccountKeeper.SetAccount(ctx, sellerAccount)
	seller := sellerAccount.GetAddress()

	negotiation := types.NewNegotiation(buyer, seller, assetPegHash)
	negotiation.SetBid(1000)
	negotiation.SetTime(ctx.BlockHeight() + 1)

	signBytes := negotiationTypes.NewSignNegotiationBody(buyer, seller, assetPegHash, negotiation.GetBid(), negotiation.GetTime()).GetSignBytes()

	buyerSign, _, _ := kb.Sign("buyer", "1234567890", signBytes)
	negotiation.SetBuyerSignature(buyerSign)
	sellerSign, _, _ := kb.Sign("seller", "1234567890", signBytes)
	negotiation.SetSellerSignature(sellerSign)

	aclAccount := acl.BaseACLAccount{
		Address:        buyer,
		ZoneID:         acl.DefaultZoneID,
		OrganizationID: acl.DefaultOrganizationID,
		ACL:            acl.ACL{ConfirmBuyerBid: false, ConfirmSellerBid: false},
	}
	confirmBid1 := negotiationTypes.NewConfirmBid(negotiation)

	negotiation2 := types.NewNegotiation(buyer, seller, assetPegHash)
	negotiation2.SetBid(1100)
	confirmBid2 := negotiationTypes.NewConfirmBid(negotiation2)

	negotiation3 := types.NewNegotiation(buyer, seller, assetPegHash)
	negotiation3.SetBid(negotiation.GetBid())
	negotiation3.SetSellerSignature([]byte("Seller Signature"))
	confirmBid3 := negotiationTypes.NewConfirmBid(negotiation3)

	negotiation4 := types.NewNegotiation(buyer, seller, assetPegHash)
	negotiation4.SetBid(negotiation.GetBid())
	negotiation4.SetBuyerSignature([]byte("Buyer Signature"))
	confirmBid4 := negotiationTypes.NewConfirmBid(negotiation4)

	type args struct {
		ctx        cTypes.Context
		confirmBid negotiationTypes.ConfirmBid
	}

	arg1 := args{ctx, confirmBid1}
	arg2 := args{ctx, confirmBid2}
	arg3 := args{ctx, confirmBid3}
	arg4 := args{ctx, confirmBid4}
	tests := []struct {
		name string
		args args
		pre  prerequisites
		want cTypes.Error
	}{
		{"ACL not defined for buyer.",
			arg1,
			func() {},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"ACL not defined for seller.",
			arg1,
			func() {
				aclAccount.SetAddress(buyer)
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"Negotiation in ACL defined set to false.",
			arg1,
			func() {
				aclAccount.SetAddress(seller)
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInternal("Unauthorized transaction")},
		{"Seller does not have given asset peg.",
			arg1,
			func() {
				aclAccount.SetAddress(seller)
				aclAccount.SetACL(acl.ACL{ConfirmBuyerBid: true, ConfirmSellerBid: true})
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
				aclAccount.SetAddress(buyer)
				app.ACLKeeper.SetACLAccount(ctx, &aclAccount)
			},
			cTypes.ErrInvalidCoins("Asset peg not found!")},
		{"Buyer address is not taker address.",
			arg1,
			func() {
				sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.AccountKeeper.SetAccount(ctx, sellerAccount)
			},
			cTypes.ErrInternal(fmt.Sprintf("Transaction is not permitted with %s", arg1.confirmBid.Negotiation.GetBuyerAddress().String()))},
		{"Buyer and seller both confirming the bid.",
			arg1,
			func() {
				assetPeg.SetTakerAddress(buyer)
				sellerAccount.SetAssetPegWallet(types.AssetPegWallet{assetPeg})
				app.AccountKeeper.SetAccount(ctx, sellerAccount)
			},
			nil},
		{"Trying to confirm the negotiation again.",
			arg1,
			func() {
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			negotiationTypes.ErrVerifySignature(negotiationTypes.DefaultCodeSpace, "Already Exist the signatures")},
		{"Confirming negotiation bid is different than old negotiating bid.",
			arg2,
			func() {
				negotiation.SetBuyerSignature(nil)
				negotiation.SetSellerSignature(nil)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			negotiationTypes.ErrInvalidBid(negotiationTypes.DefaultCodeSpace, "Buyer and Seller must confirm with same bid amount")},
		{"Invalid seller signature.",
			arg3,
			func() {},
			negotiationTypes.ErrVerifySignature(negotiationTypes.DefaultCodeSpace, "Seller signature verification failed")},
		{"Invalid buyer signature.",
			arg4,
			func() {
				negotiation.SetSellerSignature(sellerSign)
				app.NegotiationKeeper.SetNegotiation(ctx, negotiation)
			},
			negotiationTypes.ErrVerifySignature(negotiationTypes.DefaultCodeSpace, "Buyer signature verification failed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pre()
			if got := app.NegotiationKeeper.ConfirmNegotiationBidWithACL(tt.args.ctx, tt.args.confirmBid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.ConfirmNegotiationBidWithACL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_GetNegotiations(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)
	var negotiations []types.Negotiation

	negotiation := types.NewNegotiation(cTypes.AccAddress("buyer"), cTypes.AccAddress("seller"), []byte("30"))
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	type args struct {
		ctx cTypes.Context
	}
	tests := []struct {
		name             string
		args             args
		wantNegotiations []types.Negotiation
	}{
		{"Get all negotiations.",
			args{ctx},
			append(negotiations, negotiation)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNegotiations := app.NegotiationKeeper.GetNegotiations(tt.args.ctx); !reflect.DeepEqual(gotNegotiations, tt.wantNegotiations) {
				t.Errorf("Keeper.GetNegotiations() = %v, want %v", gotNegotiations, tt.wantNegotiations)
			}
		})
	}
}

func TestKeeper_GetNegotiationDetails(t *testing.T) {
	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress("buyer")
	seller := cTypes.AccAddress("seller")
	pegHash := []byte("30")
	negotiation := types.NewNegotiation(buyer, seller, pegHash)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation)

	type args struct {
		ctx           cTypes.Context
		buyerAddress  cTypes.AccAddress
		sellerAddress cTypes.AccAddress
		hash          types.PegHash
	}
	tests := []struct {
		name  string
		args  args
		want  types.Negotiation
		want1 cTypes.Error
	}{
		{"Get negotiation details.",
			args{ctx, buyer, seller, pegHash},
			negotiation,
			nil},
		{"Querying negotiation which is not present",
			args{ctx, seller, buyer, pegHash},
			nil,
			negotiationTypes.ErrInvalidNegotiationID(negotiationTypes.DefaultCodeSpace, "negotiation not found.")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := app.NegotiationKeeper.GetNegotiationDetails(tt.args.ctx, tt.args.buyerAddress, tt.args.sellerAddress, tt.args.hash)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetNegotiationDetails() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Keeper.GetNegotiationDetails() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKeeper_IterateNegotiations(t *testing.T) {

	app, ctx := simApp.CreateTestApp(false)

	buyer := cTypes.AccAddress("buyer")
	seller := cTypes.AccAddress("seller")
	pegHash := []byte("30")
	negotiation2 := types.NewNegotiation(buyer, seller, pegHash)
	app.NegotiationKeeper.SetNegotiation(ctx, negotiation2)

	type args struct {
		ctx     cTypes.Context
		handler func(negotiation types.Negotiation) (stop bool)
	}
	tests := []struct {
		name string
		args args
	}{
		{"Iterating over negotiations",
			args{ctx, func(negotiation types.Negotiation) (stop bool) {
				if bytes.Compare(negotiation.GetNegotiationID(), negotiation2.GetNegotiationID()) == 0 {
					return true
				}
				return true
			}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.NegotiationKeeper.IterateNegotiations(tt.args.ctx, tt.args.handler)
		})
	}
}
