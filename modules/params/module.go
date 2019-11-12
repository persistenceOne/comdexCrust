package params

import (
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/kafka"
	"github.com/commitHub/commitBlockchain/types/module"

	"github.com/commitHub/commitBlockchain/modules/params/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
)

const moduleName = "params"

// app module basics object
type AppModuleBasic struct{}

// module name
func (AppModuleBasic) Name() string {
	return moduleName
}

// register module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// default genesis state
func (AppModuleBasic) DefaultGenesis() json.RawMessage { return nil }

// module validate genesis
func (AppModuleBasic) ValidateGenesis(_ json.RawMessage) error { return nil }

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router, kafkaBool bool, kafkaState kafka.KafkaState) {
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command { return nil }
