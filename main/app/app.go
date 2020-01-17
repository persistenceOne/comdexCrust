package app

import (
	"encoding/json"
	"io"
	"os"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/acl"
	"github.com/persistenceOne/persistenceSDK/modules/assetFactory"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/bank"
	"github.com/persistenceOne/persistenceSDK/modules/crisis"
	distr "github.com/persistenceOne/persistenceSDK/modules/distribution"
	distrClient "github.com/persistenceOne/persistenceSDK/modules/distribution/client"
	"github.com/persistenceOne/persistenceSDK/modules/fiatFactory"
	"github.com/persistenceOne/persistenceSDK/modules/genaccounts"
	"github.com/persistenceOne/persistenceSDK/modules/genutil"
	"github.com/persistenceOne/persistenceSDK/modules/gov"
	"github.com/persistenceOne/persistenceSDK/modules/ibc"
	ibctransfer "github.com/persistenceOne/persistenceSDK/modules/ibc/20-transfer"
	"github.com/persistenceOne/persistenceSDK/modules/mint"
	"github.com/persistenceOne/persistenceSDK/modules/negotiation"
	"github.com/persistenceOne/persistenceSDK/modules/orders"
	"github.com/persistenceOne/persistenceSDK/modules/params"
	paramsClient "github.com/persistenceOne/persistenceSDK/modules/params/client"
	"github.com/persistenceOne/persistenceSDK/modules/reputation"
	"github.com/persistenceOne/persistenceSDK/modules/slashing"
	"github.com/persistenceOne/persistenceSDK/modules/staking"
	"github.com/persistenceOne/persistenceSDK/modules/supply"
	"github.com/persistenceOne/persistenceSDK/types/module"
	"github.com/persistenceOne/persistenceSDK/version"
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
		gov.NewAppModuleBasic(paramsClient.ProposalHandler, distrClient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		reputation.AppModuleBasic{},
		acl.AppModuleBasic{},
		negotiation.AppModuleBasic{},
		orders.AppModuleBasic{},
		ibc.AppModuleBasic{},
		assetFactory.AppModuleBasic{},
		fiatFactory.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:              nil,
		distr.ModuleName:                   nil,
		mint.ModuleName:                    {supply.Minter},
		staking.BondedPoolName:             {supply.Burner, supply.Staking},
		staking.NotBondedPoolName:          {supply.Burner, supply.Staking},
		gov.ModuleName:                     {supply.Burner},
		ibctransfer.GetModuleAccountName(): {supply.Minter, supply.Burner},
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

	keyIBC   *cTypes.KVStoreKey
	keyAsset *cTypes.KVStoreKey
	keyFiat  *cTypes.KVStoreKey

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

	ibcKeeper   ibc.Keeper
	assetKeeper assetFactory.Keeper
	fiatKeeper  fiatFactory.Keeper

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

		keyIBC:   cTypes.NewKVStoreKey(ibc.StoreKey),
		keyAsset: cTypes.NewKVStoreKey(assetFactory.StoreKey),
		keyFiat:  cTypes.NewKVStoreKey(fiatFactory.StoreKey),
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

	app.aclKeeper = acl.NewKeeper(app.keyACL, app.accountKeeper, app.cdc)
	app.orderKeeper = orders.NewKeeper(app.keyOrder, app.cdc, app.accountKeeper)
	app.reputationKeeper = reputation.NewKeeper(cdc, app.keyReputation, app.orderKeeper)
	app.negotiationKeeper = negotiation.NewKeeper(app.keyNegotiation, app.accountKeeper, app.aclKeeper, app.reputationKeeper, app.cdc)

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

	app.assetKeeper = assetFactory.NewKeeper(app.cdc, app.keyAsset, app.accountKeeper)
	app.fiatKeeper = fiatFactory.NewKeeper(app.cdc, app.keyFiat, app.accountKeeper)

	app.ibcKeeper = ibc.NewKeeper(app.cdc, app.keyIBC, ibc.DefaultCodespace, app.bankKeeper, app.supplyKeeper, app.assetKeeper, app.aclKeeper)

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
		orders.NewAppModule(app.orderKeeper),
		negotiation.NewAppModule(app.negotiationKeeper),
		reputation.NewAppModule(app.reputationKeeper),

		ibc.NewAppModule(app.ibcKeeper),
		assetFactory.NewAppModule(app.assetKeeper, app.accountKeeper),
		fiatFactory.NewAppModule(app.fiatKeeper, app.accountKeeper),
	)

	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName)

	app.mm.SetOrderInitGenesis(genaccounts.ModuleName, distr.ModuleName,
		staking.ModuleName, auth.ModuleName, bank.ModuleName, slashing.ModuleName,
		gov.ModuleName, mint.ModuleName, supply.ModuleName, crisis.ModuleName, genutil.ModuleName,
		acl.ModuleName, orders.ModuleName, negotiation.ModuleName, reputation.ModuleName, ibc.ModuleName, assetFactory.ModuleName, fiatFactory.ModuleName)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.MountStores(app.keyMain, app.keyAccount, app.keySupply, app.keyStaking,
		app.keyMint, app.keyDistribution, app.keySlashing, app.keyGov, app.keyParams,
		app.tkeyParams, app.tkeyStaking, app.tkeyDistribution, app.keyACL, app.keyOrder, app.keyNegotiation, app.keyReputation, app.keyIBC, app.keyAsset, app.keyFiat)

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
func (app *MainApp) BeginBlocker(ctx cTypes.Context, req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *MainApp) EndBlocker(ctx cTypes.Context, req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *MainApp) InitChainer(ctx cTypes.Context, req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
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
) (appState json.RawMessage, validators []tmTypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abciTypes.Header{Height: app.LastBlockHeight()})

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.stakingKeeper)

	return appState, validators, nil
}
