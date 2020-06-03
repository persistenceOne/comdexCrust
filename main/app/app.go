package app

import (
	"encoding/json"
	"io"
	"os"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/bank"
	"github.com/commitHub/commitBlockchain/modules/crisis"
	distr "github.com/commitHub/commitBlockchain/modules/distribution"
	distrclient "github.com/commitHub/commitBlockchain/modules/distribution/client"
	"github.com/commitHub/commitBlockchain/modules/genaccounts"
	"github.com/commitHub/commitBlockchain/modules/genutil"
	"github.com/commitHub/commitBlockchain/modules/gov"
	"github.com/commitHub/commitBlockchain/modules/mint"
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	"github.com/commitHub/commitBlockchain/modules/orders"
	"github.com/commitHub/commitBlockchain/modules/params"
	paramsclient "github.com/commitHub/commitBlockchain/modules/params/client"
	"github.com/commitHub/commitBlockchain/modules/reputation"
	"github.com/commitHub/commitBlockchain/modules/slashing"
	"github.com/commitHub/commitBlockchain/modules/staking"
	"github.com/commitHub/commitBlockchain/modules/supply"
	"github.com/commitHub/commitBlockchain/types/module"
	"github.com/commitHub/commitBlockchain/version"
)

const (
	appName        = "Main App"
	DefaultKeyPass = "1234567890"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.maincli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.maind")

	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distrclient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		reputation.AppModuleBasic{},
		acl.AppModuleBasic{},
		negotiation.AppModuleBasic{},
		orders.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	cTypes.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type MainApp struct {
	*bam.BaseApp
	cdc            *codec.Codec
	invCheckPeriod uint

	keyMain *cTypes.KVStoreKey

	keyAccount      *cTypes.KVStoreKey
	keySupply       *cTypes.KVStoreKey
	keyStaking      *cTypes.KVStoreKey
	keyDistribution *cTypes.KVStoreKey
	keySlashing     *cTypes.KVStoreKey
	keyGov          *cTypes.KVStoreKey
	keyMint         *cTypes.KVStoreKey
	keyParams       *cTypes.KVStoreKey

	keyACL         *cTypes.KVStoreKey
	keyOrder       *cTypes.KVStoreKey
	keyNegotiation *cTypes.KVStoreKey
	keyReputation  *cTypes.KVStoreKey

	tkeyStaking      *cTypes.TransientStoreKey
	tkeyDistribution *cTypes.TransientStoreKey
	tkeyParams       *cTypes.TransientStoreKey

	accountKeeper      auth.AccountKeeper
	bankKeeper         bank.Keeper
	supplyKeeper       supply.Keeper
	stakingKeeper      staking.Keeper
	distributionKeeper distr.Keeper
	slashingKeeper     slashing.Keeper
	govKeeper          gov.Keeper
	mintKeeper         mint.Keeper
	crisisKeeper       crisis.Keeper
	paramsKeeper       params.Keeper

	aclKeeper         acl.Keeper
	orderKeeper       orders.Keeper
	negotiationKeeper negotiation.Keeper
	reputationKeeper  reputation.Keeper

	mm *module.Manager
}

func NewMainApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *MainApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	app := &MainApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,

		keyMain:          cTypes.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:       cTypes.NewKVStoreKey(auth.ModuleName),
		keySupply:        cTypes.NewKVStoreKey(supply.ModuleName),
		keyStaking:       cTypes.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      cTypes.NewTransientStoreKey(staking.TStoreKey),
		keyDistribution:  cTypes.NewKVStoreKey(distr.StoreKey),
		tkeyDistribution: cTypes.NewTransientStoreKey(distr.TStoreKey),
		keySlashing:      cTypes.NewKVStoreKey(slashing.ModuleName),
		keyGov:           cTypes.NewKVStoreKey(gov.ModuleName),
		keyMint:          cTypes.NewKVStoreKey(mint.ModuleName),
		keyParams:        cTypes.NewKVStoreKey(params.ModuleName),
		tkeyParams:       cTypes.NewTransientStoreKey(params.TStoreKey),

		keyACL:         cTypes.NewKVStoreKey(acl.ModuleName),
		keyNegotiation: cTypes.NewKVStoreKey(negotiation.ModuleName),
		keyOrder:       cTypes.NewKVStoreKey(orders.ModuleName),
		keyReputation:  cTypes.NewKVStoreKey(reputation.ModuleName),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams, params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace)
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	app.accountKeeper = auth.NewAccountKeeper(app.cdc, app.keyAccount, authSubspace, auth.ProtoBaseAccount)
	app.negotiationKeeper = negotiation.NewKeeper(app.keyNegotiation, app.accountKeeper, app.cdc)
	app.aclKeeper = acl.NewKeeper(app.keyACL, app.accountKeeper, app.cdc)
	app.orderKeeper = orders.NewKeeper(app.keyOrder, app.cdc, app.negotiationKeeper, app.aclKeeper, app.accountKeeper)
	app.reputationKeeper = reputation.NewKeeper(cdc, app.keyReputation, app.orderKeeper)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, app.negotiationKeeper, app.aclKeeper, app.orderKeeper, app.reputationKeeper, bankSubspace, bank.DefaultCodespace)
	app.supplyKeeper = supply.NewKeeper(app.cdc, app.keySupply, app.accountKeeper, app.bankKeeper, supply.DefaultCodespace, maccPerms)
	stakingKeeper := staking.NewKeeper(app.cdc, app.keyStaking, app.tkeyStaking,
		app.supplyKeeper, stakingSubspace, staking.DefaultCodespace)
	app.mintKeeper = mint.NewKeeper(app.cdc, app.keyMint, mintSubspace, &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName)
	app.distributionKeeper = distr.NewKeeper(app.cdc, app.keyDistribution, distrSubspace, &stakingKeeper,
		app.supplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName)
	app.slashingKeeper = slashing.NewKeeper(app.cdc, app.keySlashing, &stakingKeeper,
		slashingSubspace, slashing.DefaultCodespace)
	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distributionKeeper))
	app.govKeeper = gov.NewKeeper(app.cdc, app.keyGov, app.paramsKeeper, govSubspace,
		app.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter)

	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distributionKeeper.Hooks(), app.slashingKeeper.Hooks()))

	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distributionKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.distributionKeeper, app.accountKeeper, app.supplyKeeper),

		acl.NewAppModule(app.aclKeeper, app.accountKeeper),
		orders.NewAppModule(app.orderKeeper, app.negotiationKeeper),
		negotiation.NewAppModule(app.negotiationKeeper),
		reputation.NewAppModule(app.reputationKeeper),
	)

	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName)

	app.mm.SetOrderInitGenesis(genaccounts.ModuleName, distr.ModuleName,
		staking.ModuleName, auth.ModuleName, bank.ModuleName, slashing.ModuleName,
		gov.ModuleName, mint.ModuleName, supply.ModuleName, crisis.ModuleName, genutil.ModuleName,
		acl.ModuleName, orders.ModuleName, negotiation.ModuleName, reputation.ModuleName)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.MountStores(app.keyMain, app.keyAccount, app.keySupply, app.keyStaking,
		app.keyMint, app.keyDistribution, app.keySlashing, app.keyGov, app.keyParams,
		app.tkeyParams, app.tkeyStaking, app.tkeyDistribution, app.keyACL, app.keyOrder, app.keyNegotiation, app.keyReputation)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keyMain)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}
	return app
}

// application updates every begin block
func (app *MainApp) BeginBlocker(ctx cTypes.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *MainApp) EndBlocker(ctx cTypes.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *MainApp) InitChainer(ctx cTypes.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *MainApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *MainApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[app.supplyKeeper.GetModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *MainApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.stakingKeeper)

	return appState, validators, nil
}
