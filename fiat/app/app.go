package app

import (
	"encoding/json"
	"io"
	"os"
	
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	
	bam "github.com/comdex-blockchain/baseapp"
	"github.com/comdex-blockchain/rest"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/bank"
	"github.com/comdex-blockchain/x/fiatFactory"
	"github.com/comdex-blockchain/x/ibc"
	"github.com/comdex-blockchain/x/order"
	"github.com/comdex-blockchain/x/slashing"
	"github.com/comdex-blockchain/x/stake"
)

const (
	appName = "FiatApp"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.fiatcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.fiatd")
)

// FiatApp : Extended ABCI application
type FiatApp struct {
	*bam.BaseApp
	cdc *wire.Codec
	
	keyMain       *sdk.KVStoreKey
	keyAccount    *sdk.KVStoreKey
	keyIBC        *sdk.KVStoreKey
	keyStake      *sdk.KVStoreKey
	keySlashing   *sdk.KVStoreKey
	keyFiatStore  *sdk.KVStoreKey
	keyOrderStore *sdk.KVStoreKey
	
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	stakeKeeper         stake.Keeper
	slashingKeeper      slashing.Keeper
	fiatMapper          fiatFactory.FiatPegMapper
	fiatKeeper          fiatFactory.Keeper
	orderMapper         order.Mapper
	orderKeeer          order.Keeper
}

// NewFiatApp : returns new  fiat app
func NewFiatApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *FiatApp {
	cdc := MakeCodec()
	
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	
	var app = &FiatApp{
		BaseApp:       bApp,
		cdc:           cdc,
		keyMain:       sdk.NewKVStoreKey("main"),
		keyAccount:    sdk.NewKVStoreKey("acc"),
		keyIBC:        sdk.NewKVStoreKey("ibc"),
		keyStake:      sdk.NewKVStoreKey("stake"),
		keyFiatStore:  sdk.NewKVStoreKey("fiat"),
		keyOrderStore: sdk.NewKVStoreKey("order"),
	}
	
	// define the accountMapper
	app.accountMapper = auth.NewAccountMapper(
		app.cdc,
		app.keyAccount,        // target store
		auth.ProtoBaseAccount, // prototype
	)
	app.fiatMapper = fiatFactory.NewFiatPegMapper(
		app.cdc,
		app.keyFiatStore,     // target store
		sdk.ProtoBaseFiatPeg, // prototype
	)
	app.orderMapper = order.NewMapper(
		app.cdc,
		app.keyOrderStore,
		sdk.ProtoBaseOrder,
	)
	
	// add handlers
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.fiatKeeper = fiatFactory.NewKeeper(app.fiatMapper)
	app.orderKeeer = order.NewKeeper(app.orderMapper)
	
	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewFiatHandler(app.ibcMapper, app.coinKeeper, app.fiatKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("fiatFactory", fiatFactory.NewHandler(app.fiatKeeper))
	
	// initialize BaseApp
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStake, app.keyFiatStore)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}
	
	return app
}

// MakeCodec : custom tx codec
func MakeCodec() *wire.Codec {
	
	var cdc = wire.NewCodec()
	
	ibc.RegisterWire(cdc)
	bank.RegisterWire(cdc)
	stake.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	sdk.RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	fiatFactory.RegisterWire(cdc)
	fiatFactory.RegisterFiatPeg(cdc)
	order.RegisterWire(cdc)
	rest.RegisterWire(cdc)
	
	return cdc
}

// BeginBlocker : application updates every end block
func (app *FiatApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker : application updates every end block
func (app *FiatApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)
	
	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
	}
}

// custom logic for fiat initialization
func (app *FiatApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	
	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}
	
	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}
	
	for _, gfiat := range genesisState.Fiats {
		fiat := gfiat.ToFiatPeg()
		app.fiatMapper.SetFiatPeg(ctx, fiat)
	}
	
	// load the initial stake information
	stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	
	return abci.ResponseInitChain{}
}

// ExportAppStateAndValidators : export the state of fiats for a genesis file
func (app *FiatApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	
	// iterate to get the accounts
	accounts := []GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)
	
	// iterate to get fiats
	fiats := []GenesisFiatPeg{}
	appendFiatPeg := func(fiat sdk.FiatPeg) (stop bool) {
		fiatPeg := NewGenesisFiatPegI(fiat)
		fiats = append(fiats, fiatPeg)
		return false
	}
	app.fiatMapper.IterateFiats(ctx, appendFiatPeg)
	
	genState := GenesisState{
		Accounts:  accounts,
		Fiats:     fiats,
		StakeData: stake.WriteGenesis(ctx, app.stakeKeeper),
	}
	appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	validators = stake.WriteValidators(ctx, app.stakeKeeper)
	return appState, validators, nil
}
