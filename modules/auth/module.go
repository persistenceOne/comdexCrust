package auth

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/kafka"
	"github.com/persistenceOne/comdexCrust/types/module"

	"github.com/persistenceOne/comdexCrust/modules/auth/client/cli"
	"github.com/persistenceOne/comdexCrust/modules/auth/client/rest"
	"github.com/persistenceOne/comdexCrust/modules/auth/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := types.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router, kafkaBool bool, kafkaState kafka.KafkaState) {
	rest.RegisterRoutes(ctx, rtr, types.StoreKey)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

// ___________________________
// app module object
type AppModule struct {
	AppModuleBasic
	accountKeeper AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(accountKeeper AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		accountKeeper:  accountKeeper,
	}
}

// module name
func (AppModule) Name() string {
	return types.ModuleName
}

// register invariants
func (AppModule) RegisterInvariants(_ cTypes.InvariantRegistry) {}

// module message route name
func (AppModule) Route() string { return "" }

// module handler
func (AppModule) NewHandler() cTypes.Handler { return nil }

// module querier route name
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// module querier
func (am AppModule) NewQuerierHandler() cTypes.Querier {
	return NewQuerier(am.accountKeeper)
}

// module init-genesis
func (am AppModule) InitGenesis(ctx cTypes.Context, data json.RawMessage) []abciTypes.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.accountKeeper, genesisState)
	return []abciTypes.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx cTypes.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.accountKeeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// module begin-block
func (AppModule) BeginBlock(_ cTypes.Context, _ abciTypes.RequestBeginBlock) {}

// module end-block
func (AppModule) EndBlock(_ cTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}
