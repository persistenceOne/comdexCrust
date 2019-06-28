package bank

import (
	"fmt"
	"strconv"
	"testing"
	
	"github.com/comdex-blockchain/x/reputation"
	
	"github.com/comdex-blockchain/store"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/acl"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/negotiation"
	"github.com/comdex-blockchain/x/order"
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
		// accountMapper.SetAccount(ctx, baseAccount[i])
	}
	
	var coins = []sdk.Coins{
		{sdk.NewInt64Coin("atom", 10)},
		{sdk.NewInt64Coin("atom", 5)},
		{sdk.NewInt64Coin("atom", 10)},
		{sdk.NewInt64Coin("atom", -5)}, // This needs to be checked, -5 should not be added and -5 should not be subtracted
	}
	var inputs = []Input{
		{addr[0], coins[0]},
		{addr[1], coins[1]},
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
	// require.Equal(t, addr[3].String(), string(res.Tags[0].Value))
	require.Equal(t, sdk.Tags(nil), res.Tags)
}

func TestMsgBankIssueAssets(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, _, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
		sdk.AccAddress([]byte("addr2")),
		sdk.AccAddress([]byte("addr3")),
	}
	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address: addr[1],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAssets: true}},
		{
			Address: addr[2],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAssets: false}},
		{
			Address: addr[3],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAssets: true}},
	}
	
	var baseAccount = []auth.BaseAccount{}
	
	for i, address := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		baseAccount[i].AccountNumber = int64(i)
		accountMapper.SetAccount(ctx, &baseAccount[i])
		if i <= len(aclAccount) && i != 0 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i-1])
		}
	}
	
	var assets = []sdk.BaseAssetPeg{
		{
			AssetType:     "gold",
			AssetQuantity: 10,
		},
	}
	
	for i := 0; i < 1; i++ {
		pegHash, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
		baseAccount[0].AssetPegWallet = append(baseAccount[0].AssetPegWallet, sdk.NewBaseAssetPegWithPegHash(pegHash))
	}
	accountMapper.SetAccount(ctx, &baseAccount[0])
	pegHash, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(0)))
	
	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	defineZone := NewDefineZone(addr[1], addr[1], zoneID)
	msgDefineZone := NewMsgDefineZones([]DefineZone{defineZone})
	var res = fun(ctx, msgDefineZone)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	defineZone = NewDefineZone(addr[0], addr[0], zoneID)
	msgDefineZone = NewMsgDefineZones([]DefineZone{defineZone})
	res = fun(ctx, msgDefineZone)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	issueAsset := NewIssueAsset(addr[0], addr[1], &assets[0])
	var issueAssets = []IssueAsset{issueAsset}
	msgBankIssueAssets := NewMsgBankIssueAssets(issueAssets)
	res = fun(ctx, msgBankIssueAssets)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
	baseAccount0 := accountMapper.GetAccount(ctx, baseAccount[0].Address)
	baseAccount1 := accountMapper.GetAccount(ctx, baseAccount[1].Address)
	require.Equal(t, sdk.AssetPegWallet(nil), baseAccount0.GetAssetPegWallet())
	require.Equal(t, pegHash, baseAccount1.GetAssetPegWallet()[0].PegHash)
	require.Equal(t, assets[0].AssetType, baseAccount1.GetAssetPegWallet()[0].AssetType)
	require.Equal(t, assets[0].AssetQuantity, baseAccount1.GetAssetPegWallet()[0].AssetQuantity)
	require.Equal(t, true, baseAccount1.GetAssetPegWallet()[0].Locked)
	
	issueAsset = NewIssueAsset(addr[1], addr[2], &assets[0])
	issueAssets = []IssueAsset{issueAsset}
	msgBankIssueAssets = NewMsgBankIssueAssets(issueAssets)
	res = fun(ctx, msgBankIssueAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	issueAsset = NewIssueAsset(addr[0], addr[2], &assets[0])
	issueAssets = []IssueAsset{issueAsset}
	msgBankIssueAssets = NewMsgBankIssueAssets(issueAssets)
	res = fun(ctx, msgBankIssueAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	issueAsset = NewIssueAsset(addr[0], addr[3], &assets[0])
	issueAssets = []IssueAsset{issueAsset}
	msgBankIssueAssets = NewMsgBankIssueAssets(issueAssets)
	res = fun(ctx, msgBankIssueAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
}

func TestMsgBankRedeemAssets(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, _, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
		sdk.AccAddress([]byte("addr2")),
		sdk.AccAddress([]byte("addr3")),
	}
	
	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address: addr[1],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAssets: true, RedeemAssets: true}},
		{
			Address: addr[2],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAssets: true, RedeemAssets: false}},
		{
			Address: addr[3],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{RedeemAssets: true}},
	}
	
	var baseAccount = []auth.BaseAccount{}
	
	for i, address := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		baseAccount[i].AccountNumber = int64(i)
		accountMapper.SetAccount(ctx, &baseAccount[i])
		if i <= len(aclAccount) && i != 0 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i-1])
		}
	}
	
	var assets = []sdk.BaseAssetPeg{
		{
			AssetType:     "sona",
			AssetQuantity: 10,
		},
	}
	
	for i := 0; i < 1; i++ {
		pegHash, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
		baseAccount[0].AssetPegWallet = append(baseAccount[0].AssetPegWallet, sdk.NewBaseAssetPegWithPegHash(pegHash))
	}
	accountMapper.SetAccount(ctx, &baseAccount[0])
	pegHash, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(0)))
	
	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	defineZone := NewDefineZone(addr[0], addr[0], zoneID)
	msgDefineZone := NewMsgDefineZones([]DefineZone{defineZone})
	var res = fun(ctx, msgDefineZone)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	issueAsset := NewIssueAsset(addr[0], addr[1], &assets[0])
	var issueAssets = []IssueAsset{issueAsset}
	msgBankIssueAssets := NewMsgBankIssueAssets(issueAssets)
	res = fun(ctx, msgBankIssueAssets)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
	baseAccount0 := accountMapper.GetAccount(ctx, baseAccount[0].Address)
	baseAccount1 := accountMapper.GetAccount(ctx, baseAccount[1].Address)
	require.Equal(t, sdk.AssetPegWallet(nil), baseAccount0.GetAssetPegWallet())
	require.Equal(t, pegHash, baseAccount1.GetAssetPegWallet()[0].PegHash)
	require.Equal(t, assets[0].AssetType, baseAccount1.GetAssetPegWallet()[0].AssetType)
	require.Equal(t, assets[0].AssetQuantity, baseAccount1.GetAssetPegWallet()[0].AssetQuantity)
	require.Equal(t, true, baseAccount1.GetAssetPegWallet()[0].Locked)
	
	redeemAsset := NewRedeemAsset(addr[0], addr[1], sdk.PegHash([]byte("10")))
	redeemAssets := []RedeemAsset{redeemAsset}
	msgBankRedeemAssets := NewMsgBankRedeemAssets(redeemAssets)
	res = fun(ctx, msgBankRedeemAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	redeemAsset = NewRedeemAsset(addr[0], addr[1], pegHash)
	redeemAssets = []RedeemAsset{redeemAsset}
	msgBankRedeemAssets = NewMsgBankRedeemAssets(redeemAssets)
	res = fun(ctx, msgBankRedeemAssets)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
	baseAccount0 = accountMapper.GetAccount(ctx, baseAccount[0].Address)
	baseAccount1 = accountMapper.GetAccount(ctx, baseAccount[1].Address)
	require.Equal(t, pegHash, baseAccount0.GetAssetPegWallet()[0].PegHash)
	require.Equal(t, "", baseAccount0.GetAssetPegWallet()[0].AssetType)
	require.Equal(t, int64(0), baseAccount0.GetAssetPegWallet()[0].AssetQuantity)
	require.Equal(t, false, baseAccount0.GetAssetPegWallet()[0].Locked)
	require.Equal(t, sdk.AssetPegWallet(nil), baseAccount1.GetAssetPegWallet())
	
	redeemAsset = NewRedeemAsset(addr[1], addr[2], assets[0].GetPegHash())
	redeemAssets = []RedeemAsset{redeemAsset}
	msgBankRedeemAssets = NewMsgBankRedeemAssets(redeemAssets)
	res = fun(ctx, msgBankRedeemAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	redeemAsset = NewRedeemAsset(addr[0], addr[2], assets[0].GetPegHash())
	redeemAssets = []RedeemAsset{redeemAsset}
	msgBankRedeemAssets = NewMsgBankRedeemAssets(redeemAssets)
	res = fun(ctx, msgBankRedeemAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	redeemAsset = NewRedeemAsset(addr[0], addr[3], assets[0].GetPegHash())
	redeemAssets = []RedeemAsset{redeemAsset}
	msgBankRedeemAssets = NewMsgBankRedeemAssets(redeemAssets)
	res = fun(ctx, msgBankRedeemAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
}

func TestMsgBankIssueFiat(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, _, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
		sdk.AccAddress([]byte("addr2")),
	}
	
	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address: addr[1],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiats: true}},
		{
			Address: addr[2],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiats: false}},
	}
	
	var baseAccount = []auth.BaseAccount{}
	
	for i, address := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		baseAccount[i].AccountNumber = int64(i)
		accountMapper.SetAccount(ctx, &baseAccount[i])
		if i <= len(aclAccount) && i != 0 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i-1])
		}
	}
	
	var fiats = []sdk.BaseFiatPeg{
		{
			TransactionID:     "AA",
			TransactionAmount: 10,
		},
	}
	
	for i := 0; i < 1; i++ {
		pegHash, _ := sdk.GetFiatPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
		baseAccount[0].FiatPegWallet = append(baseAccount[0].FiatPegWallet, sdk.NewBaseFiatPegWithPegHash(pegHash))
	}
	accountMapper.SetAccount(ctx, &baseAccount[0])
	pegHash, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(0)))
	
	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	defineZone := NewDefineZone(addr[0], addr[0], zoneID)
	msgDefineZone := NewMsgDefineZones([]DefineZone{defineZone})
	var res = fun(ctx, msgDefineZone)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	issueFiat := NewIssueFiat(addr[0], addr[1], &fiats[0])
	issueFiats := []IssueFiat{issueFiat}
	msgBankIssueFiat := NewMsgBankIssueFiats(issueFiats)
	res = fun(ctx, msgBankIssueFiat)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
	baseAccount0 := accountMapper.GetAccount(ctx, baseAccount[0].Address)
	baseAccount1 := accountMapper.GetAccount(ctx, baseAccount[1].Address)
	require.Equal(t, sdk.FiatPegWallet(nil), baseAccount0.GetFiatPegWallet())
	require.Equal(t, pegHash, baseAccount1.GetFiatPegWallet()[0].PegHash)
	require.Equal(t, fiats[0].TransactionAmount, baseAccount1.GetFiatPegWallet()[0].TransactionAmount)
	require.Equal(t, fiats[0].TransactionID, baseAccount1.GetFiatPegWallet()[0].TransactionID)
	
	issueFiat = NewIssueFiat(addr[1], addr[2], &fiats[0])
	issueFiats = []IssueFiat{issueFiat}
	msgBankIssueFiat = NewMsgBankIssueFiats(issueFiats)
	res = fun(ctx, msgBankIssueFiat)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	issueFiat = NewIssueFiat(addr[0], addr[2], &fiats[0])
	issueFiats = []IssueFiat{issueFiat}
	msgBankIssueFiat = NewMsgBankIssueFiats(issueFiats)
	res = fun(ctx, msgBankIssueFiat)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	issueFiat = NewIssueFiat(addr[0], addr[1], &fiats[0])
	issueFiats = []IssueFiat{issueFiat}
	msgBankIssueFiat = NewMsgBankIssueFiats(issueFiats)
	res = fun(ctx, msgBankIssueFiat)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
}

