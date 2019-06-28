package fiatFactory

import (
	"testing"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/stretchr/testify/require"
)

var issuerAddress = sdk.AccAddress([]byte("issuerAddress"))
var toAddress = sdk.AccAddress([]byte("toAddress"))
var fiatPeg = sdk.BaseFiatPeg{
	PegHash:           sdk.PegHash([]byte("pegHash")),
	TransactionID:     "ABC123",
	TransactionAmount: 100,
	RedeemedAmount:    50,
	Owners: []sdk.Owner{
		{
			OwnerAddress: sdk.AccAddress([]byte("issuer")),
			Amount:       2000,
		},
	},
}

var TestIssueFiat = IssueFiat{
	IssuerAddress: issuerAddress,
	ToAddress:     toAddress,
	FiatPeg:       &fiatPeg,
}

var addresses = []sdk.AccAddress{
	sdk.AccAddress([]byte("issuer")),
	sdk.AccAddress([]byte("to2")),
	sdk.AccAddress([]byte("from")),
	sdk.AccAddress([]byte("")),
	sdk.AccAddress(nil),
	sdk.AccAddress([]byte("efuhrgjkklirjgnlekgopkelfjgrselkmhijngekfjigorjnsekijgrjnekijorseknijorselfpojigobijkafnijgrjneafkmgrs")),
}

type TestCase struct {
	from              sdk.AccAddress
	to                sdk.AccAddress
	pegHash           sdk.PegHash
	transactionID     string
	transactionAmount int64
	redeemedAmount    int64
	owners            []sdk.Owner
	expectedResult    bool
}

var listTests = []TestCase{
	{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("peghash")), transactionID: "documentHash", transactionAmount: 10, redeemedAmount: 5, owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}, expectedResult: true},
	{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("")), transactionID: "documentHash", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}, expectedResult: true},
	{from: addresses[1], to: addresses[2], pegHash: sdk.PegHash([]byte("peghash")), transactionID: "documentHash", transactionAmount: 20, redeemedAmount: 4, owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}, expectedResult: true},
	{from: addresses[2], to: addresses[1], pegHash: sdk.PegHash([]byte("")), transactionID: "", transactionAmount: 30, redeemedAmount: 5, owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}, expectedResult: true},
	{from: addresses[3], to: addresses[3], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "", transactionAmount: 4, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}, expectedResult: true},
}

func genIssueFiat(testCase TestCase) (sdk.BaseFiatPeg, IssueFiat) {
	oneFiatPeg := sdk.BaseFiatPeg{
		PegHash:           testCase.pegHash,
		TransactionID:     testCase.transactionID,
		TransactionAmount: testCase.transactionAmount,
		RedeemedAmount:    testCase.redeemedAmount,
		Owners:            testCase.owners,
	}
	oneIssueFiat := IssueFiat{
		IssuerAddress: testCase.from,
		ToAddress:     testCase.to,
		FiatPeg:       &oneFiatPeg,
	}
	return oneFiatPeg, oneIssueFiat
}

// Issue Fiat Test Cases
func TestNewIssueFiat(t *testing.T) {
	for _, testCase := range listTests {
		oneFiatPeg, oneIssueFiat := genIssueFiat(testCase)
		testIssueFiat := NewIssueFiat(testCase.from, testCase.to, &oneFiatPeg)
		if testCase.expectedResult {
			require.Equal(t, oneIssueFiat, testIssueFiat)
		} else {
			require.NotEqual(t, oneIssueFiat, testIssueFiat)
		}
	}
}

