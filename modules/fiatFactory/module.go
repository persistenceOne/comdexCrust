package fiatFactory

import (
	"encoding/json"
	"github.com/commitHub/commitBlockchain/kafka"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/types/module"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"

	"github.com/commitHub/commitBlockchain/modules/fiatFactory/client/cli"
	"github.com/commitHub/commitBlockchain/modules/fiatFactory/client/rest"
	fiatFactoryTypes "github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
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

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, router *mux.Router, kafkaBool bool, kafkaState kafka.KafkaState) {
	rest.RegisterRoutes(ctx, router)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	fiatTxCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "fiat transaction sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	fiatTxCmd.AddCommand(client.PostCommands(
		cli.IssueFiatCmd(cdc),
		cli.SendFiatCmd(cdc),
		cli.ExecuteFiatCmd(cdc),
		cli.RedeemFiatCmd(cdc),
	)...)

	return fiatTxCmd
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {

	fiatQueryCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "fiat query sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	fiatQueryCmd.AddCommand(client.GetCommands(
		cli.QueryFiatCmd(),
	)...)

	return fiatQueryCmd
}

type AppModule struct {
	AppModuleBasic
	keeper        Keeper
	accountKeeper fiatFactoryTypes.AccountKeeper
}

func NewModule(keeper Keeper, accountKeeper fiatFactoryTypes.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
	}
}

func (AppModule) Name() string { return ModuleName }

func (am AppModule) RegisterInvariants(ir cTypes.InvariantRegistry) {}

func (AppModule) InitGenesis(ctx cTypes.Context, data json.RawMessage) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}

func (AppModule) ExportGenesis(ctx cTypes.Context) json.RawMessage {
	return nil
}

func (AppModule) BeginBlock(cTypes.Context, abciTypes.RequestBeginBlock) {}

func (AppModule) EndBlock(_ cTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}

func (AppModule) Route() string { return ModuleName }

func (AppModule) QuerierRoute() string { return QuerierRoute }

func (AppModule) NewQuerierHandler() cTypes.Querier { return nil }

func (am AppModule) NewHandler() cTypes.Handler { return NewHandler(am.keeper) }
