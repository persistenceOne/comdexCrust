package assetFactory

import (
	"github.com/persistenceOne/comdexCrust/modules/assetFactory/internal/keeper"
	"github.com/persistenceOne/comdexCrust/modules/assetFactory/internal/types"
)

const (
	ModuleName       = types.ModuleName
	RouterKey        = types.RouterKey
	QuerierRoute     = types.QuerierRoute
	DefaultCodeSpace = types.DefaultCodeSpace
	StoreKey         = types.StoreKey
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	DefaultGenesisState  = types.DefaultGenesisState
	ValidateGenesis      = types.ValidateGenesis
	AssetPegHashStoreKey = types.AssetPegHashStoreKey

	NewQuerier = keeper.NewQuerier
	NewKeeper  = keeper.NewKeeper

	EventTypeAssetFactoryIssueAsset   = types.EventTypeAssetFactoryIssueAsset
	EventTypeAssetFactoryRedeemAsset  = types.EventTypeAssetFactoryRedeemAsset
	EventTypeAssetFactorySendAsset    = types.EventTypeAssetFactorySendAsset
	EventTypeAssetFactoryExecuteAsset = types.EventTypeAssetFactoryExecuteAsset

	NewIssueAsset             = types.NewIssueAsset
	NewMsgFactoryIssueAssets  = types.NewMsgFactoryIssueAssets
	NewRedeemAsset            = types.NewRedeemAsset
	NewMsgFactoryRedeemAssets = types.NewMsgFactoryRedeemAssets
	NewSendAsset              = types.NewSendAsset
	NewMsgFactorySendAssets   = types.NewMsgFactorySendAssets
)

type (
	GenesisState = types.GenesisState
	Keeper       = keeper.Keeper

	AccountKeeper           = types.AccountKeeper
	MsgFactoryIssueAssets   = types.MsgFactoryIssueAssets
	MsgFactoryRedeemAssets  = types.MsgFactoryRedeemAssets
	MsgFactorySendAssets    = types.MsgFactorySendAssets
	MsgFactoryExecuteAssets = types.MsgFactoryExecuteAssets

	IssueAsset  = types.IssueAsset
	RedeemAsset = types.RedeemAsset
	SendAsset   = types.SendAsset
)
