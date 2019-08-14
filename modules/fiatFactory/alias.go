package fiatFactory

import (
	"github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/fiatFactory/internal/types"
)

const (
	ModuleName               = types.ModuleName
	QuerierRoute             = types.QuerierRoute
	DefaultCodeSpace         = types.DefaultCodeSpace
	CodeInvalidAmount        = types.CodeInvalidAmount
	CodeInvalidString        = types.CodeInvalidString
	CodeInvalidInputsOutputs = types.CodeInvalidInputsOutputs
	CodeInvalidPegHash       = types.CodeInvalidPegHash
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc
	
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	FiatPegHashStoreKey = types.FiatPegHashStoreKey
	
	EventTypeFiatFactoryAssignFiat  = types.EventTypeFiatFactoryAssignFiat
	EventTypeFiatFactoryRedeemFiat  = types.EventTypeFiatFactoryRedeemFiat
	EventTypeFiatFactorySendFiat    = types.EventTypeFiatFactorySendFiat
	EventTypeFiatFactoryExecuteFiat = types.EventTypeFiatFactoryExecuteFiat
	
	ErrNoInputs       = types.ErrNoInputs
	ErrInvalidPegHash = types.ErrInvalidPegHash
)

type (
	GenesisState = types.GenesisState
	Keeper       = keeper.Keeper
	
	MsgFactoryIssueFiats   = types.MsgFactoryIssueFiats
	MsgFactoryRedeemFiats  = types.MsgFactoryRedeemFiats
	MsgFactorySendFiats    = types.MsgFactorySendFiats
	MsgFactoryExecuteFiats = types.MsgFactoryExecuteFiats
)
