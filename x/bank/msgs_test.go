package bank

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/commitHub/commitBlockchain/types"
)

func TestNewMsgSend(t *testing.T) {}

func TestMsgSendType(t *testing.T) {
	// Construct a MsgSend
	addr1 := sdk.AccAddress([]byte("input"))
	addr2 := sdk.AccAddress([]byte("output"))
	coins := sdk.Coins{sdk.NewInt64Coin("atom", 10)}
	var msg = MsgSend{
		Inputs:  []Input{NewInput(addr1, coins)},
		Outputs: []Output{NewOutput(addr2, coins)},
	}

	// TODO some failures for bad result
	require.Equal(t, msg.Type(), "bank")
}

func TestInputValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte{1, 2})
	addr2 := sdk.AccAddress([]byte{7, 8})
	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123)}
	multiCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 20)}

	var emptyAddr sdk.AccAddress
	emptyCoins := sdk.Coins{}
	emptyCoins2 := sdk.Coins{sdk.NewInt64Coin("eth", 0)}
	someEmptyCoins := sdk.Coins{sdk.NewInt64Coin("eth", 10), sdk.NewInt64Coin("atom", 0)}
	minusCoins := sdk.Coins{sdk.NewInt64Coin("eth", -34)}
	someMinusCoins := sdk.Coins{sdk.NewInt64Coin("atom", 20), sdk.NewInt64Coin("eth", -34)}
	unsortedCoins := sdk.Coins{sdk.NewInt64Coin("eth", 1), sdk.NewInt64Coin("atom", 1)}

	cases := []struct {
		valid bool
		txIn  Input
	}{
		// auth works with different apps
		{true, NewInput(addr1, someCoins)},
		{true, NewInput(addr2, someCoins)},
		{true, NewInput(addr2, multiCoins)},

		{false, NewInput(emptyAddr, someCoins)},  // empty address
		{false, NewInput(addr1, emptyCoins)},     // invalid coins
		{false, NewInput(addr1, emptyCoins2)},    // invalid coins
		{false, NewInput(addr1, someEmptyCoins)}, // invalid coins
		{false, NewInput(addr1, minusCoins)},     // negative coins
		{false, NewInput(addr1, someMinusCoins)}, // negative coins
		{false, NewInput(addr1, unsortedCoins)},  // unsorted coins
	}

	for i, tc := range cases {
		err := tc.txIn.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestOutputValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte{1, 2})
	addr2 := sdk.AccAddress([]byte{7, 8})
	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123)}
	multiCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 20)}

	var emptyAddr sdk.AccAddress
	emptyCoins := sdk.Coins{}
	emptyCoins2 := sdk.Coins{sdk.NewInt64Coin("eth", 0)}
	someEmptyCoins := sdk.Coins{sdk.NewInt64Coin("eth", 10), sdk.NewInt64Coin("atom", 0)}
	minusCoins := sdk.Coins{sdk.NewInt64Coin("eth", -34)}
	someMinusCoins := sdk.Coins{sdk.NewInt64Coin("atom", 20), sdk.NewInt64Coin("eth", -34)}
	unsortedCoins := sdk.Coins{sdk.NewInt64Coin("eth", 1), sdk.NewInt64Coin("atom", 1)}

	cases := []struct {
		valid bool
		txOut Output
	}{
		// auth works with different apps
		{true, NewOutput(addr1, someCoins)},
		{true, NewOutput(addr2, someCoins)},
		{true, NewOutput(addr2, multiCoins)},

		{false, NewOutput(emptyAddr, someCoins)},  // empty address
		{false, NewOutput(addr1, emptyCoins)},     // invalid coins
		{false, NewOutput(addr1, emptyCoins2)},    // invalid coins
		{false, NewOutput(addr1, someEmptyCoins)}, // invalid coins
		{false, NewOutput(addr1, minusCoins)},     // negative coins
		{false, NewOutput(addr1, someMinusCoins)}, // negative coins
		{false, NewOutput(addr1, unsortedCoins)},  // unsorted coins
	}

	for i, tc := range cases {
		err := tc.txOut.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestMsgSendValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte{1, 2})
	addr2 := sdk.AccAddress([]byte{7, 8})
	atom123 := sdk.Coins{sdk.NewInt64Coin("atom", 123)}
	atom124 := sdk.Coins{sdk.NewInt64Coin("atom", 124)}
	eth123 := sdk.Coins{sdk.NewInt64Coin("eth", 123)}
	atom123eth123 := sdk.Coins{sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 123)}

	input1 := NewInput(addr1, atom123)
	input2 := NewInput(addr1, eth123)
	output1 := NewOutput(addr2, atom123)
	output2 := NewOutput(addr2, atom124)
	output3 := NewOutput(addr2, eth123)
	outputMulti := NewOutput(addr2, atom123eth123)

	var emptyAddr sdk.AccAddress

	cases := []struct {
		valid bool
		tx    MsgSend
	}{
		{false, MsgSend{}},                           // no input or output
		{false, MsgSend{Inputs: []Input{input1}}},    // just input
		{false, MsgSend{Outputs: []Output{output1}}}, // just output
		{false, MsgSend{
			Inputs:  []Input{NewInput(emptyAddr, atom123)}, // invalid input
			Outputs: []Output{output1}}},
		{false, MsgSend{
			Inputs:  []Input{input1},
			Outputs: []Output{{emptyAddr, atom123}}}, // invalid output
		},
		{false, MsgSend{
			Inputs:  []Input{input1},
			Outputs: []Output{output2}}, // amounts dont match
		},
		{false, MsgSend{
			Inputs:  []Input{input1},
			Outputs: []Output{output3}}, // amounts dont match
		},
		{false, MsgSend{
			Inputs:  []Input{input1},
			Outputs: []Output{outputMulti}}, // amounts dont match
		},
		{false, MsgSend{
			Inputs:  []Input{input2},
			Outputs: []Output{output1}}, // amounts dont match
		},

		{true, MsgSend{
			Inputs:  []Input{input1},
			Outputs: []Output{output1}},
		},
		{true, MsgSend{
			Inputs:  []Input{input1, input2},
			Outputs: []Output{outputMulti}},
		},
	}

	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestMsgSendGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("input"))
	addr2 := sdk.AccAddress([]byte("output"))
	coins := sdk.Coins{sdk.NewInt64Coin("atom", 10)}
	var msg = MsgSend{
		Inputs:  []Input{NewInput(addr1, coins)},
		Outputs: []Output{NewOutput(addr2, coins)},
	}
	res := msg.GetSignBytes()

	expected := `{"inputs":[{"address":"cosmos1d9h8qat57ljhcm","coins":[{"amount":"10","denom":"atom"}]}],"outputs":[{"address":"cosmos1da6hgur4wsmpnjyg","coins":[{"amount":"10","denom":"atom"}]}]}`
	require.Equal(t, expected, string(res))
}

