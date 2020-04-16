package fiatFactory

import (
	"github.com/persistenceOne/comdexCrust/modules/fiatFactory/internal/keeper"
	"github.com/persistenceOne/comdexCrust/modules/fiatFactory/internal/types"
)

const (
	ModuleName       = types.ModuleName
	RouterKey        = types.RouterKey
	QuerierRoute     = types.QuerierRoute
	DefaultCodeSpace = types.DefaultCodeSpace
	StoreKey         = types.StoreKey

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

	NewQuerier = keeper.NewQuerier
	NewKeeper  = keeper.NewKeeper

	EventTypeFiatFactoryAssignFiat  = types.EventTypeFiatFactoryAssignFiat
	EventTypeFiatFactoryRedeemFiat  = types.EventTypeFiatFactoryRedeemFiat
	EventTypeFiatFactorySendFiat    = types.EventTypeFiatFactorySendFiat
	EventTypeFiatFactoryExecuteFiat = types.EventTypeFiatFactoryExecuteFiat

	NewIssueFiat             = types.NewIssueFiat
	NewMsgFactoryIssueFiats  = types.NewMsgFactoryIssueFiats
	NewRedeemFiat            = types.NewRedeemFiat
	NewMsgFactoryRedeemFiats = types.NewMsgFactoryRedeemFiats
	NewSendFiat              = types.NewSendFiat
	NewMsgFactorySendFiats   = types.NewMsgFactorySendFiats
)

type (
	GenesisState = types.GenesisState
	Keeper       = keeper.Keeper

	MsgFactoryIssueFiats   = types.MsgFactoryIssueFiats
	MsgFactoryRedeemFiats  = types.MsgFactoryRedeemFiats
	MsgFactorySendFiats    = types.MsgFactorySendFiats
	MsgFactoryExecuteFiats = types.MsgFactoryExecuteFiats

	IssueFiat  = types.IssueFiat
	RedeemFiat = types.RedeemFiat
	SendFiat   = types.SendFiat
)
