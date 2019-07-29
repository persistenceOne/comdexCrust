package assetFactory

import (
	"testing"

	sdk "github.com/commitHub/commitBlockchain/types"

	"github.com/stretchr/testify/require"
)

var issuerAddr = sdk.AccAddress([]byte("issuerAddr"))
var toAddr = sdk.AccAddress([]byte("toAddr"))
var assetPeg = sdk.BaseAssetPeg{
	PegHash:       sdk.PegHash([]byte("pegHash")),
	DocumentHash:  "ABC123",
	AssetType:     "assetType",
	AssetQuantity: 5,
	QuantityUnit:  "quantityUnit",
	OwnerAddress:  sdk.AccAddress([]byte("ownerAddr")),
}
var issueAsset = IssueAsset{
	IssuerAddress: issuerAddr,
	ToAddress:     toAddr,
	AssetPeg:      &assetPeg,
}

var addrs = []sdk.AccAddress{
	sdk.AccAddress([]byte("issuer")),
	sdk.AccAddress([]byte("from")),
	sdk.AccAddress([]byte("to")),
	sdk.AccAddress([]byte("")),
	sdk.AccAddress(nil),
	sdk.AccAddress([]byte("efuhrgjkklirjgnlekgopkelfjgrselkmhijngekfjigorjnsekijgrjnekijorseknijorselfpojigobijkafnijgrjneafkmgrs")),
}

type TestCase struct {
	from           sdk.AccAddress
	to             sdk.AccAddress
	pegHash        sdk.PegHash
	documentHash   string
	assetType      string
	assetQuantity  int64
	quantityUnit   string
	ownerAddress   sdk.AccAddress
	expectedResult bool
}

var listTests = []TestCase{
	{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "documentHash", "assetType", 5, "quantityUnit", addrs[2], true},
	{addrs[0], addrs[1], sdk.PegHash([]byte("")), "documentHash", "assetType", 5, "quantityUnit", addrs[2], true},
	{addrs[1], addrs[2], sdk.PegHash([]byte("peghash")), "documentHash", "assetType", 558463, "quantityUnit", addrs[3], true},
	{addrs[2], addrs[1], sdk.PegHash([]byte("")), "", "", 0, "", addrs[2], true},
	{addrs[3], addrs[4], sdk.PegHash([]byte("pegHash")), "", "", -34, "", addrs[3], true},
}

func genIssueAsset(testCase TestCase) (sdk.BaseAssetPeg, IssueAsset) {
	oneAssetPeg := sdk.BaseAssetPeg{
		PegHash:       testCase.pegHash,
		DocumentHash:  testCase.documentHash,
		AssetType:     testCase.assetType,
		AssetQuantity: testCase.assetQuantity,
		QuantityUnit:  testCase.quantityUnit,
		OwnerAddress:  testCase.ownerAddress,
	}
	oneIssueAsset := IssueAsset{
		IssuerAddress: testCase.from,
		ToAddress:     testCase.to,
		AssetPeg:      &oneAssetPeg,
	}
	return oneAssetPeg, oneIssueAsset
}

//--------------------Issye Asset
func TestNewIssueAsset(t *testing.T) {
	for _, testCase := range listTests {
		oneAssetPeg, oneIssueAsset := genIssueAsset(testCase)
		testIssueAsset := NewIssueAsset(testCase.from, testCase.to, &oneAssetPeg)
		if testCase.expectedResult {
			require.Equal(t, oneIssueAsset, testIssueAsset)
		} else {
			require.NotEqual(t, oneIssueAsset, testIssueAsset)
		}
	}
}
func TestIssueAssetGetSignBytes(t *testing.T) {
	for _, testCase := range listTests {
		_, oneIssueAsset := genIssueAsset(testCase)
		if testCase.expectedResult {
			require.NotNil(t, oneIssueAsset.GetSignBytes())
		} else {
			require.Nil(t, oneIssueAsset.GetSignBytes())
		}
	}

}

