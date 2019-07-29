package fiatFactory

import (
	"testing"

	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

// TestHandleMsgFactoryIssueFiat handles test cases for message factory issue fiat
func TestHandleMsgFactoryIssueFiat(t *testing.T) {

	ms, capkey, _ := setupMultiStore()
	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)

	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	keeper := NewKeeper(newfiatpegmapper)

	_, addrs, _, _ := mock.CreateGenAccounts(6, sdk.Coins{})
	randFiatPeg1 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("prashant")))
	randFiatPeg2 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("92151818")))
	randFiatPeg3 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("*&^%$#")))
	randFiatPeg4 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("")))

	fiatpeg1 := sdk.ToFiatPeg(randFiatPeg1)
	fiatpeg2 := sdk.ToFiatPeg(randFiatPeg2)
	fiatpeg3 := sdk.ToFiatPeg(randFiatPeg3)
	fiatpeg4 := sdk.ToFiatPeg(randFiatPeg4)

	newfiatpegmapper.SetFiatPeg(ctx, fiatpeg1)
	newfiatpegmapper.SetFiatPeg(ctx, fiatpeg2)
	newfiatpegmapper.SetFiatPeg(ctx, fiatpeg3)

	issuFiat1 := NewIssueFiat(addrs[0], addrs[1], fiatpeg1)
	issuFiat2 := NewIssueFiat(addrs[2], addrs[3], fiatpeg2)
	issuFiat3 := NewIssueFiat(addrs[4], addrs[5], fiatpeg3)
	issuFiat4 := NewIssueFiat(addrs[0], addrs[1], fiatpeg4)

	issuFiats := []IssueFiat{issuFiat1, issuFiat2, issuFiat3}
	issuFiats2 := []IssueFiat{issuFiat4}

	msgfactoryissuefiat := NewMsgFactoryIssueFiats(issuFiats)
	msgfactoryissuefiat2 := NewMsgFactoryIssueFiats(issuFiats2)

	result := handleMsgFactoryIssueFiat(ctx, keeper, msgfactoryissuefiat)
	result2 := handleMsgFactoryIssueFiat(ctx, keeper, msgfactoryissuefiat2)

	require.Equal(t, result2.Tags, sdk.Tags(nil))
	require.Equal(t, addrs[0].String(), string(result.Tags[1].Value))
	require.Equal(t, addrs[1].String(), string(result.Tags[0].Value))
	require.Equal(t, addrs[2].String(), string(result.Tags[4].Value))
	require.Equal(t, addrs[3].String(), string(result.Tags[3].Value))
	require.Equal(t, addrs[4].String(), string(result.Tags[7].Value))
	require.Equal(t, addrs[5].String(), string(result.Tags[6].Value))

}

// TestHandleMsgFactorySendFiats handles test cases for message factory send fiats
func TestHandleMsgFactorySendFiats(t *testing.T) {
	ms, capkey, _ := setupMultiStore()
	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)

	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	keeper := NewKeeper(newfiatpegmapper)

	var testSendFiat = []SendFiat{
		SendFiat{
			RelayerAddress: sdk.AccAddress([]byte("relayer")),
			FromAddress:    sdk.AccAddress([]byte("from")),
			ToAddress:      sdk.AccAddress([]byte("to")),
			PegHash:        sdk.PegHash([]byte("pegHash")),
			FiatPegWallet: sdk.FiatPegWallet{
				sdk.BaseFiatPeg{
					PegHash:           sdk.PegHash([]byte("pegHash1")),
					TransactionID:     "FB8AE3A02BBCD2",
					TransactionAmount: 5000,
					RedeemedAmount:    0,
					Owners: []sdk.Owner{
						sdk.Owner{
							OwnerAddress: sdk.AccAddress([]byte("from")),
							Amount:       2000,
						},
						sdk.Owner{
							OwnerAddress: sdk.AccAddress([]byte("relayer")),
							Amount:       3000,
						},
					},
				},
			},
		},
	}
	var testSendFiat2 = []SendFiat{
		SendFiat{
			RelayerAddress: sdk.AccAddress([]byte("relayer1")),
			FromAddress:    sdk.AccAddress([]byte("from1")),
			ToAddress:      sdk.AccAddress([]byte("to1")),
			PegHash:        sdk.PegHash([]byte("sdds")),
			FiatPegWallet: sdk.FiatPegWallet{
				sdk.BaseFiatPeg{
					PegHash:           sdk.PegHash([]byte("pegHash122")),
					TransactionID:     "FB8AE3A02BBCDAD2",
					TransactionAmount: 1000,
					RedeemedAmount:    300,
					Owners:            nil,
				},
			},
		},
	}
	fiatPegWallet := testSendFiat[0].FiatPegWallet

	for _, fiat := range fiatPegWallet {
		newfiatpegmapper.SetFiatPeg(ctx, &fiat)
	}
	fiatPegWallet2 := testSendFiat2[0].FiatPegWallet

	for _, fiat2 := range fiatPegWallet2 {
		newfiatpegmapper.SetFiatPeg(ctx, &fiat2)
	}
	msgfactorysendfiat := NewMsgFactorySendFiats(testSendFiat)

	msgfactorysendfiat2 := NewMsgFactorySendFiats(testSendFiat2)
	result := handleMsgFactorySendFiats(ctx, keeper, msgfactorysendfiat)
	result2 := handleMsgFactorySendFiats(ctx, keeper, msgfactorysendfiat2)

	if true {
		require.Equal(t, testSendFiat[0].FromAddress.String(), string(result.Tags[1].Value))
		require.Equal(t, sdk.AccAddress((sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash))).String(), string(result.Tags[0].Value))
		require.Equal(t, result2.Tags, sdk.Tags(nil))
	}

}

