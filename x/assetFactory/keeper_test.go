package assetFactory

import (
	"encoding/json"
	"testing"

	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
)

func initiateSetupMultiStore() (*wire.Codec, sdk.Context, AssetPegMapper, Keeper) {
	ms, assetKey := setupMultiStore()

	cdc := wire.NewCodec()
	RegisterAssetPeg(cdc)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	assetMapper := NewAssetPegMapper(cdc, assetKey, sdk.ProtoBaseAssetPeg)
	assetKeeper := NewKeeper(assetMapper)
	return cdc, ctx, assetMapper, assetKeeper
}

func vars() ([]IssueAsset, []SendAsset) {
	var issueAssetsKeeper = []IssueAsset{
		{ //for owner == issuer
			IssuerAddress: sdk.AccAddress([]byte("issuer")),
			ToAddress:     sdk.AccAddress([]byte("to")),
			AssetPeg: &sdk.BaseAssetPeg{
				PegHash:      sdk.PegHash([]byte("pegHash")),
				OwnerAddress: sdk.AccAddress([]byte("issuer")),
			},
		},
		{ //for from == owner
			IssuerAddress: sdk.AccAddress([]byte("issuer")),
			ToAddress:     sdk.AccAddress([]byte("to")),
			AssetPeg: &sdk.BaseAssetPeg{
				PegHash:      sdk.PegHash([]byte("pegHash")),
				OwnerAddress: sdk.AccAddress([]byte("from")),
			},
		},
		{ //for from ,to,peghash== owner
			IssuerAddress: sdk.AccAddress([]byte("issuer")),
			ToAddress:     sdk.AccAddress([]byte("to")),
			AssetPeg: &sdk.BaseAssetPeg{
				PegHash:      sdk.PegHash([]byte("pegHash")),
				OwnerAddress: sdk.AccAddress([]byte("tofrompegHash")),
			},
		},
	}

	var testSendAsset = []SendAsset{
		{
			RelayerAddress: sdk.AccAddress([]byte("relayer")),
			FromAddress:    sdk.AccAddress([]byte("from")),
			ToAddress:      sdk.AccAddress([]byte("to")),
			PegHash:        sdk.PegHash([]byte("pegHash")),
		},
	}
	return issueAssetsKeeper, testSendAsset
}

func TestInstantiateAndAssignAsset(t *testing.T) {
	issueAssetsKeeper, _ := vars()
	cdc, ctx, assetMapper, _ := initiateSetupMultiStore()
	t.Log(string(assetPeg.GetPegHash()), assetMapper)
	newAsset, tags, _ := instantiateAndAssignAsset(ctx, assetMapper, issueAssetsKeeper[0].IssuerAddress, issueAssetsKeeper[0].ToAddress, issueAssetsKeeper[0].AssetPeg)
	t.Log(newAsset, tags)
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	data, err := cdc.MarshalJSONIndent(tags, "", "")
	if err != nil {
		(panic(err))
	}
	json.Unmarshal(data, &Tags)
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}
	t.Log(tagList)
	require.Nil(t, tagList)

	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[1].AssetPeg)
	_, tag1, _ := instantiateAndAssignAsset(ctx, assetMapper, issueAssetsKeeper[1].IssuerAddress, issueAssetsKeeper[1].ToAddress, issueAssetsKeeper[1].AssetPeg)
	require.Nil(t, tag1)
}

func TestMainissueAssetsToWallets(t *testing.T) {
	issueAssetsKeeper, _ := vars()
	cdc, ctx, assetMapper, _ := initiateSetupMultiStore()
	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[0].AssetPeg)
	listAssets := []IssueAsset{issueAssetsKeeper[0]}
	tags, err := issueAssetsToWallets(ctx, assetMapper, listAssets)
	if err != nil {
		(panic(err))
	}
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList = []jsonData{}
	data, err1 := cdc.MarshalJSONIndent(tags, "", "")
	if err1 != nil {
		(panic(err1))
	}
	json.Unmarshal(data, &Tags)
	var tagCheck = []string{"recepient", "issuer", "asset"}
	for _, tag := range Tags {
		tagList = append(tagList, jsonData{string(tag.Key), string(tag.Value)})
		require.Equal(t, tagCheck[0], tagList[0].key)
	}
	t.Log(tagList)
	require.Equal(t, "recepient", tagList[0].key)
	require.Equal(t, "issuer", tagList[1].key)
	require.Equal(t, "asset", tagList[2].key)
	require.Equal(t, issueAssetsKeeper[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, issueAssetsKeeper[0].IssuerAddress.String(), tagList[1].value)
	require.Equal(t, issueAssetsKeeper[0].AssetPeg.GetPegHash().String(), tagList[2].value)

	tag1, _ := issueAssetsToWallets(ctx, assetMapper, issueAssetsKeeper[1:2])
	require.Nil(t, tag1)
}

