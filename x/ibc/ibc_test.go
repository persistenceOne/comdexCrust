package ibc

import (
	"testing"

	"github.com/commitHub/commitBlockchain/x/acl"
	"github.com/commitHub/commitBlockchain/x/negotiation"
	"github.com/commitHub/commitBlockchain/x/order"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/commitHub/commitBlockchain/store"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"

	"github.com/commitHub/commitBlockchain/x/auth"
	"github.com/commitHub/commitBlockchain/x/bank"
)

// AccountMapper(/Keeper) and IBCMapper should use different StoreKey later
func defaultContext(key sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func newAddress() sdk.AccAddress {
	return sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
}

func getCoins(ck bank.Keeper, ctx sdk.Context, addr sdk.AccAddress) (sdk.Coins, sdk.Error) {
	zero := sdk.Coins(nil)
	coins, _, err := ck.AddCoins(ctx, addr, zero)
	return coins, err
}

func makeCodec() *wire.Codec {
	var cdc = wire.NewCodec()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(bank.MsgSend{}, "test/ibc/Send", nil)
	cdc.RegisterConcrete(bank.MsgIssue{}, "test/ibc/Issue", nil)
	cdc.RegisterConcrete(IBCTransferMsg{}, "test/ibc/IBCTransferMsg", nil)
	cdc.RegisterConcrete(IBCReceiveMsg{}, "test/ibc/IBCReceiveMsg", nil)
	cdc.RegisterInterface((*sdk.AssetPeg)(nil), nil)
	cdc.RegisterConcrete(&sdk.BaseAssetPeg{}, "commit-blockchain/AssetPeg", nil)
	cdc.RegisterInterface((*sdk.FiatPeg)(nil), nil)
	cdc.RegisterConcrete(&sdk.BaseFiatPeg{}, "commit-blockchain/FiatPeg", nil)
	// Register AppAccount
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/ibc/Account", nil)
	wire.RegisterCrypto(cdc)

	cdc.Seal()

	return cdc
}

func TestIBC(t *testing.T) {
	cdc := makeCodec()
	key := sdk.NewKVStoreKey("ibc")
	ctx := defaultContext(key)
	authKey := sdk.NewKVStoreKey("authKey")
	orderKey := sdk.NewKVStoreKey("orderKey")
	negoKey := sdk.NewKVStoreKey("negoKey")

	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	am := auth.NewAccountMapper(cdc, key, auth.ProtoBaseAccount)
	aclMapper := acl.NewACLMapper(cdc, key, sdk.ProtoBaseACLAccount)
	orderMapper := order.NewMapper(cdc, orderKey, sdk.ProtoBaseOrder)
	negoMapper := negotiation.NewMapper(cdc, negoKey, sdk.ProtoBaseNegotiation)

	ck := bank.NewKeeper(am)
	aclKeeper := acl.NewKeeper(aclMapper)
	orderKeeper := order.NewKeeper(orderMapper)
	negoKeeper := negotiation.NewKeeper(negoMapper, accountMapper)

	src := newAddress()
	dest := newAddress()
	chainid := "ibcchain"
	zero := sdk.Coins(nil)
	mycoins := sdk.Coins{sdk.NewInt64Coin("mycoin", 10)}

	coins, _, err := ck.AddCoins(ctx, src, mycoins)
	require.Nil(t, err)
	require.Equal(t, mycoins, coins)

	ibcm := NewMapper(cdc, key, DefaultCodespace)
	h := NewHandler(ibcm, ck, aclKeeper, negoKeeper, orderKeeper)
	packet := IBCPacket{
		SrcAddr:   src,
		DestAddr:  dest,
		Coins:     mycoins,
		SrcChain:  chainid,
		DestChain: chainid,
	}

	store := ctx.KVStore(key)

	var msg sdk.Msg
	var res sdk.Result
	var egl int64
	var igs int64

	egl = ibcm.getEgressLength(store, chainid)
	require.Equal(t, egl, int64(0))

	msg = IBCTransferMsg{
		IBCPacket: packet,
	}
	res = h(ctx, msg)
	require.True(t, res.IsOK())

	coins, err = getCoins(ck, ctx, src)
	require.Nil(t, err)
	require.Equal(t, zero, coins)

	egl = ibcm.getEgressLength(store, chainid)
	require.Equal(t, egl, int64(1))

	igs = ibcm.GetIngressSequence(ctx, chainid)
	require.Equal(t, igs, int64(0))

	msg = IBCReceiveMsg{
		IBCPacket: packet,
		Relayer:   src,
		Sequence:  0,
	}
	res = h(ctx, msg)
	require.True(t, res.IsOK())

	coins, err = getCoins(ck, ctx, dest)
	require.Nil(t, err)
	require.Equal(t, mycoins, coins)

	igs = ibcm.GetIngressSequence(ctx, chainid)
	require.Equal(t, igs, int64(1))

	res = h(ctx, msg)
	require.False(t, res.IsOK())

	igs = ibcm.GetIngressSequence(ctx, chainid)
	require.Equal(t, igs, int64(1))

	asset := MsgIssueAssets{[]IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)}}
	err = ibcm.PostIBCMsgIssueAssetsPacket(ctx, asset)
	require.Nil(t, err)

	packet1 := IBCTransferMsg{
		IBCPacket{
			SrcAddr:   src,
			DestAddr:  dest,
			Coins:     mycoins,
			SrcChain:  chainid,
			DestChain: chainid,
		},
	}
	err = ibcm.PostIBCTransferMsg(ctx, packet1)
	require.Nil(t, err)

	fiatPeg := MsgIssueFiats{[]IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)}}
	err = ibcm.PostIBCMsgIssueFiatsPacket(ctx, fiatPeg)
	require.Nil(t, err)
	/*
		order := MsgRedeemOrders{[]RedeemOrder{NewRedeemOrder(mediatorAddress, issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)}}
		err = ibcm.PostIBCMsgRedeemOrders(ctx, order)
	*/
	require.Nil(t, err)
}

