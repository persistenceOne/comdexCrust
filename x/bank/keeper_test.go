package bank

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/comdex-blockchain/store"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	
	"github.com/comdex-blockchain/x/auth"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey
}

func TestKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()
	
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)
	
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	coinKeeper := NewKeeper(accountMapper)
	
	addr := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addr3 := sdk.AccAddress([]byte("addr3"))
	acc := accountMapper.NewAccountWithAddress(ctx, addr)
	
	// Test GetCoins/SetCoins
	accountMapper.SetAccount(ctx, acc)
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))
	
	coinKeeper.SetCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	
	// Test HasCoins
	require.True(t, coinKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, coinKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	require.False(t, coinKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)}))
	require.False(t, coinKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 5)}))
	
	// Test AddCoins
	coinKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 25)}))
	
	coinKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 15)})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 15), sdk.NewInt64Coin("foocoin", 25)}))
	
	// Test SubtractCoins
	coinKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	coinKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 5)})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 15)}))
	
	coinKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 11)})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 15)}))
	
	coinKeeper.SubtractCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 10)})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 15)}))
	require.False(t, coinKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 1)}))
	
	// Test SendCoins
	coinKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, coinKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	
	_, err2 := coinKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 50)})
	assert.Implements(t, (*sdk.Error)(nil), err2)
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, coinKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	
	coinKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 30)})
	coinKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 5)})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 5)}))
	require.True(t, coinKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 10)}))
	
	// Test InputOutputCoins
	input1 := NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 2)})
	output1 := NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 2)})
	coinKeeper.InputOutputCoins(ctx, []Input{input1}, []Output{output1})
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 7)}))
	require.True(t, coinKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 8)}))
	
	inputs := []Input{
		NewInput(addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 3)}),
		NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("barcoin", 3), sdk.NewInt64Coin("foocoin", 2)}),
	}
	
	outputs := []Output{
		NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 1)}),
		NewOutput(addr3, sdk.Coins{sdk.NewInt64Coin("barcoin", 2), sdk.NewInt64Coin("foocoin", 5)}),
	}
	coinKeeper.InputOutputCoins(ctx, inputs, outputs)
	require.True(t, coinKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 21), sdk.NewInt64Coin("foocoin", 4)}))
	require.True(t, coinKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 7), sdk.NewInt64Coin("foocoin", 6)}))
	require.True(t, coinKeeper.GetCoins(ctx, addr3).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 2), sdk.NewInt64Coin("foocoin", 5)}))
	
}

func TestSendKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()
	
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)
	
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	coinKeeper := NewKeeper(accountMapper)
	sendKeeper := NewSendKeeper(accountMapper)
	
	addr := sdk.AccAddress([]byte("addr1"))
	addr2 := sdk.AccAddress([]byte("addr2"))
	addr3 := sdk.AccAddress([]byte("addr3"))
	acc := accountMapper.NewAccountWithAddress(ctx, addr)
	
	// Test GetCoins/SetCoins
	accountMapper.SetAccount(ctx, acc)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))
	
	coinKeeper.SetCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	
	// Test HasCoins
	require.True(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	require.False(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)}))
	require.False(t, sendKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 5)}))
	
	coinKeeper.SetCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)})
	
	// Test SendCoins
	sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	
	_, err2 := sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 50)})
	assert.Implements(t, (*sdk.Error)(nil), err2)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	
	coinKeeper.AddCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 30)})
	sendKeeper.SendCoins(ctx, addr, addr2, sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 5)})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 5)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 10)}))
	
	// Test InputOutputCoins
	input1 := NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("foocoin", 2)})
	output1 := NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 2)})
	sendKeeper.InputOutputCoins(ctx, []Input{input1}, []Output{output1})
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 20), sdk.NewInt64Coin("foocoin", 7)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 10), sdk.NewInt64Coin("foocoin", 8)}))
	
	inputs := []Input{
		NewInput(addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 3)}),
		NewInput(addr2, sdk.Coins{sdk.NewInt64Coin("barcoin", 3), sdk.NewInt64Coin("foocoin", 2)}),
	}
	
	outputs := []Output{
		NewOutput(addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 1)}),
		NewOutput(addr3, sdk.Coins{sdk.NewInt64Coin("barcoin", 2), sdk.NewInt64Coin("foocoin", 5)}),
	}
	sendKeeper.InputOutputCoins(ctx, inputs, outputs)
	require.True(t, sendKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 21), sdk.NewInt64Coin("foocoin", 4)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr2).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 7), sdk.NewInt64Coin("foocoin", 6)}))
	require.True(t, sendKeeper.GetCoins(ctx, addr3).IsEqual(sdk.Coins{sdk.NewInt64Coin("barcoin", 2), sdk.NewInt64Coin("foocoin", 5)}))
	
}