// TestIssueFiatGetSignBytes
func TestIssueFiatGetSignBytes(t *testing.T) {
	for _, testCase := range listTests {
		_, oneIssueFiat := genIssueFiat(testCase)
		
		if testCase.expectedResult {
			require.NotNil(t, oneIssueFiat.GetSignBytes())
		} else {
			require.Nil(t, oneIssueFiat.GetSignBytes())
		}
	}
}
func TestIssueFiatValidateBasics(t *testing.T) {
	var testValidateIssueFiats = []TestCase{
		
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[2], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[0], pegHash: sdk.PegHash([]byte("pegHashpegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[0], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[2], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCD123", transactionAmount: 10, redeemedAmount: 20, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCD123", transactionAmount: 10, redeemedAmount: -10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: false},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCD123", transactionAmount: 0, redeemedAmount: 0, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: false},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCD123", transactionAmount: -5, redeemedAmount: 0, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: false},
	}
	for _, testCase := range testValidateIssueFiats {
		_, oneIssueFiat := genIssueFiat(testCase)
		if testCase.expectedResult {
			require.Nil(t, MsgFactoryIssueFiats{[]IssueFiat{oneIssueFiat}}.ValidateBasic())
		} else {
			require.NotNil(t, MsgFactoryIssueFiats{[]IssueFiat{oneIssueFiat}}.ValidateBasic())
		}
	}
}

// ------------------------------------------
func TestNewMsgFactoryIssueFiats(t *testing.T) {
	
	for _, testCase := range listTests {
		oneFiatPeg, _ := genIssueFiat(testCase)
		issueFiats := []IssueFiat{NewIssueFiat(testCase.from, testCase.to, &oneFiatPeg)}
		var msg = MsgFactoryIssueFiats{issueFiats}
		newMsgFactoryIssueFiats := NewMsgFactoryIssueFiats(issueFiats)
		if testCase.expectedResult {
			require.Equal(t, msg, newMsgFactoryIssueFiats)
		} else {
			require.NotEqual(t, msg, newMsgFactoryIssueFiats)
		}
	}
	
}

func TestMsgFactoryIssueFiatsType(t *testing.T) {
	
	for _, testCase := range listTests {
		oneFiatPeg, _ := genIssueFiat(testCase)
		issueFiats := []IssueFiat{NewIssueFiat(testCase.from, testCase.to, &oneFiatPeg)}
		var msg = MsgFactoryIssueFiats{issueFiats}
		
		if testCase.expectedResult {
			require.Equal(t, "fiatFactory", msg.Type())
		} else {
			require.NotEqual(t, "fiatFactory", msg.Type())
		}
	}
}

func TestMsgFactoryIssueFiatsValidateBasic(t *testing.T) {
	var testValidateIssueFiats = []TestCase{
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[2], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[0], pegHash: sdk.PegHash([]byte("pegHashpegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[0], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[2], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("")), transactionID: "ABCDEF12345", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCD123", transactionAmount: 10, redeemedAmount: 20, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCD123", transactionAmount: 10, redeemedAmount: -10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: true},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "", transactionAmount: 100, redeemedAmount: 10, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: false},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCD123", transactionAmount: 0, redeemedAmount: 0, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: false},
		{from: addresses[0], to: addresses[1], pegHash: sdk.PegHash([]byte("pegHash")), transactionID: "ABCD123", transactionAmount: -5, redeemedAmount: 0, owners: []sdk.Owner{{OwnerAddress: addresses[0], Amount: 12}}, expectedResult: false},
	}
	for _, testCase := range testValidateIssueFiats {
		oneFiatPeg, _ := genIssueFiat(testCase)
		issueFiats := []IssueFiat{NewIssueFiat(testCase.from, testCase.to, &oneFiatPeg)}
		var msg = MsgFactoryIssueFiats{issueFiats}
		
		if testCase.expectedResult {
			require.Nil(t, msg.ValidateBasic())
		} else {
			require.NotNil(t, msg.ValidateBasic())
		}
	}
}
func TestMsgFactoryIssueFiatsGetSignBytes(t *testing.T) {
	issuerAddress := sdk.AccAddress([]byte("issuer"))
	toAddress := sdk.AccAddress([]byte("to"))
	fiatPeg := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("pegHash")))
	issuerAddress1 := sdk.AccAddress([]byte("issuer1"))
	toAddress1 := sdk.AccAddress([]byte("to1"))
	fiatPeg1 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("pegHash1")))
	var msg = MsgFactoryIssueFiats{
		IssueFiats: []IssueFiat{
			NewIssueFiat(issuerAddress, toAddress, &fiatPeg),
			NewIssueFiat(issuerAddress1, toAddress1, &fiatPeg1),
		},
	}
	res := msg.GetSignBytes()
	
	expected := `{"issueFiats":[{"issuerAddress":"cosmos1d9ehxat9wgjjln07","toAddress":"cosmos1w3hsjttrfq","fiatPeg":{"type":"comdex-blockchain/FiatPeg","value":{"pegHash":"70656748617368","transactionID":"","transactionAmount":"0","redeemedAmount":"0","owners":null}}},{"issuerAddress":"cosmos1d9ehxat9wgcs3mnw0e","toAddress":"cosmos1w3hnz0y6kfc","fiatPeg":{"type":"comdex-blockchain/FiatPeg","value":{"pegHash":"7065674861736831","transactionID":"","transactionAmount":"0","redeemedAmount":"0","owners":null}}}]}`
	require.Equal(t, expected, string(res))
}