func TestIssueAssetsToWallets(t *testing.T) {
	issueAssetsKeeper, _ := vars()

	cdc, ctx, assetMapper, assetKeeper := initiateSetupMultiStore()

	assetPeg := issueAssetsKeeper[0].AssetPeg
	assetMapper.SetAssetPeg(ctx, assetPeg)
	tags, _ := assetKeeper.IssueAssetsToWallets(ctx, issueAssetsKeeper[:1])

	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	data, _ := cdc.MarshalJSONIndent(tags, "", "")
	json.Unmarshal(data, &Tags)
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}

	t.Log(tagList)
	require.Equal(t, "recepient", tagList[0].key)
	require.Equal(t, "issuer", tagList[1].key)
	require.Equal(t, "asset", tagList[2].key)
	require.Equal(t, issueAssetsKeeper[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, issueAssetsKeeper[0].IssuerAddress.String(), tagList[1].value)
	require.Equal(t, issueAssetsKeeper[0].AssetPeg.GetPegHash().String(), tagList[2].value)
}

func TestSendAssetsToOrder(t *testing.T) {
	issueAssetsKeeper, testSendAsset := vars()

	cdc, ctx, assetMapper, _ := initiateSetupMultiStore()

	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[1].AssetPeg)

	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData

	asset, tags, _ := sendAssetToOrder(ctx, assetMapper, testSendAsset[0].FromAddress, testSendAsset[0].ToAddress, testSendAsset[0].PegHash)
	t.Log(asset)
	data, _ := cdc.MarshalJSONIndent(tags, "", "")
	json.Unmarshal(data, &Tags)
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}

	t.Log(tagList)
	require.Equal(t, "recepient", tagList[0].key)
	require.Equal(t, "sender", tagList[1].key)
	require.Equal(t, "asset", tagList[2].key)
	require.Equal(t, testSendAsset[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, testSendAsset[0].FromAddress.String(), tagList[1].value)
	require.Equal(t, string(testSendAsset[0].PegHash), tagList[2].value)
}

func TestSendAssetToOrders(t *testing.T) {
	issueAssetsKeeper, testSendAsset := vars()

	cdc, ctx, assetMapper, _ := initiateSetupMultiStore()

	_, tag1, _ := sendAssetToOrder(ctx, assetMapper, testSendAsset[0].FromAddress, testSendAsset[0].ToAddress, testSendAsset[0].PegHash)
	t.Log(tag1)
	require.Nil(t, tag1)

	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[1].AssetPeg)

	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	tags, _ := sendAssetsToOrders(ctx, assetMapper, testSendAsset)
	data, _ := cdc.MarshalJSONIndent(tags, "", "")
	json.Unmarshal(data, &Tags)
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}
	t.Log(tagList)
	require.Equal(t, "recepient", tagList[0].key)
	require.Equal(t, "sender", tagList[1].key)
	require.Equal(t, "asset", tagList[2].key)
	require.Equal(t, testSendAsset[0].FromAddress.String(), tagList[1].value)
	require.Equal(t, string(testSendAsset[0].PegHash), tagList[2].value)
}