func TestMsgBankRedeemFiat(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, _, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
		sdk.AccAddress([]byte("addr2")),
		sdk.AccAddress([]byte("addr3")),
		sdk.AccAddress([]byte("addr4")),
	}
	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address: addr[1],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiats: true, RedeemFiats: true}},
		{
			Address: addr[2],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiats: true, RedeemFiats: true}},
		{
			Address: addr[3],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiats: true, RedeemFiats: true}},
		{
			Address: addr[4],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiats: true, RedeemFiats: false}},
	}
	
	var baseAccount = []auth.BaseAccount{}
	
	for i, address := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		baseAccount[i].AccountNumber = int64(i)
		accountMapper.SetAccount(ctx, &baseAccount[i])
		if i <= len(aclAccount) && i != 0 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i-1])
		}
	}
	var fiats = []sdk.BaseFiatPeg{
		{
			TransactionID:     "A1",
			TransactionAmount: 1000,
		},
		{
			TransactionID:     "A2",
			TransactionAmount: 1000,
		},
		{
			TransactionID:     "B1",
			TransactionAmount: 100,
		},
		{
			TransactionID:     "B2",
			TransactionAmount: 100,
		},
		{
			TransactionID:     "C1",
			TransactionAmount: 100,
		},
		{
			TransactionID:     "C2",
			TransactionAmount: 100,
		},
	}
	
	pegHash := make([]sdk.PegHash, 0, 10)
	for i := 0; i < 10; i++ {
		peg, _ := sdk.GetFiatPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
		pegHash = append(pegHash, peg)
		baseAccount[0].FiatPegWallet = append(baseAccount[0].FiatPegWallet, sdk.NewBaseFiatPegWithPegHash(pegHash[i]))
	}
	accountMapper.SetAccount(ctx, &baseAccount[0])
	
	fun := NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	defineZone := NewDefineZone(addr[0], addr[0], zoneID)
	msgDefineZone := NewMsgDefineZones([]DefineZone{defineZone})
	var res = fun(ctx, msgDefineZone)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	{
		issueFiat0 := NewIssueFiat(addr[0], addr[1], &fiats[0])
		issueFiat1 := NewIssueFiat(addr[0], addr[1], &fiats[1])
		issueFiats := []IssueFiat{issueFiat0, issueFiat1}
		msgBankIssueFiat := NewMsgBankIssueFiats(issueFiats)
		res = fun(ctx, msgBankIssueFiat)
		require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
		require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
		baseAccount0 := accountMapper.GetAccount(ctx, baseAccount[0].Address)
		baseAccount1 := accountMapper.GetAccount(ctx, baseAccount[1].Address)
		require.Equal(t, int(8), len(baseAccount0.GetFiatPegWallet()))
		require.Equal(t, int(2), len(baseAccount1.GetFiatPegWallet()))
		require.Equal(t, pegHash[9], baseAccount1.GetFiatPegWallet()[0].PegHash)
		require.Equal(t, pegHash[8], baseAccount1.GetFiatPegWallet()[1].PegHash)
		require.Equal(t, fiats[0].TransactionAmount, baseAccount1.GetFiatPegWallet()[0].TransactionAmount)
		require.Equal(t, fiats[1].TransactionAmount, baseAccount1.GetFiatPegWallet()[1].TransactionAmount)
		require.Equal(t, fiats[0].TransactionID, baseAccount1.GetFiatPegWallet()[0].TransactionID)
		require.Equal(t, fiats[1].TransactionID, baseAccount1.GetFiatPegWallet()[1].TransactionID)
		
		redeemFiat := NewRedeemFiat(addr[1], addr[0], 500)
		redeemFiats := []RedeemFiat{redeemFiat}
		msgBankRedeemFiats := NewMsgBankRedeemFiats(redeemFiats)
		res = fun(ctx, msgBankRedeemFiats)
		require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
		baseAccount0 = accountMapper.GetAccount(ctx, baseAccount[0].Address)
		baseAccount1 = accountMapper.GetAccount(ctx, baseAccount[1].Address)
		require.Equal(t, int(8), len(baseAccount0.GetFiatPegWallet()))
		require.Equal(t, int(2), len(baseAccount1.GetFiatPegWallet()))
		require.Equal(t, pegHash[9], baseAccount1.GetFiatPegWallet()[0].PegHash)
		require.Equal(t, pegHash[8], baseAccount1.GetFiatPegWallet()[1].PegHash)
		require.Equal(t, int64(500), baseAccount1.GetFiatPegWallet()[0].RedeemedAmount)
		require.Equal(t, int64(0), baseAccount1.GetFiatPegWallet()[1].RedeemedAmount)
		require.Equal(t, int64(500), baseAccount1.GetFiatPegWallet()[0].TransactionAmount)
		require.Equal(t, fiats[1].TransactionAmount, baseAccount1.GetFiatPegWallet()[1].TransactionAmount)
		require.Equal(t, fiats[0].TransactionID, baseAccount1.GetFiatPegWallet()[0].TransactionID)
		require.Equal(t, fiats[1].TransactionID, baseAccount1.GetFiatPegWallet()[1].TransactionID)
		
		redeemFiat = NewRedeemFiat(addr[1], addr[0], 1250)
		redeemFiats = []RedeemFiat{redeemFiat}
		msgBankRedeemFiats = NewMsgBankRedeemFiats(redeemFiats)
		res = fun(ctx, msgBankRedeemFiats)
		require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
		baseAccount0 = accountMapper.GetAccount(ctx, baseAccount[0].Address)
		baseAccount1 = accountMapper.GetAccount(ctx, baseAccount[1].Address)
		require.Equal(t, int(8), len(baseAccount0.GetFiatPegWallet()))
		require.Equal(t, int(1), len(baseAccount1.GetFiatPegWallet()))
		require.Equal(t, pegHash[8], baseAccount1.GetFiatPegWallet()[0].PegHash)
		require.Equal(t, int64(750), baseAccount1.GetFiatPegWallet()[0].RedeemedAmount)
		require.Equal(t, int64(250), baseAccount1.GetFiatPegWallet()[0].TransactionAmount)
		require.Equal(t, fiats[1].TransactionID, baseAccount1.GetFiatPegWallet()[0].TransactionID)
	}
	
	{
		issueFiat0 := NewIssueFiat(addr[0], addr[2], &fiats[2])
		issueFiat1 := NewIssueFiat(addr[0], addr[2], &fiats[3])
		issueFiats := []IssueFiat{issueFiat0, issueFiat1}
		msgBankIssueFiat := NewMsgBankIssueFiats(issueFiats)
		res = fun(ctx, msgBankIssueFiat)
		require.Equal(t, addr[2].String(), string(res.Tags[0].Value))
		require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
		baseAccount0 := accountMapper.GetAccount(ctx, baseAccount[0].Address)
		baseAccount2 := accountMapper.GetAccount(ctx, baseAccount[2].Address)
		require.Equal(t, int(6), len(baseAccount0.GetFiatPegWallet()))
		require.Equal(t, int(2), len(baseAccount2.GetFiatPegWallet()))
		require.Equal(t, pegHash[7], baseAccount2.GetFiatPegWallet()[0].PegHash)
		require.Equal(t, pegHash[6], baseAccount2.GetFiatPegWallet()[1].PegHash)
		require.Equal(t, fiats[2].TransactionAmount, baseAccount2.GetFiatPegWallet()[0].TransactionAmount)
		require.Equal(t, fiats[3].TransactionAmount, baseAccount2.GetFiatPegWallet()[1].TransactionAmount)
		require.Equal(t, fiats[2].TransactionID, baseAccount2.GetFiatPegWallet()[0].TransactionID)
		require.Equal(t, fiats[3].TransactionID, baseAccount2.GetFiatPegWallet()[1].TransactionID)
		
		redeemFiat := NewRedeemFiat(addr[2], addr[0], 100)
		redeemFiats := []RedeemFiat{redeemFiat}
		msgBankRedeemFiats := NewMsgBankRedeemFiats(redeemFiats)
		res = fun(ctx, msgBankRedeemFiats)
		require.Equal(t, addr[2].String(), string(res.Tags[0].Value))
		baseAccount0 = accountMapper.GetAccount(ctx, baseAccount[0].Address)
		baseAccount2 = accountMapper.GetAccount(ctx, baseAccount[2].Address)
		require.Equal(t, int(6), len(baseAccount0.GetFiatPegWallet()))
		require.Equal(t, int(1), len(baseAccount2.GetFiatPegWallet()))
		require.Equal(t, pegHash[6], baseAccount2.GetFiatPegWallet()[0].PegHash)
		require.Equal(t, fiats[3].TransactionAmount, baseAccount2.GetFiatPegWallet()[0].TransactionAmount)
		require.Equal(t, fiats[3].TransactionID, baseAccount2.GetFiatPegWallet()[0].TransactionID)
	}
	
	{
		issueFiat0 := NewIssueFiat(addr[0], addr[3], &fiats[4])
		issueFiat1 := NewIssueFiat(addr[0], addr[3], &fiats[5])
		issueFiats := []IssueFiat{issueFiat0, issueFiat1}
		msgBankIssueFiat := NewMsgBankIssueFiats(issueFiats)
		res = fun(ctx, msgBankIssueFiat)
		require.Equal(t, addr[3].String(), string(res.Tags[0].Value))
		require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
		baseAccount0 := accountMapper.GetAccount(ctx, baseAccount[0].Address)
		baseAccount3 := accountMapper.GetAccount(ctx, baseAccount[3].Address)
		require.Equal(t, int(4), len(baseAccount0.GetFiatPegWallet()))
		require.Equal(t, int(2), len(baseAccount3.GetFiatPegWallet()))
		require.Equal(t, pegHash[5], baseAccount3.GetFiatPegWallet()[0].PegHash)
		require.Equal(t, pegHash[4], baseAccount3.GetFiatPegWallet()[1].PegHash)
		require.Equal(t, fiats[4].TransactionAmount, baseAccount3.GetFiatPegWallet()[0].TransactionAmount)
		require.Equal(t, fiats[5].TransactionAmount, baseAccount3.GetFiatPegWallet()[1].TransactionAmount)
		require.Equal(t, fiats[4].TransactionID, baseAccount3.GetFiatPegWallet()[0].TransactionID)
		require.Equal(t, fiats[5].TransactionID, baseAccount3.GetFiatPegWallet()[1].TransactionID)
		
		redeemFiat := NewRedeemFiat(addr[3], addr[0], 200)
		redeemFiats := []RedeemFiat{redeemFiat}
		msgBankRedeemFiats := NewMsgBankRedeemFiats(redeemFiats)
		res = fun(ctx, msgBankRedeemFiats)
		require.Equal(t, addr[3].String(), string(res.Tags[0].Value))
		baseAccount0 = accountMapper.GetAccount(ctx, baseAccount[0].Address)
		baseAccount3 = accountMapper.GetAccount(ctx, baseAccount[3].Address)
		require.Equal(t, int(4), len(baseAccount0.GetFiatPegWallet()))
		require.Equal(t, int(0), len(baseAccount3.GetFiatPegWallet()))
	}
	
	{
		redeemFiat := NewRedeemFiat(addr[2], addr[0], 120)
		redeemFiats := []RedeemFiat{redeemFiat}
		msgBankRedeemFiats := NewMsgBankRedeemFiats(redeemFiats)
		res = fun(ctx, msgBankRedeemFiats)
		require.Equal(t, sdk.Tags(nil), res.Tags)
		
		redeemFiat = NewRedeemFiat(addr[4], addr[1], 100)
		redeemFiats = []RedeemFiat{redeemFiat}
		msgBankRedeemFiats = NewMsgBankRedeemFiats(redeemFiats)
		res = fun(ctx, msgBankRedeemFiats)
		require.Equal(t, sdk.Tags(nil), res.Tags)
		
		redeemFiat = NewRedeemFiat(addr[4], addr[0], 120)
		redeemFiats = []RedeemFiat{redeemFiat}
		msgBankRedeemFiats = NewMsgBankRedeemFiats(redeemFiats)
		res = fun(ctx, msgBankRedeemFiats)
		require.Equal(t, sdk.Tags(nil), res.Tags)
		
	}
	
}