func TestMsgFactoryIssueFiatsGetSigners(t *testing.T) {
	for _, testCase := range listTests {
		_, oneIssueFiat := genIssueFiat(testCase)
		var issueFiats = []IssueFiat{TestIssueFiat, oneIssueFiat}
		var msgFactoryIssueFiats = MsgFactoryIssueFiats{issueFiats}
		var issuers = []sdk.AccAddress{issuerAddress, oneIssueFiat.IssuerAddress}
		
		if testCase.expectedResult {
			require.Equal(t, msgFactoryIssueFiats.GetSigners(), issuers)
		} else {
			require.NotEqual(t, msgFactoryIssueFiats.GetSigners(), issuers)
		}
	}
}

func TestBuildIssueFiatMsg(t *testing.T) {
	
	for _, testCase := range listTests {
		oneFiatPeg, oneIssueFiat := genIssueFiat(testCase)
		var issueFiats = []IssueFiat{oneIssueFiat}
		var msgFactoryIssueFiats = MsgFactoryIssueFiats{issueFiats}
		if testCase.expectedResult {
			require.Equal(t, msgFactoryIssueFiats, BuildIssueFiatMsg(testCase.from, testCase.to, &oneFiatPeg))
		} else {
			require.NotEqual(t, msgFactoryIssueFiats, BuildIssueFiatMsg(testCase.from, testCase.to, &oneFiatPeg))
		}
	}
}

// -------------------------------------REDEEM FIAT-------------------------------------

type TestCaseRedeem struct {
	relayer        sdk.AccAddress
	redeemer       sdk.AccAddress
	amount         int64
	fiatPegWallet  sdk.FiatPegWallet
	expectedResult bool
}