func TestViewKeeper(t *testing.T) {
	ms, authKey := setupMultiStore()
	
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)
	
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	coinKeeper := NewKeeper(accountMapper)
	viewKeeper := NewViewKeeper(accountMapper)
	
	addr := sdk.AccAddress([]byte("addr1"))
	acc := accountMapper.NewAccountWithAddress(ctx, addr)
	
	// Test GetCoins/SetCoins
	accountMapper.SetAccount(ctx, acc)
	require.True(t, viewKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{}))
	
	coinKeeper.SetCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)})
	require.True(t, viewKeeper.GetCoins(ctx, addr).IsEqual(sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	
	// Test HasCoins
	require.True(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}))
	require.True(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 5)}))
	require.False(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("foocoin", 15)}))
	require.False(t, viewKeeper.HasCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin("barcoin", 5)}))
}

// ------------------------------
func TestGetAssetWallet(t *testing.T) {
	ms, authKey := setupMultiStore()
	cdc := wire.NewCodec()
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	auth.RegisterBaseAccount(cdc)
	// assetKeeper := NewKeeper(accountMapper)
	
	addr := sdk.AccAddress([]byte("addr1"))
	require.Equal(t, sdk.AssetPegWallet{}, getAssetWallet(ctx, accountMapper, addr))
}

