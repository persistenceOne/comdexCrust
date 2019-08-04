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
	PegHashKey          = types.PegHashKey
	FiatPegHashStoreKey = types.FiatPegHashStoreKey

	EventTypeFiatFactoryAssignFiat  = types.EventTypeFiatFactoryAssignFiat
	EventTypeFiatFactoryRedeemFiat  = types.EventTypeFiatFactoryRedeemFiat
	EventTypeFiatFactorySendFiat    = types.EventTypeFiatFactorySendFiat
	EventTypeFiatFactoryExecuteFiat = types.EventTypeFiatFactoryExecuteFiat

	ErrInvalidAmount  = types.ErrInvalidAmount
	ErrInvalidString  = types.ErrInvalidString
	ErrNoInputs       = types.ErrNoInputs
	ErrInvalidPegHash = types.ErrInvalidPegHash

	BuildIssueFiatMsg   = types.BuildIssueFiatMsg
	BuildSendFiatMsg    = types.BuildSendFiatMsg
	BuildExecuteFiatMsg = types.BuildExecuteFiatMsg
	BuildRedeemFiatMsg  = types.BuildRedeemFiatMsg
)

type (
	GenesisState = types.GenesisState
	Keeper       = keeper.Keeper

	MsgFactoryIssueFiats   = types.MsgFactoryIssueFiats
	MsgFactoryRedeemFiats  = types.MsgFactoryRedeemFiats
	MsgFactorySendFiats    = types.MsgFactorySendFiats
	MsgFactoryExecuteFiats = types.MsgFactoryExecuteFiats
)
