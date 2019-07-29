package order

import (
	"fmt"
	"strconv"
	"testing"

	key1 "github.com/commitHub/commitBlockchain/crypto/keys"
	"github.com/commitHub/commitBlockchain/store"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/auth"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	auth.RegisterWire(cdc)
	RegisterOrder(cdc)
	return cdc
}

var (
	cdc      = MakeCodec()
	ctx      = sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	authKey  = sdk.NewKVStoreKey("authKey")
	orderKey = sdk.NewKVStoreKey("order")
)

func setUpMultiStore() sdk.MultiStore {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(orderKey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()
	return ms
}

func accountMapper() auth.AccountMapper {
	am := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	acc := am.NewAccountWithAddress(ctx, sdk.AccAddress("address"))
	am.SetAccount(ctx, acc)
	return am
}

func mapper() (Mapper, Keeper) {
	om := NewMapper(cdc, orderKey, sdk.ProtoBaseOrder)
	orderKeeper := NewKeeper(om)
	return om, orderKeeper
}

var keybase key1.Keybase
var (
	ms              = setUpMultiStore()
	om, orderKeeper = mapper()
	am              = accountMapper()
	fromAddress     = sdk.AccAddress([]byte("FromAddress"))
	toAddress       = sdk.AccAddress([]byte("ToAddress"))
	peghash         = sdk.PegHash([]byte(fmt.Sprintf("%x", strconv.Itoa(1))))
	peghash1        = sdk.PegHash([]byte(fmt.Sprintf("%x", strconv.Itoa(2))))
	assetPeg        = []sdk.BaseAssetPeg{
		{
			PegHash:       peghash,
			DocumentHash:  "DocumentHash",
			AssetType:     "AssetType",
			AssetQuantity: 100,
			AssetPrice:    10,
			QuantityUnit:  "MT",
			OwnerAddress:  nil,
			Locked:        true,
			Moderated:     false,
		},
		{
			PegHash:       peghash,
			DocumentHash:  "DocumentHash",
			AssetType:     "AssetType",
			AssetQuantity: 100,
			AssetPrice:    10,
			QuantityUnit:  "MT",
			OwnerAddress:  nil,
			Locked:        false,
			Moderated:     true,
		},
	}
	fiatPegWallet = sdk.FiatPegWallet{
		{
			PegHash:           peghash,
			TransactionID:     "TansactionID",
			TransactionAmount: 10,
			RedeemedAmount:    10,
			Owners:            nil,
		},
	}
)

func Test_SendAssetsToOrder(t *testing.T) {
	err := orderKeeper.SendAssetsToOrder(ctx, fromAddress, toAddress, &assetPeg[0])
	require.Nil(t, err)
}

func Test_SendFiatToOrder(t *testing.T) {
	err := orderKeeper.SendFiatsToOrder(ctx, fromAddress, toAddress, peghash, fiatPegWallet)
	require.Nil(t, err)
}

func Test_SendFiatsFromOrder(t *testing.T) {
	order := om.NewOrder(fromAddress, toAddress, peghash)

	order.SetAssetPegWallet(sdk.AssetPegWallet{assetPeg[0]})
	order.SetFiatPegWallet(fiatPegWallet)

	om.SetOrder(ctx, order)

	fiatPeg := sdk.FiatPegWallet{}

	fiatWallet := orderKeeper.SendFiatsFromOrder(ctx, fromAddress, toAddress, peghash, fiatPegWallet)
	fmt.Println(om.GetOrder(ctx, sdk.NegotiationID(append(append(fromAddress.Bytes(), toAddress.Bytes()...), peghash.Bytes()...))))
	require.Equal(t, fiatPeg, fiatWallet)
}

func Test_SendAssetFromOrder(t *testing.T) {
	order := om.NewOrder(fromAddress, toAddress, peghash)

	order.SetAssetPegWallet(sdk.AssetPegWallet{assetPeg[0]})
	order.SetFiatPegWallet(fiatPegWallet)
	om.SetOrder(ctx, order)

	assetPeg1 := sdk.AssetPegWallet{}
	assetPegWallet := orderKeeper.SendAssetFromOrder(ctx, fromAddress, toAddress, &assetPeg[0])
	require.Equal(t, assetPeg1, assetPegWallet)
}

func Test_GetOrderDetails(t *testing.T) {
	order := om.NewOrder(fromAddress, toAddress, peghash)

	order.SetAssetPegWallet(sdk.AssetPegWallet{assetPeg[0]})
	order.SetFiatPegWallet(fiatPegWallet)
	order.SetFiatProofHash("FiatProofHash")
	order.SetAWBProofHash("AWBProofHash")
	om.SetOrder(ctx, order)

	err, assetWallet, fiatWallet, fiatProof, awbProof := orderKeeper.GetOrderDetails(ctx, fromAddress, toAddress, peghash)
	err1, _, _, _, _ := orderKeeper.GetOrderDetails(ctx, sdk.AccAddress("address"), fromAddress, peghash)

	require.Nil(t, err)
	require.Equal(t, assetWallet, sdk.AssetPegWallet{assetPeg[0]})
	require.Equal(t, fiatWallet, fiatPegWallet)
	require.Equal(t, fiatProof, "FiatProofHash")
	require.Equal(t, awbProof, "AWBProofHash")
	require.NotNil(t, err1)
}

func Test_SetOrderFiatProofHash(t *testing.T) {
	order := om.NewOrder(fromAddress, toAddress, peghash)
	om.SetOrder(ctx, order)
	orderKeeper.SetOrderFiatProofHash(ctx, fromAddress, toAddress, peghash, "fiatProofHash")
}

func Test_SetOrderAWBProofHash(t *testing.T) {
	order := om.NewOrder(fromAddress, toAddress, peghash)
	om.SetOrder(ctx, order)
	orderKeeper.SetOrderAWBProofHash(ctx, fromAddress, toAddress, peghash, "AWBProofHash")
}