func TestMsgBankSendAsset(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, negoMapper, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
		sdk.AccAddress([]byte("addr2")),
	}
	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address: addr[1],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAssets: true, SendAssets: false, ReleaseAssets: false}},
	}
	
	var baseAccount = []auth.BaseAccount{}
	
	for i, address := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		baseAccount[i].AccountNumber = int64(i)
		accountMapper.SetAccount(ctx, &baseAccount[i])
		if i <= len(aclAccount) && i != 0 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i-1])
		}
	}
	
	var assets = []sdk.BaseAssetPeg{
		{
			AssetType:     "gold",
			AssetQuantity: 10,
		},
	}
	
	pegHash, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(0)))
	baseAccount[0].AssetPegWallet = append(baseAccount[0].AssetPegWallet, sdk.NewBaseAssetPegWithPegHash(pegHash))
	accountMapper.SetAccount(ctx, &baseAccount[0])
	
	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	nego := sdk.BaseNegotiation{
		NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(addr[2], addr[1], pegHash)),
		BuyerAddress:    addr[2],
		SellerAddress:   addr[1],
		PegHash:         pegHash,
		Bid:             0,
		Time:            -1,
		BuyerSignature:  []byte("addr2"),
		SellerSignature: []byte("addr1"),
	}
	
	defineZone := NewDefineZone(addr[0], addr[0], zoneID)
	msgDefineZone := NewMsgDefineZones([]DefineZone{defineZone})
	var res = fun(ctx, msgDefineZone)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	sendAsset := NewSendAsset(addr[2], addr[1], pegHash)
	sendAssets := []SendAsset{sendAsset}
	msgBankSendAsset := NewMsgBankSendAssets(sendAssets)
	res = fun(ctx, msgBankSendAsset)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	sendAsset = NewSendAsset(addr[1], addr[2], pegHash)
	sendAssets = []SendAsset{sendAsset}
	msgBankSendAsset = NewMsgBankSendAssets(sendAssets)
	res = fun(ctx, msgBankSendAsset)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	aclAccount[0].ACL.SendAssets = true
	aclMapper.SetAccount(ctx, addr[1], &aclAccount[0])
	
	sendAsset = NewSendAsset(addr[1], addr[2], pegHash)
	sendAssets = []SendAsset{sendAsset}
	msgBankSendAsset = NewMsgBankSendAssets(sendAssets)
	res = fun(ctx, msgBankSendAsset)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	negoMapper.SetNegotiation(ctx, &nego)
	
	sendAsset = NewSendAsset(addr[1], addr[2], pegHash)
	sendAssets = []SendAsset{sendAsset}
	msgBankSendAsset = NewMsgBankSendAssets(sendAssets)
	res = fun(ctx, msgBankSendAsset)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	nego.Time = 0
	negoMapper.SetNegotiation(ctx, &nego)
	
	sendAsset = NewSendAsset(addr[1], addr[2], pegHash)
	sendAssets = []SendAsset{sendAsset}
	msgBankSendAsset = NewMsgBankSendAssets(sendAssets)
	res = fun(ctx, msgBankSendAsset)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	releaseAsset := NewReleaseAsset(addr[1], addr[1], pegHash)
	releaseAssets := []ReleaseAsset{releaseAsset}
	msgBankReleaseAssets := NewMsgBankReleaseAssets(releaseAssets)
	res = fun(ctx, msgBankReleaseAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	releaseAsset = NewReleaseAsset(addr[0], addr[1], pegHash)
	releaseAssets = []ReleaseAsset{releaseAsset}
	msgBankReleaseAssets = NewMsgBankReleaseAssets(releaseAssets)
	res = fun(ctx, msgBankReleaseAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	aclAccount[0].ACL.ReleaseAssets = true
	aclMapper.SetAccount(ctx, addr[1], &aclAccount[0])
	
	releaseAsset = NewReleaseAsset(addr[0], addr[1], pegHash)
	releaseAssets = []ReleaseAsset{releaseAsset}
	msgBankReleaseAssets = NewMsgBankReleaseAssets(releaseAssets)
	res = fun(ctx, msgBankReleaseAssets)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	issueAsset := NewIssueAsset(addr[0], addr[1], &assets[0])
	issueAssets := []IssueAsset{issueAsset}
	msgBankIssueAssets := NewMsgBankIssueAssets(issueAssets)
	res = fun(ctx, msgBankIssueAssets)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
	
	sendAsset = NewSendAsset(addr[1], addr[2], pegHash)
	sendAssets = []SendAsset{sendAsset}
	msgBankSendAsset = NewMsgBankSendAssets(sendAssets)
	res = fun(ctx, msgBankSendAsset)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	releaseAsset = NewReleaseAsset(addr[0], addr[1], pegHash)
	releaseAssets = []ReleaseAsset{releaseAsset}
	msgBankReleaseAssets = NewMsgBankReleaseAssets(releaseAssets)
	res = fun(ctx, msgBankReleaseAssets)
	baseAccount1 := accountMapper.GetAccount(ctx, baseAccount[1].Address)
	require.Equal(t, false, baseAccount1.GetAssetPegWallet()[0].Locked)
	
	sendAsset = NewSendAsset(addr[1], addr[2], pegHash)
	sendAssets = []SendAsset{sendAsset}
	msgBankSendAsset = NewMsgBankSendAssets(sendAssets)
	res = fun(ctx, msgBankSendAsset)
	require.Equal(t, addr[2].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[1].String(), string(res.Tags[1].Value))
	baseAccount1 = accountMapper.GetAccount(ctx, baseAccount[1].Address)
	require.Equal(t, sdk.AssetPegWallet(nil), baseAccount1.GetAssetPegWallet())
	_, assetPegWallet, _, _, _ := orderKeeper.GetOrderDetails(ctx, addr[2], addr[1], pegHash)
	require.Equal(t, pegHash, assetPegWallet[0].PegHash)
	require.Equal(t, assets[0].AssetType, assetPegWallet[0].AssetType)
	require.Equal(t, assets[0].AssetQuantity, assetPegWallet[0].AssetQuantity)
}

func TestMsgBankSendFiat(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, negoMapper, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
		sdk.AccAddress([]byte("addr2")),
	}
	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address: addr[2],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiats: true, SendFiats: false}},
	}
	
	var baseAccount = []auth.BaseAccount{}
	
	for i := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		accountMapper.SetAccount(ctx, &baseAccount[i])
	}
	aclMapper.SetAccount(ctx, addr[2], &aclAccount[0])
	
	var fiats = []sdk.BaseFiatPeg{
		{
			TransactionID:     "A1",
			TransactionAmount: 1000,
		},
	}
	
	for i := 0; i < 1; i++ {
		pegHash, _ := sdk.GetFiatPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
		baseAccount[0].FiatPegWallet = append(baseAccount[0].FiatPegWallet, sdk.NewBaseFiatPegWithPegHash(pegHash))
	}
	accountMapper.SetAccount(ctx, &baseAccount[0])
	negoHash, _ := sdk.GetFiatPegHashHex(fmt.Sprintf("%x", strconv.Itoa(0)))
	
	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	nego := sdk.BaseNegotiation{
		NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(addr[2], addr[1], negoHash)),
		BuyerAddress:    addr[2],
		SellerAddress:   addr[1],
		PegHash:         negoHash,
		Bid:             0,
		Time:            -1,
		BuyerSignature:  []byte("addr2"),
		SellerSignature: []byte("addr1"),
	}
	
	defineZone := NewDefineZone(addr[0], addr[0], zoneID)
	msgDefineZone := NewMsgDefineZones([]DefineZone{defineZone})
	var res = fun(ctx, msgDefineZone)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	sendFiat := NewSendFiat(addr[1], addr[2], negoHash, 750)
	sendFiats := []SendFiat{sendFiat}
	msgBankSendFiat := NewMsgBankSendFiats(sendFiats)
	res = fun(ctx, msgBankSendFiat)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	sendFiat = NewSendFiat(addr[2], addr[1], negoHash, 750)
	sendFiats = []SendFiat{sendFiat}
	msgBankSendFiat = NewMsgBankSendFiats(sendFiats)
	res = fun(ctx, msgBankSendFiat)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	aclAccount[0].ACL.SendFiats = true
	aclMapper.SetAccount(ctx, addr[2], &aclAccount[0])
	
	sendFiat = NewSendFiat(addr[2], addr[1], negoHash, 750)
	sendFiats = []SendFiat{sendFiat}
	msgBankSendFiat = NewMsgBankSendFiats(sendFiats)
	res = fun(ctx, msgBankSendFiat)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	negoMapper.SetNegotiation(ctx, &nego)
	
	sendFiat = NewSendFiat(addr[2], addr[1], negoHash, 750)
	sendFiats = []SendFiat{sendFiat}
	msgBankSendFiat = NewMsgBankSendFiats(sendFiats)
	res = fun(ctx, msgBankSendFiat)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	nego.Time = 0
	negoMapper.SetNegotiation(ctx, &nego)
	
	sendFiat = NewSendFiat(addr[2], addr[1], negoHash, 750)
	sendFiats = []SendFiat{sendFiat}
	msgBankSendFiat = NewMsgBankSendFiats(sendFiats)
	res = fun(ctx, msgBankSendFiat)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	issueFiat := NewIssueFiat(addr[0], addr[2], &fiats[0])
	issueFiats := []IssueFiat{issueFiat}
	msgBankIssueFiats := NewMsgBankIssueFiats(issueFiats)
	res = fun(ctx, msgBankIssueFiats)
	require.Equal(t, addr[2].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
	
	sendFiat = NewSendFiat(addr[2], addr[1], negoHash, 750)
	sendFiats = []SendFiat{sendFiat}
	msgBankSendFiat = NewMsgBankSendFiats(sendFiats)
	res = fun(ctx, msgBankSendFiat)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[2].String(), string(res.Tags[1].Value))
	baseAccount2 := accountMapper.GetAccount(ctx, baseAccount[2].Address)
	require.Equal(t, int64(250), baseAccount2.GetFiatPegWallet()[0].TransactionAmount)
	_, _, fiatPegWallet, _, _ := orderKeeper.GetOrderDetails(ctx, addr[2], addr[1], negoHash)
	require.Equal(t, negoHash, fiatPegWallet[0].PegHash)
	require.Equal(t, int64(750), fiatPegWallet[0].TransactionAmount)
	require.Equal(t, "A1", fiatPegWallet[0].TransactionID)
}

