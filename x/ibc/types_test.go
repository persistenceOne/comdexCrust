package ibc

import (
	"fmt"
	"testing"
	
	"github.com/stretchr/testify/require"
	
	sdk "github.com/comdex-blockchain/types"
	"github.com/tendermint/tendermint/libs/common"
)

// --------------------------------
// IBCPacket Tests

func TestIBCPacketValidation(t *testing.T) {
	cases := []struct {
		valid  bool
		packet IBCPacket
	}{
		{true, constructIBCPacket(true)},
		{false, constructIBCPacket(false)},
	}
	
	for i, tc := range cases {
		err := tc.packet.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

// -------------------------------
// IBCTransferMsg Tests

func TestIBCTransferMsg(t *testing.T) {
	packet := constructIBCPacket(true)
	msg := IBCTransferMsg{packet}
	
	require.Equal(t, msg.Type(), "ibc")
}

func TestIBCTransferMsgValidation(t *testing.T) {
	validPacket := constructIBCPacket(true)
	invalidPacket := constructIBCPacket(false)
	
	cases := []struct {
		valid bool
		msg   IBCTransferMsg
	}{
		{true, IBCTransferMsg{validPacket}},
		{false, IBCTransferMsg{invalidPacket}},
	}
	
	for i, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

// -------------------------------
// IBCReceiveMsg Tests

func TestIBCReceiveMsg(t *testing.T) {
	packet := constructIBCPacket(true)
	msg := IBCReceiveMsg{packet, sdk.AccAddress([]byte("relayer")), 0}
	
	require.Equal(t, msg.Type(), "ibc")
}

func TestIBCReceiveMsgValidation(t *testing.T) {
	validPacket := constructIBCPacket(true)
	invalidPacket := constructIBCPacket(false)
	
	cases := []struct {
		valid bool
		msg   IBCReceiveMsg
	}{
		{true, IBCReceiveMsg{validPacket, sdk.AccAddress([]byte("relayer")), 0}},
		{false, IBCReceiveMsg{invalidPacket, sdk.AccAddress([]byte("relayer")), 0}},
	}
	
	for i, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

// -------------------------------
// Helpers

func constructIBCPacket(valid bool) IBCPacket {
	srcAddr := sdk.AccAddress([]byte("source"))
	destAddr := sdk.AccAddress([]byte("destination"))
	coins := sdk.Coins{sdk.NewInt64Coin("atom", 10)}
	srcChain := "source-chain"
	destChain := "dest-chain"
	
	if valid {
		return NewIBCPacket(srcAddr, destAddr, coins, srcChain, destChain)
	}
	return NewIBCPacket(srcAddr, destAddr, coins, srcChain, srcChain)
}

var issuerAddress = sdk.AccAddress([]byte("IssuerAddress"))
var fromAddress = sdk.AccAddress([]byte("FromAddress"))
var toAddress = sdk.AccAddress([]byte("toAddress"))
var redeemerAddress = sdk.AccAddress([]byte("RedeemerAddress"))
var relayerAddress = sdk.AccAddress([]byte("RelayerAddress"))
var mediatorAddress = sdk.AccAddress([]byte("MediatorAddress"))
var srcChain = "comdex-main"
var destAssetChain = "comdex-asset"
var destFiatChain = "comdex-fiat"
var amount int64 = 1234
var pegHash common.HexBytes = sdk.PegHash([]byte("pegHash1"))
var pegHash2 common.HexBytes = sdk.PegHash([]byte(""))

var fiatPegWallet = sdk.FiatPegWallet{sdk.BaseFiatPeg{PegHash: sdk.PegHash([]byte("1")), TransactionID: "FB8AE3A02BBCD2", TransactionAmount: 1000, RedeemedAmount: 0, Owners: []sdk.Owner{{OwnerAddress: nil, Amount: 500}, {OwnerAddress: sdk.AccAddress([]byte("relayer")), Amount: 500}}}}
var fiatPegWallet2 sdk.FiatPegWallet
var assetPeg = &sdk.BaseAssetPeg{PegHash: sdk.PegHash([]byte("1")), DocumentHash: "ABCD123", AssetType: "FGHJK", AssetQuantity: 1234, AssetPrice: 1234, QuantityUnit: "MT", OwnerAddress: nil, Locked: false}
var peg = &sdk.BaseAssetPeg{PegHash: sdk.PegHash([]byte("1")), DocumentHash: "ABCD123fghj", AssetType: "FGHJK", AssetQuantity: 1234, AssetPrice: 1234, QuantityUnit: "MT", OwnerAddress: nil, Locked: false}
var fiatPeg = &sdk.BaseFiatPeg{PegHash: sdk.PegHash([]byte("1")), TransactionID: "ASDFGHJKL", TransactionAmount: 12345, RedeemedAmount: 23456, Owners: []sdk.Owner{{OwnerAddress: nil, Amount: 12}}}
var peg2 = &sdk.BaseFiatPeg{PegHash: sdk.PegHash([]byte("1")), TransactionID: "dfghj", TransactionAmount: 0, RedeemedAmount: 23456, Owners: []sdk.Owner{{OwnerAddress: nil, Amount: 12}}}

// *****MsgIssueAssets
func TestIssueAssetType(t *testing.T) {
	var msg = MsgIssueAssets{
		IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)},
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestIssueAssetsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgIssueAssets
	}{
		{true, MsgIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)}}},
		{false, MsgIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(nil, toAddress, assetPeg, srcChain, destAssetChain)}}},
		{false, MsgIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, nil, assetPeg, srcChain, destAssetChain)}}},
		{false, MsgIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, peg, srcChain, destAssetChain)}}},
		{false, MsgIssueAssets{IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, srcChain)}}},
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