func TestIBCHandlerAssets(t *testing.T) {
	ctx, ibcMapper, accountMapper, coinKeeper, _, orderKeeper, _, negoKeeper, aclMapper, aclKeeper, assetMapper, assetKeeper, _, _ := setup()

	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("addr0")),
		sdk.AccAddress([]byte("addr1")),
	}
	zoneID, _ := sdk.GetZoneIDFromString("zoneID")
	aclMapper.SetZone(ctx, addr[0], zoneID)

	var aclAccount = []sdk.BaseACLAccount{
		sdk.BaseACLAccount{
			Address: addr[0],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{MainIssueAssets: true, MainRedeemAssets: true}},
		sdk.BaseACLAccount{
			Address: addr[1],
			ACL:     sdk.ACL{MainIssueAssets: true, MainRedeemAssets: true}},
	}
	setcoins := []int64{0, 0}

	var baseAccount = []auth.Account{}
	//var aclAccounts = []sdk.ACLAccount{}

	for i, address := range addr {
		setupSetCoins(ctx, coinKeeper, address, "atom", setcoins[i])
		baseAccount = append(baseAccount, accountMapper.GetAccount(ctx, address))
		accountMapper.SetAccount(ctx, baseAccount[i])
		if i <= len(aclAccount)-1 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i])
			//aclAccounts = append(aclAccounts, aclMapper.GetAccount(ctx, address))
		}
	}

	chainid := "ibcchain"
	assets := []sdk.BaseAssetPeg{
		sdk.BaseAssetPeg{
			PegHash:       sdk.PegHash([]byte("test")),
			DocumentHash:  "AA",
			AssetType:     "sona",
			AssetQuantity: 10,
			OwnerAddress:  addr[0],
		},
		sdk.BaseAssetPeg{
			PegHash:       sdk.PegHash([]byte("tip2")),
			AssetType:     "silver",
			AssetQuantity: 5,
			OwnerAddress:  addr[1],
		},
	}
	assetMapper.SetAssetPeg(ctx, &assets[0])
	var assetPegWallets = [][]sdk.BaseAssetPeg{{assets[0]}, {assets[1]}}
	for i, assetPeg := range assetPegWallets {
		baseAccount[i].SetAssetPegWallet(assetPeg)
		accountMapper.SetAccount(ctx, baseAccount[i])
	}

	issueAsset := NewIssueAsset(addr[0], addr[1], &assets[0], chainid, chainid)
	redeemAsset := NewRedeemAsset(addr[0], addr[1], assets[0].GetPegHash(), chainid, chainid)
	sendAsset := NewSendAsset(addr[0], addr[1], assets[0].GetPegHash(), chainid, chainid)
	//issueFiat := NewIssueFiat(addr[0], addr[1], &fiats[0], chainid, chainid)

	fun := NewHandler(ibcMapper, coinKeeper, aclKeeper, negoKeeper, orderKeeper)

	/*
		var issueAssets = []IssueAsset{issueAsset}
		msgIssueAssets := NewMsgIssueAssets(issueAssets)

		fun := NewHandler(ibcMapper, coinKeeper, aclKeeper, negoKeeper, orderKeeper)
		var res = fun(ctx, msgIssueAssets)
		baseAccount[0] = accountMapper.GetAccount(ctx, addr[0])
		baseAccount[1] = accountMapper.GetAccount(ctx, addr[1])
		require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	*/
	//RelayIssueAsset
	fun = NewAssetHandler(ibcMapper, coinKeeper, assetKeeper)
	seq := ibcMapper.GetIngressSequence(ctx, chainid)
	msgRelayIssueAssets := NewMsgRelayIssueAssets([]IssueAsset{issueAsset}, addr[0], seq)
	res := fun(ctx, msgRelayIssueAssets)
	require.True(t, res.IsOK())
	/*
		fun = NewHandler(ibcMapper, coinKeeper, aclKeeper, negoKeeper, orderKeeper)
		msgRedeemAsset := NewMsgRedeemAssets([]RedeemAsset{redeemAsset})
		res = fun(ctx, msgRedeemAsset)
		require.True(t, res.IsOK())
	*/
	//RelayRedeemAsset
	fun = NewAssetHandler(ibcMapper, coinKeeper, assetKeeper)
	seq = ibcMapper.GetIngressSequence(ctx, chainid)
	msgRelayRedeemAssets := NewMsgRelayRedeemAssets([]RedeemAsset{redeemAsset}, addr[0], seq)
	res = fun(ctx, msgRelayRedeemAssets)
	require.True(t, res.IsOK())

	//RelaySendAsset
	fun = NewAssetHandler(ibcMapper, coinKeeper, assetKeeper)
	seq = ibcMapper.GetIngressSequence(ctx, chainid)
	msgRelaySendAssets := NewMsgRelaySendAssets([]SendAsset{sendAsset}, addr[0], seq)
	res = fun(ctx, msgRelaySendAssets)
	require.True(t, res.IsOK())
}
func TestIBCHandlerFiats(t *testing.T) {
	ctx, ibcMapper, accountMapper, coinKeeper, _, _, _, _, aclMapper, _, _, _, fiatMapper, fiatKeeper := setup()

	var addr = []sdk.AccAddress{
		sdk.AccAddress([]byte("addr0")),
		sdk.AccAddress([]byte("addr1")),
	}
	zoneID, _ := sdk.GetZoneIDFromString("zoneID")
	aclMapper.SetZone(ctx, addr[0], zoneID)

	var aclAccount = []sdk.BaseACLAccount{
		sdk.BaseACLAccount{
			Address: addr[0],
			ZoneID:  zoneID,
			ACL:     sdk.ACL{MainIssueFiats: true, MainRedeemFiats: true}},
		sdk.BaseACLAccount{
			Address: addr[1],
			ACL:     sdk.ACL{MainIssueFiats: true, MainRedeemFiats: true}},
	}
	setcoins := []int64{0, 0}

	var baseAccount = []auth.Account{}
	//var aclAccounts = []sdk.ACLAccount{}

	for i, address := range addr {
		setupSetCoins(ctx, coinKeeper, address, "atom", setcoins[i])
		baseAccount = append(baseAccount, accountMapper.GetAccount(ctx, address))
		accountMapper.SetAccount(ctx, baseAccount[i])
		if i <= len(aclAccount)-1 {
			aclMapper.SetAccount(ctx, address, &aclAccount[i])
			//aclAccounts = append(aclAccounts, aclMapper.GetAccount(ctx, address))
		}
	}

	chainid := "ibcchain"
	var fiats = []sdk.BaseFiatPeg{
		sdk.BaseFiatPeg{
			PegHash:           sdk.PegHash([]byte("test")),
			TransactionID:     "one",
			TransactionAmount: 5,
			RedeemedAmount:    0,
			Owners:            []sdk.Owner{{OwnerAddress: addr[0], Amount: 5}},
		},
		sdk.BaseFiatPeg{
			PegHash:           sdk.PegHash([]byte("tip")),
			TransactionID:     "two",
			TransactionAmount: 5,
			RedeemedAmount:    0,
			Owners:            []sdk.Owner{{OwnerAddress: addr[1], Amount: 5}},
		},
	}
	fiatMapper.SetFiatPeg(ctx, &fiats[0])
	fiatMapper.SetFiatPeg(ctx, &fiats[1])
	var fiatPegWallets = [][]sdk.BaseFiatPeg{{fiats[0]}}
	for i, fiatPeg := range fiatPegWallets {
		baseAccount[i].SetFiatPegWallet(fiatPeg)
		accountMapper.SetAccount(ctx, baseAccount[i])
	}

	var fiatPegWallet sdk.FiatPegWallet = sdk.FiatPegWallet{
		sdk.BaseFiatPeg{
			PegHash:           sdk.PegHash([]byte("test")),
			TransactionID:     "one",
			TransactionAmount: 5,
			RedeemedAmount:    1,
			Owners:            []sdk.Owner{{OwnerAddress: addr[0], Amount: 5}},
		},
	}

	issueFiat := NewIssueFiat(addr[0], addr[1], &fiats[0], chainid, chainid)
	redeemFiat := NewRedeemFiat(addr[1], addr[0], 5, fiatPegWallet, chainid, chainid)
	sendFiat := NewSendFiat(addr[0], addr[1], fiats[0].GetPegHash(), 5, fiatPegWallet, chainid, chainid)

	/*
		var issueFiats = []IssueFiat{issueFiat}
		msgIssueFiats := NewMsgIssueFiats(issueFiats)

		fun := NewHandler(ibcMapper, coinKeeper, aclKeeper, negoKeeper, orderKeeper)
		var res = fun(ctx, msgIssueFiats)
		baseAccount[0] = accountMapper.GetAccount(ctx, addr[0])
		baseAccount[1] = accountMapper.GetAccount(ctx, addr[1])
		require.Equal(t, addr[1].String(), string(res.Tags[0].Value))
	*/

	//RelayIssueFiat
	fun := NewFiatHandler(ibcMapper, coinKeeper, fiatKeeper)
	seq := ibcMapper.GetIngressSequence(ctx, chainid)
	msgRelayIssueFiats := NewMsgRelayIssueFiats([]IssueFiat{issueFiat}, addr[0], seq)
	res := fun(ctx, msgRelayIssueFiats)
	require.True(t, res.IsOK())
	//RelayRedeemFiat
	fun = NewFiatHandler(ibcMapper, coinKeeper, fiatKeeper)
	seq = ibcMapper.GetIngressSequence(ctx, chainid)
	msgRelayRedeemFiats := NewMsgRelayRedeemFiats([]RedeemFiat{redeemFiat}, addr[0], seq)
	res = fun(ctx, msgRelayRedeemFiats)
	require.True(t, res.IsOK())

	/*

		fun = NewHandler(ibcMapper, coinKeeper, aclKeeper, negoKeeper, orderKeeper)
		msgRedeemFiat := NewMsgRedeemFiats([]RedeemFiat{redeemFiat})
		res = fun(ctx, msgRedeemFiat)
		require.True(t, res.IsOK())
	*/

	//RelaySendFiat
	fun2 := NewFiatHandler(ibcMapper, coinKeeper, fiatKeeper)
	seq = ibcMapper.GetIngressSequence(ctx, chainid)
	msgRelaySendFiats := NewMsgRelaySendFiats([]SendFiat{sendFiat}, addr[0], seq)
	res = fun2(ctx, msgRelaySendFiats)
	require.True(t, res.IsOK())
}