func TestMainSendAssetsToOrders(t *testing.T) {
	issueAssetsKeeper, testSendAsset := vars()
	cdc, ctx, assetMapper, assetKeeper := initiateSetupMultiStore()

	tag1, _ := assetKeeper.SendAssetsToOrders(ctx, testSendAsset)
	t.Log(tag1)
	require.Nil(t, tag1)
	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[1].AssetPeg)
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	tags, _ := assetKeeper.SendAssetsToOrders(ctx, testSendAsset)
	data, _ := cdc.MarshalJSONIndent(tags, "", "")
	json.Unmarshal(data, &Tags)
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}

	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[2].AssetPeg)
	tag2, _ := assetKeeper.SendAssetsToOrders(ctx, testSendAsset)
	require.Nil(t, tag2)

	t.Log(tagList)
	require.Equal(t, "recepient", tagList[0].key)
	require.Equal(t, "sender", tagList[1].key)
	require.Equal(t, "asset", tagList[2].key)
	require.Equal(t, testSendAsset[0].FromAddress.String(), tagList[1].value)
	require.Equal(t, string(testSendAsset[0].PegHash), tagList[2].value)
}

func TestSendAssetFromOrder(t *testing.T) {
	issueAssetsKeeper, testSendAsset := vars()
	cdc, ctx, assetMapper, _ := initiateSetupMultiStore()
	t.Log(string(assetPeg.GetPegHash()), assetMapper)

	asset, tag1, _ := sendAssetFromOrder(ctx, assetMapper, testSendAsset[0].FromAddress, testSendAsset[0].ToAddress, testSendAsset[0].PegHash)
	t.Log(asset, tag1)

	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[2].AssetPeg)
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData

	_, tags, _ := sendAssetFromOrder(ctx, assetMapper, testSendAsset[0].FromAddress, testSendAsset[0].ToAddress, testSendAsset[0].PegHash)
	data, _ := cdc.MarshalJSONIndent(tags, "", "")
	t.Log(tags)
	json.Unmarshal(data, &Tags)
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}

	t.Log(tagList)

	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[1].AssetPeg)
	_, tag2, _ := sendAssetFromOrder(ctx, assetMapper, testSendAsset[0].FromAddress, testSendAsset[0].ToAddress, testSendAsset[0].PegHash)
	require.Nil(t, tag2)

	require.Equal(t, "recepient", tagList[0].key, nil)
	require.Equal(t, "sender", tagList[1].key, nil)
	require.Equal(t, "asset", tagList[2].key, nil)
	require.Equal(t, testSendAsset[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, string(testSendAsset[0].PegHash), tagList[2].value)
}

func TestMainExecuteAssetOrders(t *testing.T) {
	issueAssetsKeeper, testSendAsset := vars()

	cdc, ctx, assetMapper, _ := initiateSetupMultiStore()
	tags1, _ := executeAssetOrders(ctx, assetMapper, testSendAsset)
	require.Nil(t, tags1)

	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[2].AssetPeg)

	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData

	tags, _ := executeAssetOrders(ctx, assetMapper, testSendAsset)
	data, _ := cdc.MarshalJSONIndent(tags, "", "")
	json.Unmarshal(data, &Tags)
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}
	t.Log(tagList)
	require.Equal(t, "recepient", tagList[0].key)
	require.Equal(t, "sender", tagList[1].key)
	require.Equal(t, "asset", tagList[2].key)
	require.Equal(t, testSendAsset[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, string(testSendAsset[0].PegHash), tagList[2].value)
}

func TestExecuteAssetOrders(t *testing.T) {
	issueAssetsKeeper, testSendAsset := vars()
	cdc, ctx, assetMapper, assetKeeper := initiateSetupMultiStore()

	tag1, _ := assetKeeper.SendAssetsToOrders(ctx, testSendAsset)
	t.Log(tag1)
	require.Nil(t, tag1)
	assetMapper.SetAssetPeg(ctx, issueAssetsKeeper[2].AssetPeg)

	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData

	tags, _ := assetKeeper.ExecuteAssetOrders(ctx, testSendAsset)
	data, _ := cdc.MarshalJSONIndent(tags, "", "")
	json.Unmarshal(data, &Tags)
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}
	t.Log(tagList)
	require.Equal(t, "recepient", tagList[0].key)
	require.Equal(t, "sender", tagList[1].key)
	require.Equal(t, "asset", tagList[2].key)
	require.Equal(t, testSendAsset[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, string(testSendAsset[0].PegHash), tagList[2].value)
}