func TestIssueAssetGetSignBytes(t *testing.T) {
	var msg = MsgIssueAssets{
		IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)},
	}
	expected := `{"issueAssets":[{"issuerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","toAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","assetPeg":{"type":"comdex-blockchain/AssetPeg","value":{"pegHash":"31","documentHash":"ABCD123","assetType":"FGHJK","assetQuantity":"1234","assetPrice":"1234","quantityUnit":"MT","ownerAddress":"cosmos1550dq7","locked":false}},"sourceChain":"comdex-main","destinationChain":"comdex-asset"}]}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestIssueAssetGetSigners(t *testing.T) {
	var msg = MsgIssueAssets{
		IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)},
	}
	expected := "[49737375657241646472657373]"
	require.Equal(t, expected, fmt.Sprintf("%v", msg.GetSigners()))
}

// #####MsgIssueAssets

// *****MsgRelayIssueAssets
func TestRelayIssueAssetType(t *testing.T) {
	var msg = MsgRelayIssueAssets{
		Relayer:     relayerAddress,
		IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)},
		Sequence:    0,
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestRelayIssueAssetsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgRelayIssueAssets
	}{
		{true, MsgRelayIssueAssets{[]IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueAssets{[]IssueAsset{NewIssueAsset(nil, toAddress, assetPeg, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueAssets{[]IssueAsset{NewIssueAsset(issuerAddress, nil, assetPeg, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueAssets{[]IssueAsset{NewIssueAsset(issuerAddress, toAddress, peg, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueAssets{[]IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, srcChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueAssets{[]IssueAsset{NewIssueAsset(issuerAddress, toAddress, peg, srcChain, destAssetChain)}, nil, 0}},
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

func TestRelayIssueAssetGetSignBytes(t *testing.T) {
	var msg = MsgRelayIssueAssets{
		IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)},
		Relayer:     relayerAddress,
		Sequence:    0,
	}
	expected := `{"issueAssets":[{"issuerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","toAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","assetPeg":{"type":"comdex-blockchain/AssetPeg","value":{"pegHash":"31","documentHash":"ABCD123","assetType":"FGHJK","assetQuantity":"1234","assetPrice":"1234","quantityUnit":"MT","ownerAddress":"cosmos1550dq7","locked":false}},"sourceChain":"comdex-main","destinationChain":"comdex-asset"}],"relayer":"cosmos12fjkcctev4eyzerywfjhxucryd09t","sequence":"0"}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestRelayIssueAssetGetSigners(t *testing.T) {
	var msg = MsgRelayIssueAssets{
		IssueAssets: []IssueAsset{NewIssueAsset(issuerAddress, toAddress, assetPeg, srcChain, destAssetChain)},
		Relayer:     relayerAddress,
		Sequence:    0,
	}
	expected := "[52656C6179657241646472657373]"
	require.Equal(t, expected, fmt.Sprintf("%v", msg.GetSigners()))
}