var listRedeemTests = []TestCaseRedeem{
	{relayer: addresses[0], redeemer: addresses[1], amount: 200, fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
	{relayer: addresses[0], redeemer: addresses[1], amount: 0, fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
	{relayer: addresses[1], redeemer: addresses[2], amount: 200, fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
	{relayer: addresses[0], redeemer: addresses[1], amount: 200, fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
}

func genRedeemFiat(testCase TestCaseRedeem) RedeemFiat {
	return RedeemFiat{
		RelayerAddress:  testCase.relayer,
		RedeemerAddress: testCase.redeemer,
		Amount:          testCase.amount,
		FiatPegWallet:   testCase.fiatPegWallet,
	}
}

func TestNewRedeemFiat(t *testing.T) {
	for _, testCase := range listRedeemTests {
		oneRedeemFiat := genRedeemFiat(testCase)
		if testCase.expectedResult {
			require.Equal(t, oneRedeemFiat, NewRedeemFiat(testCase.relayer, testCase.redeemer, testCase.amount, testCase.fiatPegWallet))
		} else {
			require.NotEqual(t, oneRedeemFiat, NewRedeemFiat(testCase.relayer, testCase.redeemer, testCase.amount, testCase.fiatPegWallet))
		}
	}
}

func TestRedeemFiatGetSignBytes(t *testing.T) {
	for _, testCase := range listRedeemTests {
		oneRedeemFiat := genRedeemFiat(testCase)
		if testCase.expectedResult {
			require.NotNil(t, oneRedeemFiat.GetSignBytes())
		} else {
			require.Nil(t, oneRedeemFiat.GetSignBytes())
		}
	}
}

// ---------------Msg Factory Redeem Fiats

func TestNewMsgFactoryRedeemFiats(t *testing.T) {
	for _, testCase := range listRedeemTests {
		oneRedeemFiat := genRedeemFiat(testCase)
		redeemFiats := []RedeemFiat{oneRedeemFiat}
		msgFactoryRedeemFiats := MsgFactoryRedeemFiats{redeemFiats}
		if testCase.expectedResult {
			require.Equal(t, msgFactoryRedeemFiats, NewMsgFactoryRedeemFiats(redeemFiats))
		} else {
			require.NotEqual(t, msgFactoryRedeemFiats, NewMsgFactoryRedeemFiats(redeemFiats))
		}
	}
}

func TestMsgFactoryRedeemFiatsType(t *testing.T) {
	for _, testCase := range listRedeemTests {
		oneRedeemFiat := genRedeemFiat(testCase)
		redeemFiats := []RedeemFiat{oneRedeemFiat}
		msgFactoryRedeemFiats := MsgFactoryRedeemFiats{redeemFiats}
		if testCase.expectedResult {
			require.Equal(t, msgFactoryRedeemFiats.Type(), "fiatFactory")
		} else {
			require.NotEqual(t, msgFactoryRedeemFiats.Type(), "fiatFactory")
		}
	}
}

// -------------------------------------SEND FIAT-------------------------------------

type TestCaseSend struct {
	relayer        sdk.AccAddress
	from           sdk.AccAddress
	to             sdk.AccAddress
	pegHash        sdk.PegHash
	fiatPegWallet  sdk.FiatPegWallet
	expectedResult bool
}

var listSendTests = []TestCaseSend{
	{relayer: addresses[0], from: addresses[1], to: addresses[2], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
	{relayer: addresses[0], from: addresses[1], to: addresses[2], pegHash: sdk.PegHash([]byte("")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
	{relayer: addresses[1], from: addresses[2], to: addresses[2], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
	{relayer: addresses[0], from: addresses[1], to: addresses[3], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
}

func genSendFiat(testCase TestCaseSend) SendFiat {
	return SendFiat{
		RelayerAddress: testCase.relayer,
		FromAddress:    testCase.from,
		ToAddress:      testCase.to,
		PegHash:        testCase.pegHash,
		FiatPegWallet:  testCase.fiatPegWallet,
	}
}

func TestNewSendFiat(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		if testCase.expectedResult {
			require.Equal(t, oneSendFiat, NewSendFiat(testCase.relayer, testCase.from, testCase.to, testCase.pegHash, testCase.fiatPegWallet))
		} else {
			require.NotEqual(t, oneSendFiat, NewSendFiat(testCase.relayer, testCase.from, testCase.to, testCase.pegHash, testCase.fiatPegWallet))
		}
	}
}

func TestSendFiatGetSignBytes(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		if testCase.expectedResult {
			require.NotNil(t, oneSendFiat.GetSignBytes())
		} else {
			require.Nil(t, oneSendFiat.GetSignBytes())
		}
	}
}

func TestMsgFactoryRedeemFiatsValidateBasic(t *testing.T) {
	var listRedeemTestsValidate = []TestCaseRedeem{
		{relayer: addresses[0], redeemer: addresses[1], amount: 20, fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 1000, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 1000}}}}, expectedResult: true},
		{relayer: addresses[0], redeemer: addresses[1], amount: 40, fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 0, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
		{relayer: addresses[1], redeemer: addresses[2], amount: 50, fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
		{relayer: addresses[0], redeemer: addresses[1], amount: 90, fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
	}
	for _, testCase := range listRedeemTestsValidate {
		oneRedeemFiat := genRedeemFiat(testCase)
		redeemFiats := []RedeemFiat{oneRedeemFiat}
		msgFactoryRedeemFiats := MsgFactoryRedeemFiats{redeemFiats}
		if testCase.expectedResult {
			require.Nil(t, msgFactoryRedeemFiats.ValidateBasic())
		} else {
			require.NotNil(t, msgFactoryRedeemFiats.ValidateBasic())
		}
	}
}

func TestMsgFactoryRedeemFiatsGetSignBytes(t *testing.T) {
	for _, testCase := range listRedeemTests {
		oneRedeemFiat := genRedeemFiat(testCase)
		redeemFiats := []RedeemFiat{oneRedeemFiat}
		msgFactoryRedeemFiats := MsgFactoryRedeemFiats{redeemFiats}
		if testCase.expectedResult {
			require.NotNil(t, msgFactoryRedeemFiats.GetSignBytes())
		} else {
			require.Nil(t, msgFactoryRedeemFiats.GetSignBytes())
		}
	}
}

func TestBuildRedeemFiatMsg(t *testing.T) {
	for _, testCase := range listRedeemTests {
		oneRedeemFiat := genRedeemFiat(testCase)
		redeemFiats := []RedeemFiat{oneRedeemFiat}
		msgFactoryRedeemFiats := MsgFactoryRedeemFiats{redeemFiats}
		buildedMsg := BuildRedeemFiatMsg(testCase.relayer, testCase.redeemer, testCase.amount, testCase.fiatPegWallet)
		if testCase.expectedResult {
			require.Equal(t, buildedMsg, msgFactoryRedeemFiats)
		} else {
			require.NotEqual(t, buildedMsg, msgFactoryRedeemFiats)
		}
	}
}

// ---------------Msg Factory Send Fiats

func TestNewMsgFactorySendFiats(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactorySendFiats := MsgFactorySendFiats{sendFiats}
		if testCase.expectedResult {
			require.Equal(t, msgFactorySendFiats, NewMsgFactorySendFiats(sendFiats))
		} else {
			require.NotEqual(t, msgFactorySendFiats, NewMsgFactorySendFiats(sendFiats))
		}
	}
}

func TestMsgFactorySendFiatsType(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactorySendFiats := MsgFactorySendFiats{sendFiats}
		if testCase.expectedResult {
			require.Equal(t, msgFactorySendFiats.Type(), "fiatFactory")
		} else {
			require.NotEqual(t, msgFactorySendFiats.Type(), "fiatFactory")
		}
	}
}

func TestMsgFactorySendFiatsValidateBasic(t *testing.T) {
	var listSendTestsValidate = []TestCaseSend{
		{relayer: addresses[0], from: addresses[1], to: addresses[2], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
		{relayer: addresses[0], from: addresses[1], to: addresses[2], pegHash: sdk.PegHash([]byte("")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
		{relayer: addresses[1], from: addresses[2], to: addresses[4], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: false},
		{relayer: addresses[0], from: addresses[1], to: addresses[3], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: false},
	}
	for _, testCase := range listSendTestsValidate {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactorySendFiats := MsgFactorySendFiats{sendFiats}
		if testCase.expectedResult {
			require.Nil(t, msgFactorySendFiats.ValidateBasic())
		} else {
			require.NotNil(t, msgFactorySendFiats.ValidateBasic())
		}
	}
}

func TestMsgFactorySendFiatsGetSignBytes(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactorySendFiats := MsgFactorySendFiats{sendFiats}
		if testCase.expectedResult {
			require.NotNil(t, msgFactorySendFiats.GetSignBytes())
		} else {
			require.Nil(t, msgFactorySendFiats.GetSignBytes())
		}
	}
}

func TestMsgFactorySendFiatsGetSigners(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendAssets := []SendFiat{oneSendFiat}
		msgFactorySendFiats := MsgFactorySendFiats{sendAssets}
		signers := []sdk.AccAddress{testCase.relayer}
		if testCase.expectedResult {
			require.Equal(t, msgFactorySendFiats.GetSigners(), signers)
		} else {
			require.NotEqual(t, msgFactorySendFiats.GetSigners(), signers)
		}
	}
}

func TestBuildSendFiatMsg(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactorySendFiats := MsgFactorySendFiats{sendFiats}
		buildedMsg := BuildSendFiatMsg(testCase.relayer, testCase.from, testCase.to, testCase.pegHash, testCase.fiatPegWallet)
		if testCase.expectedResult {
			require.Equal(t, buildedMsg, msgFactorySendFiats)
		} else {
			require.NotEqual(t, buildedMsg, msgFactorySendFiats)
		}
	}
}

// --------------------Msg Factory Execute Assets

func TestNewMsgFactoryExecuteFiats(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactoryExecuteFiats := MsgFactoryExecuteFiats{sendFiats}
		toTest := NewMsgFactoryExecuteFiats(sendFiats)
		if testCase.expectedResult {
			require.Equal(t, toTest, msgFactoryExecuteFiats)
		} else {
			require.NotEqual(t, toTest, msgFactoryExecuteFiats)
		}
	}
}

func TestMsgFactoryExecuteFiatsType(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactoryExecuteFiats := MsgFactoryExecuteFiats{sendFiats}
		toTest := "fiatFactory"
		if testCase.expectedResult {
			require.Equal(t, toTest, msgFactoryExecuteFiats.Type())
		} else {
			require.NotEqual(t, toTest, msgFactoryExecuteFiats.Type())
		}
	}
}

func TestMsgFactoryExecuteFiatsValidateBasic(t *testing.T) {
	var listSendTestsValidate = []TestCaseSend{
		{relayer: addresses[0], from: addresses[1], to: addresses[2], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
		{relayer: addresses[0], from: addresses[1], to: addresses[2], pegHash: sdk.PegHash([]byte("")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: true},
		{relayer: addresses[1], from: addresses[2], to: addresses[4], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: false},
		{relayer: addresses[0], from: addresses[1], to: addresses[3], pegHash: sdk.PegHash([]byte("pegHash")), fiatPegWallet: []sdk.BaseFiatPeg{{PegHash: sdk.PegHash([]byte("pegHash")), TransactionID: "ABCDEF12345", TransactionAmount: 100, RedeemedAmount: 10, Owners: []sdk.Owner{{OwnerAddress: addresses[1], Amount: 12}}}}, expectedResult: false},
	}
	for _, testCase := range listSendTestsValidate {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactoryExecuteFiats := MsgFactoryExecuteFiats{sendFiats}
		if testCase.expectedResult {
			require.Nil(t, msgFactoryExecuteFiats.ValidateBasic())
		} else {
			require.NotNil(t, msgFactoryExecuteFiats.ValidateBasic())
		}
	}
}

func TestMsgFactoryExecuteFiatsGetSignBytes(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactoryExecuteFiats := MsgFactoryExecuteFiats{sendFiats}
		if testCase.expectedResult {
			require.NotNil(t, msgFactoryExecuteFiats.GetSignBytes())
		} else {
			require.Nil(t, msgFactoryExecuteFiats.GetSignBytes())
		}
	}
}

func TestMsgFactoryExecuteFiatsGetSigners(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactoryExecuteFiats := MsgFactoryExecuteFiats{sendFiats}
		signers := []sdk.AccAddress{testCase.relayer}
		if testCase.expectedResult {
			require.Equal(t, msgFactoryExecuteFiats.GetSigners(), signers)
		} else {
			require.NotEqual(t, msgFactoryExecuteFiats.GetSigners(), signers)
		}
	}
}

func TestBuildExecuteAssetMsg(t *testing.T) {
	for _, testCase := range listSendTests {
		oneSendFiat := genSendFiat(testCase)
		sendFiats := []SendFiat{oneSendFiat}
		msgFactoryExecuteFiats := MsgFactoryExecuteFiats{sendFiats}
		toTest := BuildExecuteFiatMsg(testCase.relayer, testCase.from, testCase.to, testCase.pegHash, testCase.fiatPegWallet)
		if testCase.expectedResult {
			require.Equal(t, toTest, msgFactoryExecuteFiats)
		} else {
			require.NotEqual(t, toTest, msgFactoryExecuteFiats)
		}
	}
}