func TestSetAssetWallet(t *testing.T) {
	ms, authKey := setupMultiStore()
	cdc := wire.NewCodec()
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	auth.RegisterBaseAccount(cdc)
	addr := sdk.AccAddress([]byte("addr1"))
	acc := accountMapper.NewAccountWithAddress(ctx, addr)
	accountMapper.SetAccount(ctx, acc)
	asset := sdk.AssetPegWallet{sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "ADFC", AssetType: "sdf", AssetQuantity: 12, QuantityUnit: "asd", OwnerAddress: nil, Locked: false}}
	require.Nil(t, setAssetWallet(ctx, accountMapper, addr, asset))
	require.Equal(t, asset, getAssetWallet(ctx, accountMapper, addr))
	
}
func setupMultiStore1() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	authKey := sdk.NewKVStoreKey("authkey")
	orderKey := sdk.NewKVStoreKey("order")
	negotiationKey := sdk.NewKVStoreKey("negotitation")
	aclKey := sdk.NewKVStoreKey("acl")
	reputationKey := sdk.NewKVStoreKey("reputation")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(negotiationKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(orderKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(aclKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(reputationKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, authKey, orderKey, negotiationKey, aclKey, reputationKey
}

/*
func TestBankKeeper(t *testing.T) {
	ms, authKey, orderKey, negotiationKey, aclKey, reputationKey := setupMultiStore1()
	cdc := wire.NewCodec()
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	aclMapper := acl.NewACLMapper(cdc, aclKey, sdk.ProtoBaseACLAccount)
	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	orderMapper := order.NewMapper(cdc, orderKey, sdk.ProtoBaseOrder)
	negotiationMapper := negotiation.NewMapper(cdc, negotiationKey, sdk.ProtoBaseNegotiation)

	auth.RegisterBaseAccount(cdc)
	acl.RegisterWire(cdc)
	acl.RegisterACLAccount(cdc)
	order.RegisterOrder(cdc)
	negotiation.RegisterNegotiation(cdc)

	keeper := NewKeeper(accountMapper)
	aclKeeper := acl.NewKeeper(aclMapper)
	orderKeeper := order.NewKeeper(orderMapper)
	negotiationKeeper := negotiation.NewKeeper(negotiationMapper, accountMapper)
	reputationKeeper := reputation.NewKeeper(reputation.NewMapper(cdc, reputationKey, sdk.ProtoBaseAccountReputation))
	addr1 := sdk.AccAddress([]byte("IssuerAddress"))
	addr2 := sdk.AccAddress([]byte("SellerAddress"))
	addr3 := sdk.AccAddress([]byte("BuyerAddress"))
	addr4 := sdk.AccAddress([]byte("ZoneAddress"))
	addr5 := sdk.AccAddress([]byte("OwnerAddress"))

	assetPeg := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "ADEC", AssetType: "asdf", AssetQuantity: 234, QuantityUnit: "MT", OwnerAddress: addr5, Locked: false}
	asset := sdk.AssetPegWallet{sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "ADEC", AssetType: "asdf", AssetQuantity: 234, QuantityUnit: "MT", OwnerAddress: addr5, Locked: true}}
	releasedAsset := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "ADEC", AssetType: "asdf", AssetQuantity: 234, QuantityUnit: "MT", OwnerAddress: addr5, Locked: false}
	fiatPeg := &sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "ADEC", TransactionAmount: 123, RedeemedAmount: 0, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 0}}}
	fiat := sdk.FiatPegWallet{sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "ADEC", TransactionAmount: 123, RedeemedAmount: 0, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 0}}}}
	fiatt := sdk.FiatPegWallet{sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "ADEC", TransactionAmount: 50, RedeemedAmount: 0, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 0}}}}

	defineZone := []DefineZone{
		DefineZone{
			From:   addr1,
			To:     addr4,
			ZoneID: []byte("ASDF"),
		},
	}
	keeper.DefineZones(ctx, aclKeeper, defineZone)
	acl1 := &sdk.BaseACLAccount{
		Address:        addr4,
		ACL:            sdk.ACL{MainIssueAssets: true, MainIssueFiats: true, SendAssets: true, SendFiats: true, BuyerExecuteOrder: true},
		OrganizationID: nil,
		ZoneID:         []byte("ASDF"),
	}
	defineACL := []DefineACL{
		DefineACL{
			From:       addr4,
			To:         addr1,
			ACLAccount: acl1,
		},
		DefineACL{
			From:       addr4,
			To:         addr2,
			ACLAccount: acl1,
		},
		DefineACL{
			From:       addr4,
			To:         addr3,
			ACLAccount: acl1,
		},
		DefineACL{
			From:       addr4,
			To:         addr5,
			ACLAccount: acl1,
		},
	}

	keeper.DefineACLs(ctx, aclKeeper, defineACL)

	acc := accountMapper.NewAccountWithAddress(ctx, addr1)
	accountMapper.SetAccount(ctx, acc)
	setAssetWallet(ctx, accountMapper, addr1, asset)
	fmt.Println(getAssetWallet(ctx, accountMapper, addr1))
	setFiatWallet(ctx, accountMapper, addr1, fiat)

	keeper.IssueAssetsToWallets(ctx, []IssueAsset{NewIssueAsset(addr1, addr2, assetPeg)}, aclKeeper)
	require.Equal(t, asset, getAssetWallet(ctx, accountMapper, addr2))
	_, err, _ := keeper.IssueAssetsToWallets(ctx, []IssueAsset{NewIssueAsset(addr1, addr3, assetPeg)}, aclKeeper)
	require.NotNil(t, err)
	fmt.Println(getAssetWallet(ctx, accountMapper, addr2))

	keeper.IssueFiatsToWallets(ctx, []IssueFiat{NewIssueFiat(addr1, addr3, fiatPeg)}, aclKeeper)
	require.Equal(t, fiat, getFiatWallet(ctx, accountMapper, addr3))
	require.NotEqual(t, fiat, getFiatWallet(ctx, accountMapper, addr1))
	fmt.Println(getFiatWallet(ctx, accountMapper, addr3))

	_, err, _ = keeper.IssueFiatsToWallets(ctx, []IssueFiat{NewIssueFiat(addr1, addr2, fiatPeg)}, aclKeeper)
	require.NotNil(t, err)

	releaseAsset := []ReleaseAsset{
		ReleaseAsset{
			ZoneAddress:  addr4,
			OwnerAddress: addr2,
			PegHash:      getPegHash(1),
		},
	}
	_, _ = keeper.ReleaseLockedAssets(ctx, releaseAsset, aclKeeper)
	//send Asset to order
	negotiation := sdk.BaseNegotiation{
		NegotiationID: sdk.NegotiationID(append(append(addr3.Bytes(), addr2.Bytes()...), getPegHash(1).Bytes()...)),
		BuyerAddress:  addr3,
		SellerAddress: addr2,
	}
	negotiationMapper.SetNegotiation(ctx, &negotiation)
	_, negotiation1 := negotiationKeeper.GetNegotiation(ctx, addr3, addr2, getPegHash(1))

	order := sdk.BaseOrder{
		NegotiationID: sdk.NegotiationID(append(append(addr3.Bytes(), addr2.Bytes()...), getPegHash(1).Bytes()...)),
	}
	orderMapper.SetOrder(ctx, &order)

	sendAsset := NewSendAsset(addr2, addr3, getPegHash(1))
	_, _, AssetPeg := keeper.SendAssetsToWallets(ctx, orderKeeper, negotiationKeeper, []SendAsset{sendAsset}, aclKeeper, reputationKeeper)
	_, err, _ = keeper.SendAssetsToWallets(ctx, orderKeeper, negotiationKeeper, []SendAsset{NewSendAsset(addr2, addr3, getPegHash(1))}, aclKeeper, reputationKeeper)
	require.Equal(t, releasedAsset, AssetPeg[0])
	require.NotNil(t, err)

	//query order
	orderAcc := orderMapper.GetOrder(ctx, negotiation1.GetNegotiationID())
	require.NotEqual(t, asset, orderAcc.GetAssetPegWallet())

	//send fiat to order
	sendFiat := NewSendFiat(addr3, addr2, getPegHash(1), 50)
	_, _, FiatPeg := keeper.SendFiatsToWallets(ctx, orderKeeper, negotiationKeeper, []SendFiat{sendFiat}, aclKeeper, reputationKeeper)
	_, err, _ = keeper.SendFiatsToWallets(ctx, orderKeeper, negotiationKeeper, []SendFiat{NewSendFiat(addr3, addr2, getPegHash(1), 100)}, aclKeeper, reputationKeeper)
	require.Equal(t, fiatt, FiatPeg[0])
	require.NotNil(t, err)

	//query order
	orderAcc = orderMapper.GetOrder(ctx, negotiation1.GetNegotiationID())
	require.Equal(t, fiatt, orderAcc.GetFiatPegWallet())

	//execute order
	executeOrder := NewBuyerExecuteOrder(addr1, addr3, addr2, getPegHash(1), "")
	keeper.BuyerExecuteTradeOrders(ctx, negotiationKeeper, orderKeeper, []BuyerExecuteOrder{executeOrder}, aclKeeper, reputationKeeper)
	require.Equal(t, fiatt, getFiatWallet(ctx, accountMapper, addr2))

	executeSellerOrder := NewSellerExecuteOrder(addr1, addr3, addr2, getPegHash(1), "AWEProof")
	keeper.SellerExecuteTradeOrders(ctx, negotiationKeeper, orderKeeper, []SellerExecuteOrder{executeSellerOrder}, aclKeeper, reputationKeeper)
	require.Equal(t, sdk.AssetPegWallet{*releasedAsset}, getAssetWallet(ctx, accountMapper, addr3))

	orderAcc = orderMapper.GetOrder(ctx, negotiation1.GetNegotiationID())
	fmt.Println(orderAcc)
}
*/
