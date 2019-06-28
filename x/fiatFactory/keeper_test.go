package fiatFactory

import (
	"encoding/json"
	"testing"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/stretchr/testify/require"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// TestInstantiateAndAssignFiat tests function instantiateAndAssignFiat
func TestInstantiateAndAssignFiat(t *testing.T) {
	cdc, ctx, fiatMapper, _ := initiateSetupMultiStore()
	fiatPeg := issueFiat[0].FiatPeg
	fiatMapper.SetFiatPeg(ctx, fiatPeg)
	_, tags, _ := instantiateAndAssignFiat(ctx, fiatMapper, issueFiat[0].IssuerAddress, issueFiat[0].ToAddress, issueFiat[0].FiatPeg)
	t.Log(tags)
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, issueFiat[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, issueFiat[0].IssuerAddress.String(), tagList[1].value)
	require.Equal(t, issueFiat[0].FiatPeg.GetPegHash().String(), tagList[2].value)
}

// TestMainIssueFiatsToWallets tests function issueFiatsToWallets
func TestMainIssueFiatsToWallets(t *testing.T) {
	cdc, ctx, fiatMapper, _ := initiateSetupMultiStore()
	
	fiatPeg := issueFiat[0].FiatPeg
	fiatMapper.SetFiatPeg(ctx, fiatPeg)
	tags, _ := issueFiatsToWallets(ctx, fiatMapper, issueFiat)
	
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, issueFiat[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, issueFiat[0].IssuerAddress.String(), tagList[1].value)
	require.Equal(t, issueFiat[0].FiatPeg.GetPegHash().String(), tagList[2].value)
}

// TestIssueFiatsToWallets tests function IssueFiatsToWallets
func TestIssueFiatsToWallets(t *testing.T) {
	cdc, ctx, fiatMapper, fiatKeeper := initiateSetupMultiStore()
	
	fiatPeg := issueFiat[0].FiatPeg
	fiatMapper.SetFiatPeg(ctx, fiatPeg)
	tags, _ := fiatKeeper.IssueFiatsToWallets(ctx, issueFiat)
	
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	data, _ := cdc.MarshalJSONIndent(tags, "", "")
	json.Unmarshal(data, &Tags)
	// Tags.Sort()
	var tagData jsonData
	for _, tag := range Tags {
		tagData.key = string(tag.Key)
		tagData.value = string(tag.Value)
		tagList = append(tagList, tagData)
	}
	
	t.Log(tagList)
	require.Equal(t, "recepient", tagList[0].key)
	require.Equal(t, "issuer", tagList[1].key)
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, issueFiat[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, issueFiat[0].IssuerAddress.String(), tagList[1].value)
	require.Equal(t, issueFiat[0].FiatPeg.GetPegHash().String(), tagList[2].value)
}

// TestIssueFiatsToWalletsError handles error case for the IssueFiatsToWallets
func TestIssueFiatsToWalletsError(t *testing.T) {
	_, ctx, _, fiatKeeper := initiateSetupMultiStore()
	
	tags, _ := fiatKeeper.IssueFiatsToWallets(ctx, issueFiat)
	t.Log(tags)
	require.Nil(t, tags)
}

// TestRedeemFiatsFromWallet handles test case for the RedeemFiatsFromWallets
func TestRedeemFiatsFromWallet(t *testing.T) {
	cdc, ctx, fiatMapper, fiatKeeper := initiateSetupMultiStore()
	fiatPegWallet := testRedeemFiat[0].FiatPegWallet
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	tags, _ := fiatKeeper.RedeemFiatsFromWallets(ctx, testRedeemFiat)
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
	require.Equal(t, testRedeemFiat[0].RedeemerAddress.String(), tagList[0].value)
	require.Equal(t, testRedeemFiat[0].RedeemerAddress.String(), tagList[1].value)
}

// TestSendFiats tests function sendFiats
func TestSendFiats(t *testing.T) {
	cdc, ctx, fiatMapper, _ := initiateSetupMultiStore()
	
	fiatPegWallet := testSendFiat[0].FiatPegWallet
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	
	tags, _ := sendFiats(testSendFiat[0].FiatPegWallet, fiatMapper, ctx, testSendFiat[0].PegHash, testSendFiat[0].FromAddress, testSendFiat[0].ToAddress)
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, testSendFiat[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, testSendFiat[0].FromAddress.String(), tagList[1].value)
	require.Equal(t, string(testSendFiat[0].PegHash), tagList[2].value)
}

// TestSendFiatToOrder tests function sendFiatToOrder
func TestSendFiatToOrder(t *testing.T) {
	cdc, ctx, fiatMapper, _ := initiateSetupMultiStore()
	fiatPegWallet := testSendFiat[0].FiatPegWallet
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	tags, _ := sendFiatToOrder(ctx, fiatMapper, testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash, testSendFiat[0].FiatPegWallet)
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash)).String(), tagList[0].value)
	require.Equal(t, testSendFiat[0].FromAddress.String(), tagList[1].value)
	require.Equal(t, string(testSendFiat[0].PegHash), tagList[2].value)
}

