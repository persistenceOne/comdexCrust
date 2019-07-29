package bank

import (
	"fmt"
	"testing"

	"github.com/commitHub/commitBlockchain/x/reputation"

	"github.com/commitHub/commitBlockchain/store"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/acl"
	"github.com/commitHub/commitBlockchain/x/auth"
	"github.com/commitHub/commitBlockchain/x/negotiation"
	"github.com/commitHub/commitBlockchain/x/order"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setup() (*wire.Codec, sdk.Context, auth.AccountMapper, Keeper, order.Mapper, order.Keeper, negotiation.Mapper, negotiation.Keeper, acl.Mapper, acl.Keeper, reputation.Keeper) {
	db := dbm.NewMemDB()

	authKey := sdk.NewKVStoreKey("authKey")
	orderKey := sdk.NewKVStoreKey("orderKey")
	negoKey := sdk.NewKVStoreKey("negoKey")
	aclKey := sdk.NewKVStoreKey("aclKey")
	reputationKey := sdk.NewKVStoreKey("reputationKey")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(orderKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(negoKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(aclKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(reputationKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)
	negotiation.RegisterNegotiation(cdc)
	order.RegisterOrder(cdc)
	acl.RegisterWire(cdc)
	acl.RegisterACLAccount(cdc)
	reputation.RegisterReputation(cdc)
	reputation.RegisterWire(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	coinKeeper := NewKeeper(accountMapper)

	orderMapper := order.NewMapper(cdc, orderKey, sdk.ProtoBaseOrder)
	orderKeeper := order.NewKeeper(orderMapper)

	negoMapper := negotiation.NewMapper(cdc, negoKey, sdk.ProtoBaseNegotiation)
	negoKeeper := negotiation.NewKeeper(negoMapper, accountMapper)

	aclMapper := acl.NewACLMapper(cdc, aclKey, sdk.ProtoBaseACLAccount)
	aclKeeper := acl.NewKeeper(aclMapper)

	reputationMapper := reputation.NewMapper(cdc, reputationKey, sdk.ProtoBaseAccountReputation)
	reputationKeeper := reputation.NewKeeper(reputationMapper)

	return cdc, ctx, accountMapper, coinKeeper, orderMapper, orderKeeper, negoMapper, negoKeeper, aclMapper, aclKeeper, reputationKeeper
}

func setAccount(ctx sdk.Context, accountMapper auth.AccountMapper, baseAccount auth.Account, addr sdk.AccAddress) {
	baseAccount = accountMapper.GetAccount(ctx, addr)
	accountMapper.SetAccount(ctx, baseAccount)
	return

}

func setupSetCoins(ctx sdk.Context, coinKeeper Keeper, addr sdk.AccAddress, denom string, coins int64) {
	coinKeeper.SetCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin(denom, coins)})
	return
}

func TestMsgSend(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, _, negoKeeper, _, aclKeeper, reputationKeeper := setup()

	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1")),
		sdk.AccAddress([]byte("addr2")),
		sdk.AccAddress([]byte("addr3")),
		sdk.AccAddress([]byte("addr4")),
	}

	setcoins := []int64{100, 0, 25, 1}

	var baseAccount = []auth.Account{}

	for i, address := range addr {
		setupSetCoins(ctx, coinKeeper, address, "atom", setcoins[i])
		baseAccount = append(baseAccount, accountMapper.GetAccount(ctx, addr[i]))
		//accountMapper.SetAccount(ctx, baseAccount[i])
	}

	var coins = []sdk.Coins{
		sdk.Coins{sdk.NewInt64Coin("atom", 10)},
		sdk.Coins{sdk.NewInt64Coin("atom", 5)},
		sdk.Coins{sdk.NewInt64Coin("atom", 10)},
		sdk.Coins{sdk.NewInt64Coin("atom", -5)}, // This needs to be checked, -5 should not be added and -5 should not be subtracted
	}
	var inputs = []Input{
		Input{addr[0], coins[0]},
		Input{addr[1], coins[1]},
	}

	output1 := Output{addr[2], coins[2]}
	output2 := Output{addr[3], coins[3]}

	msgSend1 := NewMsgSend([]Input{inputs[0]}, []Output{output1})
	msgSend2 := NewMsgSend([]Input{inputs[1]}, []Output{})
	msgSend3 := NewMsgSend([]Input{}, []Output{output2})

	fun := NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)

	var res = fun(ctx, msgSend1)
	require.Equal(t, addr[2].String(), string(res.Tags[1].Value))

	res = fun(ctx, msgSend2)
	require.Equal(t, sdk.Tags(nil), res.Tags)

	res = fun(ctx, msgSend3)
	//require.Equal(t, addr[3].String(), string(res.Tags[0].Value))
	require.Equal(t, sdk.Tags(nil), res.Tags)
}

func TestMsgBankIssueAssets(t *testing.T) {
	_, ctx, am, coinKeeper, om, orderKeeper, negoMapper, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()

	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)

	var seller = []sdk.AccAddress{
		sdk.AccAddress([]byte("seller1")),
		sdk.AccAddress([]byte("seller2")),
		sdk.AccAddress([]byte("seller3")),
		sdk.AccAddress([]byte("seller4")),
		sdk.AccAddress([]byte("seller5")),
		sdk.AccAddress([]byte("seller6")),
		sdk.AccAddress([]byte("seller7")),
	}

	var buyer = []sdk.AccAddress{
		sdk.AccAddress([]byte("buyer1")),
		sdk.AccAddress([]byte("buyer2")),
		sdk.AccAddress([]byte("buyer3")),
		sdk.AccAddress([]byte("buyer4")),
		sdk.AccAddress([]byte("buyer5")),
		sdk.AccAddress([]byte("buyer6")),
		sdk.AccAddress([]byte("buyer7")),
	}

	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	orgID, _ := sdk.GetOrganizationIDFromString("AB12")

	var aclAccount = []sdk.BaseACLAccount{
		sdk.BaseACLAccount{
			Address: seller[0],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAsset: true, RedeemFiat: true, RedeemAsset: true, SendAsset: true, SellerExecuteOrder: true}},
		sdk.BaseACLAccount{
			Address: seller[1],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAsset: false, RedeemFiat: false, RedeemAsset: false}},
		sdk.BaseACLAccount{
			Address: seller[2],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAsset: true}},
		sdk.BaseACLAccount{
			Address: buyer[0],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiat: true, SendFiat: true, BuyerExecuteOrder: true}},
		sdk.BaseACLAccount{
			Address: buyer[1],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiat: false, BuyerExecuteOrder: true}},
		sdk.BaseACLAccount{
			Address: buyer[2],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiat: true}},
		sdk.BaseACLAccount{
			Address: seller[3],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAsset: true, SendAsset: true, ReleaseAsset: true}},
		sdk.BaseACLAccount{
			Address: buyer[4],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAsset: true, SendAsset: true, BuyerExecuteOrder: true}},
		sdk.BaseACLAccount{
			Address: buyer[5],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAsset: true, SendAsset: true, BuyerExecuteOrder: true}},
		sdk.BaseACLAccount{
			Address: buyer[5],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAsset: true, SendAsset: true, BuyerExecuteOrder: true}},
		sdk.BaseACLAccount{
			Address: seller[6],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAsset: true, SendAsset: true, BuyerExecuteOrder: true, SellerExecuteOrder: true, ReleaseAsset: true}},
	}

	pegHash := []sdk.PegHash{
		sdk.PegHash("30"),
		sdk.PegHash("31"),
		sdk.PegHash("32"),
		sdk.PegHash("33"),
		sdk.PegHash("34"),
	}

	var fiats = []sdk.BaseFiatPeg{
		{
			PegHash:           pegHash[0],
			TransactionID:     "ABCD",
			TransactionAmount: 100,
			RedeemedAmount:    10,
			Owners:            []sdk.Owner{{OwnerAddress: sdk.AccAddress("address"), Amount: 10}},
		},
		{
			PegHash:           pegHash[0],
			TransactionID:     "TransactionID",
			TransactionAmount: 100,
			RedeemedAmount:    0,
		},
	}

	var assets = []sdk.BaseAssetPeg{
		sdk.BaseAssetPeg{
			PegHash:       pegHash[0],
			AssetType:     "gold",
			AssetQuantity: 10,
		},
		sdk.BaseAssetPeg{
			PegHash:       pegHash[0],
			AssetType:     "gold",
			AssetQuantity: 10,
			Moderated:     true,
		},
		sdk.BaseAssetPeg{
			PegHash:       pegHash[0],
			AssetType:     "gold",
			AssetQuantity: 10,
			Moderated:     true,
			Locked:        false,
			TakerAddress:  nil,
			QuantityUnit:  "MT",
			AssetPrice:    10,
			OwnerAddress:  nil,
			DocumentHash:  "DOC",
		},
		sdk.BaseAssetPeg{
			PegHash:       pegHash[3],
			AssetType:     "gold",
			AssetQuantity: 10,
			Locked:        true,
			Moderated:     true,
		},
		sdk.BaseAssetPeg{
			PegHash:       pegHash[2],
			AssetType:     "gold",
			AssetQuantity: 10,
			Locked:        false,
			Moderated:     true,
		},
		sdk.BaseAssetPeg{
			PegHash:       pegHash[1],
			AssetType:     "gold",
			AssetQuantity: 10,
			Locked:        false,
			Moderated:     false,
		},
	}

	aclMapper.SetAccount(ctx, seller[0], &aclAccount[0])
	aclMapper.SetAccount(ctx, seller[1], &aclAccount[1])
	aclMapper.SetAccount(ctx, seller[2], &aclAccount[0])
	aclMapper.SetAccount(ctx, buyer[0], &aclAccount[3])
	aclMapper.SetAccount(ctx, buyer[1], &aclAccount[4])
	aclMapper.SetAccount(ctx, seller[3], &aclAccount[6])
	aclMapper.SetAccount(ctx, buyer[2], &aclAccount[5])
	aclMapper.SetAccount(ctx, buyer[4], &aclAccount[7])
	aclMapper.SetAccount(ctx, buyer[5], &aclAccount[8])
	aclMapper.SetAccount(ctx, seller[6], &aclAccount[10])

	acc := am.NewAccountWithAddress(ctx, seller[2])
	acc1 := am.NewAccountWithAddress(ctx, seller[3])
	accB := am.NewAccountWithAddress(ctx, buyer[0])
	acc1.SetAssetPegWallet(sdk.AssetPegWallet{assets[3]})
	accB.SetAssetPegWallet(sdk.AssetPegWallet{assets[4]})
	acc.SetAccountNumber(0)

	am.SetAccount(ctx, acc1)
	am.SetAccount(ctx, acc)
	am.SetAccount(ctx, accB)
	fmt.Println(am.GetAccount(ctx, seller[3]))
	acc = am.GetAccount(ctx, seller[2])
	acc.SetFiatPegWallet(sdk.FiatPegWallet{fiats[0]})
	acc.SetAssetPegWallet(sdk.AssetPegWallet{assets[0]})
	am.SetAccount(ctx, acc)

	gen := am.NewAccountWithAddress(ctx, sdk.AccAddress("main"))
	am.SetAccount(ctx, gen)

	var zoneAddr = sdk.AccAddress("Zone")
	var orgAddr = sdk.AccAddress("Organization")

	zn := am.NewAccountWithAddress(ctx, zoneAddr)
	zn.SetAccountNumber(0)
	am.SetAccount(ctx, zn)

	negotiation := negoMapper.NewNegotiation(buyer[0], seller[0], pegHash[1])
	nego := []sdk.BaseNegotiation{
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[0].Bytes(), seller[0].Bytes()...), pegHash[2].Bytes()...)),
			BuyerAddress:    buyer[0],
			SellerAddress:   seller[0],
			PegHash:         pegHash[2],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[0].Bytes(), seller[0].Bytes()...), pegHash[3].Bytes()...)),
			BuyerAddress:    buyer[0],
			SellerAddress:   seller[0],
			PegHash:         pegHash[3],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[0].Bytes(), seller[3].Bytes()...), pegHash[3].Bytes()...)),
			BuyerAddress:    buyer[0],
			SellerAddress:   seller[3],
			PegHash:         pegHash[3],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[4].Bytes(), seller[4].Bytes()...), pegHash[0].Bytes()...)),
			BuyerAddress:    buyer[4],
			SellerAddress:   seller[4],
			PegHash:         pegHash[0],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[4].Bytes(), seller[4].Bytes()...), pegHash[0].Bytes()...)),
			BuyerAddress:    buyer[4],
			SellerAddress:   seller[4],
			PegHash:         pegHash[0],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[4].Bytes(), seller[4].Bytes()...), pegHash[1].Bytes()...)),
			BuyerAddress:    buyer[4],
			SellerAddress:   seller[4],
			PegHash:         pegHash[1],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[5].Bytes(), seller[4].Bytes()...), pegHash[1].Bytes()...)),
			BuyerAddress:    buyer[5],
			SellerAddress:   seller[4],
			PegHash:         pegHash[1],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[6].Bytes(), seller[6].Bytes()...), pegHash[2].Bytes()...)),
			BuyerAddress:    buyer[6],
			SellerAddress:   seller[6],
			PegHash:         pegHash[2],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
		{
			NegotiationID:   sdk.NegotiationID(append(append(buyer[5].Bytes(), seller[6].Bytes()...), pegHash[1].Bytes()...)),
			BuyerAddress:    buyer[5],
			SellerAddress:   seller[6],
			PegHash:         pegHash[1],
			Bid:             100,
			Time:            100,
			BuyerSignature:  []byte("buyerSignature"),
			SellerSignature: []byte("sellerSignature"),
		},
	}
	negoMapper.SetNegotiation(ctx, negotiation)
	negoMapper.SetNegotiation(ctx, &nego[0])
	negoMapper.SetNegotiation(ctx, &nego[1])
	negoMapper.SetNegotiation(ctx, &nego[2])
	negoMapper.SetNegotiation(ctx, &nego[3])
	negoMapper.SetNegotiation(ctx, &nego[4])
	negoMapper.SetNegotiation(ctx, &nego[5])
	negoMapper.SetNegotiation(ctx, &nego[6])
	negoMapper.SetNegotiation(ctx, &nego[7])
	negoMapper.SetNegotiation(ctx, &nego[8])

	order := om.NewOrder(buyer[0], seller[0], pegHash[0])
	order1 := om.NewOrder(buyer[1], seller[0], pegHash[2])
	order2 := om.NewOrder(buyer[0], seller[1], pegHash[2])
	order3 := om.NewOrder(buyer[0], seller[0], pegHash[2])
	order4 := om.NewOrder(buyer[0], seller[0], pegHash[3])
	order5 := om.NewOrder(buyer[3], seller[0], pegHash[0])
	order6 := om.NewOrder(buyer[2], seller[0], pegHash[0])
	order7 := om.NewOrder(buyer[4], seller[4], pegHash[0])
	order8 := om.NewOrder(buyer[4], seller[4], pegHash[1])
	order9 := om.NewOrder(buyer[5], seller[4], pegHash[1])
	order10 := om.NewOrder(buyer[0], seller[0], pegHash[4])
	order11 := om.NewOrder(buyer[0], seller[1], pegHash[0])
	order12 := om.NewOrder(buyer[6], seller[6], pegHash[2])
	order13 := om.NewOrder(buyer[5], seller[6], pegHash[1])
	order.SetAssetPegWallet(sdk.AssetPegWallet{assets[0]})
	order1.SetAssetPegWallet(sdk.AssetPegWallet{assets[4]})
	order2.SetAssetPegWallet(sdk.AssetPegWallet{assets[4]})
	order3.SetAssetPegWallet(sdk.AssetPegWallet{assets[4]})
	order4.SetAssetPegWallet(sdk.AssetPegWallet{assets[3]})
	order4.SetAssetPegWallet(sdk.AssetPegWallet{assets[3]})
	order5.SetAssetPegWallet(sdk.AssetPegWallet{assets[1]})
	order6.SetAssetPegWallet(sdk.AssetPegWallet{assets[0]})
	order7.SetAssetPegWallet(sdk.AssetPegWallet{assets[0]})
	order8.SetAssetPegWallet(sdk.AssetPegWallet{assets[5]})
	order9.SetAssetPegWallet(sdk.AssetPegWallet{assets[5]})
	order9.SetFiatPegWallet(sdk.FiatPegWallet{fiats[0]})
	order11.SetAssetPegWallet(sdk.AssetPegWallet{assets[0]})
	order12.SetAssetPegWallet(sdk.AssetPegWallet{assets[4]})
	order12.SetFiatPegWallet(sdk.FiatPegWallet{fiats[0]})
	order13.SetAssetPegWallet(sdk.AssetPegWallet{assets[5]})
	order9.SetAWBProofHash("awbProof")
	order4.SetAWBProofHash("awbproof")
	om.SetOrder(ctx, order)
	om.SetOrder(ctx, order1)
	om.SetOrder(ctx, order2)
	om.SetOrder(ctx, order3)
	om.SetOrder(ctx, order4)
	om.SetOrder(ctx, order5)
	om.SetOrder(ctx, order6)
	om.SetOrder(ctx, order7)
	om.SetOrder(ctx, order8)
	om.SetOrder(ctx, order9)
	om.SetOrder(ctx, order10)
	om.SetOrder(ctx, order11)
	om.SetOrder(ctx, order12)
	om.SetOrder(ctx, order13)

	aclMapper.SetZone(ctx, zoneAddr, zoneID)

	msgIssueAsset := []sdk.Msg{
		NewMsgBankIssueAssets([]IssueAsset{NewIssueAsset(sdk.AccAddress("zone"), seller[0], &assets[0])}),
		NewMsgBankIssueAssets([]IssueAsset{NewIssueAsset(zoneAddr, seller[0], &assets[0])}),
		NewMsgBankIssueAssets([]IssueAsset{NewIssueAsset(zoneAddr, seller[1], &assets[0])}),
		NewMsgBankIssueAssets([]IssueAsset{NewIssueAsset(zoneAddr, seller[0], &assets[1])}),
		NewMsgBankIssueAssets([]IssueAsset{NewIssueAsset(seller[0], seller[0], &assets[1])}),
		NewMsgBankIssueAssets([]IssueAsset{NewIssueAsset(seller[0], seller[0], &assets[2])}),
		NewMsgBankIssueAssets([]IssueAsset{NewIssueAsset(nil, seller[0], &assets[2])}),
		NewMsgBankIssueAssets([]IssueAsset{NewIssueAsset(seller[0], nil, &assets[2])}),
		NewMsgBankIssueAssets([]IssueAsset{}),
	}

	msgRedeemFiat := []sdk.Msg{
		NewMsgBankRedeemFiats([]RedeemFiat{NewRedeemFiat(seller[0], sdk.AccAddress("zone"), 10)}),
		NewMsgBankRedeemFiats([]RedeemFiat{NewRedeemFiat(seller[1], zoneAddr, 10)}),
		NewMsgBankRedeemFiats([]RedeemFiat{NewRedeemFiat(seller[0], zoneAddr, 10)}),
		NewMsgBankRedeemFiats([]RedeemFiat{NewRedeemFiat(seller[2], zoneAddr, 10)}),
		NewMsgBankRedeemFiats([]RedeemFiat{}),
		NewMsgBankRedeemFiats([]RedeemFiat{NewRedeemFiat(nil, zoneAddr, 10)}),
		NewMsgBankRedeemFiats([]RedeemFiat{NewRedeemFiat(seller[2], nil, 10)}),
		NewMsgBankRedeemFiats([]RedeemFiat{NewRedeemFiat(seller[2], zoneAddr, 0)}),
	}

	msgRedeemAsset := []sdk.Msg{
		NewMsgBankRedeemAssets([]RedeemAsset{NewRedeemAsset(sdk.AccAddress("zone"), seller[0], pegHash[0])}),
		NewMsgBankRedeemAssets([]RedeemAsset{NewRedeemAsset(zoneAddr, seller[1], pegHash[0])}),
		NewMsgBankRedeemAssets([]RedeemAsset{NewRedeemAsset(zoneAddr, seller[0], pegHash[0])}),
		NewMsgBankRedeemAssets([]RedeemAsset{NewRedeemAsset(zoneAddr, seller[0], pegHash[0])}),
		NewMsgBankRedeemAssets([]RedeemAsset{NewRedeemAsset(zoneAddr, seller[2], pegHash[0])}),
		NewMsgBankRedeemAssets([]RedeemAsset{NewRedeemAsset(nil, seller[2], pegHash[0])}),
		NewMsgBankRedeemAssets([]RedeemAsset{NewRedeemAsset(zoneAddr, nil, pegHash[0])}),
		NewMsgBankRedeemAssets([]RedeemAsset{NewRedeemAsset(zoneAddr, seller[2], nil)}),
		NewMsgBankRedeemAssets([]RedeemAsset{}),
	}

	msgIssueFiats := []sdk.Msg{
		NewMsgBankIssueFiats([]IssueFiat{NewIssueFiat(sdk.AccAddress("zone"), buyer[0], &fiats[0])}),
		NewMsgBankIssueFiats([]IssueFiat{NewIssueFiat(zoneAddr, buyer[0], &fiats[0])}),
		NewMsgBankIssueFiats([]IssueFiat{NewIssueFiat(zoneAddr, buyer[1], &fiats[0])}),
		NewMsgBankIssueFiats([]IssueFiat{}),
		NewMsgBankIssueFiats([]IssueFiat{NewIssueFiat(nil, buyer[1], &fiats[0])}),
		NewMsgBankIssueFiats([]IssueFiat{NewIssueFiat(zoneAddr, nil, &fiats[0])}),
		NewMsgBankIssueFiats([]IssueFiat{NewIssueFiat(zoneAddr, buyer[1], &fiats[0])}),
		NewMsgBankIssueFiats([]IssueFiat{NewIssueFiat(zoneAddr, buyer[1], &fiats[1])}),
	}

	msgSendAsset := []sdk.Msg{
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(sdk.AccAddress("zone"), buyer[0], pegHash[0])}),
		NewMsgBankSendAssets([]SendAsset{}),
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(nil, buyer[0], pegHash[0])}),
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(seller[0], nil, pegHash[0])}),
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(seller[0], buyer[0], nil)}),
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(seller[1], buyer[0], pegHash[0])}),
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(seller[0], buyer[0], pegHash[0])}),
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(seller[0], buyer[0], pegHash[1])}),
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(seller[0], buyer[0], pegHash[2])}),
		NewMsgBankSendAssets([]SendAsset{NewSendAsset(seller[3], buyer[0], pegHash[3])}),
	}

	msgSendFiat := []sdk.Msg{
		NewMsgBankSendFiats([]SendFiat{NewSendFiat(sdk.AccAddress("zone"), buyer[0], pegHash[0], 10)}),
		NewMsgBankSendFiats([]SendFiat{}),
		NewMsgBankSendFiats([]SendFiat{NewSendFiat(nil, buyer[0], pegHash[0], 10)}),
		NewMsgBankSendFiats([]SendFiat{NewSendFiat(seller[0], nil, pegHash[0], 10)}),
		NewMsgBankSendFiats([]SendFiat{NewSendFiat(seller[0], buyer[0], nil, 10)}),
		NewMsgBankSendFiats([]SendFiat{NewSendFiat(buyer[0], seller[0], pegHash[0], 0)}),
		NewMsgBankSendFiats([]SendFiat{NewSendFiat(buyer[0], seller[0], pegHash[0], 10)}),
		NewMsgBankSendFiats([]SendFiat{NewSendFiat(buyer[1], seller[0], pegHash[0], 10)}),
		NewMsgBankSendFiats([]SendFiat{NewSendFiat(buyer[0], seller[0], pegHash[2], 10)}),
	}

	msgBuyerExecuteOrder := []sdk.Msg{
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(sdk.AccAddress("zone"), buyer[3], seller[0], pegHash[0], "")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(nil, buyer[0], seller[0], pegHash[0], "")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[0], nil, pegHash[0], "")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, nil, seller[0], pegHash[0], "")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[0], seller[0], nil, "")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[0], seller[0], pegHash[0], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[1], seller[0], pegHash[2], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[0], seller[1], pegHash[2], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(buyer[1], buyer[1], seller[0], pegHash[2], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(buyer[0], buyer[0], seller[0], pegHash[2], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(buyer[0], buyer[0], seller[0], pegHash[3], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(sdk.AccAddress("zone"), buyer[0], seller[0], pegHash[0], "")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(buyer[3], buyer[3], seller[0], pegHash[0], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[0], seller[0], pegHash[3], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[2], seller[0], pegHash[0], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[4], seller[4], pegHash[0], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[4], seller[4], pegHash[1], "fiat")}),
		NewMsgBankBuyerExecuteOrders([]BuyerExecuteOrder{NewBuyerExecuteOrder(zoneAddr, buyer[5], seller[4], pegHash[1], "fiat")}),
	}

	msgSellerExecuteOrder := []sdk.Msg{
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(sdk.AccAddress("zone"), buyer[0], seller[0], pegHash[0], "")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(nil, buyer[0], seller[0], pegHash[0], "")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(zoneAddr, nil, seller[0], pegHash[0], "")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(zoneAddr, buyer[0], nil, pegHash[0], "")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(zoneAddr, buyer[0], seller[0], nil, "")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(zoneAddr, buyer[0], seller[0], pegHash[0], "asset")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(zoneAddr, buyer[0], seller[0], pegHash[4], "asset")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(zoneAddr, buyer[0], seller[1], pegHash[0], "asset")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(zoneAddr, buyer[3], seller[0], pegHash[0], "asset")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(seller[0], buyer[3], seller[0], pegHash[0], "asset")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(seller[6], buyer[6], seller[6], pegHash[2], "")}),
		NewMsgBankSellerExecuteOrders([]SellerExecuteOrder{NewSellerExecuteOrder(zoneAddr, buyer[5], seller[6], pegHash[1], "")}),
	}

	msgReleaseAsset := []sdk.Msg{
		NewMsgBankReleaseAssets([]ReleaseAsset{NewReleaseAsset(sdk.AccAddress("zone"), seller[0], pegHash[0])}),
		NewMsgBankReleaseAssets([]ReleaseAsset{}),
		NewMsgBankReleaseAssets([]ReleaseAsset{NewReleaseAsset(nil, seller[0], pegHash[0])}),
		NewMsgBankReleaseAssets([]ReleaseAsset{NewReleaseAsset(zoneAddr, nil, pegHash[0])}),
		NewMsgBankReleaseAssets([]ReleaseAsset{NewReleaseAsset(zoneAddr, seller[0], nil)}),
		NewMsgBankReleaseAssets([]ReleaseAsset{NewReleaseAsset(zoneAddr, seller[0], nil)}),
		NewMsgBankReleaseAssets([]ReleaseAsset{NewReleaseAsset(zoneAddr, seller[6], nil)}),
		NewMsgBankReleaseAssets([]ReleaseAsset{NewReleaseAsset(zoneAddr, seller[3], pegHash[3])}),
	}

	msgDefineZone := []sdk.Msg{
		BuildMsgDefineZone(sdk.AccAddress("main"), zoneAddr, zoneID),
		NewMsgDefineZones([]DefineZone{}),
		BuildMsgDefineZone(nil, zoneAddr, zoneID),
		BuildMsgDefineZone(sdk.AccAddress("main"), nil, zoneID),
		BuildMsgDefineZone(sdk.AccAddress("main"), zoneAddr, nil),
		BuildMsgDefineZone(seller[2], zoneAddr, nil),
	}

	msgDefineOrganization := []sdk.Msg{
		BuildMsgDefineOrganization(zoneAddr, orgAddr, orgID, zoneID),
		NewMsgDefineOrganizations([]DefineOrganization{}),
		BuildMsgDefineOrganization(nil, orgAddr, orgID, zoneID),
		BuildMsgDefineOrganization(zoneAddr, nil, orgID, zoneID),
		BuildMsgDefineOrganization(zoneAddr, orgAddr, nil, zoneID),
		BuildMsgDefineOrganization(zoneAddr, orgAddr, orgID, nil),
		BuildMsgDefineOrganization(orgAddr, orgAddr, orgID, nil),
	}

	msgDefineACL := []sdk.Msg{
		BuildMsgDefineACL(orgAddr, seller[0], &aclAccount[0]),
		NewMsgDefineACLs([]DefineACL{}),
		BuildMsgDefineACL(nil, seller[0], &aclAccount[0]),
		BuildMsgDefineACL(orgAddr, nil, &aclAccount[0]),
		BuildMsgDefineACL(sdk.AccAddress("main"), seller[0], &aclAccount[0]),
		BuildMsgDefineACL(zoneAddr, seller[0], &aclAccount[0]),
	}

	//IssueAsset
	require.Equal(t, msgIssueAsset[0].Type(), "bank")
	require.NotNil(t, msgIssueAsset[0].ValidateBasic())
	require.Nil(t, msgIssueAsset[5].ValidateBasic())
	require.NotNil(t, msgIssueAsset[6].ValidateBasic())
	require.NotNil(t, msgIssueAsset[7].ValidateBasic())
	require.NotNil(t, msgIssueAsset[8].ValidateBasic())
	require.NotNil(t, msgIssueAsset[0].GetSignBytes())
	require.Equal(t, msgIssueAsset[1].GetSigners()[0], zoneAddr)

	//RedeemAsset
	require.Equal(t, msgRedeemAsset[0].Type(), "bank")
	require.Nil(t, msgRedeemAsset[0].ValidateBasic())
	require.NotNil(t, msgRedeemAsset[5].ValidateBasic())
	require.NotNil(t, msgRedeemAsset[6].ValidateBasic())
	require.NotNil(t, msgRedeemAsset[7].ValidateBasic())
	require.NotNil(t, msgRedeemAsset[8].ValidateBasic())
	require.NotNil(t, msgRedeemAsset[0].GetSignBytes())
	require.NotNil(t, msgRedeemAsset[1].GetSigners()[0], zoneAddr)

	//IssueFiat
	require.Equal(t, msgIssueFiats[0].Type(), "bank")
	require.NotNil(t, msgIssueFiats[3].ValidateBasic())
	require.NotNil(t, msgIssueFiats[4].ValidateBasic())
	require.NotNil(t, msgIssueFiats[5].ValidateBasic())
	require.Nil(t, msgIssueFiats[6].ValidateBasic())
	require.NotNil(t, msgIssueFiats[7].ValidateBasic())
	require.NotNil(t, msgIssueFiats[0].GetSignBytes())
	require.Equal(t, msgIssueFiats[1].GetSigners()[0], zoneAddr)

	//RedeemFiat
	require.Equal(t, msgRedeemFiat[0].Type(), "bank")
	require.NotNil(t, msgRedeemFiat[4].ValidateBasic())
	require.NotNil(t, msgRedeemFiat[5].ValidateBasic())
	require.NotNil(t, msgRedeemFiat[6].ValidateBasic())
	require.NotNil(t, msgRedeemFiat[7].ValidateBasic())
	require.Nil(t, msgRedeemFiat[2].ValidateBasic())
	require.NotNil(t, msgRedeemFiat[0].GetSignBytes())
	require.Equal(t, msgRedeemFiat[1].GetSigners()[0], seller[1])

	//SendAsset
	require.Equal(t, msgSendAsset[0].Type(), "bank")
	require.NotNil(t, msgSendAsset[1].ValidateBasic())
	require.NotNil(t, msgSendAsset[2].ValidateBasic())
	require.NotNil(t, msgSendAsset[3].ValidateBasic())
	require.NotNil(t, msgSendAsset[4].ValidateBasic())
	require.Nil(t, msgSendAsset[5].ValidateBasic())
	require.NotNil(t, msgSendAsset[5].GetSignBytes())
	require.Equal(t, msgSendAsset[5].GetSigners()[0], seller[1])

	//sendFiat
	require.Equal(t, msgSendFiat[0].Type(), "bank")
	require.NotNil(t, msgSendFiat[1].ValidateBasic())
	require.NotNil(t, msgSendFiat[2].ValidateBasic())
	require.NotNil(t, msgSendFiat[3].ValidateBasic())
	require.NotNil(t, msgSendFiat[4].ValidateBasic())
	require.NotNil(t, msgSendFiat[5].ValidateBasic())
	require.Nil(t, msgSendFiat[6].ValidateBasic())
	require.NotNil(t, msgSendFiat[6].GetSignBytes())
	require.Equal(t, msgSendFiat[6].GetSigners()[0], buyer[0])

	//buyerExecuteOrder
	require.Equal(t, msgBuyerExecuteOrder[0].Type(), "bank")
	require.NotNil(t, msgBuyerExecuteOrder[1].ValidateBasic())
	require.NotNil(t, msgBuyerExecuteOrder[2].ValidateBasic())
	require.NotNil(t, msgBuyerExecuteOrder[3].ValidateBasic())
	require.NotNil(t, msgBuyerExecuteOrder[4].ValidateBasic())
	require.NotNil(t, msgBuyerExecuteOrder[5].ValidateBasic())
	require.NotNil(t, msgBuyerExecuteOrder[0].ValidateBasic())
	require.Nil(t, msgBuyerExecuteOrder[6].ValidateBasic())
	require.NotNil(t, msgBuyerExecuteOrder[6].GetSignBytes())
	require.Equal(t, msgBuyerExecuteOrder[6].GetSigners()[0], zoneAddr)

	//sellerExecuteOrder
	require.Equal(t, msgSellerExecuteOrder[0].Type(), "bank")
	require.NotNil(t, msgSellerExecuteOrder[1].ValidateBasic())
	require.NotNil(t, msgSellerExecuteOrder[2].ValidateBasic())
	require.NotNil(t, msgSellerExecuteOrder[3].ValidateBasic())
	require.NotNil(t, msgSellerExecuteOrder[4].ValidateBasic())
	require.NotNil(t, msgSellerExecuteOrder[5].ValidateBasic())
	require.NotNil(t, msgSellerExecuteOrder[0].ValidateBasic())
	require.Nil(t, msgSellerExecuteOrder[6].ValidateBasic())
	require.NotNil(t, msgSellerExecuteOrder[6].GetSignBytes())
	require.Equal(t, msgSellerExecuteOrder[6].GetSigners()[0], zoneAddr)

	//buyerExecuteOrder
	require.Equal(t, msgReleaseAsset[0].Type(), "bank")
	require.NotNil(t, msgReleaseAsset[1].ValidateBasic())
	require.NotNil(t, msgReleaseAsset[2].ValidateBasic())
	require.NotNil(t, msgReleaseAsset[3].ValidateBasic())
	require.NotNil(t, msgReleaseAsset[4].ValidateBasic())
	require.Nil(t, msgReleaseAsset[0].ValidateBasic())
	require.NotNil(t, msgReleaseAsset[0].GetSignBytes())
	require.Equal(t, msgReleaseAsset[0].GetSigners()[0], sdk.AccAddress("zone"))

	//definezone
	require.Equal(t, msgDefineZone[0].Type(), "bank")
	require.NotNil(t, msgDefineZone[1].ValidateBasic())
	require.Nil(t, msgDefineZone[0].ValidateBasic())
	require.NotNil(t, msgDefineZone[2].ValidateBasic())
	require.NotNil(t, msgDefineZone[3].ValidateBasic())
	require.NotNil(t, msgDefineZone[4].ValidateBasic())
	require.NotNil(t, msgDefineZone[0].GetSignBytes())
	require.Equal(t, msgDefineZone[0].GetSigners()[0], sdk.AccAddress("main"))

	//defineorganization
	require.Equal(t, msgDefineOrganization[0].Type(), "bank")
	require.Nil(t, msgDefineOrganization[0].ValidateBasic())
	require.NotNil(t, msgDefineOrganization[1].ValidateBasic())
	require.NotNil(t, msgDefineOrganization[2].ValidateBasic())
	require.NotNil(t, msgDefineOrganization[3].ValidateBasic())
	require.NotNil(t, msgDefineOrganization[4].ValidateBasic())
	require.NotNil(t, msgDefineOrganization[5].ValidateBasic())
	require.NotNil(t, msgDefineOrganization[0].GetSignBytes())
	require.Equal(t, msgDefineOrganization[0].GetSigners()[0], zoneAddr)

	//defineAcl
	require.Equal(t, msgDefineACL[0].Type(), "bank")
	require.Nil(t, msgDefineACL[0].ValidateBasic())
	require.NotNil(t, msgDefineACL[1].ValidateBasic())
	require.NotNil(t, msgDefineACL[2].ValidateBasic())
	require.NotNil(t, msgDefineACL[3].ValidateBasic())
	require.Nil(t, msgDefineACL[0].ValidateBasic())
	require.NotNil(t, msgDefineACL[0].GetSignBytes())
	require.Equal(t, msgDefineACL[0].GetSigners()[0], orgAddr)

	testIssueAsset := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgIssueAsset[0]},
		{true, msgIssueAsset[1]},
		{false, msgIssueAsset[2]},
		{false, msgIssueAsset[3]},
		{true, msgIssueAsset[4]},
	}

	testRedeemFiat := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgRedeemFiat[0]},
		{false, msgRedeemFiat[1]},
		{false, msgRedeemFiat[2]},
		{true, msgRedeemFiat[3]},
	}

	testRedeemAsset := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgRedeemAsset[0]},
		{false, msgRedeemAsset[1]},
		{false, msgRedeemAsset[2]},
		{false, msgRedeemAsset[3]},
		{true, msgRedeemAsset[4]},
	}

	testIssueFiat := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgIssueFiats[0]},
		{false, msgIssueFiats[2]},
		{true, msgIssueFiats[1]},
	}

	testSendAsset := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgSendAsset[0]},
		{false, msgSendAsset[5]},
		{false, msgSendAsset[6]},
		{false, msgSendAsset[7]},
		{false, msgSendAsset[8]},
		{true, msgSendAsset[9]},
	}

	testSendFiat := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgSendFiat[0]},
		{false, msgSendFiat[5]},
		{false, msgSendFiat[7]},
		{true, msgSendFiat[8]},
	}

	testBuyerExecuteOrder := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgBuyerExecuteOrder[0]},
		{false, msgBuyerExecuteOrder[12]},
		{false, msgBuyerExecuteOrder[6]},
		{false, msgBuyerExecuteOrder[7]},
		{false, msgBuyerExecuteOrder[8]},
		{false, msgBuyerExecuteOrder[9]},
		{true, msgBuyerExecuteOrder[10]},
		{false, msgBuyerExecuteOrder[13]},
		{false, msgBuyerExecuteOrder[14]},
		{true, msgBuyerExecuteOrder[11]},
		{false, msgBuyerExecuteOrder[15]},
		{true, msgBuyerExecuteOrder[16]},
		{true, msgBuyerExecuteOrder[18]},
	}

	testSellerExecuteOrder := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgSellerExecuteOrder[0]},
		{false, msgSellerExecuteOrder[7]},
		{false, msgSellerExecuteOrder[6]},
		{false, msgSellerExecuteOrder[8]},
		{false, msgSellerExecuteOrder[9]},
		{false, msgSellerExecuteOrder[10]},
		{true, msgSellerExecuteOrder[11]},
		{true, msgSellerExecuteOrder[12]},
	}

	testReleaseAsset := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgReleaseAsset[0]},
		{false, msgReleaseAsset[4]},
		{false, msgReleaseAsset[6]},
		{true, msgReleaseAsset[7]},
	}

	testDefineZone := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgDefineZone[0]},
		{true, msgDefineZone[5]},
	}

	testDefineOrg := []struct {
		value bool
		msg   sdk.Msg
	}{
		{true, msgDefineOrganization[0]},
		{false, msgDefineOrganization[6]},
	}

	testDefineACL := []struct {
		value bool
		msg   sdk.Msg
	}{
		{false, msgDefineACL[4]},
		{true, msgDefineACL[5]},
	}

	for _, tt := range testDefineACL {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testDefineOrg {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testDefineZone {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testReleaseAsset {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testSellerExecuteOrder {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testIssueAsset {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testRedeemFiat {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testRedeemAsset {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testIssueFiat {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testSendAsset {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testSendFiat {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}

	for _, tt := range testBuyerExecuteOrder {
		if tt.value == true {
			res := fun(ctx, tt.msg)
			require.NotNil(t, res.Tags)
		} else {
			res := fun(ctx, tt.msg)
			require.Nil(t, res.Tags)
		}
	}
}
