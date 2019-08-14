package assetFactory

import (
	"encoding/json"
	
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/kafka"
	"github.com/commitHub/commitBlockchain/types/module"
	
	abciTypes "github.com/tendermint/tendermint/abci/types"
	
	"github.com/commitHub/commitBlockchain/modules/assetFactory/client/cli"
	"github.com/commitHub/commitBlockchain/modules/assetFactory/client/rest"
	"github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return ModuleName }

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState)
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
	assetTxCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "asset transaction sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	
	assetTxCmd.AddCommand(client.PostCommands(
		cli.IssueAssetCmd(cdc),
		cli.SendAssetCmd(cdc),
		cli.ExecuteAssetCmd(cdc),
		cli.RedeemAssetCmd(cdc),
	)...)
	
	return assetTxCmd
}
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	assetQueryCmd := &cobra.Command{
		Use:                        ModuleName,
		Short:                      "asset query sub commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	
	assetQueryCmd.AddCommand(client.GetCommands(
		cli.QueryAssetCmd(),
	)...)
	
	return assetQueryCmd
}

type AppModule struct {
	AppModuleBasic
	keeper        Keeper
	accountKeeper types.AccountKeeper
}

func NewModule(keeper Keeper, accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		accountKeeper:  accountKeeper,
	}
}

func (AppModule) Name() string {
	return types.ModuleName
}

func (am AppModule) RegisterInvariants(ir cTypes.InvariantRegistry) {}

func (AppModule) Route() string { return RouterKey }

func (am AppModule) NewHandler() cTypes.Handler { return NewHandler(am.keeper) }

func (AppModule) QuerierRoute() string { return QuerierRoute }

func (AppModule) NewQuerierHandler() cTypes.Querier { return nil }

func (am AppModule) InitGenesis(ctx cTypes.Context, data json.RawMessage) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}

func (AppModule) ExportGenesis(ctx cTypes.Context) json.RawMessage {
	return nil
}

func (AppModule) BeginBlock(cTypes.Context, abciTypes.RequestBeginBlock) {}

func (AppModule) EndBlock(_ cTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}