func TestMsgSendGetSigners(t *testing.T) {
	var msg = MsgSend{
		Inputs: []Input{
			NewInput(sdk.AccAddress([]byte("input1")), nil),
			NewInput(sdk.AccAddress([]byte("input2")), nil),
			NewInput(sdk.AccAddress([]byte("input3")), nil),
		},
	}
	res := msg.GetSigners()
	// TODO: fix this !
	require.Equal(t, fmt.Sprintf("%v", res), "[696E70757431 696E70757432 696E70757433]")
}

/*
// what to do w/ this test?
func TestMsgSendSigners(t *testing.T) {
	signers := []sdk.AccAddress{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	someCoins := sdk.Coins{sdk.NewInt64Coin("atom", 123)}
	inputs := make([]Input, len(signers))
	for i, signer := range signers {
		inputs[i] = NewInput(signer, someCoins)
	}
	tx := NewMsgSend(inputs, nil)

	require.Equal(t, signers, tx.Signers())
}
*/

// ----------------------------------------
// MsgIssue Tests
func TestNewMsgIssue(t *testing.T) {
	// TODO
}

func TestMsgIssueType(t *testing.T) {
	// Construct an MsgIssue
	addr := sdk.AccAddress([]byte("loan-from-bank"))
	coins := sdk.Coins{sdk.NewInt64Coin("atom", 10)}
	var msg = MsgIssue{
		Banker:  sdk.AccAddress([]byte("input")),
		Outputs: []Output{NewOutput(addr, coins)},
	}

	// TODO some failures for bad result
	require.Equal(t, msg.Type(), "bank")
}