// #####MsgRelayIssueAssets

// *****MsgRedeemAssets
func TestRedeemAssetType(t *testing.T) {
	var msg = MsgRedeemAssets{
		RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)},
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestRedeemAssetsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgRedeemAssets
	}{
		{true, MsgRedeemAssets{RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)}}},
		{false, MsgRedeemAssets{RedeemAssets: []RedeemAsset{NewRedeemAsset(nil, redeemerAddress, pegHash, srcChain, destAssetChain)}}},
		{false, MsgRedeemAssets{RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, nil, pegHash, srcChain, destAssetChain)}}},
		{false, MsgRedeemAssets{RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash2, srcChain, destAssetChain)}}},
		{false, MsgRedeemAssets{RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, srcChain)}}},
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

func TestRedeemAssetGetSignBytes(t *testing.T) {
	var msg = MsgRedeemAssets{
		RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)},
	}
	expected := `{"redeemAssets":[{"issuerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","redeemeraddress":"cosmos12fjkget9d4jhystyv3ex2umnshydxs","pegHash":"7065674861736831","sourceChain":"comdex-main","destinationChain":"comdex-asset"}]}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestRedeemAssetGetSigners(t *testing.T) {
	var msg = MsgRedeemAssets{
		RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)},
	}
	expected := "[49737375657241646472657373]"
	require.Equal(t, expected, fmt.Sprintf("%v", msg.GetSigners()))
}

// #####MsgRedeemAssets

// *****MsgRelayRedeemAssets
func TestRelayRedeemAssetType(t *testing.T) {
	var msg = MsgRelayRedeemAssets{
		Relayer:      relayerAddress,
		RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)},
		Sequence:     0,
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestRelayRedeemAssetsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgRelayRedeemAssets
	}{
		{true, MsgRelayRedeemAssets{[]RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemAssets{[]RedeemAsset{NewRedeemAsset(nil, redeemerAddress, pegHash, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemAssets{[]RedeemAsset{NewRedeemAsset(issuerAddress, nil, pegHash, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemAssets{[]RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash2, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemAssets{[]RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, srcChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemAssets{[]RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)}, nil, 0}},
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

func TestRelayRedeemAssetGetSignBytes(t *testing.T) {
	var msg = MsgRelayRedeemAssets{
		RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)},
		Relayer:      relayerAddress,
		Sequence:     0,
	}
	expected := `{"redeemAssets":[{"issuerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","redeemeraddress":"cosmos12fjkget9d4jhystyv3ex2umnshydxs","pegHash":"7065674861736831","sourceChain":"comdex-main","destinationChain":"comdex-asset"}],"relayer":"cosmos12fjkcctev4eyzerywfjhxucryd09t","sequence":"0"}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestRelayRedeemAssetGetSigners(t *testing.T) {
	var msg = MsgRelayRedeemAssets{
		RedeemAssets: []RedeemAsset{NewRedeemAsset(issuerAddress, redeemerAddress, pegHash, srcChain, destAssetChain)},
		Relayer:      relayerAddress,
		Sequence:     0,
	}
	expected := "[52656C6179657241646472657373]"
	require.Equal(t, expected, fmt.Sprintf("%v", msg.GetSigners()))
}

// #####MsgRelayRedeemAssets

// *****MsgIssueFiats
func TestIssueFiatType(t *testing.T) {
	var msg = MsgIssueFiats{
		IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)},
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestIssueFiatsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgIssueFiats
	}{
		{true, MsgIssueFiats{IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)}}},
		{false, MsgIssueFiats{IssueFiats: []IssueFiat{NewIssueFiat(nil, toAddress, fiatPeg, srcChain, destFiatChain)}}},
		{false, MsgIssueFiats{IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, nil, fiatPeg, srcChain, destFiatChain)}}},
		{false, MsgIssueFiats{IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, peg2, srcChain, destFiatChain)}}},
		{false, MsgIssueFiats{IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, srcChain)}}},
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

