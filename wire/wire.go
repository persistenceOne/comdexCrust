package wire

import (
	"bytes"
	"encoding/json"

	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// Codec amino codec to marshal/unmarshal
type Codec = amino.Codec

// NewCodec : returns new codec
func NewCodec() *Codec {
	cdc := amino.NewCodec()
	return cdc
}

// RegisterCrypto : Register the go-crypto to the codec
func RegisterCrypto(cdc *Codec) {
	cryptoAmino.RegisterAmino(cdc)
}

// MarshalJSONIndent : attempt to make some pretty json
func MarshalJSONIndent(cdc *Codec, obj interface{}) ([]byte, error) {
	bz, err := cdc.MarshalJSON(obj)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = json.Indent(&out, bz, "", "  ")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

//__________________________________________________________________

// Cdc : generic sealed codec to be used throughout sdk
var Cdc *Codec

func init() {
	cdc := NewCodec()
	RegisterCrypto(cdc)
	Cdc = cdc.Seal()
}
