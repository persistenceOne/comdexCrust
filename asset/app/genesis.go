package app

import (
	"encoding/json"
	"errors"
	"fmt"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/comdex-blockchain/server"
	"github.com/comdex-blockchain/server/config"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/stake"
	
	"github.com/spf13/pflag"
	
	"strconv"
	
	keyss "github.com/comdex-blockchain/client/keys"
	"github.com/spf13/viper"
	
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

// DefaultKeyPass :
const DefaultKeyPass = "1234567890"

var (
	flagName        = "name"
	flagClientHome  = "home-client"
	flagOWK         = "owk"
	freeFermionVal  = int64(100)
	freeFermionsAcc = sdk.NewInt(50)
)

// GenesisState to Unmarshal
type GenesisState struct {
	Accounts  []GenesisAccount   `json:"accounts"`
	Assets    []GenesisAssetPeg  `json:"assets"`
	StakeData stake.GenesisState `json:"stake"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewGenesisAccount returns GenesisAccount
func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
	}
}

// NewGenesisAccountI returns a new GenesisAccount
func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address: acc.GetAddress(),
		Coins:   acc.GetCoins(),
	}
}

// ToAccount convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
}

// GenesisAssetPeg : genesis state of assets
type GenesisAssetPeg struct {
	PegHash       sdk.PegHash    `json:"pegHash"`
	DocumentHash  string         `json:"documentHash"`
	AssetType     string         `json:"assetType"`
	AssetQuantity int64          `json:"assetQuantity"`
	QuantityUnit  string         `json:"quantityUnit"`
	OwnerAddress  sdk.AccAddress `json:"ownerAddress"`
}

// NewGenesisAssetPegI : returns a new genesis state account
func NewGenesisAssetPegI(assetPeg sdk.AssetPeg) GenesisAssetPeg {
	return GenesisAssetPeg{
		PegHash:       assetPeg.GetPegHash(),
		DocumentHash:  assetPeg.GetDocumentHash(),
		AssetType:     assetPeg.GetAssetType(),
		AssetQuantity: assetPeg.GetAssetQuantity(),
		QuantityUnit:  assetPeg.GetQuantityUnit(),
		OwnerAddress:  assetPeg.GetOwnerAddress(),
	}
}

// NewGenesisAssetPeg : returns a new genesis state account
func NewGenesisAssetPeg(assetPeg *sdk.BaseAssetPeg) GenesisAssetPeg {
	return GenesisAssetPeg{
		PegHash:       assetPeg.PegHash,
		DocumentHash:  assetPeg.DocumentHash,
		AssetType:     assetPeg.AssetType,
		AssetQuantity: assetPeg.AssetQuantity,
		QuantityUnit:  assetPeg.QuantityUnit,
		OwnerAddress:  assetPeg.OwnerAddress,
	}
}

// ToAssetPeg : convert GenesisAssetPeg to sdk.BaseAssetPeg
func (ga *GenesisAssetPeg) ToAssetPeg() (asset *sdk.BaseAssetPeg) {
	return &sdk.BaseAssetPeg{
		PegHash:       ga.PegHash,
		DocumentHash:  ga.DocumentHash,
		AssetType:     ga.AssetType,
		AssetQuantity: ga.AssetQuantity,
		QuantityUnit:  ga.QuantityUnit,
		OwnerAddress:  ga.OwnerAddress,
	}
}

// AssetAppInit : get app init parameters for server init command
func AssetAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)
	
	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagName, "", "validator moniker, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, true, "overwrite the accounts created")
	
	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         AssetAppGenTx,
		AppGenState:      AssetAppGenStateJSON,
	}
}

// AssetGenTx : simple genesis tx
type AssetGenTx struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
	PubKey  string         `json:"pub_key"`
}

// AssetAppGenTx generates a Asset genesis transaction.
func AssetAppGenTx(
	cdc *wire.Codec, pk crypto.PubKey, genTxConfig config.GenTx,
) (appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	
	if genTxConfig.Name == "" {
		return nil, nil, tmtypes.GenesisValidator{}, errors.New("Must specify --name (validator moniker)")
	}
	
	var addr sdk.AccAddress
	var kb keys.Keybase
	kb, err = keyss.GetKeyBaseFromDir(genTxConfig.CliRoot)
	if err != nil {
		return appGenTx, cliPrint, validator, err
	}
	
	secret := viper.GetString(client.FlagSeed)
	
	if secret == "" {
		addr, secret, err = server.GenerateSaveCoinKey(
			genTxConfig.CliRoot,
			genTxConfig.Name,
			DefaultKeyPass,
			genTxConfig.Overwrite,
		)
		if err != nil {
			return appGenTx, cliPrint, validator, err
		}
	} else {
		info, err := kb.CreateKey(genTxConfig.Name, secret, DefaultKeyPass)
		addr = sdk.AccAddress(info.GetPubKey().Address())
		
		if err != nil {
			return appGenTx, cliPrint, validator, err
		}
	}
	mm := map[string]string{"secret": secret}
	bz, err := cdc.MarshalJSON(mm)
	if err != nil {
		return appGenTx, cliPrint, validator, err
	}
	
	cliPrint = json.RawMessage(bz)
	appGenTx, _, validator, err = AssetAppGenTxNF(cdc, pk, addr, genTxConfig.Name)
	
	return appGenTx, cliPrint, validator, err
}

// AssetAppGenTxNF : Generate a asset genesis transaction without flags
func AssetAppGenTxNF(cdc *wire.Codec, pk crypto.PubKey, addr sdk.AccAddress, name string) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	
	var bz []byte
	assetGenTx := AssetGenTx{
		Name:    name,
		Address: addr,
		PubKey:  sdk.MustBech32ifyConsPub(pk),
	}
	bz, err = wire.MarshalJSONIndent(cdc, assetGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)
	
	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  freeFermionVal,
	}
	return
}

// AssetAppGenState : Create the core parameters for genesis initialization for main
// note that the pubkey input is this machines pubkey
func AssetAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {
	
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}
	
	stakeData := stake.DefaultGenesisState()
	
	genaccs := make([]GenesisAccount, len(appGenTxs))
	var genesisAssetPegs []GenesisAssetPeg
	for i, appGenTx := range appGenTxs {
		
		var genTx AssetGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}
		
		accAuth := auth.NewBaseAccountWithAddress(genTx.Address)
		accAuth.Coins = sdk.Coins{
			{"comdex", sdk.NewInt(1000)},
			{"steak", freeFermionsAcc},
		}
		
		acc := NewGenesisAccount(&accAuth)
		genaccs[i] = acc
		stakeData.Pool.LooseTokens = stakeData.Pool.LooseTokens.Add(sdk.NewDecFromInt(freeFermionsAcc)) // increase the supply
		
		if len(genTx.Name) > 0 {
			desc := stake.NewDescription(genTx.Name, "", "", "")
			validator := stake.NewValidator(
				sdk.ValAddress(genTx.Address), sdk.MustGetConsPubKeyBech32(genTx.PubKey), desc,
			)
			
			stakeData.Pool.LooseTokens = stakeData.Pool.LooseTokens.Add(sdk.NewDec(freeFermionVal)) // increase the supply
			
			var issuedDelShares sdk.Dec
			validator, stakeData.Pool, issuedDelShares = validator.AddTokensFromDel(stakeData.Pool, sdk.NewInt(freeFermionVal))
			stakeData.Validators = append(stakeData.Validators, validator)
			
			delegation := stake.Delegation{
				DelegatorAddr: sdk.AccAddress(validator.Operator),
				ValidatorAddr: validator.Operator,
				Shares:        issuedDelShares,
				Height:        0,
			}
			
			stakeData.Bonds = append(stakeData.Bonds, delegation)
		}
		
		for i := 0; i < 10000; i++ {
			pegHash, err := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
			if err == nil {
				genesisAssetPegs = append(genesisAssetPegs, NewGenesisAssetPeg(&sdk.BaseAssetPeg{
					PegHash:      pegHash,
					OwnerAddress: acc.Address,
				}))
			}
		}
	}
	
	genesisState = GenesisState{
		Accounts:  genaccs,
		Assets:    genesisAssetPegs,
		StakeData: stakeData,
	}
	return
}

// AssetAppGenStateJSON  but with JSON
func AssetAppGenStateJSON(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	
	genesisState, err := AssetAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	
	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	if err != nil {
		return nil, err
	}
	return
}