func TestIssueFiatGetSignBytes(t *testing.T) {
	var msg = MsgIssueFiats{
		IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)},
	}
	expected := `{"issueFiats":[{"issuerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","toAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","fiatPeg":{"type":"comdex-blockchain/FiatPeg","value":{"pegHash":"31","transactionID":"ASDFGHJKL","transactionAmount":"12345","redeemedAmount":"23456","owners":[{"ownerAddress":"cosmos1550dq7","amount":"12"}]}},"sourceChain":"comdex-main","destinationChain":"comdex-fiat"}]}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestIssueFiatGetSigners(t *testing.T) {
	var msg = MsgIssueFiats{
		IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)},
	}
	expected := "[49737375657241646472657373]"
	require.Equal(t, expected, fmt.Sprintf("%v", msg.GetSigners()))
}

// #####MsgIssueFiat

// *****MsgRelayIssueFiats
func TestRelayIssueFiatType(t *testing.T) {
	var msg = MsgRelayIssueFiats{
		IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)},
		Relayer:    relayerAddress,
		Sequence:   0,
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestRelayIssueFiatsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgRelayIssueFiats
	}{
		{true, MsgRelayIssueFiats{[]IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueFiats{[]IssueFiat{NewIssueFiat(nil, toAddress, fiatPeg, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueFiats{[]IssueFiat{NewIssueFiat(issuerAddress, nil, fiatPeg, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueFiats{[]IssueFiat{NewIssueFiat(issuerAddress, toAddress, peg2, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueFiats{[]IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, srcChain)}, relayerAddress, 0}},
		{false, MsgRelayIssueFiats{[]IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)}, nil, 0}},
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

func TestRelayIssueFiatGetSignBytes(t *testing.T) {
	var msg = MsgRelayIssueFiats{
		IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)},
		Relayer:    relayerAddress,
		Sequence:   0,
	}
	expected := `{"issueFiats":[{"issuerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","toAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","fiatPeg":{"type":"comdex-blockchain/FiatPeg","value":{"pegHash":"31","transactionID":"ASDFGHJKL","transactionAmount":"12345","redeemedAmount":"23456","owners":[{"ownerAddress":"cosmos1550dq7","amount":"12"}]}},"sourceChain":"comdex-main","destinationChain":"comdex-fiat"}],"relayer":"cosmos12fjkcctev4eyzerywfjhxucryd09t","sequence":"0"}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestRelayIssueFiatGetSigners(t *testing.T) {
	var msg = MsgRelayIssueFiats{
		IssueFiats: []IssueFiat{NewIssueFiat(issuerAddress, toAddress, fiatPeg, srcChain, destFiatChain)},
		Relayer:    relayerAddress,
		Sequence:   0,
	}
	expected := `[52656C6179657241646472657373]`
	require.Equal(t, expected, fmt.Sprintf("%d", msg.GetSigners()))
}

// #####MsgRelayIssueFiats

