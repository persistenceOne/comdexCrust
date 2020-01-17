package negotiation

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/kafka"
	"github.com/persistenceOne/persistenceSDK/modules/negotiation/client/cli"
	"github.com/persistenceOne/persistenceSDK/modules/negotiation/client/rest"
	"github.com/persistenceOne/persistenceSDK/types/module"
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
	rest.RegisterRoutes(ctx, rtr, kafkaBool, kafkaState)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	negotiationTxCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "negotiation transaction sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	negotiationTxCmd.AddCommand(client.PostCommands(
		cli.ChangeBuyerBidCmd(cdc),
		cli.ChangeSellerBidCmd(cdc),
		cli.ConfirmBuyerBidCmd(cdc),
		cli.ConfirmSellerBidCmd(cdc),
	)...)

	return negotiationTxCmd
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	negotiationQueryCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "negotiation query sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	negotiationQueryCmd.AddCommand(client.GetCommands(
		cli.GetNegotiationCmd(cdc),
	)...)

	return negotiationQueryCmd
}

type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir cTypes.InvariantRegistry) {

}

func (AppModule) Route() string { return RouterKey }

func (am AppModule) NewHandler() cTypes.Handler { return NewHandler(am.keeper) }

func (am AppModule) QuerierRoute() string { return QuerierRoute }

func (am AppModule) NewQuerierHandler() cTypes.Querier { return NewQuerier(am.keeper) }

func (am AppModule) InitGenesis(ctx cTypes.Context, data json.RawMessage) []abciTypes.ValidatorUpdate {
	var genesisState GenesisState

	_ = ModuleCdc.UnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)

	return []abciTypes.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx cTypes.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)

	return ModuleCdc.MustMarshalJSON(gs)
}

func (AppModule) BeginBlock(cTypes.Context, abciTypes.RequestBeginBlock) {}

func (AppModule) EndBlock(_ cTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}