// TestHandleMsgFactoryExecuteFiats handles test cases for message factory execute fiats
func TestHandleMsgFactoryExecuteFiats(t *testing.T) {

	ms, capkey, _ := setupMultiStore()
	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)

	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	keeper := NewKeeper(newfiatpegmapper)

	var testSendFiat = []SendFiat{
		SendFiat{
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
						sdk.Owner{
							OwnerAddress: sdk.AccAddress([]byte("from")),
							Amount:       2000,
						},
						sdk.Owner{
							OwnerAddress: sdk.AccAddress([]byte("relayer")),
							Amount:       3000,
						},
					},
				},
			},
		},
	}
	var testSendFiat2 = []SendFiat{
		SendFiat{
			RelayerAddress: sdk.AccAddress([]byte("relayer1")),
			FromAddress:    sdk.AccAddress([]byte("from1")),
			ToAddress:      sdk.AccAddress([]byte("to1")),
			PegHash:        sdk.PegHash([]byte("sdds")),
			FiatPegWallet: sdk.FiatPegWallet{
				sdk.BaseFiatPeg{
					PegHash:           sdk.PegHash([]byte("pegHash122")),
					TransactionID:     "FB8AE3A02BBCDAD2",
					TransactionAmount: 1000,
					RedeemedAmount:    300,
					Owners:            nil,
				},
			},
		},
	}
	fiatPegWallet := testSendFiat[0].FiatPegWallet
	testSendFiat[0].FiatPegWallet[0].Owners[0].OwnerAddress = sdk.AccAddress((sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash)))
	for _, fiat := range fiatPegWallet {
		newfiatpegmapper.SetFiatPeg(ctx, &fiat)
	}
	fiatPegWallet2 := testSendFiat2[0].FiatPegWallet

	for _, fiat2 := range fiatPegWallet2 {
		newfiatpegmapper.SetFiatPeg(ctx, &fiat2)
	}

	msgfactoryexecutefiat := NewMsgFactoryExecuteFiats(testSendFiat)
	msgfactoryexecutefiat2 := NewMsgFactoryExecuteFiats(testSendFiat2)

	result := handleMsgFactoryExecuteFiats(ctx, keeper, msgfactoryexecutefiat)
	result2 := handleMsgFactoryExecuteFiats(ctx, keeper, msgfactoryexecutefiat2)

	if true {

		require.Equal(t, testSendFiat[0].ToAddress.String(), string(result.Tags[0].Value))
		require.Equal(t, sdk.AccAddress((sdk.GenerateNegotiationIDBytes(testSendFiat[0].FromAddress, testSendFiat[0].ToAddress, testSendFiat[0].PegHash))).String(), string(result.Tags[1].Value))
		require.Equal(t, result2.Tags, sdk.Tags(nil))
	}

}

// TestNewHandler tests function NewHandler
func TestNewHandler(t *testing.T) {

	ms, capkey, _ := setupMultiStore()
	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)

	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	keeper := NewKeeper(newfiatpegmapper)

	issufiat := []IssueFiat{}
	sendfiat := []SendFiat{}
	_, addrs, _, _ := mock.CreateGenAccounts(1, sdk.Coins{})
	msg5 := sdk.NewTestMsg(addrs[0])
	msg := NewMsgFactoryIssueFiats(issufiat)
	msg2 := NewMsgFactorySendFiats(sendfiat)
	msg3 := NewMsgFactoryExecuteFiats(sendfiat)
	handler := NewHandler(keeper)

	result := handler(ctx, msg)
	result2 := handler(ctx, msg2)
	result3 := handler(ctx, msg3)
	result5 := handler(ctx, msg5)
	require.True(t, result.IsOK())
	require.Equal(t, result5.Tags, sdk.Tags(nil))
	require.True(t, result2.IsOK())
	require.True(t, result3.IsOK())
}
