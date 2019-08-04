package genutil

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commitHub/commitBlockchain/codec"

	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/staking"
)

// generic sealed codec to be used throughout this module
var moduleCdc *codec.Codec

// TODO abstract genesis transactions registration back to staking
// required for genesis transactions
func init() {
	moduleCdc = codec.New()
	staking.RegisterCodec(moduleCdc)
	auth.RegisterCodec(moduleCdc)
	sdk.RegisterCodec(moduleCdc)
	codec.RegisterCrypto(moduleCdc)
	moduleCdc.Seal()
}
