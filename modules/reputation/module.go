package reputation

import (
	"encoding/json"
	"github.com/commitHub/commitBlockchain/kafka"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/types/module"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commitHub/commitBlockchain/modules/reputation/client/cli"
	"github.com/commitHub/commitBlockchain/modules/reputation/client/rest"
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

	reputationTxCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "reputation transaction sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	reputationTxCmd.AddCommand(client.PostCommands(
		cli.SubmitBuyerFeedbackCmd(cdc),
		cli.SubmitSellerFeedbackCmd(cdc),
	)...)

	return reputationTxCmd
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	reputationQueryCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "reputation query sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	reputationQueryCmd.AddCommand(client.GetCommands(
		cli.GetReputationCmd(cdc),
	)...)

	return reputationQueryCmd
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

func (am AppModule) InitGenesis(ctx cTypes.Context, data json.RawMessage) []abci.ValidatorUpdate {
	// TODO
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx cTypes.Context) json.RawMessage {
	// TODO
	return nil
}

func (AppModule) BeginBlock(cTypes.Context, abci.RequestBeginBlock) {}

func (AppModule) EndBlock(_ cTypes.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
