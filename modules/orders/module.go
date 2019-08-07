package orders

import (
	"encoding/json"
	
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/types/module"
	
	abciTypes "github.com/tendermint/tendermint/abci/types"
	
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	"github.com/commitHub/commitBlockchain/modules/orders/client/cli"
	"github.com/commitHub/commitBlockchain/modules/orders/client/rest"
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

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, r *mux.Router) {
	rest.RegisterRoutes(ctx, r)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil // TODO
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	orderQueryCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "order query sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	
	orderQueryCmd.AddCommand(client.GetCommands(
		cli.GetOrderCmd(cdc),
	)...)
	
	return orderQueryCmd
}

type AppModule struct {
	AppModuleBasic
	keeper            Keeper
	negotiationKeeper negotiation.Keeper
}

func NewAppModule(keeper Keeper, negotiationKeeper negotiation.Keeper) AppModule {
	return AppModule{
		AppModuleBasic:    AppModuleBasic{},
		keeper:            keeper,
		negotiationKeeper: negotiationKeeper,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir cTypes.InvariantRegistry) {

}

func (AppModule) Route() string { return RouterKey }

func (am AppModule) NewHandler() cTypes.Handler { return nil }

func (am AppModule) QuerierRoute() string { return QuerierRoute }

func (am AppModule) NewQuerierHandler() cTypes.Querier { return NewQuerier(am.keeper) }

func (am AppModule) InitGenesis(ctx cTypes.Context, data json.RawMessage) []abciTypes.ValidatorUpdate {
	var gs GenesisState
	
	_ = ModuleCdc.UnmarshalJSON(data, gs)
	InitGenesis(ctx, am.keeper, gs)
	
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
