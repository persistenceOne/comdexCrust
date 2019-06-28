package assetFactory

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

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("name")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, capKey
}

type testStruct = struct {
	key            sdk.StoreKey
	proto          func() sdk.AssetPeg
	cdc            *wire.Codec
	assetPeg       sdk.AssetPeg
	expectedResult bool
}

var assetPegs = []sdk.BaseAssetPeg{
	sdk.NewBaseAssetPegWithPegHash(sdk.PegHash([]byte("PegHash0"))),
	sdk.NewBaseAssetPegWithPegHash(nil),
	{},
	sdk.NewBaseAssetPegWithPegHash(sdk.PegHash([]byte("PegHash1"))),
	sdk.NewBaseAssetPegWithPegHash(sdk.PegHash([]byte("PegHash2"))),
}

var listTestsMapper = []testStruct{
	{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
	{sdk.NewKVStoreKey(""), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
	{nil, sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
	{sdk.NewKVStoreKey("name"), nil, wire.NewCodec(), &assetPegs[0], true},
	{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[1], true},
	{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[2], true},
}

func genAssetPegMapper(testCase testStruct) AssetPegMapper {
	return AssetPegMapper{
		key:   testCase.key,
		proto: testCase.proto,
		cdc:   testCase.cdc,
	}
}

func TestNewAssetPegMapper(t *testing.T) {
	for _, testCase := range listTestsMapper[:] {
		oneAssetPegMapper := genAssetPegMapper(testCase)
		toTest := NewAssetPegMapper(testCase.cdc, testCase.key, testCase.proto)
		if testCase.expectedResult {
			require.Equal(t, oneAssetPegMapper.key, toTest.key)
			require.Equal(t, oneAssetPegMapper.cdc, toTest.cdc)
			
		} else {
			require.NotEqual(t, oneAssetPegMapper.key, toTest.key)
			require.NotEqual(t, oneAssetPegMapper.cdc, toTest.cdc)
			require.NotEqual(t, oneAssetPegMapper.proto(), toTest.proto())
		}
	}
}

func TestAssetPegHashStoreKey(t *testing.T) {
	onePegHashStoreKey := AssetPegHashStoreKey(sdk.PegHash([]byte("pegHash")))
	twoPegHashStoreKey := AssetPegHashStoreKey(sdk.PegHash([]byte("")))
	threePegHashStoreKey := AssetPegHashStoreKey(nil)
	t.Logf("%v,%v", threePegHashStoreKey[7:], []byte("PegHash:"))
	require.Equal(t, onePegHashStoreKey, []byte("PegHash:pegHash"))
	require.Equal(t, twoPegHashStoreKey, []byte("PegHash:"))
	require.Equal(t, threePegHashStoreKey[8:], []byte(""))
	
}

func TestAssetPegMapperencodeAssetPeg(t *testing.T) {
	var listTestsMapperCustom = []testStruct{
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey(""), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, nil, &assetPegs[0], false},
		{nil, sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), nil, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[1], true},
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[2], true},
	}
	
	for _, testCase := range listTestsMapperCustom {
		if testCase.expectedResult {
			oneAssetPegMapper := genAssetPegMapper(testCase)
			RegisterAssetPeg(oneAssetPegMapper.cdc)
			
			toTest := oneAssetPegMapper.encodeAssetPeg(testCase.assetPeg)
			testByte, _ := oneAssetPegMapper.cdc.MarshalBinaryBare(testCase.assetPeg)
			
			require.Equal(t, toTest, testByte)
			
		} else {
			defer func() {
				if r := recover(); r == nil {
					t.Logf("The code did not panic %v \n", "no")
				} else {
					t.Logf("There had been an error")
				}
			}()
			
			oneAssetPegMapper := genAssetPegMapper(testCase)
			RegisterAssetPeg(oneAssetPegMapper.cdc)
			
			toTest := oneAssetPegMapper.encodeAssetPeg(testCase.assetPeg)
			testByte, _ := oneAssetPegMapper.cdc.MarshalBinaryBare(testCase.assetPeg)
			
			require.NotEqual(t, toTest, testByte)
		}
	}
}

