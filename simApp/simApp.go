package simApp

import (
	"encoding/json"
	app2 "github.com/commitHub/commitBlockchain/main/app"
	"io"
	"os"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/bank"
	"github.com/commitHub/commitBlockchain/modules/crisis"
	distr "github.com/commitHub/commitBlockchain/modules/distribution"
	distrClient "github.com/commitHub/commitBlockchain/modules/distribution/client"
	"github.com/commitHub/commitBlockchain/modules/genaccounts"
	"github.com/commitHub/commitBlockchain/modules/genutil"
	"github.com/commitHub/commitBlockchain/modules/gov"
	"github.com/commitHub/commitBlockchain/modules/mint"
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	"github.com/commitHub/commitBlockchain/modules/orders"
	"github.com/commitHub/commitBlockchain/modules/params"
	paramsClient "github.com/commitHub/commitBlockchain/modules/params/client"
	"github.com/commitHub/commitBlockchain/modules/reputation"
	"github.com/commitHub/commitBlockchain/modules/slashing"
	"github.com/commitHub/commitBlockchain/modules/staking"
	"github.com/commitHub/commitBlockchain/modules/supply"
	"github.com/commitHub/commitBlockchain/types/module"
	"github.com/commitHub/commitBlockchain/version"

	"github.com/commitHub/commitBlockchain/modules/assetFactory"
)