func TestMsgBankExecuteOrders(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, orderMapper, orderKeeper, negoMapper, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
		sdk.AccAddress([]byte("addr2")),
	}
	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address: addr[1],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueAssets: true, SendAssets: true, ReleaseAssets: true, SellerExecuteOrder: false}},
		{
			Address: addr[2],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{IssueFiats: true, SendFiats: true, BuyerExecuteOrder: false}},
	}
	var baseAccount = []auth.BaseAccount{}
	
	for i := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		accountMapper.SetAccount(ctx, &baseAccount[i])
	}
	aclMapper.SetAccount(ctx, addr[1], &aclAccount[0])
	aclMapper.SetAccount(ctx, addr[2], &aclAccount[1])
	
	var assets = []sdk.BaseAssetPeg{
		{
			AssetType:     "gold",
			AssetQuantity: 10,
		},
		{
			AssetType:     "silver",
			AssetQuantity: 100,
		},
	}
	
	var fiats = []sdk.BaseFiatPeg{
		{
			TransactionID:     "A1",
			TransactionAmount: 1000,
		},
		{
			TransactionID:     "B1",
			TransactionAmount: 1000,
		},
	}
	pegHash := make([]sdk.PegHash, 0, 2)
	for i := 0; i < 2; i++ {
		peg1, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
		peg2, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i+5)))
		pegHash = append(pegHash, peg1)
		baseAccount[0].AssetPegWallet = append(baseAccount[0].AssetPegWallet, sdk.NewBaseAssetPegWithPegHash(peg1))
		baseAccount[0].FiatPegWallet = append(baseAccount[0].FiatPegWallet, sdk.NewBaseFiatPegWithPegHash(peg2))
	}
	accountMapper.SetAccount(ctx, &baseAccount[0])
	
	nego := sdk.BaseNegotiation{
		NegotiationID:   sdk.NegotiationID(sdk.GenerateNegotiationIDBytes(addr[2], addr[1], pegHash[1])),
		BuyerAddress:    addr[2],
		SellerAddress:   addr[1],
		PegHash:         pegHash[1],
		Bid:             25000,
		Time:            0,
		BuyerSignature:  []byte("addr2"),
		SellerSignature: []byte("addr1"),
	}
	
	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	defineZone := NewDefineZone(addr[0], addr[0], zoneID)
	msgDefineZone := NewMsgDefineZones([]DefineZone{defineZone})
	var res = fun(ctx, msgDefineZone)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	{
		buyerExecuteOrder := NewBuyerExecuteOrder(addr[1], addr[2], addr[1], pegHash[1], "fiatProofHash")
		buyerExecuteOrders := []BuyerExecuteOrder{buyerExecuteOrder}
		msgBankBuyerExecuteOrders := NewMsgBankBuyerExecuteOrders(buyerExecuteOrders)
		res = fun(ctx, msgBankBuyerExecuteOrders)
		require.Equal(t, sdk.Tags(nil), res.Tags)
		sellerExecuteOrder := NewSellerExecuteOrder(addr[1], addr[2], addr[1], pegHash[1], "awbProofHash")
		sellerExecuteOrders := []SellerExecuteOrder{sellerExecuteOrder}
		msgBankSellerExecuteOrders := NewMsgBankSellerExecuteOrders(sellerExecuteOrders)
		res = fun(ctx, msgBankSellerExecuteOrders)
		require.Equal(t, sdk.Tags(nil), res.Tags)
	}
	
	{
		buyerExecuteOrder := NewBuyerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "fiatProofHash")
		buyerExecuteOrders := []BuyerExecuteOrder{buyerExecuteOrder}
		msgBankBuyerExecuteOrders := NewMsgBankBuyerExecuteOrders(buyerExecuteOrders)
		res = fun(ctx, msgBankBuyerExecuteOrders)
		require.Equal(t, sdk.Tags(nil), res.Tags)
		sellerExecuteOrder := NewSellerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "awbProofHash")
		sellerExecuteOrders := []SellerExecuteOrder{sellerExecuteOrder}
		msgBankSellerExecuteOrders := NewMsgBankSellerExecuteOrders(sellerExecuteOrders)
		res = fun(ctx, msgBankSellerExecuteOrders)
		require.Equal(t, sdk.Tags(nil), res.Tags)
	}
	
	aclAccount[0].ACL.SellerExecuteOrder = true
	aclAccount[1].ACL.BuyerExecuteOrder = true
	aclMapper.SetAccount(ctx, addr[1], &aclAccount[0])
	aclMapper.SetAccount(ctx, addr[2], &aclAccount[1])
	
	{
		buyerExecuteOrder := NewBuyerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "fiatProofHash")
		buyerExecuteOrders := []BuyerExecuteOrder{buyerExecuteOrder}
		msgBankBuyerExecuteOrders := NewMsgBankBuyerExecuteOrders(buyerExecuteOrders)
		res = fun(ctx, msgBankBuyerExecuteOrders)
		require.Equal(t, sdk.Tags(nil), res.Tags)
		sellerExecuteOrder := NewSellerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "awbProofHash")
		sellerExecuteOrders := []SellerExecuteOrder{sellerExecuteOrder}
		msgBankSellerExecuteOrders := NewMsgBankSellerExecuteOrders(sellerExecuteOrders)
		res = fun(ctx, msgBankSellerExecuteOrders)
		require.Equal(t, sdk.Tags(nil), res.Tags)
	}
	
	negoMapper.SetNegotiation(ctx, &nego)
	
	issueFiat1 := NewIssueFiat(addr[0], addr[2], &fiats[0])
	issueFiat2 := NewIssueFiat(addr[0], addr[2], &fiats[1])
	issueFiats := []IssueFiat{issueFiat1, issueFiat2}
	msgBankIssueFiats := NewMsgBankIssueFiats(issueFiats)
	res = fun(ctx, msgBankIssueFiats)
	require.Equal(t, addr[2].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
	
	sendFiat1 := NewSendFiat(addr[2], addr[1], pegHash[1], 750)
	sendFiat2 := NewSendFiat(addr[2], addr[1], pegHash[1], 600)
	sendFiats := []SendFiat{sendFiat1, sendFiat2}
	msgBankSendFiat := NewMsgBankSendFiats(sendFiats)
	res = fun(ctx, msgBankSendFiat)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[2].String(), string(res.Tags[1].Value))
	
	_, assetPegWallet, fiatPegWallet, awbProofHash, fiatProofHash := orderKeeper.GetOrderDetails(ctx, addr[2], addr[1], pegHash[1])
	order := orderMapper.GetOrder(ctx, nego.NegotiationID)
	
	buyerExecuteOrder := NewBuyerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "")
	buyerExecuteOrders := []BuyerExecuteOrder{buyerExecuteOrder}
	msgBankBuyerExecuteOrders := NewMsgBankBuyerExecuteOrders(buyerExecuteOrders)
	res = fun(ctx, msgBankBuyerExecuteOrders)
	require.Equal(t, "false", string(res.Tags[3].Value))
	sellerExecuteOrder := NewSellerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "")
	sellerExecuteOrders := []SellerExecuteOrder{sellerExecuteOrder}
	msgBankSellerExecuteOrders := NewMsgBankSellerExecuteOrders(sellerExecuteOrders)
	res = fun(ctx, msgBankSellerExecuteOrders)
	require.Equal(t, "false", string(res.Tags[3].Value))
	
	issueAsset := NewIssueAsset(addr[0], addr[1], &assets[0])
	issueAssets := []IssueAsset{issueAsset}
	msgBankIssueAssets := NewMsgBankIssueAssets(issueAssets)
	res = fun(ctx, msgBankIssueAssets)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[0].String(), string(res.Tags[1].Value))
	
	releaseAsset := NewReleaseAsset(addr[0], addr[1], pegHash[1])
	releaseAssets := []ReleaseAsset{releaseAsset}
	msgBankReleaseAssets := NewMsgBankReleaseAssets(releaseAssets)
	res = fun(ctx, msgBankReleaseAssets)
	baseAccount1 := accountMapper.GetAccount(ctx, baseAccount[1].Address)
	require.Equal(t, false, baseAccount1.GetAssetPegWallet()[0].Locked)
	
	sendAsset := NewSendAsset(addr[1], addr[2], pegHash[1])
	sendAssets := []SendAsset{sendAsset}
	msgBankSendAsset := NewMsgBankSendAssets(sendAssets)
	res = fun(ctx, msgBankSendAsset)
	require.Equal(t, addr[2].String(), string(res.Tags[0].Value))
	require.Equal(t, addr[1].String(), string(res.Tags[1].Value))
	
	_, assetPegWallet, fiatPegWallet, awbProofHash, fiatProofHash = orderKeeper.GetOrderDetails(ctx, addr[2], addr[1], pegHash[1])
	fmt.Println(assetPegWallet, fiatPegWallet, awbProofHash, fiatProofHash)
	order.SetAssetPegWallet(assetPegWallet)
	orderMapper.SetOrder(ctx, order)
	_, assetPegWallet, fiatPegWallet, awbProofHash, fiatProofHash = orderKeeper.GetOrderDetails(ctx, addr[2], addr[1], pegHash[1])
	
	buyerExecuteOrder = NewBuyerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "")
	buyerExecuteOrders = []BuyerExecuteOrder{buyerExecuteOrder}
	msgBankBuyerExecuteOrders = NewMsgBankBuyerExecuteOrders(buyerExecuteOrders)
	res = fun(ctx, msgBankBuyerExecuteOrders)
	require.Equal(t, "false", string(res.Tags[3].Value))
	sellerExecuteOrder = NewSellerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "")
	sellerExecuteOrders = []SellerExecuteOrder{sellerExecuteOrder}
	msgBankSellerExecuteOrders = NewMsgBankSellerExecuteOrders(sellerExecuteOrders)
	res = fun(ctx, msgBankSellerExecuteOrders)
	require.Equal(t, "false", string(res.Tags[3].Value))
	
	nego.Time = -10
	negoMapper.SetNegotiation(ctx, &nego)
	
	buyerExecuteOrder = NewBuyerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "")
	buyerExecuteOrders = []BuyerExecuteOrder{buyerExecuteOrder}
	msgBankBuyerExecuteOrders = NewMsgBankBuyerExecuteOrders(buyerExecuteOrders)
	res = fun(ctx, msgBankBuyerExecuteOrders)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	sellerExecuteOrder = NewSellerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "")
	sellerExecuteOrders = []SellerExecuteOrder{sellerExecuteOrder}
	msgBankSellerExecuteOrders = NewMsgBankSellerExecuteOrders(sellerExecuteOrders)
	res = fun(ctx, msgBankSellerExecuteOrders)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	nego.Time = 0
	negoMapper.SetNegotiation(ctx, &nego)
	
	nego.Bid = 250
	negoMapper.SetNegotiation(ctx, &nego)
	orderMapper.SetOrder(ctx, order)
	
	buyerExecuteOrder = NewBuyerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "fiatProofHash")
	buyerExecuteOrders = []BuyerExecuteOrder{buyerExecuteOrder}
	msgBankBuyerExecuteOrders = NewMsgBankBuyerExecuteOrders(buyerExecuteOrders)
	res = fun(ctx, msgBankBuyerExecuteOrders)
	require.Equal(t, "false", string(res.Tags[3].Value))
	sellerExecuteOrder = NewSellerExecuteOrder(addr[0], addr[2], addr[1], pegHash[1], "awbProofHash")
	sellerExecuteOrders = []SellerExecuteOrder{sellerExecuteOrder}
	msgBankSellerExecuteOrders = NewMsgBankSellerExecuteOrders(sellerExecuteOrders)
	res = fun(ctx, msgBankSellerExecuteOrders)
	require.Equal(t, "true", string(res.Tags[3].Value))
}