func TestMsgIssueValidation(t *testing.T) {
	addr := sdk.AccAddress([]byte("MainAddress"))
	addr1 := sdk.AccAddress([]byte("OutputAddress"))
	coins := sdk.Coins{sdk.NewInt64Coin("atom", 10)}
	cases := []struct {
		valid bool
		tx    MsgIssue
	}{
		{true, NewMsgIssue(addr, []Output{NewOutput(addr1, coins)})},
		{false, NewMsgIssue(addr, nil)},
	}
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", err)
		}
	}
}

func TestMsgIssueGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress([]byte("loan-from-bank"))
	coins := sdk.Coins{sdk.NewInt64Coin("atom", 10)}
	var msg = MsgIssue{
		Banker:  sdk.AccAddress([]byte("input")),
		Outputs: []Output{NewOutput(addr, coins)},
	}
	res := msg.GetSignBytes()

	expected := `{"banker":"cosmos1d9h8qat57ljhcm","outputs":[{"address":"cosmos1d3hkzm3dveex7mfdvfsku6cjngpcj","coins":[{"amount":"10","denom":"atom"}]}]}`
	require.Equal(t, expected, string(res))
}

func TestMsgIssueGetSigners(t *testing.T) {
	var msg = MsgIssue{
		Banker: sdk.AccAddress([]byte("onlyone")),
	}
	res := msg.GetSigners()
	require.Equal(t, fmt.Sprintf("%v", res), "[6F6E6C796F6E65]")
}

func getPegHash(i int) (peghash sdk.PegHash) {
	peghash, _ = sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(1)))
	return peghash
}
func TestMsgIssueAssetType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("issuerAddress"))
	addr2 := sdk.AccAddress([]byte("toAddress"))
	assetPeg := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "SDF123", AssetType: "RICE", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: nil, Locked: false}
	var msg = MsgBankIssueAssets{
		IssueAssets: []IssueAsset{NewIssueAsset(addr1, addr2, assetPeg)},
	}
	require.Equal(t, msg.Type(), "bank")
}

