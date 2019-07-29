package fiatFactory

import (
	"github.com/commitHub/commitBlockchain/store"
	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/wire"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setupMultiStore1() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	fiatKey := sdk.NewKVStoreKey("fiatkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(fiatKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, fiatKey
}

func initiateSetupMultiStore() (*wire.Codec, sdk.Context, FiatPegMapper, Keeper) {
	ms, fiatKey := setupMultiStore1()

	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	fiatMapper := NewFiatPegMapper(cdc, fiatKey, sdk.ProtoBaseFiatPeg)
	fiatKeeper := NewKeeper(fiatMapper)
	return cdc, ctx, fiatMapper, fiatKeeper
}

var issueFiat = []IssueFiat{
	IssueFiat{
		IssuerAddress: sdk.AccAddress([]byte("issuer")),
		ToAddress:     sdk.AccAddress([]byte("to")),
		FiatPeg: &sdk.BaseFiatPeg{
			PegHash:           sdk.PegHash([]byte("pegHash")),
			TransactionID:     "FH8GH3V02HNJG2",
			TransactionAmount: 9,
			RedeemedAmount:    3,
			Owners: []sdk.Owner{
				sdk.Owner{
					OwnerAddress: sdk.AccAddress([]byte("issuer")),
					Amount:       2000,
				},
				sdk.Owner{
					OwnerAddress: sdk.AccAddress([]byte("to")),
					Amount:       3000,
				},
			},
		},
	},
}

var testRedeemFiat = []RedeemFiat{
	RedeemFiat{
		RelayerAddress:  sdk.AccAddress([]byte("relayer")),
		RedeemerAddress: sdk.AccAddress([]byte("redeemer")),
		Amount:          60,
		FiatPegWallet: sdk.FiatPegWallet{
			sdk.BaseFiatPeg{
				PegHash:           sdk.PegHash([]byte("pegHash1")),
				TransactionID:     "FB8AE3A02BBCD2",
				TransactionAmount: 1000,
				RedeemedAmount:    0,
				Owners: []sdk.Owner{
					sdk.Owner{
						OwnerAddress: sdk.AccAddress([]byte("from")),
						Amount:       500,
					},
					sdk.Owner{
						OwnerAddress: sdk.AccAddress([]byte("relayer")),
						Amount:       500,
					},
				},
			},
		},
	},
}

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