func TestAssetPegMapperdecodeAssetPeg(t *testing.T) {
	var listTestsMapperCustom = []testStruct{
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey(""), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, nil, &assetPegs[0], false},
		{nil, sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), nil, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[1], true},
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[2], true},
	}
	
	for _, testCase := range listTestsMapperCustom {
		
		if testCase.expectedResult {
			oneAssetPegMapper := genAssetPegMapper(testCase)
			RegisterAssetPeg(oneAssetPegMapper.cdc)
			testBytes, err := oneAssetPegMapper.cdc.MarshalBinaryBare(&assetPeg)
			if err != nil {
				t.Logf("%v,test bytes failed", err)
			}
			toTest := oneAssetPegMapper.decodeAssetPeg(testBytes)
			require.Equal(t, toTest, &assetPeg)
			
		} else {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			oneAssetPegMapper := genAssetPegMapper(testCase)
			RegisterAssetPeg(oneAssetPegMapper.cdc)
			testBytes, err := oneAssetPegMapper.cdc.MarshalBinaryBare(&assetPeg)
			if err != nil {
				t.Logf("%v,test bytes failed", err)
			}
			toTest := oneAssetPegMapper.decodeAssetPeg(testBytes)
			require.NotEqual(t, toTest, &assetPeg)
		}
	}
}

func TestAssetPegMapperSetGetAssetPeg(t *testing.T) {
	var listTestsMapperCustom = []testStruct{
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey(""), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{nil, sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), nil, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[1], true},
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[2], true},
	}
	for _, testCase := range listTestsMapperCustom {
		ms, key := setupMultiStore()
		RegisterAssetPeg(testCase.cdc)
		
		ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
		oneAssetPegMapper := NewAssetPegMapper(testCase.cdc, key, testCase.proto)
		oneAssetPegMapper.SetAssetPeg(ctx, testCase.assetPeg)
		toTest := oneAssetPegMapper.GetAssetPeg(ctx, testCase.assetPeg.GetPegHash())
		if testCase.expectedResult {
			require.Equal(t, toTest, testCase.assetPeg)
		} else {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			require.NotEqual(t, toTest, testCase.assetPeg)
		}
	}
}

/*
func TestAssetPegMapperGetAssetPeg(t *testing.T){} //implemented with Set only
*/

func processTrue(sdk.AssetPeg) (stop bool) {
	return true
}

func TestAssetPegMapperIterateAssets(t *testing.T) {
	var listTestsMapperCustom = []testStruct{
		{sdk.NewKVStoreKey("name"), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey(""), sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{nil, sdk.ProtoBaseAssetPeg, wire.NewCodec(), &assetPegs[0], true},
		{sdk.NewKVStoreKey("name"), nil, wire.NewCodec(), &assetPegs[0], true},
	}
	
	for _, testCase := range listTestsMapperCustom {
		ms, key := setupMultiStore()
		RegisterAssetPeg(testCase.cdc)
		ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
		oneAssetPegMapper := NewAssetPegMapper(testCase.cdc, key, testCase.proto)
		oneAssetPegMapper.SetAssetPeg(ctx, testCase.assetPeg)
		oneAssetPegMapper.SetAssetPeg(ctx, &assetPegs[3])
		oneAssetPegMapper.SetAssetPeg(ctx, &assetPegs[4])
		var pegs []string
		someProcess := func(asset sdk.AssetPeg) (stop bool) {
			some := asset.GetPegHash()
			pegs = append(pegs, string(some))
			return false
		}
		oneAssetPegMapper.IterateAssets(ctx, processTrue)
		oneAssetPegMapper.IterateAssets(ctx, someProcess)
		require.Equal(t, pegs, []string{"PegHash0", "PegHash1", "PegHash2"})
	}
	
}