func TestMsgDefineOrganization(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, _, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	organizationID, _ := sdk.GetOrganizationIDFromString("ABCD1234")
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
	}
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address:        addr[1],
			OrganizationID: organizationID,
		},
	}
	var baseAccount = []auth.BaseAccount{}
	for i, address := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		baseAccount[i].AccountNumber = int64(i)
		accountMapper.SetAccount(ctx, &baseAccount[i])
		if i <= len(aclAccount) && i != 0 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i-1])
		}
	}
	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	defineOrganization := NewDefineOrganization(addr[1], addr[1], organizationID)
	msgDefineOrganization := NewMsgDefineOrganizations([]DefineOrganization{defineOrganization})
	var res = fun(ctx, msgDefineOrganization)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	defineOrganization = NewDefineOrganization(addr[0], addr[0], organizationID)
	msgDefineOrganization = NewMsgDefineOrganizations([]DefineOrganization{defineOrganization})
	res = fun(ctx, msgDefineOrganization)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
}

func TestMsgDefineACLs(t *testing.T) {
	_, ctx, accountMapper, coinKeeper, _, orderKeeper, _, negoKeeper, aclMapper, aclKeeper, reputationKeeper := setup()
	
	zoneID, _ := sdk.GetZoneIDFromString("ABCD1234")
	organizationID, _ := sdk.GetOrganizationIDFromString("ABCD1234")
	
	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("genesis")),
		sdk.AccAddress([]byte("addr1")),
	}
	var aclAccount = []sdk.BaseACLAccount{
		{
			Address:        addr[1],
			OrganizationID: organizationID,
			ZoneID:         zoneID,
			ACL:            sdk.ACL{IssueAssets: true},
		},
	}
	var baseAccount = []auth.BaseAccount{}
	for i, address := range addr {
		baseAccount = append(baseAccount, auth.NewBaseAccountWithAddress(addr[i]))
		baseAccount[i].AccountNumber = int64(i)
		accountMapper.SetAccount(ctx, &baseAccount[i])
		if i <= len(aclAccount) && i != 0 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i-1])
		}
	}
	var fun = NewAssetFiatHandler(coinKeeper, negoKeeper, orderKeeper, aclKeeper, reputationKeeper)
	
	defineOrganization := NewDefineOrganization(addr[0], addr[0], organizationID)
	msgDefineOrganization := NewMsgDefineOrganizations([]DefineOrganization{defineOrganization})
	res := fun(ctx, msgDefineOrganization)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	defineZone := NewDefineZone(addr[0], addr[0], zoneID)
	msgDefineZone := NewMsgDefineZones([]DefineZone{defineZone})
	res = fun(ctx, msgDefineZone)
	require.Equal(t, addr[0].String(), string(res.Tags[0].Value))
	require.Equal(t, "ABCD1234", string(res.Tags[1].Value))
	
	defineACL := NewDefineACL(addr[1], addr[1], &aclAccount[0])
	msgDefineACL := NewMsgDefineACLs([]DefineACL{defineACL})
	res = fun(ctx, msgDefineACL)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	defineACL = NewDefineACL(addr[0], addr[1], &aclAccount[0])
	msgDefineACL = NewMsgDefineACLs([]DefineACL{defineACL})
	res = fun(ctx, msgDefineACL)
	require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	
}