func TestIssueAssetValidateBasics(t *testing.T) {
	var testsValidatIssues = []TestCase{
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantityUnit", addrs[2], true},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "dfdc", "assetType", 5, "quantityUnit", addrs[2], false},
		{addrs[3], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantityUnit", addrs[2], false},
		{addrs[0], addrs[3], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantityUnit", addrs[3], false},
		{addrs[0], addrs[0], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantityUnit", addrs[0], true},
		{addrs[0], addrs[1], sdk.PegHash([]byte("")), "ABC234", "assetType", 5, "quantityUnit", addrs[2], true},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "", "assetType", 5, "quantityUnit", addrs[2], false},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "", 5, "quantityUnit", addrs[2], false},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantity Type", addrs[2], false},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 0, "", addrs[2], false},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", -7, "quantityUnit", addrs[2], false},
	}
	for _, testCase := range testsValidatIssues {
		_, oneIssueAsset := genIssueAsset(testCase)
		if testCase.expectedResult {
			require.Nil(t, MsgFactoryIssueAssets{[]IssueAsset{oneIssueAsset}}.ValidateBasic())
		} else {
			require.NotNil(t, MsgFactoryIssueAssets{[]IssueAsset{oneIssueAsset}}.ValidateBasic())
		}
	}
}

//-------------------Msg Factory Issue Assets

func TestNewMsgFactoryIssueAssets(t *testing.T) {
	for _, testCase := range listTests {
		_, oneIssueAsset := genIssueAsset(testCase)
		var issueAssets = []IssueAsset{issueAsset, oneIssueAsset}
		var msgFactoryIssueAssets = MsgFactoryIssueAssets{issueAssets}
		newMsgFactoryIssueAssets := NewMsgFactoryIssueAssets(issueAssets)
		if testCase.expectedResult {
			require.Equal(t, newMsgFactoryIssueAssets, msgFactoryIssueAssets)
		} else {
			require.NotEqual(t, newMsgFactoryIssueAssets, msgFactoryIssueAssets)
		}
	}

}

func TestMsgFactoryIssueAssetsType(t *testing.T) {
	for _, testCase := range listTests {
		_, oneIssueAsset := genIssueAsset(testCase)
		var issueAssets = []IssueAsset{issueAsset, oneIssueAsset}
		var msgFactoryIssueAssets = MsgFactoryIssueAssets{issueAssets}
		if testCase.expectedResult {
			require.Equal(t, msgFactoryIssueAssets.Type(), "assetFactory")
		} else {
			require.NotEqual(t, msgFactoryIssueAssets.Type(), "assetFactory")
		}
	}
}

func TestMsgFactoryIssueAssetsValidateBasic(t *testing.T) {
	var testsValidatIssues = []TestCase{
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantityUnit", addrs[2], true},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "dfdc", "assetType", 5, "quantityUnit", addrs[2], false},
		{addrs[3], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantityUnit", addrs[2], false},
		{addrs[0], addrs[3], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantityUnit", addrs[3], false},
		{addrs[0], addrs[0], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantityUnit", addrs[0], true},
		{addrs[0], addrs[1], sdk.PegHash([]byte("")), "ABC234", "assetType", 5, "quantityUnit", addrs[2], true},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "", "assetType", 5, "quantityUnit", addrs[2], false},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "", 5, "quantityUnit", addrs[2], false},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 5, "quantity Type", addrs[2], false},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", 0, "", addrs[2], false},
		{addrs[0], addrs[1], sdk.PegHash([]byte("peghash")), "ABC234", "assetType", -7, "quantityUnit", addrs[2], false},
	}
	for _, testCase := range testsValidatIssues {
		_, oneIssueAsset := genIssueAsset(testCase)
		var issueAssets = []IssueAsset{issueAsset, oneIssueAsset}
		var msgFactoryIssueAssets = MsgFactoryIssueAssets{issueAssets}
		if testCase.expectedResult {
			require.Nil(t, msgFactoryIssueAssets.ValidateBasic())
		} else {
			require.NotNil(t, msgFactoryIssueAssets.ValidateBasic())
		}
	}
}
func TestMsgFactoryIssueAssetsGetSignBytes(t *testing.T) {
	for _, testCase := range listTests {
		_, oneIssueAsset := genIssueAsset(testCase)
		var issueAssets = []IssueAsset{issueAsset, oneIssueAsset}
		var msgFactoryIssueAssets = MsgFactoryIssueAssets{issueAssets}
		if testCase.expectedResult {
			require.NotNil(t, msgFactoryIssueAssets.GetSignBytes())
		} else {
			require.Nil(t, msgFactoryIssueAssets.GetSignBytes())
		}
	}
}

