package types

import (
	"context"
	"sync"

	"github.com/golang/protobuf/proto"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

/*
The intent of Context is for it to be an immutable object that can be
cloned and updated cheaply with WithValue() and passed forward to the
next decorator or handler. For example,

 func MsgHandler(ctx Context, tx Tx) Result {
 	...
 	ctx = ctx.WithValue(key, value)
 	...
 }
*/

// Context : context structure
type Context struct {
	context.Context
	pst *thePast
	gen int
	// Don't add any other fields here,
	// it's probably not what you want to do.
}

// NewContext create a new context
func NewContext(ms MultiStore, header abci.Header, isCheckTx bool, logger log.Logger) Context {
	c := Context{
		Context: context.Background(),
		pst:     newThePast(),
		gen:     0,
	}
	c = c.WithMultiStore(ms)
	c = c.WithBlockHeader(header)
	c = c.WithBlockHeight(header.Height)
	c = c.WithChainID(header.ChainID)
	c = c.WithTxBytes(nil)
	c = c.WithLogger(logger)
	c = c.WithSigningValidators(nil)
	c = c.WithGasMeter(NewInfiniteGasMeter())
	return c
}

// IsZero is context nil
func (c Context) IsZero() bool {
	return c.Context == nil
}

//----------------------------------------
// Getting a value

// Value context value for the provided key
func (c Context) Value(key interface{}) interface{} {
	value := c.Context.Value(key)
	if cloner, ok := value.(cloner); ok {
		return cloner.Clone()
	}
	if message, ok := value.(proto.Message); ok {
		return proto.Clone(message)
	}
	return value
}

// KVStore fetches a KVStore from the MultiStore.
func (c Context) KVStore(key StoreKey) KVStore {
	return c.multiStore().GetKVStore(key).Gas(c.GasMeter(), cachedDefaultGasConfig)
}

// TransientStore fetches a TransientStore from the MultiStore.
func (c Context) TransientStore(key StoreKey) KVStore {
	return c.multiStore().GetKVStore(key).Gas(c.GasMeter(), cachedTransientGasConfig)
}

//----------------------------------------
// With* (setting a value)

// WithValue nolint
func (c Context) WithValue(key interface{}, value interface{}) Context {
	return c.withValue(key, value)
}

// WithCloner returns context with withCloner
func (c Context) WithCloner(key interface{}, value cloner) Context {
	return c.withValue(key, value)
}

// WithCacheWrapper returns context with withcachewrapper
func (c Context) WithCacheWrapper(key interface{}, value CacheWrapper) Context {
	return c.withValue(key, value)
}

//WithProtoMsg returns context with WithProtoMsg
func (c Context) WithProtoMsg(key interface{}, value proto.Message) Context {
	return c.withValue(key, value)
}

// WithString returns context with WithString
func (c Context) WithString(key interface{}, value string) Context {
	return c.withValue(key, value)
}

// WithInt32 returns context with WithInt32
func (c Context) WithInt32(key interface{}, value int32) Context {
	return c.withValue(key, value)
}

// WithUint32 returns context with WithUint32
func (c Context) WithUint32(key interface{}, value uint32) Context {
	return c.withValue(key, value)
}

//WithUint64 returns context with WithUint64
func (c Context) WithUint64(key interface{}, value uint64) Context {
	return c.withValue(key, value)
}

func (c Context) withValue(key interface{}, value interface{}) Context {
	c.pst.bump(Op{
		gen:   c.gen + 1,
		key:   key,
		value: value,
	}) // increment version for all relatives.

	return Context{
		Context: context.WithValue(c.Context, key, value),
		pst:     c.pst,
		gen:     c.gen + 1,
	}
}

//----------------------------------------
// Values that require no key.

type contextKey int // local to the context module

const (
	contextKeyMultiStore contextKey = iota
	contextKeyBlockHeader
	contextKeyBlockHeight
	contextKeyConsensusParams
	contextKeyChainID
	contextKeyTxBytes
	contextKeyLogger
	contextKeySigningValidators
	contextKeyGasMeter
)

// NOTE: Do not expose MultiStore.
// MultiStore exposes all the keys.
// Instead, pass the context and the store key.
func (c Context) multiStore() MultiStore {
	return c.Value(contextKeyMultiStore).(MultiStore)
}

// BlockHeader : returns blockheader
func (c Context) BlockHeader() abci.Header {
	return c.Value(contextKeyBlockHeader).(abci.Header)
}

// BlockHeight : returns blockHeight
func (c Context) BlockHeight() int64 {
	return c.Value(contextKeyBlockHeight).(int64)
}

// ConsensusParams : returns consensusparams
func (c Context) ConsensusParams() abci.ConsensusParams {
	return c.Value(contextKeyConsensusParams).(abci.ConsensusParams)
}

// ChainID : returns Chainid
func (c Context) ChainID() string {
	return c.Value(contextKeyChainID).(string)
}

// TxBytes : returns bytes
func (c Context) TxBytes() []byte {
	return c.Value(contextKeyTxBytes).([]byte)
}

// Logger : returns logger
func (c Context) Logger() log.Logger {
	return c.Value(contextKeyLogger).(log.Logger)
}

// SigningValidators : returns the validator who signs it
func (c Context) SigningValidators() []abci.SigningValidator {
	return c.Value(contextKeySigningValidators).([]abci.SigningValidator)
}

//GasMeter :
func (c Context) GasMeter() GasMeter {
	return c.Value(contextKeyGasMeter).(GasMeter)
}

// WithMultiStore returns context with multi store
func (c Context) WithMultiStore(ms MultiStore) Context {
	return c.withValue(contextKeyMultiStore, ms)
}

// WithBlockHeader : returns context with blockheader
func (c Context) WithBlockHeader(header abci.Header) Context {
	var _ proto.Message = &header // for cloning.
	return c.withValue(contextKeyBlockHeader, header)
}

// WithBlockHeight returns context with blockheight
func (c Context) WithBlockHeight(height int64) Context {
	return c.withValue(contextKeyBlockHeight, height)
}

// WithConsensusParams returns context
func (c Context) WithConsensusParams(params *abci.ConsensusParams) Context {
	if params == nil {
		return c
	}
	return c.withValue(contextKeyConsensusParams, params).
		WithGasMeter(NewGasMeter(params.TxSize.MaxGas))
}

//WithChainID : returns context WithChainID
func (c Context) WithChainID(chainID string) Context {
	return c.withValue(contextKeyChainID, chainID)
}

// WithTxBytes : returns context with txbytes
func (c Context) WithTxBytes(txBytes []byte) Context {
	return c.withValue(contextKeyTxBytes, txBytes)
}

//WithLogger : returns context with WithLogger
func (c Context) WithLogger(logger log.Logger) Context {
	return c.withValue(contextKeyLogger, logger)
}

// WithSigningValidators : returns context
func (c Context) WithSigningValidators(SigningValidators []abci.SigningValidator) Context {
	return c.withValue(contextKeySigningValidators, SigningValidators)
}

// WithGasMeter :returns context
func (c Context) WithGasMeter(meter GasMeter) Context {
	return c.withValue(contextKeyGasMeter, meter)
}

// CacheContext : Cache the multistore and return a new cached context. The cached context is
// written to the context when writeCache is called.
func (c Context) CacheContext() (cc Context, writeCache func()) {
	cms := c.multiStore().CacheMultiStore()
	cc = c.WithMultiStore(cms)
	return cc, cms.Write
}

//----------------------------------------
// thePast

// GetOp : Returns false if ver <= 0 || ver > len(c.pst.ops).
// The first operation is version 1.
func (c Context) GetOp(ver int64) (Op, bool) {
	return c.pst.getOp(ver)
}

//----------------------------------------
// Misc.

type cloner interface {
	Clone() interface{} // deep copy
}

// Op struct
type Op struct {
	// type is always 'with'
	gen   int
	key   interface{}
	value interface{}
}

type thePast struct {
	mtx sync.RWMutex
	ver int
	ops []Op
}

func newThePast() *thePast {
	return &thePast{
		ver: 0,
		ops: nil,
	}
}

func (pst *thePast) bump(op Op) {
	pst.mtx.Lock()
	pst.ver++
	pst.ops = append(pst.ops, op)
	pst.mtx.Unlock()
}

func (pst *thePast) version() int {
	pst.mtx.RLock()
	defer pst.mtx.RUnlock()
	return pst.ver
}

// Returns false if ver <= 0 || ver > len(pst.ops).
// The first operation is version 1.
func (pst *thePast) getOp(ver int64) (Op, bool) {
	pst.mtx.RLock()
	defer pst.mtx.RUnlock()
	l := int64(len(pst.ops))
	if l < ver || ver <= 0 {
		return Op{}, false
	}
	return pst.ops[ver-1], true
}