// *****MsgRedeemFiats
func TestRedeemFiatType(t *testing.T) {
	var msg = MsgRedeemFiats{
		RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)},
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestRedeemFiatsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgRedeemFiats
	}{
		{true, MsgRedeemFiats{RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgRedeemFiats{RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, nil, amount, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgRedeemFiats{RedeemFiats: []RedeemFiat{NewRedeemFiat(nil, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgRedeemFiats{RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, -44, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgRedeemFiats{RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet2, srcChain, destFiatChain)}}},
		{false, MsgRedeemFiats{RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, srcChain)}}},
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

func TestRedeemFiatGetSignBytes(t *testing.T) {
	var msg = MsgRedeemFiats{
		RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)},
	}
	expected := `{"redeemFiats":[{"redeemerAddress":"cosmos12fjkget9d4jhystyv3ex2umnshydxs","issuerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","amount":"1234","fiatPegWallet":[{"pegHash":"31","transactionID":"FB8AE3A02BBCD2","transactionAmount":"1000","redeemedAmount":"0","owners":[{"ownerAddress":"cosmos1550dq7","amount":"500"},{"ownerAddress":"cosmos1wfjkcctev4eqs3083t","amount":"500"}]}],"sourceChain":"comdex-main","destinationChain":"comdex-fiat"}]}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestRedeemFiatGetSigners(t *testing.T) {
	var msg = MsgRedeemFiats{
		RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)},
	}
	expected := "[49737375657241646472657373]"
	require.Equal(t, expected, fmt.Sprintf("%v", msg.GetSigners()))
}

// #####MsgRedeemFiats

// *****MsgRelayRedeemFiats
func TestRelayRedeemFiatType(t *testing.T) {
	var msg = MsgRelayRedeemFiats{
		Relayer:     relayerAddress,
		RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)},
		Sequence:    0,
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestRelayRedeemFiatsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgRelayRedeemFiats
	}{
		{true, MsgRelayRedeemFiats{[]RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemFiats{[]RedeemFiat{NewRedeemFiat(redeemerAddress, nil, amount, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemFiats{[]RedeemFiat{NewRedeemFiat(nil, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemFiats{[]RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, -44, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemFiats{[]RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet2, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemFiats{[]RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, srcChain)}, relayerAddress, 0}},
		{false, MsgRelayRedeemFiats{[]RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)}, nil, 0}},
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

func TestRelayRedeemFiatGetSignBytes(t *testing.T) {
	var msg = MsgRelayRedeemFiats{
		RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)},
		Relayer:     relayerAddress,
		Sequence:    0,
	}
	expected := `{"redeemFiats":[{"redeemerAddress":"cosmos12fjkget9d4jhystyv3ex2umnshydxs","issuerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","amount":"1234","fiatPegWallet":[{"pegHash":"31","transactionID":"FB8AE3A02BBCD2","transactionAmount":"1000","redeemedAmount":"0","owners":[{"ownerAddress":"cosmos1550dq7","amount":"500"},{"ownerAddress":"cosmos1wfjkcctev4eqs3083t","amount":"500"}]}],"sourceChain":"comdex-main","destinationChain":"comdex-fiat"}],"relayer":"cosmos12fjkcctev4eyzerywfjhxucryd09t","sequence":"0"}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestRelayRedeemFiatGetSigners(t *testing.T) {
	var msg = MsgRelayRedeemFiats{
		RedeemFiats: []RedeemFiat{NewRedeemFiat(redeemerAddress, issuerAddress, amount, fiatPegWallet, srcChain, destFiatChain)},
		Relayer:     relayerAddress,
		Sequence:    0,
	}
	expected := "[52656C6179657241646472657373]"
	require.Equal(t, expected, fmt.Sprintf("%v", msg.GetSigners()))
}

// #####MsgRelayRedeemFiats

// *****MsgSendAssets

func TestSendAssetType(t *testing.T) {
	var msg = MsgSendAssets{
		SendAssets: []SendAsset{NewSendAsset(issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)},
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestSendAssetsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgSendAssets
	}{
		{true, MsgSendAssets{[]SendAsset{NewSendAsset(issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)}}},
		{false, MsgSendAssets{[]SendAsset{NewSendAsset(nil, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)}}},
		{false, MsgSendAssets{[]SendAsset{NewSendAsset(issuerAddress, toAddress, nil, srcChain, "")}}},
		{false, MsgSendAssets{[]SendAsset{NewSendAsset(issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, srcChain)}}},
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

func TestSendAssetGetSignBytes(t *testing.T) {
	var msg = MsgSendAssets{
		SendAssets: []SendAsset{NewSendAsset(issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)},
	}
	expected := `{"sendAssets":[{"fromAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","toAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","pegHash":"31","sourceChain":"comdex-main","destinationChain":"comdex-asset"}]}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}
func TestSendAssetGetSigners(t *testing.T) {
	var msg = MsgSendAssets{
		SendAssets: []SendAsset{NewSendAsset(issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)},
	}
	expected := `[49737375657241646472657373]`
	require.Equal(t, expected, fmt.Sprintf("%d", msg.GetSigners()))
}

// #####MsgSendAssets

// *****MsgRelaySendAssets
func TestRelaySendAssetType(t *testing.T) {
	var msg = MsgRelaySendAssets{
		Relayer:    relayerAddress,
		SendAssets: []SendAsset{NewSendAsset(fromAddress, toAddress, pegHash, srcChain, destAssetChain)},
		Sequence:   0,
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestRelaySendAssetsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgRelaySendAssets
	}{
		{true, MsgRelaySendAssets{[]SendAsset{NewSendAsset(fromAddress, toAddress, pegHash, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelaySendAssets{[]SendAsset{NewSendAsset(nil, toAddress, pegHash, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelaySendAssets{[]SendAsset{NewSendAsset(fromAddress, nil, pegHash, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelaySendAssets{[]SendAsset{NewSendAsset(fromAddress, toAddress, pegHash2, srcChain, destAssetChain)}, relayerAddress, 0}},
		{false, MsgRelaySendAssets{[]SendAsset{NewSendAsset(fromAddress, toAddress, pegHash, srcChain, srcChain)}, relayerAddress, 0}},
		{false, MsgRelaySendAssets{[]SendAsset{NewSendAsset(fromAddress, toAddress, pegHash, srcChain, destAssetChain)}, nil, 0}},
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

func TestRelaySendAssetGetSignBytes(t *testing.T) {
	var msg = MsgRelaySendAssets{
		SendAssets: []SendAsset{NewSendAsset(fromAddress, toAddress, pegHash, srcChain, destAssetChain)},
		Relayer:    relayerAddress,
		Sequence:   0,
	}
	expected := `{"sendAssets":[{"fromAddress":"cosmos1geex7m2pv3j8yetnwvv6x74m","toAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","pegHash":"7065674861736831","sourceChain":"comdex-main","destinationChain":"comdex-asset"}],"relayer":"cosmos12fjkcctev4eyzerywfjhxucryd09t","sequence":"0"}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestRelaySendAssetGetSigners(t *testing.T) {
	var msg = MsgRelaySendAssets{
		SendAssets: []SendAsset{NewSendAsset(fromAddress, toAddress, pegHash, srcChain, destAssetChain)},
		Relayer:    relayerAddress,
		Sequence:   0,
	}
	expected := "[52656C6179657241646472657373]"
	require.Equal(t, expected, fmt.Sprintf("%v", msg.GetSigners()))
}

// #####MsgRelaySendAssets

// *****MsgSendFiats
func TestSendFiatType(t *testing.T) {
	var msg = MsgSendFiats{
		SendFiats: []SendFiat{NewSendFiat(fromAddress, toAddress, sdk.PegHash([]byte("1")), amount, fiatPegWallet, srcChain, destFiatChain)},
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestSendFiatsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgSendFiats
	}{
		{true, MsgSendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, sdk.PegHash([]byte("1")), amount, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgSendFiats{[]SendFiat{NewSendFiat(nil, toAddress, sdk.PegHash([]byte("1")), amount, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgSendFiats{[]SendFiat{NewSendFiat(fromAddress, nil, sdk.PegHash([]byte("1")), amount, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgSendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, nil, amount, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgSendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, sdk.PegHash([]byte("1")), -44, fiatPegWallet, srcChain, destFiatChain)}}},
		{false, MsgSendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, sdk.PegHash([]byte("1")), amount, nil, srcChain, destFiatChain)}}},
		{false, MsgSendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, sdk.PegHash([]byte("1")), amount, fiatPegWallet, srcChain, srcChain)}}},
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

func TestSendFiatGetSignBytes(t *testing.T) {
	var msg = MsgSendFiats{
		SendFiats: []SendFiat{NewSendFiat(fromAddress, toAddress, sdk.PegHash([]byte("1")), amount, fiatPegWallet, srcChain, destFiatChain)},
	}
	expected := `{"sendFiats":[{"fromAddress":"cosmos1geex7m2pv3j8yetnwvv6x74m","toAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","pegHash":"31","amount":"1234","fiatPegWallet":[{"pegHash":"31","transactionID":"FB8AE3A02BBCD2","transactionAmount":"1000","redeemedAmount":"0","owners":[{"ownerAddress":"cosmos1550dq7","amount":"500"},{"ownerAddress":"cosmos1wfjkcctev4eqs3083t","amount":"500"}]}],"sourceChain":"comdex-main","destinationChain":"comdex-fiat"}]}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}
func TestSendFiatGetSigners(t *testing.T) {
	var msg = MsgSendFiats{
		SendFiats: []SendFiat{NewSendFiat(fromAddress, toAddress, sdk.PegHash([]byte("1")), amount, fiatPegWallet, srcChain, destFiatChain)},
	}
	expected := `[46726F6D41646472657373]`
	require.Equal(t, expected, fmt.Sprintf("%d", msg.GetSigners()))
}

// #####MsgSendFiats

// *****MsgRelaySendFiats
func TestRelaySendFiatType(t *testing.T) {
	var msg = MsgRelaySendFiats{
		SendFiats: []SendFiat{NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet, srcChain, destFiatChain)},
		Relayer:   relayerAddress,
		Sequence:  0,
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestRelaySendFiatsValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgRelaySendFiats
	}{
		{true, MsgRelaySendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelaySendFiats{[]SendFiat{NewSendFiat(nil, toAddress, pegHash, amount, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelaySendFiats{[]SendFiat{NewSendFiat(fromAddress, nil, pegHash, amount, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelaySendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, pegHash2, amount, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelaySendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, pegHash, -44, fiatPegWallet, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelaySendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet2, srcChain, destFiatChain)}, relayerAddress, 0}},
		{false, MsgRelaySendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet, srcChain, srcChain)}, relayerAddress, 0}},
		{false, MsgRelaySendFiats{[]SendFiat{NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet, srcChain, destFiatChain)}, nil, 0}},
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

func TestRelaySendFiatGetSignBytes(t *testing.T) {
	var msg = MsgRelaySendFiats{
		SendFiats: []SendFiat{NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet, srcChain, destFiatChain)},
		Relayer:   relayerAddress,
		Sequence:  0,
	}
	expected := `{"sendFiats":[{"fromAddress":"cosmos1geex7m2pv3j8yetnwvv6x74m","toAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","pegHash":"7065674861736831","amount":"1234","fiatPegWallet":[{"pegHash":"31","transactionID":"FB8AE3A02BBCD2","transactionAmount":"1000","redeemedAmount":"0","owners":[{"ownerAddress":"cosmos1550dq7","amount":"500"},{"ownerAddress":"cosmos1wfjkcctev4eqs3083t","amount":"500"}]}],"sourceChain":"comdex-main","destinationChain":"comdex-fiat"}],"relayer":"cosmos12fjkcctev4eyzerywfjhxucryd09t","sequence":"0"}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestRelaySendFiatGetSigners(t *testing.T) {
	var msg = MsgRelaySendFiats{
		SendFiats: []SendFiat{NewSendFiat(fromAddress, toAddress, pegHash, amount, fiatPegWallet, srcChain, destFiatChain)},
		Relayer:   relayerAddress,
		Sequence:  0,
	}
	expected := `[52656C6179657241646472657373]`
	require.Equal(t, expected, fmt.Sprintf("%d", msg.GetSigners()))
}

// #####MsgRelaySendFiats

/*
func TestExecuteOrderType(t *testing.T) {
	var msg = MsgExecuteOrders{
		ExecuteOrders: []ExecuteOrder{NewExecuteOrder(mediatorAddress, issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)},
	}
	require.Equal(t, msg.Type(), "ibc")
}

func TestExecuteOrderValidateBasic(t *testing.T) {
	cases := []struct {
		valid bool
		tx    MsgExecuteOrders
	}{
		{true, MsgExecuteOrders{[]ExecuteOrder{NewExecuteOrder(mediatorAddress, issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)}}},
		{false, MsgExecuteOrders{[]ExecuteOrder{NewExecuteOrder(nil, issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)}}},
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
func TestExecuteOrderGetSigneBytes(t *testing.T) {
	var msg = MsgExecuteOrders{
		ExecuteOrders: []ExecuteOrder{NewExecuteOrder(mediatorAddress, issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)},
	}
	expected := `{"executeOrders":[{"mediatorAddress":"cosmos1f4jkg6tpw3hhystyv3ex2umn8gfvne","buyerAddress":"cosmos1f9ehxat9wfqkgerjv4ehxy08szk","sellerAddress":"cosmos1w3h5zerywfjhxuc7mfk6f","pegHash":"31","sourceChain":"comdex-main","destinationChain":"comdex-asset"}]}`
	require.Equal(t, expected, string(msg.GetSignBytes()))
}

func TestExecuteOrderGetSigners(t *testing.T) {
	var msg = MsgExecuteOrders{
		ExecuteOrders: []ExecuteOrder{NewExecuteOrder(mediatorAddress, issuerAddress, toAddress, sdk.PegHash([]byte("1")), srcChain, destAssetChain)},
	}
	expected := `[4D65646961746F7241646472657373]`
	require.Equal(t, expected, fmt.Sprintf("%d", msg.GetSigners()))
}
*/
