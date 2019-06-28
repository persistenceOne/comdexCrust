package fiatFactory

import (
	"testing"
	
	"github.com/comdex-blockchain/store"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

// TestFiatPegMapper
func TestFiatPegMapper(t *testing.T) {
	
	peghash1 := sdk.PegHash([]byte("132454"))
	peghash2 := sdk.PegHash("abcdef")
	tests := []struct {
		Peghash   sdk.PegHash
		entrypass bool
	}{
		{peghash1, true},
		{peghash2, true},
	}
	
	for _, tc := range tests {
		
		bite := FiatPegHashStoreKey(tc.Peghash)
		bite2 := bite[8:]
		
		emptybite := []byte("")
		if !tc.entrypass {
			require.Equal(t, bite2, emptybite)
		} else {
			require.NotEqual(t, bite2, emptybite)
		}
	}
	
}

func TestEncodeFiatPeg(t *testing.T) {
	
	capkey := sdk.NewKVStoreKey("capkey")
	cdc := wire.NewCodec()
	
	BfietWpeghash1 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("helloworld")))
	BfietWpeghash2 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("1234567")))
	fiatpeg1 := sdk.ToFiatPeg(BfietWpeghash1)
	fiatpeg2 := sdk.ToFiatPeg(BfietWpeghash2)
	
	tests2 := []struct {
		Fiatpeg   sdk.FiatPeg
		entrypass bool
	}{
		{fiatpeg1, true},
		{fiatpeg2, true},
	}
	
	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)
	
	for _, tc := range tests2 {
		
		bite := newfiatpegmapper.encodeFiatPeg(tc.Fiatpeg)
		var Bfiet sdk.BaseFiatPeg
		err := newfiatpegmapper.cdc.UnmarshalBinaryBare(bite, &Bfiet)
		fiet := sdk.ToFiatPeg(Bfiet)
		
		if err != nil {
			panic(err)
		}
		
		if tc.entrypass {
			require.Equal(t, fiet, tc.Fiatpeg)
		} else {
			require.NotEqual(t, fiet, tc.Fiatpeg)
		}
		
	}
	
}

func TestDecodeFiatPeg(t *testing.T) {
	
	capkey := sdk.NewKVStoreKey("capkeyy")
	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)
	
	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)
	
	BfietWpeghash1 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("gvnkjn")))
	BfietWpeghash2 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("1234567")))
	fiatpeg1 := sdk.ToFiatPeg(BfietWpeghash1)
	fiatpeg2 := sdk.ToFiatPeg(BfietWpeghash2)
	bite1, err := newfiatpegmapper.cdc.MarshalBinaryBare(fiatpeg1)
	if err != nil {
		panic(err)
	}
	bite2, err := newfiatpegmapper.cdc.MarshalBinaryBare(fiatpeg2)
	if err != nil {
		panic(err)
	}
	
	tests := []struct {
		bite      []byte
		entrypass bool
	}{
		{bite1, true},
		{bite2, true},
	}
	
	for _, tc := range tests {
		
		testfiet := newfiatpegmapper.decodeFiatPeg(tc.bite)
		
		testbite, err := newfiatpegmapper.cdc.MarshalBinaryBare(testfiet)
		if err != nil {
			
			panic(err)
		}
		
		if tc.entrypass {
			require.Equal(t, tc.bite, testbite)
		}
		
	}
	
}

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	capKey2 := sdk.NewKVStoreKey("capkey2")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(capKey2, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, capKey, capKey2
}

func TestSetFiatPeg(t *testing.T) {
	
	ms, capkey, _ := setupMultiStore()
	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)
	
	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)
	
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	
	BfietWpeghash1 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("gvnkjn")))
	BfietWpeghash2 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("1234567")))
	fiatpeg1 := sdk.ToFiatPeg(BfietWpeghash1)
	fiatpeg2 := sdk.ToFiatPeg(BfietWpeghash2)
	
	tests := []struct {
		conText   sdk.Context
		fiatpeg   sdk.FiatPeg
		entrypass bool
	}{
		
		{ctx, fiatpeg1, true},
		{ctx, fiatpeg2, true},
	}
	
	for _, tc := range tests {
		
		newfiatpegmapper.SetFiatPeg(tc.conText, tc.fiatpeg)
		// a := newfiatpegmapper.key
		store := tc.conText.KVStore(newfiatpegmapper.key)
		peghash := tc.fiatpeg.GetPegHash()
		
		TestBite := store.Get(FiatPegHashStoreKey(peghash))
		
		reqbite := newfiatpegmapper.encodeFiatPeg(tc.fiatpeg)
		
		if tc.entrypass {
			
			require.Equal(t, reqbite, TestBite)
		}
		
	}
}

func TestGetFiatPeg(t *testing.T) {
	
	ms, capkey, _ := setupMultiStore()
	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)
	
	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)
	
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	
	peghash1 := sdk.PegHash([]byte("gvnkjn"))
	peghash2 := sdk.PegHash([]byte("123456"))
	
	tests := []struct {
		conText   sdk.Context
		peghash   sdk.PegHash
		entrypass bool
	}{
		
		{ctx, peghash1, true},
		{ctx, peghash2, true},
	}
	
	for _, tc := range tests {
		
		store := ctx.KVStore(newfiatpegmapper.key)
		bite := newfiatpegmapper.encodeFiatPeg(sdk.ToFiatPeg(sdk.NewBaseFiatPegWithPegHash(tc.peghash)))
		store.Set(FiatPegHashStoreKey(tc.peghash), bite)
		
		fiet := newfiatpegmapper.GetFiatPeg(ctx, tc.peghash)
		reqfiet := sdk.ToFiatPeg(sdk.NewBaseFiatPegWithPegHash(tc.peghash))
		if tc.entrypass {
			require.Equal(t, reqfiet, fiet)
		} else {
			require.NotEqual(t, reqfiet, fiet)
		}
		
	}
}

func process1(sdk.FiatPeg) (stop bool) {
	return true
}

func process2(sdk.FiatPeg) (stop bool) {
	
	return false
}

func process3(fiet sdk.FiatPeg) (stop bool) {
	
	if fiet == nil {
		return true
	}
	return false
	
}

func TestIterateFiats(t *testing.T) {
	
	ms, capkey, _ := setupMultiStore()
	cdc := wire.NewCodec()
	RegisterFiatPeg(cdc)
	
	newfiatpegmapper := NewFiatPegMapper(cdc, capkey, sdk.ProtoBaseFiatPeg)
	
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	
	BfietWpeghash1 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("gvnkjn")))
	BfietWpeghash2 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("juygfcv")))
	BfietWpeghash3 := sdk.NewBaseFiatPegWithPegHash(sdk.PegHash([]byte("34567890")))
	fiatpeg1 := sdk.ToFiatPeg(BfietWpeghash1)
	fiatpeg2 := sdk.ToFiatPeg(BfietWpeghash2)
	fiatpeg3 := sdk.ToFiatPeg(BfietWpeghash3)
	
	newfiatpegmapper.SetFiatPeg(ctx, fiatpeg1)
	newfiatpegmapper.SetFiatPeg(ctx, fiatpeg2)
	newfiatpegmapper.SetFiatPeg(ctx, fiatpeg3)
	
	newfiatpegmapper.IterateFiats(ctx, process1)
	
	newfiatpegmapper.IterateFiats(ctx, process2)
	newfiatpegmapper.IterateFiats(ctx, process3)
	
}
