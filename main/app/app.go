package app

import (
	"encoding/json"
	"io"
	"os"
	
	bam "github.com/comdex-blockchain/baseapp"
	"github.com/comdex-blockchain/rest"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/acl"
	"github.com/comdex-blockchain/x/assetFactory"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/bank"
	"github.com/comdex-blockchain/x/fiatFactory"
	"github.com/comdex-blockchain/x/ibc"
	"github.com/comdex-blockchain/x/negotiation"
	"github.com/comdex-blockchain/x/order"
	"github.com/comdex-blockchain/x/reputation"
	"github.com/comdex-blockchain/x/stake"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName = "MainApp"
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.maincli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.maind")
)

// MainApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type MainApp struct {
	*bam.BaseApp
	cdc *wire.Codec
	
	keyMain        *sdk.KVStoreKey
	keyAccount     *sdk.KVStoreKey
	keyIBC         *sdk.KVStoreKey
	keyStake       *sdk.KVStoreKey
	keyNegotiation *sdk.KVStoreKey
	keyOrder       *sdk.KVStoreKey
	keyACL         *sdk.KVStoreKey
	keyReputation  *sdk.KVStoreKey
	
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	stakeKeeper         stake.Keeper
	negotiationMapper   negotiation.Mapper
	negotiationKeeper   negotiation.Keeper
	orderMapper         order.Mapper
	orderKeeper         order.Keeper
	aclMapper           acl.Mapper
	aclKeeper           acl.Keeper
	reputationMapper    reputation.Mapper
	reputationKeeper    reputation.Keeper
}

// NewMainApp returns a reference to a new MainApp given a logger and
// database. Internally, a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewMainApp(logger log.Logger, db dbm.DB, traceStore io.Writer, baseAppOptions ...func(*bam.BaseApp)) *MainApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()
	
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	// create your application type
	var app = &MainApp{
		cdc:            cdc,
		BaseApp:        bApp,
		keyMain:        sdk.NewKVStoreKey("main"),
		keyAccount:     sdk.NewKVStoreKey("acc"),
		keyIBC:         sdk.NewKVStoreKey("ibc"),
		keyStake:       sdk.NewKVStoreKey("stake"),
		keyNegotiation: sdk.NewKVStoreKey("negotiation"),
		keyOrder:       sdk.NewKVStoreKey("order"),
		keyACL:         sdk.NewKVStoreKey("acl"),
		keyReputation:  sdk.NewKVStoreKey("reputation"),
	}
	
	// define and attach the mappers and keepers
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.keyAccount, // target store
		auth.ProtoBaseAccount,
	)
	app.aclMapper = acl.NewACLMapper(
		app.cdc,
		app.keyACL,              // target store
		sdk.ProtoBaseACLAccount, // prototype
	)
	app.negotiationMapper = negotiation.NewMapper(
		app.cdc,
		app.keyNegotiation,       // target store
		sdk.ProtoBaseNegotiation, // prototype
	)
	
	app.orderMapper = order.NewMapper(
		app.cdc,
		app.keyOrder,       // target store
		sdk.ProtoBaseOrder, // prototype
	)
	app.reputationMapper = reputation.NewMapper(
		app.cdc,
		app.keyReputation,
		sdk.ProtoBaseAccountReputation,
	)
	app.aclKeeper = acl.NewKeeper(app.aclMapper)
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	app.negotiationKeeper = negotiation.NewKeeper(app.negotiationMapper, app.accountMapper)
	app.orderKeeper = order.NewKeeper(app.orderMapper)
	app.reputationKeeper = reputation.NewKeeper(app.reputationMapper)
	// register message routes
	app.Router().
		AddRoute("bank", bank.NewAssetFiatHandler(app.coinKeeper, app.negotiationKeeper, app.orderKeeper, app.aclKeeper, app.reputationKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper, app.aclKeeper, app.negotiationKeeper, app.orderKeeper, app.reputationKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper)).
		AddRoute("negotiation", negotiation.NewHandler(app.negotiationKeeper, app.aclKeeper, app.reputationKeeper)).
		AddRoute("reputation", reputation.NewFeedbackHandler(app.reputationKeeper, app.orderKeeper))
	
	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	
	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStake, app.keyNegotiation, app.keyOrder, app.keyACL, app.keyReputation)
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
	stake.RegisterWire(cdc)
	negotiation.RegisterWire(cdc)
	order.RegisterWire(cdc)
	assetFactory.RegisterAssetPeg(cdc)
	fiatFactory.RegisterFiatPeg(cdc)
	fiatFactory.RegisterWire(cdc)
	negotiation.RegisterNegotiation(cdc)
	order.RegisterOrder(cdc)
	acl.RegisterWire(cdc)
	rest.RegisterWire(cdc)
	acl.RegisterACLAccount(cdc)
	reputation.RegisterWire(cdc)
	reputation.RegisterReputation(cdc)
	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *MainApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *MainApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
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
func (app *MainApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	
	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}
	
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}
	
	// load the initial stake information
	stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)
	order.InitOrder(ctx, app.orderKeeper)
	negotiation.InitNegotiation(ctx, app.negotiationKeeper)
	acl.InitACL(ctx, app.aclKeeper)
	reputation.InitReputation(ctx, app.reputationKeeper)
	
	return abci.ResponseInitChain{}
}

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *MainApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	
	// iterate to get the accounts
	accounts := []GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)
	
	genState := GenesisState{
		Accounts:  accounts,
		StakeData: stake.WriteGenesis(ctx, app.stakeKeeper),
	}
	appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	validators = stake.WriteValidators(ctx, app.stakeKeeper)
	return appState, validators, nil
}