func TestMsgFactoryIssueAssetsGetSigners(t *testing.T) {
	for _, testCase := range listTests {
		_, oneIssueAsset := genIssueAsset(testCase)
		var issueAssets = []IssueAsset{issueAsset, oneIssueAsset}
		var msgFactoryIssueAssets = MsgFactoryIssueAssets{issueAssets}
		var issuers = []sdk.AccAddress{issuerAddr, oneIssueAsset.IssuerAddress}

		if testCase.expectedResult {
			require.Equal(t, msgFactoryIssueAssets.GetSigners(), issuers)
		} else {
			require.NotEqual(t, msgFactoryIssueAssets.GetSigners(), issuers)
		}
	}
}

func TestBuildIssueAssetMsg(t *testing.T) {
	for _, testCase := range listTests {
		oneAssetPeg, oneIssueAsset := genIssueAsset(testCase)
		var issueAssets = []IssueAsset{oneIssueAsset}
		var msgFactoryIssueAssets = MsgFactoryIssueAssets{issueAssets}
		if testCase.expectedResult {
			require.Equal(t, msgFactoryIssueAssets, BuildIssueAssetMsg(testCase.from, testCase.to, &oneAssetPeg))
		} else {
			require.NotEqual(t, msgFactoryIssueAssets, BuildIssueAssetMsg(testCase.from, testCase.to, &oneAssetPeg))
		}
	}
}

//------------------Send Asset
type TestCaseSend struct {
	relayer        sdk.AccAddress
	from           sdk.AccAddress
	to             sdk.AccAddress
	pegHash        sdk.PegHash
	expectedResult bool
}

var listSendTests = []TestCaseSend{
	{addrs[0], addrs[1], addrs[2], sdk.PegHash([]byte("pegHash")), true},
	{addrs[0], addrs[1], addrs[2], sdk.PegHash([]byte("")), true},
	{addrs[1], addrs[2], addrs[2], sdk.PegHash([]byte("pegHash")), true},
	{addrs[0], addrs[1], addrs[3], sdk.PegHash([]byte("pegHash")), true},
}

func genSendAsset(testCase TestCaseSend) SendAsset {
	return SendAsset{
		RelayerAddress: testCase.relayer,
		FromAddress:    testCase.from,
		ToAddress:      testCase.to,
		PegHash:        testCase.pegHash,
	}
}

func TestNewSendAsset(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		if testCase.expectedResult {
			require.Equal(t, oneSendAsset, NewSendAsset(testCase.relayer, testCase.from, testCase.to, testCase.pegHash))
		} else {
			require.NotEqual(t, oneSendAsset, NewSendAsset(testCase.relayer, testCase.from, testCase.to, testCase.pegHash))
		}
	}
}

func TestSendAssetGetSignBytes(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		if testCase.expectedResult {
			require.NotNil(t, oneSendAsset.GetSignBytes())
		} else {
			require.Nil(t, oneSendAsset.GetSignBytes())
		}
	}
}

//---------------Msg Factory Send Assets

func TestNewMsgFactorySendAssets(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactorySendAssets := MsgFactorySendAssets{sendAssets}
		if testCase.expectedResult {
			require.Equal(t, msgFactorySendAssets, NewMsgFactorySendAssets(sendAssets))
		} else {
			require.NotEqual(t, msgFactorySendAssets, NewMsgFactorySendAssets(sendAssets))
		}
	}
}
func TestMsgFactorySendAssetsType(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactorySendAssets := MsgFactorySendAssets{sendAssets}
		if testCase.expectedResult {
			require.Equal(t, msgFactorySendAssets.Type(), "assetFactory")
		} else {
			require.NotEqual(t, msgFactorySendAssets.Type(), "assetFactory")
		}
	}
}

func TestMsgFactorySendAssetsValidateBasic(t *testing.T) {
	var listSendTestsValidate = []TestCaseSend{
		{addrs[0], addrs[1], addrs[2], sdk.PegHash([]byte("pegHash")), true},
		{addrs[0], addrs[1], addrs[2], sdk.PegHash([]byte("")), true},
		{addrs[1], addrs[2], addrs[4], sdk.PegHash([]byte("pegHash")), false},
		{addrs[0], addrs[1], addrs[3], sdk.PegHash([]byte("pegHash")), false},
	}
	for _, testCase := range listSendTestsValidate {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactorySendAssets := MsgFactorySendAssets{sendAssets}
		if testCase.expectedResult {
			require.Nil(t, msgFactorySendAssets.ValidateBasic())
		} else {
			require.NotNil(t, msgFactorySendAssets.ValidateBasic())
		}
	}
}