/*
func TestValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("idfgh"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	addr3 := sdk.AccAddress([]byte("OwnerAddress"))
	var issuerAddress sdk.AccAddress
	var toAddress sdk.AccAddress
	assetPeg := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "12345", AssetType: "RICE", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: addr3, Locked: false}
	emptyPegHash := &sdk.BaseAssetPeg{PegHash: nil, DocumentHash: "AAAA", AssetType: "RICE", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: addr3, Locked: false}
	emptyDocumentHash := &sdk.BaseAssetPeg{PegHash: sdk.PegHash([]byte("")), DocumentHash: "Aasdf", AssetType: "RICE", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: addr3, Locked: false}
	emptyAssetType := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "AAAA", AssetType: "asdfg1234", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: addr3, Locked: false}
	emptyQuantitype := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "ADF", AssetType: "SDF", AssetQuantity: 12345, QuantityUnit: "234", OwnerAddress: addr3, Locked: false}
	emptyAddress := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "AFF", AssetType: "ASD12345F", AssetQuantity: 1234, QuantityUnit: "ASDFG", OwnerAddress: nil, Locked: false}
	cases := []struct {
		valid bool
		tx    MsgBankIssueAssets
	}{
		{true, MsgBankIssueAssets{[]IssueAsset{NewIssueAsset(addr1, addr2, assetPeg)}}},
		{false, MsgBankIssueAssets{[]IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg)}}},
		{false, MsgBankIssueAssets{[]IssueAsset{NewIssueAsset(addr1, toAddress, assetPeg)}}},
		{true, MsgBankIssueAssets{[]IssueAsset{NewIssueAsset(addr1, addr1, emptyPegHash)}}},
		{false, MsgBankIssueAssets{[]IssueAsset{NewIssueAsset(addr1, addr2, emptyDocumentHash)}}},
		{false, MsgBankIssueAssets{[]IssueAsset{NewIssueAsset(addr1, addr2, emptyAssetType)}}},
		{false, MsgBankIssueAssets{[]IssueAsset{NewIssueAsset(addr1, addr2, emptyQuantitype)}}},
		{false, MsgBankIssueAssets{[]IssueAsset{NewIssueAsset(addr1, addr2, emptyAddress)}}},
	}
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestMsgIssueAssetValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssuerAddress"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	assetPeg := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "12345", AssetType: "rice", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: addr3, Locked: false}
	var emptyIssuerAddress sdk.AccAddress
	var emptyToAddress sdk.AccAddress
	emptyDoccumenthash := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "ABsdfgh123456", AssetType: "rice", AssetQuantity: 123, QuantityUnit: "MT", OwnerAddress: nil, Locked: false}
	InvalidAssetType := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "ABsdfgh123456", AssetType: "rice12345", AssetQuantity: 123, QuantityUnit: "MT", OwnerAddress: nil, Locked: false}

	cases := []struct {
		valid bool
		tx    MsgBankIssueAssets
	}{
		{true, MsgBankIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(addr1, addr2, assetPeg)}}},
		{false, MsgBankIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(emptyIssuerAddress, addr2, assetPeg)}}},
		{false, MsgBankIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(addr1, emptyToAddress, assetPeg)}}},
		{false, MsgBankIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(addr1, addr2, emptyDoccumenthash)}}},
		{false, MsgBankIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(addr1, addr2, InvalidAssetType)}}},
	}
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestMsgBankIssueAssetsGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssueAddress"))
	addr2 := sdk.AccAddress(([]byte("ToAddress")))
	assetPeg := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "AB1234", AssetType: "sdfg", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: nil, Locked: false}
	issueAsset := NewIssueAsset(addr1, addr2, assetPeg)
	var msg = MsgBankIssueAssets{
		IssueAssets: []IssueAsset{issueAsset},
	}
	res := msg.GetSignBytes()
	expected := `{"issueAssets":[{"issuerAddress":"cosmos1f9ehxat9g9jxgun9wdessv22y8","toAddress":"cosmos123h5zerywfjhxuclxj67f","assetPeg":{"type":"commit-blockchain/AssetPeg","value":{"pegHash":"31","documentHash":"AB1234","assetType":"sdfg","assetQuantity":"1234","quantityUnit":"MT","ownerAddress":"cosmos1550dq7","locked":false}}}]}`
	require.Equal(t, expected, string(res))
}

func TestMsgBankIssueAssetGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssueAddress"))
	addr2 := sdk.AccAddress(([]byte("ToAddress")))
	assetPeg := &sdk.BaseAssetPeg{PegHash: getPegHash(1), DocumentHash: "AB1234", AssetType: "sdfg", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: nil, Locked: false}
	assetPeg1 := &sdk.BaseAssetPeg{PegHash: getPegHash(2), DocumentHash: "AB1234", AssetType: "sdfg", AssetQuantity: 1234, QuantityUnit: "MT", OwnerAddress: nil, Locked: false}

	var msg = MsgBankIssueAssets{
		IssueAssets: []IssueAsset{
			NewIssueAsset(addr1, addr2, assetPeg),
			NewIssueAsset(addr1, addr2, assetPeg1),
		},
	}
	res := msg.GetSigners()
	require.Equal(t, "[497373756541646472657373 497373756541646472657373]", fmt.Sprintf("%v", res))
}

func TestMsgBankIssueFiatTyep(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssuerAddress"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	addr3 := sdk.AccAddress([]byte("owner"))
	fiatPeg := &sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "ASDFG2345", TransactionAmount: 12345, RedeemedAmount: 123, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 12}}}
	var msg = MsgBankIssueFiats{
		IssueFiats: []IssueFiat{
			NewIssueFiat(addr1, addr2, fiatPeg),
		},
	}
	require.Equal(t, msg.Type(), "bank")
}

func TestMsgBankIssueFiatValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssuerAddress"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	addr3 := sdk.AccAddress([]byte("OwnerAddress"))
	fiatPeg := &sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "ADFGHJ12345678", TransactionAmount: 125, RedeemedAmount: 0, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 2345}}}
	// var emptyIssuerAddress sdk.AccAddress
	var emptyToAddress sdk.AccAddress
	emptyTxID := &sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "", TransactionAmount: 125, RedeemedAmount: 0, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 2345}}}
	emptyTxAmt := &sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "ABsdfgh123456", TransactionAmount: 0, RedeemedAmount: 23, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 2345}}}
	emptyOwnerAddr := &sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "ASDF", TransactionAmount: 123, RedeemedAmount: 2345, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: nil, Amount: 2345}}}

	issueFiat := NewIssueFiat(addr1, addr2, fiatPeg)
	issueFiatInvalidTxID := NewIssueFiat(addr1, addr2, emptyTxID)
	issueFiatEmptyAdd2 := NewIssueFiat(addr1, emptyToAddress, fiatPeg)
	fiatWithInvalidAmt := NewIssueFiat(addr1, addr2, emptyTxAmt)
	fiatWithInvalidAddr := NewIssueFiat(addr1, addr2, emptyOwnerAddr)
	cases := []struct {
		valid bool
		tx    MsgBankIssueFiats
	}{
		{true, MsgBankIssueFiats{IssueFiats: []IssueFiat{issueFiat}}},
		{false, MsgBankIssueFiats{IssueFiats: []IssueFiat{issueFiatInvalidTxID}}},
		{false, MsgBankIssueFiats{IssueFiats: []IssueFiat{issueFiatEmptyAdd2}}},
		{false, MsgBankIssueFiats{IssueFiats: []IssueFiat{fiatWithInvalidAmt}}},
		{true, MsgBankIssueFiats{IssueFiats: []IssueFiat{fiatWithInvalidAddr}}},
	}
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestMsgBankIssueFiatsGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssueAddress"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	addr3 := sdk.AccAddress([]byte("OwnerAddress"))
	fiatPeg := &sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "AB1234", TransactionAmount: 1234, RedeemedAmount: 0, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 1234}}}
	issueFiat := NewIssueFiat(addr1, addr2, fiatPeg)
	var msg = MsgBankIssueFiats{
		IssueFiats: []IssueFiat{issueFiat},
	}
	res := msg.GetSignBytes()
	expected := `{"issueFiats":[{"issuerAddress":"cosmos1f9ehxat9g9jxgun9wdessv22y8","toAddress":"cosmos123h5zerywfjhxuclxj67f","fiatPeg":{"type":"commit-blockchain/FiatPeg","value":{"pegHash":"31","transactionID":"AB1234","transactionAmount":"1234","redeemedAmount":"0","owners":[{"ownerAddress":"cosmos1famkuetjg9jxgun9wdes77yfn5","amount":"1234"}]}}}]}`
	require.Equal(t, expected, string(res))
}

func TestMsgBankIssueFiatsGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssueAddress"))
	addr2 := sdk.AccAddress(([]byte("ToAddress")))
	addr3 := sdk.AccAddress(([]byte("OwnerAddress")))
	fiatPeg := &sdk.BaseFiatPeg{PegHash: getPegHash(1), TransactionID: "AB1234", TransactionAmount: 1234, RedeemedAmount: 0, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 12345}}}
	fiatPeg1 := &sdk.BaseFiatPeg{PegHash: getPegHash(2), TransactionID: "AB1234", TransactionAmount: 1234, RedeemedAmount: 0, Owners: []sdk.Owner{sdk.Owner{OwnerAddress: addr3, Amount: 12345}}}

	var msg = MsgBankIssueFiats{
		IssueFiats: []IssueFiat{
			NewIssueFiat(addr1, addr2, fiatPeg),
			NewIssueFiat(addr1, addr2, fiatPeg1),
		},
	}
	res := msg.GetSigners()
	require.Equal(t, "[497373756541646472657373 497373756541646472657373]", fmt.Sprintf("%v", res))
}

func TestMsgBankSendAssetsType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("FromAddress"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	var msg = MsgBankSendAssets{
		SendAssets: []SendAsset{
			NewSendAsset(addr1, addr2, getPegHash(1)),
		},
	}
	require.Equal(t, msg.Type(), "bank")
}

func TestMsgBankSendAssetsValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("FromAddress"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	var emptyFromAddr sdk.AccAddress
	var emptyToAddr sdk.AccAddress
	cases := []struct {
		valid bool
		tx    MsgBankSendAssets
	}{
		{true, MsgBankSendAssets{SendAssets: []SendAsset{NewSendAsset(addr1, nil, getPegHash(1))}}},
		{false, MsgBankSendAssets{SendAssets: []SendAsset{NewSendAsset(emptyFromAddr, addr2, getPegHash(1))}}},
		{false, MsgBankSendAssets{SendAssets: []SendAsset{NewSendAsset(addr1, emptyToAddr, getPegHash(1))}}},
		{false, MsgBankSendAssets{SendAssets: []SendAsset{NewSendAsset(emptyFromAddr, emptyToAddr, getPegHash(1))}}},
	}
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", err)
		}
	}
}

func TestMsgBankSendAssetsGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssueAddress"))
	addr2 := sdk.AccAddress(([]byte("ToAddress")))
	var msg = MsgBankSendAssets{
		SendAssets: []SendAsset{NewSendAsset(addr1, addr2, getPegHash(1))},
	}
	res := msg.GetSignBytes()
	expected := `{"sendAssets":[{"fromAddress":"cosmos1f9ehxat9g9jxgun9wdessv22y8","toAddress":"cosmos123h5zerywfjhxuclxj67f","pegHash":"31"}]}`
	require.Equal(t, expected, string(res))
}

func TestMsgBankSendAssetsGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssueAddress"))
	addr2 := sdk.AccAddress(([]byte("ToAddress")))
	var msg = MsgBankSendAssets{
		SendAssets: []SendAsset{
			NewSendAsset(addr1, addr2, getPegHash(1)),
			NewSendAsset(addr1, addr2, getPegHash(2)),
		},
	}
	res := msg.GetSigners()
	require.Equal(t, "[497373756541646472657373 497373756541646472657373]", fmt.Sprintf("%v", res))
}

func TestMsgBankSendFiatsType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("FromAddress"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	var msg = MsgBankSendFiats{
		SendFiats: []SendFiat{
			NewSendFiat(addr1, addr2, getPegHash(1), 1234),
		},
	}
	require.Equal(t, msg.Type(), "bank")
}

func TestMsgBankSendFiatsValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("FromAddress"))
	addr2 := sdk.AccAddress([]byte("ToAddress"))
	var emptyFromAddr sdk.AccAddress
	var emptyToAddr sdk.AccAddress
	// sendFiat := NewSendFiat(addr1, addr2, getPegHash(1), 123)
	cases := []struct {
		valid bool
		tx    MsgBankSendFiats
	}{
		{false, MsgBankSendFiats{SendFiats: []SendFiat{NewSendFiat(emptyFromAddr, addr2, getPegHash(1), 12)}}},
		{false, MsgBankSendFiats{SendFiats: []SendFiat{NewSendFiat(addr1, emptyToAddr, getPegHash(1), 123)}}},
		{false, MsgBankSendFiats{SendFiats: []SendFiat{NewSendFiat(emptyFromAddr, emptyToAddr, getPegHash(1), 12345)}}},
		{false, MsgBankSendFiats{SendFiats: []SendFiat{NewSendFiat(addr1, addr2, getPegHash(1), 0)}}},
		{true, MsgBankSendFiats{SendFiats: []SendFiat{NewSendFiat(addr1, addr2, getPegHash(1), 2)}}},
	}
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", err)
		}
	}
}

func TestMsgBankSendFiatsGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssueAddress"))
	addr2 := sdk.AccAddress(([]byte("ToAddress")))
	var msg = MsgBankSendFiats{
		SendFiats: []SendFiat{NewSendFiat(addr1, addr2, getPegHash(1), 123)},
	}
	res := msg.GetSignBytes()
	expected := `{"sendFiats":[{"fromAddress":"cosmos1f9ehxat9g9jxgun9wdessv22y8","toAddress":"cosmos123h5zerywfjhxuclxj67f","pegHash":"31","amount":"123"}]}`
	require.Equal(t, expected, string(res))
}

func TestMsgBankSendFiatsGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("IssueAddress"))
	addr2 := sdk.AccAddress(([]byte("ToAddress")))
	var msg = MsgBankSendFiats{
		SendFiats: []SendFiat{
			NewSendFiat(addr1, addr2, getPegHash(1), 123),
			NewSendFiat(addr1, addr2, getPegHash(2), 12345),
		},
	}
	res := msg.GetSigners()
	require.Equal(t, "[497373756541646472657373 497373756541646472657373]", fmt.Sprintf("%v", res))
}

func TestMsgExecuteOrderType(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("MediatorAddress"))
	addr2 := sdk.AccAddress([]byte("BuyerAddress"))
	addr3 := sdk.AccAddress([]byte("sellerAddress"))
	var msg = MsgBankBuyerExecuteOrders{
		BuyerExecuteOrders: []BuyerExecuteOrder{
			NewBuyerExecuteOrder(addr1, addr2, addr3, getPegHash(1), ""),
		},
	}
	require.Equal(t, msg.Type(), "bank")
}

func TestMsgExecuteOrderValidation(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("MediatorAddress"))
	addr2 := sdk.AccAddress([]byte("BuyerAddress"))
	addr3 := sdk.AccAddress([]byte("SellerAddress"))
	var emptyMediatorAddr sdk.AccAddress
	var emptyBuyerAddr sdk.AccAddress
	var emptySellerAddr sdk.AccAddress
	// sendFiat := NewSendFiat(addr1, addr2, getPegHash(1), 123)
	cases := []struct {
		valid bool
		tx    MsgBankBuyerExecuteOrders
	}{
		{true, MsgBankBuyerExecuteOrders{BuyerExecuteOrders: []BuyerExecuteOrder{NewBuyerExecuteOrder(addr1, addr2, addr3, getPegHash(1), "")}}},
		{false, MsgBankBuyerExecuteOrders{BuyerExecuteOrders: []BuyerExecuteOrder{NewBuyerExecuteOrder(emptyMediatorAddr, addr2, addr3, getPegHash(1), "")}}},
		{false, MsgBankBuyerExecuteOrders{BuyerExecuteOrders: []BuyerExecuteOrder{NewBuyerExecuteOrder(addr1, emptyBuyerAddr, addr3, getPegHash(1), "")}}},
		{false, MsgBankBuyerExecuteOrders{BuyerExecuteOrders: []BuyerExecuteOrder{NewBuyerExecuteOrder(addr1, addr2, emptySellerAddr, getPegHash(1), "")}}},
		{false, MsgBankBuyerExecuteOrders{BuyerExecuteOrders: []BuyerExecuteOrder{NewBuyerExecuteOrder(emptyMediatorAddr, addr2, emptySellerAddr, getPegHash(1), "")}}},
	}
	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", err)
		}
	}
}

func TestMsgBankExecuteOrdersGetSignBytes(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("MediatorAddress"))
	addr2 := sdk.AccAddress(([]byte("BuyerAddress")))
	addr3 := sdk.AccAddress(([]byte("SellerAddress")))
	var msg = MsgBankBuyerExecuteOrders{
		BuyerExecuteOrders: []BuyerExecuteOrder{NewBuyerExecuteOrder(addr1, addr2, addr3, getPegHash(1), "")},
	}
	res := msg.GetSignBytes()
	expected := `{"executeOrders":[{"mediatorAddress":"cosmos1f4jkg6tpw3hhystyv3ex2umn8gfvne","buyerAddress":"cosmos1gf6hjetjg9jxgun9wdestq703s","sellerAddress":"cosmos12djkcmr9wfqkgerjv4ehxkfq79n","pegHash":"31"}]}`
	require.Equal(t, expected, string(res))
}

func TestMsgExecuteOrderGetSigners(t *testing.T) {
	addr1 := sdk.AccAddress([]byte("MediatorAddress"))
	addr2 := sdk.AccAddress(([]byte("BuyerAddress")))
	addr3 := sdk.AccAddress(([]byte("SellerAddress")))
	var msg = MsgBankBuyerExecuteOrders{
		BuyerExecuteOrders: []BuyerExecuteOrder{NewBuyerExecuteOrder(addr1, addr2, addr3, getPegHash(1), "")},
	}
	res := msg.GetSigners()
	require.Equal(t, "[4D65646961746F7241646472657373]", fmt.Sprintf("%v", res))
}
*/