const (
	appName        = "Sim App"
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

type SimApp struct {
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

	keyAsset 		*cTypes.KVStoreKey

	tkeyStaking      *cTypes.TransientStoreKey
	tkeyDistribution *cTypes.TransientStoreKey
	tkeyParams       *cTypes.TransientStoreKey

	AccountKeeper      auth.AccountKeeper
	BankKeeper         bank.Keeper
	SupplyKeeper       supply.Keeper
	StakingKeeper      staking.Keeper
	DistributionKeeper distr.Keeper
	SlashingKeeper     slashing.Keeper
	GovKeeper          gov.Keeper
	MintKeeper         mint.Keeper
	CrisisKeeper       crisis.Keeper
	ParamsKeeper       params.Keeper

	ACLKeeper         acl.Keeper
	OrderKeeper       orders.Keeper
	NegotiationKeeper negotiation.Keeper
	ReputationKeeper  reputation.Keeper

	AssetKeeper			assetFactory.Keeper

	mm *module.Manager
}

func NewSimApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *SimApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	app := &SimApp{
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

		keyAsset: cTypes.NewKVStoreKey(assetFactory.ModuleName),
	}

	app.ParamsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams, params.DefaultCodespace)
	authSubspace := app.ParamsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.ParamsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.ParamsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.ParamsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.ParamsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.ParamsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.ParamsKeeper.Subspace(gov.DefaultParamspace)
	crisisSubspace := app.ParamsKeeper.Subspace(crisis.DefaultParamspace)

	app.AccountKeeper = auth.NewAccountKeeper(app.cdc, app.keyAccount, authSubspace, auth.ProtoBaseAccount)

	app.ACLKeeper = acl.NewKeeper(app.keyACL, app.AccountKeeper, app.cdc)
	app.OrderKeeper = orders.NewKeeper(app.keyOrder, app.cdc, app.AccountKeeper)
	app.ReputationKeeper = reputation.NewKeeper(cdc, app.keyReputation, app.OrderKeeper)
	app.NegotiationKeeper = negotiation.NewKeeper(app.keyNegotiation, app.AccountKeeper, app.ACLKeeper, app.ReputationKeeper, app.cdc)

	app.BankKeeper = bank.NewBaseKeeper(app.AccountKeeper, app.NegotiationKeeper, app.ACLKeeper, app.OrderKeeper, app.ReputationKeeper, bankSubspace, bank.DefaultCodespace)

	app.SupplyKeeper = supply.NewKeeper(app.cdc, app.keySupply, app.AccountKeeper, app.BankKeeper, supply.DefaultCodespace, maccPerms)
	StakingKeeper := staking.NewKeeper(app.cdc, app.keyStaking, app.tkeyStaking,
		app.SupplyKeeper, stakingSubspace, staking.DefaultCodespace)
	app.MintKeeper = mint.NewKeeper(app.cdc, app.keyMint, mintSubspace, &StakingKeeper, app.SupplyKeeper, auth.FeeCollectorName)
	app.DistributionKeeper = distr.NewKeeper(app.cdc, app.keyDistribution, distrSubspace, &StakingKeeper,
		app.SupplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName)
	app.SlashingKeeper = slashing.NewKeeper(app.cdc, app.keySlashing, &StakingKeeper,
		slashingSubspace, slashing.DefaultCodespace)
	app.CrisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.SupplyKeeper, auth.FeeCollectorName)

	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistributionKeeper))
	app.GovKeeper = gov.NewKeeper(app.cdc, app.keyGov, app.ParamsKeeper, govSubspace,
		app.SupplyKeeper, &StakingKeeper, gov.DefaultCodespace, govRouter)

	app.StakingKeeper = *StakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.DistributionKeeper.Hooks(), app.SlashingKeeper.Hooks()))

	app.AssetKeeper = assetFactory.NewKeeper(app.keyAsset, app.AccountKeeper, app.cdc)

	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.AccountKeeper),
		genutil.NewAppModule(app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.AccountKeeper),
		bank.NewAppModule(app.BankKeeper, app.AccountKeeper),
		crisis.NewAppModule(app.CrisisKeeper),
		supply.NewAppModule(app.SupplyKeeper, app.AccountKeeper),
		distr.NewAppModule(app.DistributionKeeper, app.SupplyKeeper),
		gov.NewAppModule(app.GovKeeper, app.SupplyKeeper),
		mint.NewAppModule(app.MintKeeper),
		slashing.NewAppModule(app.SlashingKeeper, app.StakingKeeper),
		staking.NewAppModule(app.StakingKeeper, app.DistributionKeeper, app.AccountKeeper, app.SupplyKeeper),

		acl.NewAppModule(app.ACLKeeper, app.AccountKeeper),
		orders.NewAppModule(app.OrderKeeper),
		negotiation.NewAppModule(app.NegotiationKeeper),
		reputation.NewAppModule(app.ReputationKeeper),
	)

	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName)

	app.mm.SetOrderInitGenesis(genaccounts.ModuleName, distr.ModuleName,
		staking.ModuleName, auth.ModuleName, bank.ModuleName, slashing.ModuleName,
		gov.ModuleName, mint.ModuleName, supply.ModuleName, crisis.ModuleName, genutil.ModuleName,
		acl.ModuleName, orders.ModuleName, negotiation.ModuleName, reputation.ModuleName)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.MountStores(app.keyMain, app.keyAccount, app.keySupply, app.keyStaking,
		app.keyMint, app.keyDistribution, app.keySlashing, app.keyGov, app.keyParams,
		app.tkeyParams, app.tkeyStaking, app.tkeyDistribution, app.keyACL, app.keyOrder, app.keyNegotiation, app.keyReputation, app.keyAsset)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.AccountKeeper, app.SupplyKeeper, auth.DefaultSigVerificationGasConsumer))
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
func (app *SimApp) BeginBlocker(ctx cTypes.Context, req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *SimApp) EndBlocker(ctx cTypes.Context, req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *SimApp) InitChainer(ctx cTypes.Context, req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
	var genesisState app2.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *SimApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *SimApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[app.SupplyKeeper.GetModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// GetCDC gives CDC
func (app *SimApp) GetCDC() *codec.Codec {
	return app.cdc
}

func (app *SimApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmTypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abciTypes.Header{Height: app.LastBlockHeight()})

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.StakingKeeper)

	return appState, validators, nil
}

func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

func setup(isCheckTx bool) *SimApp {
	db := dbm.NewMemDB()
	app := NewSimApp(log.NewNopLogger(), db, nil, true, 0)
	if !isCheckTx {
		genesisState := app2.NewDefaultGenesisState()
		stateBytes, err := codec.MarshalJSONIndent(app.cdc, genesisState)
		if err != nil {
			panic(err)
		}

		app.InitChain(
			abciTypes.RequestInitChain{
				Validators:    []abciTypes.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}

	return app
}

func CreateTestApp(isCheckTx bool) (*SimApp, cTypes.Context) {
	app := setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abciTypes.Header{})

	app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	app.BankKeeper.SetSendEnabled(ctx, true)

	return app, ctx
}
