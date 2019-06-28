package ibc

import (
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	
	"github.com/comdex-blockchain/store"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/acl"
	"github.com/comdex-blockchain/x/assetFactory"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/bank"
	"github.com/comdex-blockchain/x/fiatFactory"
	"github.com/comdex-blockchain/x/negotiation"
	"github.com/comdex-blockchain/x/order"
)

func setup() (sdk.Context, Mapper, auth.AccountMapper, bank.Keeper, order.Mapper, order.Keeper, negotiation.Mapper, negotiation.Keeper, acl.Mapper, acl.Keeper, assetFactory.AssetPegMapper, assetFactory.Keeper, fiatFactory.FiatPegMapper, fiatFactory.Keeper) {
	db := dbm.NewMemDB()
	
	authKey := sdk.NewKVStoreKey("authKey")
	orderKey := sdk.NewKVStoreKey("orderKey")
	negoKey := sdk.NewKVStoreKey("negoKey")
	aclKey := sdk.NewKVStoreKey("aclKey")
	ibcKey := sdk.NewKVStoreKey("ibcKey")
	assetKey := sdk.NewKVStoreKey("assetKey")
	fiatKey := sdk.NewKVStoreKey("fiatKey")
	
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(orderKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(negoKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(aclKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(assetKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fiatKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	
	cdc := wire.NewCodec()
	auth.RegisterBaseAccount(cdc)
	negotiation.RegisterNegotiation(cdc)
	order.RegisterOrder(cdc)
	acl.RegisterWire(cdc)
	assetFactory.RegisterAssetPeg(cdc)
	fiatFactory.RegisterFiatPeg(cdc)
	RegisterWire(cdc)
	
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	
	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	coinKeeper := bank.NewKeeper(accountMapper)
	
	assetMapper := assetFactory.NewAssetPegMapper(cdc, assetKey, sdk.ProtoBaseAssetPeg)
	assetKeeper := assetFactory.NewKeeper(assetMapper)
	
	fiatMapper := fiatFactory.NewFiatPegMapper(cdc, fiatKey, sdk.ProtoBaseFiatPeg)
	fiatKeeper := fiatFactory.NewKeeper(fiatMapper)
	
	orderMapper := order.NewMapper(cdc, orderKey, sdk.ProtoBaseOrder)
	orderKeeper := order.NewKeeper(orderMapper)
	
	negoMapper := negotiation.NewMapper(cdc, negoKey, sdk.ProtoBaseNegotiation)
	negoKeeper := negotiation.NewKeeper(negoMapper, accountMapper)
	
	aclMapper := acl.NewACLMapper(cdc, aclKey, sdk.ProtoBaseACLAccount)
	aclKeeper := acl.NewKeeper(aclMapper)
	
	ibcMapper := NewMapper(cdc, ibcKey, sdk.CodespaceUndefined)
	
	return ctx, ibcMapper, accountMapper, coinKeeper, orderMapper, orderKeeper, negoMapper, negoKeeper, aclMapper, aclKeeper, assetMapper, assetKeeper, fiatMapper, fiatKeeper
}

func setAccount(ctx sdk.Context, accountMapper auth.AccountMapper, baseAccount auth.Account, addr sdk.AccAddress) {
	baseAccount = accountMapper.GetAccount(ctx, addr)
	accountMapper.SetAccount(ctx, baseAccount)
	return
	
}

func setupSetCoins(ctx sdk.Context, coinKeeper bank.Keeper, addr sdk.AccAddress, denom string, coins int64) {
	coinKeeper.SetCoins(ctx, addr, sdk.Coins{sdk.NewInt64Coin(denom, coins)})
	return
}
