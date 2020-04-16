package acl

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/kafka"
	"github.com/persistenceOne/comdexCrust/types/module"

	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/persistenceOne/comdexCrust/modules/acl/client/cli"
	"github.com/persistenceOne/comdexCrust/modules/acl/client/rest"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return ModuleName }

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router, kafkaBool bool, kafkaState kafka.KafkaState) {
	rest.RegisterRoutes(ctx, rtr)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil // TODO
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	aclQueryCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "acl query sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	aclQueryCmd.AddCommand(client.GetCommands(
		cli.GetACLAccountCmd(cdc),
		cli.GetOrganizationCmd(cdc),
		cli.GetZoneCmd(cdc),
	)...)

	return aclQueryCmd
}

type AppModule struct {
	AppModuleBasic
	keeper        Keeper
	accountKeeper AccountKeeper
}

func NewAppModule(keeper Keeper, accountKeeper AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir cTypes.InvariantRegistry) {

}

func (AppModule) Route() string { return RouterKey }

func (am AppModule) NewHandler() cTypes.Handler { return NewHandler(am.keeper) }

func (AppModule) QuerierRoute() string { return QuerierRoute }

func (am AppModule) NewQuerierHandler() cTypes.Querier { return NewQuerier(am.keeper) }

func (am AppModule) InitGenesis(ctx cTypes.Context, data json.RawMessage) []abciTypes.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	_ = InitGenesis(ctx, am.keeper, genesisState)

	return []abciTypes.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx cTypes.Context) json.RawMessage {
	gs := ExportGenesisState(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

func (AppModule) BeginBlock(cTypes.Context, abciTypes.RequestBeginBlock) {}

func (AppModule) EndBlock(_ cTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}
