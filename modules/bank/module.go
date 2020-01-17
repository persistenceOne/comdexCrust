package bank

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
	"github.com/persistenceOne/persistenceSDK/modules/bank/client/cli"
	"github.com/persistenceOne/persistenceSDK/modules/bank/client/rest"
	"github.com/persistenceOne/persistenceSDK/modules/bank/internal/keeper"
	"github.com/persistenceOne/persistenceSDK/modules/bank/internal/types"
	"github.com/persistenceOne/persistenceSDK/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string { return ModuleName }

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// module validate genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router, kafkaBool bool, kafkaState kafka.KafkaState) {
	rest.RegisterRoutes(ctx, rtr, kafkaBool, kafkaState)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	bankTxCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "bank transaction sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bankTxCmd.AddCommand(client.PostCommands(
		cli.BuyerExecuteOrderCmd(cdc),
		cli.DefineACLCmd(cdc),
		cli.DefineOrganizationCmd(cdc),
		cli.DefineZoneCmd(cdc),
		cli.IssueAssetCmd(cdc),
		cli.IssueFiatCmd(cdc),
		cli.RedeemAssetCmd(cdc),
		cli.RedeemFiatCmd(cdc),
		cli.ReleaseAssetCmd(cdc),
		cli.SellerExecuteOrderCmd(cdc),
		cli.SendAssetCmd(cdc),
		cli.SendFiatCmd(cdc),
	)...)

	return bankTxCmd
}

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	bankQueryCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "bank query sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bankQueryCmd.AddCommand(client.GetCommands(
		cli.GetAssetCmd(cdc),
		cli.GetFiatCmd(cdc),
	)...)

	return bankQueryCmd
}

// ___________________________
// app module
type AppModule struct {
	AppModuleBasic
	keeper        Keeper
	accountKeeper types.AccountKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
	}
}

// module name
func (AppModule) Name() string { return ModuleName }

// register invariants
func (am AppModule) RegisterInvariants(ir cTypes.InvariantRegistry) {
	keeper.RegisterInvariants(ir, am.accountKeeper)
}

// module message route name
func (AppModule) Route() string { return RouterKey }

// module handler
func (am AppModule) NewHandler() cTypes.Handler { return NewHandler(am.keeper) }

// module querier route name
func (AppModule) QuerierRoute() string { return RouterKey }

// module querier
func (am AppModule) NewQuerierHandler() cTypes.Querier {
	return keeper.NewQuerier(am.keeper)
}

// module init-genesis
func (am AppModule) InitGenesis(ctx cTypes.Context, data json.RawMessage) []abciTypes.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abciTypes.ValidatorUpdate{}
}

// module export genesis
func (am AppModule) ExportGenesis(ctx cTypes.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// module begin-block
func (AppModule) BeginBlock(_ cTypes.Context, _ abciTypes.RequestBeginBlock) {}

// module end-block
func (AppModule) EndBlock(_ cTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}
