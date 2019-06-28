package app

import (
	"encoding/json"
	
	"github.com/comdex-blockchain/rest"
	
	"io"
	"os"
	
	bam "github.com/comdex-blockchain/baseapp"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/assetFactory"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/bank"
	"github.com/comdex-blockchain/x/ibc"
	"github.com/comdex-blockchain/x/order"
	"github.com/comdex-blockchain/x/stake"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName = "AssetApp"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.assetcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.assetd")
)

// AssetApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type AssetApp struct {
	*bam.BaseApp
	cdc *wire.Codec
	
	// keys to access the multistore
	keyMain       *sdk.KVStoreKey
	keyAccount    *sdk.KVStoreKey
	keyIBC        *sdk.KVStoreKey
	keyStake      *sdk.KVStoreKey
	keyAssetStore *sdk.KVStoreKey
	keyOrderStore *sdk.KVStoreKey
	
	// manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	stakeKeeper         stake.Keeper
	assetMapper         assetFactory.AssetPegMapper
	assetKeeper         assetFactory.Keeper
	orderMapper         order.Mapper
	orderKeeper         order.Keeper
}

// NewAssetApp returns a reference to a new AssetApp given a logger and
// database. Internally, a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewAssetApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *AssetApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()
	
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	// create your application type
	var app = &AssetApp{
		cdc:           cdc,
		BaseApp:       bApp,
		keyMain:       sdk.NewKVStoreKey("main"),
		keyAccount:    sdk.NewKVStoreKey("acc"),
		keyIBC:        sdk.NewKVStoreKey("ibc"),
		keyStake:      sdk.NewKVStoreKey("stake"),
		keyAssetStore: sdk.NewKVStoreKey("asset"),
		keyOrderStore: sdk.NewKVStoreKey("order"),
	}
	
	// define and attach the mappers and keepers
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.keyAccount, // target store
		auth.ProtoBaseAccount,
	)
	
	app.assetMapper = assetFactory.NewAssetPegMapper(
		app.cdc,
		app.keyAssetStore,     // target store
		sdk.ProtoBaseAssetPeg, // prototype
	)
	
	app.orderMapper = order.NewMapper(
		app.cdc,
		app.keyOrderStore,
		sdk.ProtoBaseOrder,
	)
	
	// app handlers
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.assetKeeper = assetFactory.NewKeeper(app.assetMapper)
	app.orderKeeper = order.NewKeeper(app.orderMapper)
	
	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewAssetHandler(app.ibcMapper, app.coinKeeper, app.assetKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("assetFactory", assetFactory.NewHandler(app.assetKeeper))
	
	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	
	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStake, app.keyAssetStore)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}
	
	app.Seal()
	
	return app
}

// MakeCodec creates a new wire codec and registers all the necessary types
// with the codec.
func MakeCodec() *wire.Codec {
	
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	bank.RegisterWire(cdc)
	ibc.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	assetFactory.RegisterWire(cdc)
	stake.RegisterWire(cdc)
	assetFactory.RegisterAssetPeg(cdc)
	order.RegisterWire(cdc)
	rest.RegisterWire(cdc)
	
	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *AssetApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *AssetApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)
	
	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
	}
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *AssetApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	
	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		// TODO: https://github.com/comdex-blockchain/issues/468
		panic(err)
	}
	
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}
	
	for _, gasset := range genesisState.Assets {
		asset := gasset.ToAssetPeg()
		app.assetMapper.SetAssetPeg(ctx, asset)
	}
	
	stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	
	return abci.ResponseInitChain{}
}

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *AssetApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	
	accounts := []GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)
	
	assets := []GenesisAssetPeg{}
	appendAssetPeg := func(asset sdk.AssetPeg) (stop bool) {
		assetPeg := NewGenesisAssetPegI(asset)
		assets = append(assets, assetPeg)
		return false
	}
	app.assetMapper.IterateAssets(ctx, appendAssetPeg)
	
	genState := GenesisState{
		Accounts:  accounts,
		Assets:    assets,
		StakeData: stake.WriteGenesis(ctx, app.stakeKeeper),
	}
	appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	validators = stake.WriteValidators(ctx, app.stakeKeeper)
	
	return appState, validators, err
}
