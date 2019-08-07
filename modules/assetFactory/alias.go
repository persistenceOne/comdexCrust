package assetFactory

import (
	"github.com/commitHub/commitBlockchain/modules/assetFactory/internal/keeper"
	"github.com/commitHub/commitBlockchain/modules/assetFactory/internal/types"
)

const (
	ModuleName       = types.ModuleName
	RouterKey        = types.RouterKey
	QuerierRoute     = types.QuerierRoute
	DefaultCodeSpace = types.DefaultCodeSpace
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc
	PegHashKey    = types.PegHashKey
	
	AssetPegHashStoreKey = types.AssetPegHashStoreKey
	
	DefaultGenesisState               = types.DefaultGenesisState
	ValidateGenesis                   = types.ValidateGenesis
	EventTypeAssetFactoryIssueAsset   = types.EventTypeAssetFactoryIssueAsset
	EventTypeAssetFactoryRedeemAsset  = types.EventTypeAssetFactoryRedeemAsset
	EventTypeAssetFactorySendAsset    = types.EventTypeAssetFactorySendAsset
	EventTypeAssetFactoryExecuteAsset = types.EventTypeAssetFactoryExecuteAsset
	
	BuildIssueAssetMsg   = types.BuildIssueAssetMsg
	BuildSendAssetMsg    = types.BuildSendAssetMsg
	BuildExecuteAssetMsg = types.BuildExecuteAssetMsg
	BuildRedeemAssetMsg  = types.BuildRedeemAssetMsg
)

type (
	GenesisState = types.GenesisState
	Keeper = keeper.Keeper
	
	AccountKeeper = types.AccountKeeper
	MsgFactoryIssueAssets = types.MsgFactoryIssueAssets
	MsgFactoryRedeemAssets = types.MsgFactoryRedeemAssets
	MsgFactorySendAssets = types.MsgFactorySendAssets
	MsgFactoryExecuteAssets = types.MsgFactoryExecuteAssets
	
	IssueAsset = types.IssueAsset
)