func TestMsgFactorySendAssetsGetSignBytes(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactorySendAssets := MsgFactorySendAssets{sendAssets}
		if testCase.expectedResult {
			require.NotNil(t, msgFactorySendAssets.GetSignBytes())
		} else {
			require.Nil(t, msgFactorySendAssets.GetSignBytes())
		}
	}
}

func TestMsgFactorySendAssetsGetSigners(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactorySendAssets := MsgFactorySendAssets{sendAssets}
		signers := []sdk.AccAddress{testCase.relayer}
		if testCase.expectedResult {
			require.Equal(t, msgFactorySendAssets.GetSigners(), signers)
		} else {
			require.NotEqual(t, msgFactorySendAssets.GetSigners(), signers)
		}
	}
}
func TestBuildSendAssetMsg(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactorySendAssets := MsgFactorySendAssets{sendAssets}
		buildedMsg := BuildSendAssetMsg(testCase.relayer, testCase.from, testCase.to, testCase.pegHash)
		if testCase.expectedResult {
			require.Equal(t, buildedMsg, msgFactorySendAssets)
		} else {
			require.NotEqual(t, buildedMsg, msgFactorySendAssets)
		}
	}
}

//--------------------Msg Factory Execute Assets

func TestNewMsgFactoryExecuteAssets(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactoryExecuteAssets := MsgFactoryExecuteAssets{sendAssets}
		toTest := NewMsgFactoryExecuteAssets(sendAssets)
		if testCase.expectedResult {
			require.Equal(t, toTest, msgFactoryExecuteAssets)
		} else {
			require.NotEqual(t, toTest, msgFactoryExecuteAssets)
		}
	}
}

func TestMsgFactoryExecuteAssetsType(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactoryExecuteAssets := MsgFactoryExecuteAssets{sendAssets}
		toTest := "assetFactory"
		if testCase.expectedResult {
			require.Equal(t, toTest, msgFactoryExecuteAssets.Type())
		} else {
			require.NotEqual(t, toTest, msgFactoryExecuteAssets.Type())
		}
	}
}

func TestMsgFactoryExecuteAssetsValidateBasic(t *testing.T) {
	var listSendTestsValidate = []TestCaseSend{
		{addrs[0], addrs[1], addrs[2], sdk.PegHash([]byte("pegHash")), true},
		{addrs[0], addrs[1], addrs[2], sdk.PegHash([]byte("")), true},
		{addrs[1], addrs[2], addrs[4], sdk.PegHash([]byte("pegHash")), false},
		{addrs[0], addrs[1], addrs[3], sdk.PegHash([]byte("pegHash")), false},
	}
	for _, testCase := range listSendTestsValidate {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactoryExecuteAssets := MsgFactoryExecuteAssets{sendAssets}
		if testCase.expectedResult {
			require.Nil(t, msgFactoryExecuteAssets.ValidateBasic())
		} else {
			require.NotNil(t, msgFactoryExecuteAssets.ValidateBasic())
		}
	}
}

func TestMsgFactoryExecuteAssetsGetSignBytes(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactoryExecuteAssets := MsgFactoryExecuteAssets{sendAssets}
		if testCase.expectedResult {
			require.NotNil(t, msgFactoryExecuteAssets.GetSignBytes())
		} else {
			require.Nil(t, msgFactoryExecuteAssets.GetSignBytes())
		}
	}
}

func TestMsgFactoryExecuteAssetsGetSigners(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactoryExecuteAssets := MsgFactoryExecuteAssets{sendAssets}
		signers := []sdk.AccAddress{testCase.relayer}
		if testCase.expectedResult {
			require.Equal(t, msgFactoryExecuteAssets.GetSigners(), signers)
		} else {
			require.NotEqual(t, msgFactoryExecuteAssets.GetSigners(), signers)
		}
	}
}

func TestBuildExecuteAssetMsg(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendAsset := genSendAsset(testCase)
		sendAssets := []SendAsset{oneSendAsset}
		msgFactoryExecuteAssets := MsgFactoryExecuteAssets{sendAssets}
		toTest := BuildExecuteAssetMsg(testCase.relayer, testCase.from, testCase.to, testCase.pegHash)
		if testCase.expectedResult {
			require.Equal(t, toTest, msgFactoryExecuteAssets)
		} else {
			require.NotEqual(t, toTest, msgFactoryExecuteAssets)
		}
	}
}
