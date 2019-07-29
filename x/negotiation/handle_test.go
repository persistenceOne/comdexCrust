package negotiation

import (
	"testing"

	key1 "github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/store"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/acl"
	"github.com/commitHub/commitBlockchain/x/auth"
	"github.com/commitHub/commitBlockchain/x/order"
	"github.com/commitHub/commitBlockchain/x/reputation"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setup() (*wire.Codec, sdk.Context, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()

	authKey := sdk.NewKVStoreKey("authKey")
	negotiationKey := sdk.NewKVStoreKey("negotiationKey")
	reputationKey := sdk.NewKVStoreKey("reputation")
	aclKey := sdk.NewKVStoreKey("aclKey")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(reputationKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(negotiationKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(aclKey, sdk.StoreTypeIAVL, db)

	ms.LoadLatestVersion()

	cdc := wire.NewCodec()
	order.RegisterOrder(cdc)
	sdk.RegisterWire(cdc)
	RegisterWire(cdc)
	RegisterNegotiation(cdc)
	acl.RegisterWire(cdc)
	acl.RegisterACLAccount(cdc)
	auth.RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	reputation.RegisterReputation(cdc)
	reputation.RegisterWire(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	return cdc, ctx, negotiationKey, reputationKey, authKey, aclKey
}

var (
	cdc, ctx, negotiationKey, reputationKey, authKey, aclKey = setup()
	am                                                       = auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	nm                                                       = NewMapper(cdc, negotiationKey, sdk.ProtoBaseNegotiation)
	nk                                                       = NewKeeper(nm, am)
	aclMapper                                                = acl.NewACLMapper(cdc, aclKey, sdk.ProtoBaseACLAccount)
	aclKeeper                                                = acl.NewKeeper(aclMapper)
	rm                                                       = reputation.NewMapper(cdc, reputationKey, sdk.ProtoBaseAccountReputation)
	rk                                                       = reputation.NewKeeper(rm)
)

func Test_NewHandler(t *testing.T) {
	handler := NewHandler(nk, aclKeeper, rk)

	var buyer = []sdk.AccAddress{
		sdk.AccAddress("buyer1"),
		sdk.AccAddress("buyer2"),
		sdk.AccAddress("buyer3"),
		sdk.AccAddress("buyer4"),
		sdk.AccAddress("buyer5"),
		sdk.AccAddress("buyer6"),
		sdk.AccAddress(""),
	}

	var seller = []sdk.AccAddress{
		sdk.AccAddress("seller1"),
		sdk.AccAddress("seller2"),
		sdk.AccAddress("seller3"),
		sdk.AccAddress("seller4"),
		sdk.AccAddress("seller5"),
		sdk.AccAddress("seller6"),
		sdk.AccAddress("seller7"),
	}

	var pegHash = []sdk.PegHash{
		sdk.PegHash("31"),
		sdk.PegHash("32"),
		sdk.PegHash("33"),
	}

	asset := []sdk.BaseAssetPeg{
		{
			PegHash:       pegHash[0],
			DocumentHash:  "DocHash",
			AssetType:     "rice",
			AssetQuantity: 100,
			AssetPrice:    10,
			QuantityUnit:  "MT",
			TakerAddress:  nil,
			Locked:        true,
			Moderated:     false,
		},
		{
			PegHash:       pegHash[0],
			DocumentHash:  "DocHash",
			AssetType:     "rice",
			AssetQuantity: 100,
			AssetPrice:    10,
			QuantityUnit:  "MT",
			TakerAddress:  buyer[0],
			Locked:        true,
			Moderated:     false,
		},
	}
	var zoneID = []sdk.ZoneID{
		sdk.ZoneID("zone1"),
		sdk.ZoneID("zone2"),
	}
	var organizationID = []sdk.OrganizationID{
		sdk.OrganizationID("org1"),
		sdk.OrganizationID("org2"),
	}

	buyerName := "buyerName"
	sellerName := "sellerName"
	buyerName1 := "buyerName1"
	sellerName1 := "sellerName1"
	cstore := key1.New(
		dbm.NewMemDB(),
	)
	algo := key1.Secp256k1
	pass := "1234567890"
	i1, _, _ := cstore.CreateMnemonic(buyerName, key1.English, pass, algo)
	i2, _, _ := cstore.CreateMnemonic(sellerName, key1.English, pass, algo)
	i3, _, _ := cstore.CreateMnemonic(buyerName1, key1.English, pass, algo)
	i4, _, _ := cstore.CreateMnemonic(sellerName1, key1.English, pass, algo)
	buyerAddr := sdk.AccAddress(i1.GetPubKey().Address())
	sellerAddr := sdk.AccAddress(i2.GetPubKey().Address())
	buyerAddr1 := sdk.AccAddress(i3.GetPubKey().Address())
	sellerAddr1 := sdk.AccAddress(i4.GetPubKey().Address())

	signBytesBuyer := NewSignNegotiationBody(buyerAddr, sellerAddr, pegHash[0], 100, 100).GetSignBytes()
	signBytesSeller := NewSignNegotiationBody(buyerAddr, sellerAddr, pegHash[0], 100, 100).GetSignBytes()
	signBytesBuyer1 := NewSignNegotiationBody(buyerAddr1, sellerAddr1, pegHash[0], 100, 100).GetSignBytes()
	signBytesSeller1 := NewSignNegotiationBody(buyerAddr1, sellerAddr1, pegHash[0], 100, 100).GetSignBytes()
	signatureB, _, _ := cstore.Sign(buyerName, pass, signBytesBuyer)
	signatureS, _, _ := cstore.Sign(sellerName, pass, signBytesSeller)
	signatureB1, _, _ := cstore.Sign(buyerName1, pass, signBytesBuyer1)
	signatureS1, _, _ := cstore.Sign(sellerName1, pass, signBytesSeller1)

	aclAccount := acl.DefaultACLAccount(zoneID[0], organizationID[0], buyer[2])
	aclAccountS := acl.DefaultACLAccount(zoneID[0], organizationID[0], seller[3])
	aclAccountS4 := acl.DefaultACLAccount(zoneID[0], organizationID[0], seller[4])
	aclAccountS5 := acl.DefaultACLAccount(zoneID[0], organizationID[0], seller[5])
	aclAccountS6 := acl.DefaultACLAccount(zoneID[0], organizationID[0], seller[6])
	aclAccountB6 := acl.DefaultACLAccount(zoneID[0], organizationID[0], buyer[5])
	aclAccountBuyer := acl.DefaultACLAccount(zoneID[0], organizationID[0], buyerAddr)
	aclAccountSeller := acl.DefaultACLAccount(zoneID[0], organizationID[0], sellerAddr)
	aclAccountBuyer1 := acl.DefaultACLAccount(zoneID[0], organizationID[0], buyerAddr1)
	aclAccountSeller1 := acl.DefaultACLAccount(zoneID[0], organizationID[0], sellerAddr1)

	negotiation1 := nm.NewNegotiation(buyer[0], seller[0], pegHash[0])
	negotiation2 := nm.NewNegotiation(buyer[1], seller[1], pegHash[1])
	negotiation3 := nm.NewNegotiation(buyer[2], seller[2], pegHash[1])
	negotiation4 := nm.NewNegotiation(buyer[2], seller[3], pegHash[1])
	negotiation5 := nm.NewNegotiation(buyer[2], seller[4], pegHash[0])
	negotiation6 := nm.NewNegotiation(buyer[2], seller[5], pegHash[0])
	negotiation7 := nm.NewNegotiation(buyer[2], seller[6], pegHash[0])
	negotiation8 := nm.NewNegotiation(buyer[5], seller[6], pegHash[0])
	nego := []sdk.BaseNegotiation{
		{
			NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(buyerAddr, sellerAddr, pegHash[0])),
			BuyerAddress:    buyerAddr,
			SellerAddress:   sellerAddr,
			PegHash:         pegHash[0],
			Bid:             100,
			Time:            100,
			BuyerSignature:  signatureB,
			SellerSignature: signatureS,
		},
		{
			NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(buyerAddr1, sellerAddr1, pegHash[0])),
			BuyerAddress:    buyerAddr1,
			SellerAddress:   sellerAddr1,
			PegHash:         pegHash[0],
			Bid:             10,
			Time:            100,
			BuyerSignature:  signatureB1,
			SellerSignature: signatureS1,
		},
		{
			NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(buyerAddr1, sellerAddr1, pegHash[0])),
			BuyerAddress:    buyerAddr1,
			SellerAddress:   sellerAddr1,
			PegHash:         pegHash[0],
			Bid:             100,
			Time:            10,
			BuyerSignature:  signatureS1,
			SellerSignature: nil,
		},
		{
			NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(buyerAddr1, sellerAddr1, pegHash[0])),
			BuyerAddress:    nil,
			SellerAddress:   sellerAddr1,
			PegHash:         pegHash[0],
			Bid:             100,
			Time:            10,
			BuyerSignature:  signatureS1,
			SellerSignature: nil,
		},
		{
			NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(buyerAddr1, sellerAddr1, pegHash[0])),
			BuyerAddress:    buyerAddr1,
			SellerAddress:   nil,
			PegHash:         pegHash[0],
			Bid:             100,
			Time:            10,
			BuyerSignature:  signatureS1,
			SellerSignature: nil,
		},
		{
			NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(buyerAddr1, sellerAddr1, pegHash[0])),
			BuyerAddress:    buyerAddr1,
			SellerAddress:   sellerAddr1,
			PegHash:         nil,
			Bid:             100,
			Time:            10,
			BuyerSignature:  signatureS1,
			SellerSignature: nil,
		},
		{
			NegotiationID:   nil,
			BuyerAddress:    buyerAddr1,
			SellerAddress:   sellerAddr1,
			PegHash:         pegHash[0],
			Bid:             100,
			Time:            10,
			BuyerSignature:  signatureS1,
			SellerSignature: nil,
		},
		{
			NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(buyerAddr1, sellerAddr1, pegHash[0])),
			BuyerAddress:    buyerAddr1,
			SellerAddress:   sellerAddr1,
			PegHash:         pegHash[0],
			Bid:             100,
			Time:            -1,
			BuyerSignature:  signatureS1,
			SellerSignature: nil,
		},
	}
	ac := am.NewAccountWithAddress(ctx, buyerAddr)
	acS := am.NewAccountWithAddress(ctx, sellerAddr)
	ac1 := am.NewAccountWithAddress(ctx, buyerAddr1)
	acS1 := am.NewAccountWithAddress(ctx, sellerAddr1)
	ac.SetPubKey(i1.GetPubKey())
	acS.SetPubKey(i2.GetPubKey())
	ac1.SetPubKey(i3.GetPubKey())
	acS1.SetPubKey(i4.GetPubKey())
	am.SetAccount(ctx, ac)
	am.SetAccount(ctx, acS)
	am.SetAccount(ctx, ac1)
	am.SetAccount(ctx, acS1)

	aclAcc := aclAccountS.GetACL()
	aclAcc.Negotiation = false
	aclAcc.ConfirmSellerBid = false
	aclAccountS.SetACL(aclAcc)

	acc := am.NewAccountWithAddress(ctx, seller[4])
	acc1 := am.NewAccountWithAddress(ctx, seller[5])
	acc2 := am.NewAccountWithAddress(ctx, seller[6])
	am.SetAccount(ctx, acc)
	am.SetAccount(ctx, acc1)
	am.SetAccount(ctx, acc2)

	acc = am.GetAccount(ctx, seller[4])
	acc1 = am.GetAccount(ctx, seller[5])
	acc.SetAssetPegWallet(sdk.AssetPegWallet{asset[0]})
	acc1.SetAssetPegWallet(sdk.AssetPegWallet{asset[1]})

	am.SetAccount(ctx, acc)
	am.SetAccount(ctx, acc1)
	am.SetAccount(ctx, acc2)
	aclMapper.SetAccount(ctx, buyer[2], aclAccount)
	aclMapper.SetAccount(ctx, seller[3], aclAccountS)
	aclMapper.SetAccount(ctx, seller[4], aclAccountS4)
	aclMapper.SetAccount(ctx, seller[5], aclAccountS5)
	aclMapper.SetAccount(ctx, seller[6], aclAccountS6)
	aclMapper.SetAccount(ctx, buyer[5], aclAccountB6)
	aclMapper.SetAccount(ctx, buyerAddr, aclAccountBuyer)
	aclMapper.SetAccount(ctx, sellerAddr, aclAccountSeller)
	aclMapper.SetAccount(ctx, buyerAddr1, aclAccountBuyer1)
	aclMapper.SetAccount(ctx, sellerAddr1, aclAccountSeller1)
	negotiation5.SetSellerSignature(signatureS)
	negotiation5.SetBuyerSignature(signatureB)
	nm.SetNegotiation(ctx, negotiation1)
	nm.SetNegotiation(ctx, negotiation2)
	nm.SetNegotiation(ctx, negotiation5)

	msgChangeBid := []sdk.Msg{
		NewMsgChangeBuyerBids([]ChangeBid{NewChangeBid(negotiation1)}),
		NewMsgChangeBuyerBids([]ChangeBid{NewChangeBid(negotiation2)}),
		BuildMsgChangeSellerBid(negotiation1),
		NewMsgChangeSellerBids([]ChangeBid{NewChangeBid(negotiation2)}),
		NewMsgChangeBuyerBids([]ChangeBid{NewChangeBid(negotiation3)}),
		NewMsgChangeBuyerBids([]ChangeBid{NewChangeBid(negotiation4)}),
		NewMsgChangeBuyerBids([]ChangeBid{NewChangeBid(negotiation5)}),
		NewMsgChangeBuyerBids([]ChangeBid{NewChangeBid(negotiation6)}),
		NewMsgChangeBuyerBids([]ChangeBid{NewChangeBid(negotiation7)}),
		NewMsgChangeSellerBids([]ChangeBid{NewChangeBid(negotiation5)}),
		BuildMsgChangeBuyerBid(&nego[0]),
		BuildMsgChangeBuyerBid(&nego[3]),
		BuildMsgChangeBuyerBid(&nego[4]),
		BuildMsgChangeBuyerBid(&nego[5]),
		BuildMsgChangeBuyerBid(&nego[6]),
		BuildMsgChangeBuyerBid(&nego[7]),
		BuildMsgChangeSellerBid(&nego[0]),
		BuildMsgChangeSellerBid(&nego[3]),
		BuildMsgChangeSellerBid(&nego[4]),
		BuildMsgChangeSellerBid(&nego[5]),
		BuildMsgChangeSellerBid(&nego[6]),
		BuildMsgChangeSellerBid(&nego[7]),
	}
	msgConfirmBid := []sdk.Msg{
		BuildMsgConfirmBuyerBid(negotiation1),
		NewMsgConfirmBuyerBids([]ConfirmBid{NewConfirmBid(negotiation2)}),
		BuildMsgConfirmSellerBid(negotiation1),
		NewMsgConfirmSellerBids([]ConfirmBid{NewConfirmBid(negotiation2)}),
		NewMsgConfirmBuyerBids([]ConfirmBid{NewConfirmBid(negotiation3)}),
		NewMsgConfirmBuyerBids([]ConfirmBid{NewConfirmBid(negotiation5)}),
		NewMsgConfirmBuyerBids([]ConfirmBid{NewConfirmBid(negotiation4)}),
		NewMsgConfirmBuyerBids([]ConfirmBid{NewConfirmBid(negotiation8)}),
		NewMsgConfirmBuyerBids([]ConfirmBid{NewConfirmBid(&nego[0])}),
		NewMsgConfirmBuyerBids([]ConfirmBid{NewConfirmBid(&nego[1])}),
		NewMsgConfirmBuyerBids([]ConfirmBid{NewConfirmBid(&nego[2])}),
		NewMsgConfirmSellerBids([]ConfirmBid{NewConfirmBid(negotiation5)}),
		NewMsgConfirmSellerBids([]ConfirmBid{NewConfirmBid(negotiation7)}),
		BuildMsgConfirmSellerBid(&nego[0]),
		BuildMsgConfirmSellerBid(&nego[3]),
		BuildMsgConfirmSellerBid(&nego[4]),
		BuildMsgConfirmSellerBid(&nego[5]),
		BuildMsgConfirmSellerBid(&nego[6]),
		BuildMsgConfirmSellerBid(&nego[7]),
		BuildMsgConfirmBuyerBid(&nego[0]),
		BuildMsgConfirmBuyerBid(&nego[3]),
		BuildMsgConfirmBuyerBid(&nego[4]),
		BuildMsgConfirmBuyerBid(&nego[5]),
		BuildMsgConfirmBuyerBid(&nego[6]),
		BuildMsgConfirmBuyerBid(&nego[7]),
	}
	signNegotiationBody := NewSignNegotiationBody(buyer[0], seller[0], pegHash[0], 100, 100)

	require.Equal(t, msgChangeBid[0].Type(), "negotiation")
	require.Equal(t, msgChangeBid[2].Type(), "negotiation")
	require.Equal(t, msgConfirmBid[0].Type(), "negotiation")
	require.Equal(t, msgConfirmBid[2].Type(), "negotiation")

	require.NotNil(t, msgChangeBid[0].ValidateBasic())
	require.Nil(t, msgChangeBid[10].ValidateBasic())
	require.NotNil(t, msgChangeBid[11].ValidateBasic())
	require.NotNil(t, msgChangeBid[12].ValidateBasic())
	require.NotNil(t, msgChangeBid[13].ValidateBasic())
	require.NotNil(t, msgChangeBid[14].ValidateBasic())
	require.NotNil(t, msgChangeBid[15].ValidateBasic())
	require.Nil(t, msgChangeBid[16].ValidateBasic())
	require.NotNil(t, msgChangeBid[17].ValidateBasic())
	require.NotNil(t, msgChangeBid[18].ValidateBasic())
	require.NotNil(t, msgChangeBid[19].ValidateBasic())
	require.NotNil(t, msgChangeBid[20].ValidateBasic())
	require.NotNil(t, msgChangeBid[21].ValidateBasic())

	require.NotNil(t, msgConfirmBid[0].ValidateBasic())
	require.Nil(t, msgConfirmBid[13].ValidateBasic())
	require.NotNil(t, msgConfirmBid[14].ValidateBasic())
	require.NotNil(t, msgConfirmBid[15].ValidateBasic())
	require.NotNil(t, msgConfirmBid[16].ValidateBasic())
	require.NotNil(t, msgConfirmBid[17].ValidateBasic())
	require.NotNil(t, msgConfirmBid[18].ValidateBasic())
	require.Nil(t, msgConfirmBid[19].ValidateBasic())
	require.NotNil(t, msgConfirmBid[20].ValidateBasic())
	require.NotNil(t, msgConfirmBid[21].ValidateBasic())
	require.NotNil(t, msgConfirmBid[22].ValidateBasic())
	require.NotNil(t, msgConfirmBid[23].ValidateBasic())
	require.NotNil(t, msgConfirmBid[24].ValidateBasic())

	require.NotNil(t, msgChangeBid[0].GetSignBytes())
	require.NotNil(t, msgChangeBid[2].GetSignBytes())
	require.NotNil(t, msgConfirmBid[0].GetSignBytes())
	require.NotNil(t, msgConfirmBid[2].GetSignBytes())

	require.Equal(t, msgChangeBid[0].GetSigners()[0], buyer[0])
	require.Equal(t, msgChangeBid[2].GetSigners()[0], seller[0])
	require.Equal(t, msgConfirmBid[0].GetSigners()[0], buyer[0])
	require.Equal(t, msgConfirmBid[2].GetSigners()[0], seller[0])
	require.NotNil(t, signNegotiationBody.GetSignBytes())
	require.Nil(t, handler(ctx, BuildMsgChangeBuyerBid(negotiation3)).Tags)

	testBid := []struct {
		result bool
		msg    sdk.Msg
	}{
		{false, msgChangeBid[0]},
		{false, msgChangeBid[2]},
		{false, msgChangeBid[4]},
		{false, msgChangeBid[5]},
		{true, msgChangeBid[6]},
		{true, msgChangeBid[7]},
		{true, msgChangeBid[9]},
	}

	testConfirmBid := []struct {
		result bool
		msg    sdk.Msg
	}{
		{false, msgConfirmBid[0]},
		{false, msgConfirmBid[2]},
		{false, msgConfirmBid[4]},
		{false, msgConfirmBid[6]},
		{true, msgConfirmBid[7]},
		{true, msgConfirmBid[5]},
		{true, msgConfirmBid[8]},
		{true, msgConfirmBid[9]},
		{true, msgConfirmBid[10]},
		{true, msgConfirmBid[11]},
		{true, msgConfirmBid[12]},
	}
	for _, testMsg := range testBid {
		if testMsg.result == true {
			tags := handler(ctx, testMsg.msg)
			require.NotNil(t, tags)
		} else {
			tags := handler(ctx, testMsg.msg)
			require.Nil(t, tags.Tags)
		}
	}
	for _, test := range testConfirmBid {
		if test.result == true {
			tags := handler(ctx, test.msg)
			require.NotNil(t, tags)
		} else {
			tags := handler(ctx, test.msg)
			require.Nil(t, tags.Tags)
		}
	}

	_, negotiation := nk.GetNegotiation(ctx, buyer[0], seller[0], pegHash[0])
	err, _ := nk.GetNegotiation(ctx, buyer[0], seller[1], pegHash[1])
	require.NotNil(t, negotiation)
	require.NotNil(t, err)

	nm.IterateNegotiations(ctx, func(negotiation1 sdk.Negotiation) bool {
		return false
	})

	require.Nil(t, InitNegotiation(ctx, nk))
}