// TestMainSendFiatsToOrders tests function sendFiatsToOrders
func TestMainSendFiatsToOrders(t *testing.T) {
	cdc, ctx, fiatMapper, _ := initiateSetupMultiStore()
	fiatPegWallet := testSendFiat[0].FiatPegWallet
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	tags, _ := sendFiatsToOrders(ctx, fiatMapper, testSendFiat)
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash)).String(), tagList[0].value)
	require.Equal(t, testSendFiat[0].FromAddress.String(), tagList[1].value)
	require.Equal(t, string(testSendFiat[0].PegHash), tagList[2].value)
}

// TestSendFiatsToOrdersError
func TestSendFiatsToOrdersError(t *testing.T) {
	_, ctx, _, fiatKeeper := initiateSetupMultiStore()
	
	tags, _ := fiatKeeper.SendFiatsToOrders(ctx, testSendFiat)
	require.Nil(t, tags)
}

func TestSendFiatsToOrders(t *testing.T) {
	cdc, ctx, fiatMapper, fiatKeeper := initiateSetupMultiStore()
	fiatPegWallet := testSendFiat[0].FiatPegWallet
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	tags, _ := fiatKeeper.SendFiatsToOrders(ctx, testSendFiat)
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash)).String(), tagList[0].value)
	require.Equal(t, testSendFiat[0].FromAddress.String(), tagList[1].value)
	require.Equal(t, string(testSendFiat[0].PegHash), tagList[2].value)
}

func TestSendFiatFromOrder(t *testing.T) {
	cdc, ctx, fiatMapper, _ := initiateSetupMultiStore()
	fiatPegWallet := testSendFiat[0].FiatPegWallet
	testSendFiat[0].FiatPegWallet[0].Owners[0].OwnerAddress = sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash))
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	
	tags, _ := sendFiatFromOrder(ctx, fiatMapper, testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash, testSendFiat[0].FiatPegWallet)
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, testSendFiat[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash)).String(), tagList[1].value)
	require.Equal(t, string(testSendFiat[0].PegHash), tagList[2].value)
}

func TestMainExecuteFiatOrdersError(t *testing.T) {
	var testSendFiat1 = []SendFiat{
		{
			RelayerAddress: sdk.AccAddress([]byte("relayer")),
			FromAddress:    sdk.AccAddress([]byte("from")),
			ToAddress:      sdk.AccAddress([]byte("to")),
			PegHash:        sdk.PegHash([]byte("pegHash")),
			FiatPegWallet: sdk.FiatPegWallet{
				sdk.BaseFiatPeg{
					PegHash:           sdk.PegHash([]byte("pegHash1")),
					TransactionID:     "FB8AE3A02BBCD2",
					TransactionAmount: 9,
					RedeemedAmount:    3,
					Owners: []sdk.Owner{
						{
							OwnerAddress: sdk.AccAddress([]byte("from")),
							Amount:       2000,
						},
						{
							OwnerAddress: sdk.AccAddress([]byte("relayer")),
							Amount:       3000,
						},
					},
				},
			},
		},
	}
	
	_, ctx, fiatMapper, _ := initiateSetupMultiStore()
	fiatPegWallet := testSendFiat1[0].FiatPegWallet
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	
	tags, _ := executeFiatOrders(ctx, fiatMapper, testSendFiat1)
	t.Log(tags)
	require.Nil(t, tags)
}
func TestMainExecuteFiatOrders(t *testing.T) {
	cdc, ctx, fiatMapper, _ := initiateSetupMultiStore()
	fiatPegWallet := testSendFiat[0].FiatPegWallet
	testSendFiat[0].FiatPegWallet[0].Owners[0].OwnerAddress = sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash))
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	
	tags, _ := executeFiatOrders(ctx, fiatMapper, testSendFiat)
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, testSendFiat[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash)).String(), tagList[1].value)
	require.Equal(t, string(testSendFiat[0].PegHash), tagList[2].value)
}

func TestExecuteFiatOrders(t *testing.T) {
	cdc, ctx, fiatMapper, fiatKeeper := initiateSetupMultiStore()
	fiatPegWallet := testSendFiat[0].FiatPegWallet
	testSendFiat[0].FiatPegWallet[0].Owners[0].OwnerAddress = sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash))
	for _, fiat := range fiatPegWallet {
		fiatMapper.SetFiatPeg(ctx, &fiat)
	}
	var Tags cmn.KVPairs
	type jsonData struct {
		key   string
		value string
	}
	var tagList []jsonData
	
	tags, _ := fiatKeeper.ExecuteFiatOrders(ctx, testSendFiat)
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
	require.Equal(t, "fiat", tagList[2].key)
	require.Equal(t, testSendFiat[0].ToAddress.String(), tagList[0].value)
	require.Equal(t, sdk.AccAddress(sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash)).String(), tagList[1].value)
	require.Equal(t, string(testSendFiat[0].PegHash), tagList[2].value)
}
