package assetFactory

import (
	"testing"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	var testsHandlerIssues = []TestCase{
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC123", "assetType", 5, "quantityUnit", addrs[0], true},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peg")), "ABC123", "assetType", 5, "quantityUnit", sdk.AccAddress([]byte("from")), true},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peg2")), "ABC123", "assetType", 5, "quantityUnit", sdk.AccAddress([]byte("tofrompeg2")), true},
	}
	var testsHandlerSends = []TestCaseSend{
		{addrs[0], addrs[1], addrs[2], sdk.PegHash([]byte("peg")), true},
		{addrs[0], addrs[1], addrs[2], sdk.PegHash([]byte("peg2")), true},
	}
	
	cdc, ctx, mapper, keeper := initiateSetupMultiStore()
	t.Logf("\n%v \n\n%v \n\n%v \n\n%v \n\n", cdc, ctx, mapper, keeper)
	testingFunc := NewHandler(keeper)
	oneAssetPeg, oneIssueAsset := genIssueAsset(testsHandlerIssues[0])
	issueAssets := []IssueAsset{oneIssueAsset}
	
	msgIssue := NewMsgFactoryIssueAssets(issueAssets)
	require.Equal(t, testingFunc(ctx, msgIssue).Tags, sdk.Tags(nil))
	
	oneSendAsset := genSendAsset(testsHandlerSends[0])
	sendAssets := []SendAsset{oneSendAsset}
	oneAssetPegSend, _ := genIssueAsset(testsHandlerIssues[1])
	msgSend := NewMsgFactorySendAssets(sendAssets)
	require.Equal(t, testingFunc(ctx, msgSend).Tags, sdk.Tags(nil))
	
	oneAssetExecute := genSendAsset(testsHandlerSends[1])
	sendAssetsExecute := []SendAsset{oneAssetExecute}
	oneAssetPegExecute, _ := genIssueAsset(testsHandlerIssues[2])
	msgExecute := NewMsgFactoryExecuteAssets(sendAssetsExecute)
	require.Equal(t, testingFunc(ctx, msgExecute).Tags, sdk.Tags(nil))
	
	msgWrong := sdk.NewTestMsg(addrs[0])
	require.Equal(t, testingFunc(ctx, msgWrong).Tags, sdk.Tags(nil))
	
	mapper.SetAssetPeg(ctx, &oneAssetPeg)
	require.Equal(t, string(testingFunc(ctx, msgIssue).Tags.ToKVPairs()[0].Key), "recepient")
	mapper.SetAssetPeg(ctx, &oneAssetPegSend)
	require.Equal(t, string(testingFunc(ctx, msgSend).Tags.ToKVPairs()[0].Key), "recepient")
	mapper.SetAssetPeg(ctx, &oneAssetPegExecute)
	require.Equal(t, string(testingFunc(ctx, msgExecute).Tags.ToKVPairs()[0].Key), "recepient")
	
}

func TestRedeemAsset(t *testing.T) {
	_, ctx, mapper, keeper := initiateSetupMultiStore()
	
	fun := NewHandler(keeper)
	
	var testCase = []TestCase{
		{addrs[0], addrs[1], sdk.PegHash([]byte("test")), "ABC123", "assetType", 5, "quantityUnit", addrs[0], true},
	}
	
	assetPeg, issueAsset := genIssueAsset(testCase[0])
	assetPeg2 := sdk.BaseAssetPeg{
		PegHash:       sdk.PegHash([]byte("gta")),
		AssetType:     "sona",
		AssetQuantity: 10,
		OwnerAddress:  addrs[0],
	}
	mapper.SetAssetPeg(ctx, &assetPeg)
	
	issueAssets := []IssueAsset{issueAsset}
	msgIssueAssets := NewMsgFactoryIssueAssets(issueAssets)
	
	res := fun(ctx, msgIssueAssets)
	
	redeemAsset := []RedeemAsset{
		{
			OwnerAddress: addrs[1],
			ToAddress:    addrs[0],
			PegHash:      assetPeg.GetPegHash(),
		},
		{
			OwnerAddress: addrs[1],
			ToAddress:    addrs[0],
			PegHash:      assetPeg2.GetPegHash(),
		},
	}
	
	var redeemAssets = [][]RedeemAsset{{redeemAsset[0]}, {redeemAsset[1]}}
	msgAssetRedeemAsset := NewMsgFactoryRedeemAssets(redeemAssets[0])
	res = fun(ctx, msgAssetRedeemAsset)
	require.Equal(t, addrs[0].String(), string(res.Tags[0].Value))
	
	msgAssetRedeemAsset = NewMsgFactoryRedeemAssets(redeemAssets[1])
	res = fun(ctx, msgAssetRedeemAsset)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
	mapper.SetAssetPeg(ctx, &assetPeg2)
	res = fun(ctx, msgAssetRedeemAsset)
	require.Equal(t, sdk.Tags(nil), res.Tags)
	
}
