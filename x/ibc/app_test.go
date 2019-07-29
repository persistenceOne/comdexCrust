package ibc

import (
	"testing"

	"github.com/commitHub/commitBlockchain/x/acl"
	"github.com/commitHub/commitBlockchain/x/negotiation"
	"github.com/commitHub/commitBlockchain/x/order"

	"github.com/stretchr/testify/require"

	"github.com/commitHub/commitBlockchain/store"
	sdk "github.com/commitHub/commitBlockchain/types"

	"github.com/commitHub/commitBlockchain/x/auth"
	"github.com/commitHub/commitBlockchain/x/bank"
	"github.com/commitHub/commitBlockchain/x/mock"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	dbm "github.com/tendermint/tendermint/libs/db"
)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	aclKey := sdk.NewKVStoreKey("aclkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(aclKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, aclKey
}

// initialize the mock application for this module
func getMockApp(t *testing.T) *mock.App {
	mapp := mock.NewApp()
	cdc := makeCodec()
	authKey := sdk.NewKVStoreKey("authKey")
	orderKey := sdk.NewKVStoreKey("orderKey")
	negoKey := sdk.NewKVStoreKey("negoKey")

	RegisterWire(mapp.Cdc)
	keyIBC := sdk.NewKVStoreKey("ibc")
	_, aclKey := setupMultiStore()

	accountMapper := auth.NewAccountMapper(cdc, authKey, auth.ProtoBaseAccount)
	ibcMapper := NewMapper(mapp.Cdc, keyIBC, mapp.RegisterCodespace(DefaultCodespace))
	aclMapper := acl.NewACLMapper(mapp.Cdc, aclKey, sdk.ProtoBaseACLAccount)
	orderMapper := order.NewMapper(cdc, orderKey, sdk.ProtoBaseOrder)
	negoMapper := negotiation.NewMapper(cdc, negoKey, sdk.ProtoBaseNegotiation)

	coinKeeper := bank.NewKeeper(mapp.AccountMapper)
	aclKeeper := acl.NewKeeper(aclMapper)
	orderKeeper := order.NewKeeper(orderMapper)
	negoKeeper := negotiation.NewKeeper(negoMapper, accountMapper)

	mapp.Router().AddRoute("ibc", NewHandler(ibcMapper, coinKeeper, aclKeeper, negoKeeper, orderKeeper))

	require.NoError(t, mapp.CompleteSetup([]*sdk.KVStoreKey{keyIBC}))
	return mapp
}

func TestIBCMsgs(t *testing.T) {
	mapp := getMockApp(t)

	sourceChain := "source-chain"
	destChain := "dest-chain"

	priv1 := ed25519.GenPrivKey()
	addr1 := sdk.AccAddress(priv1.PubKey().Address())
	coins := sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}
	var emptyCoins sdk.Coins

	acc := &auth.BaseAccount{
		Address: addr1,
		Coins:   coins,
	}
	accs := []auth.Account{acc}

	mock.SetGenesis(mapp, accs)

	// A checkTx context (true)
	ctxCheck := mapp.BaseApp.NewContext(true, abci.Header{})
	res1 := mapp.AccountMapper.GetAccount(ctxCheck, addr1)
	require.Equal(t, acc, res1)

	packet := IBCPacket{
		SrcAddr:   addr1,
		DestAddr:  addr1,
		Coins:     coins,
		SrcChain:  sourceChain,
		DestChain: destChain,
	}

	transferMsg := IBCTransferMsg{
		IBCPacket: packet,
	}

	receiveMsg := IBCReceiveMsg{
		IBCPacket: packet,
		Relayer:   addr1,
		Sequence:  0,
	}

	mock.SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{transferMsg}, []int64{0}, []int64{0}, true, true, priv1)
	mock.CheckBalance(t, mapp, addr1, emptyCoins)
	mock.SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{transferMsg}, []int64{0}, []int64{1}, false, false, priv1)
	mock.SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{receiveMsg}, []int64{0}, []int64{2}, true, true, priv1)
	mock.CheckBalance(t, mapp, addr1, coins)
	mock.SignCheckDeliver(t, mapp.BaseApp, []sdk.Msg{receiveMsg}, []int64{0}, []int64{2}, false, false, priv1)
}
