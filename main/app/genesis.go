package app

import (
	"encoding/json"
	"errors"
	
	"github.com/comdex-blockchain/client"
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/comdex-blockchain/server"
	"github.com/comdex-blockchain/server/config"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/gov"
	"github.com/comdex-blockchain/x/stake"
	
	keyss "github.com/comdex-blockchain/client/keys"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

// DefaultKeyPass contains the default key password for genesis transactions
const DefaultKeyPass = "1234567890"

var (
	flagName       = "name"
	flagClientHome = "home-client"
	flagOWK        = "owk"
	
	// bonded tokens given to genesis validators/accounts
	freeFermionVal  = int64(100)
	freeFermionsAcc = sdk.NewInt(50)
)

// GenesisState to Unmarshal
type GenesisState struct {
	Accounts  []GenesisAccount   `json:"accounts"`
	StakeData stake.GenesisState `json:"stake"`
	GovData   gov.GenesisState   `json:"gov"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address        sdk.AccAddress     `json:"address"`
	Coins          sdk.Coins          `json:"coins"`
	FiatPegWallet  sdk.FiatPegWallet  `json:"fiatPegWallet"`
	AssetPegWallet sdk.AssetPegWallet `json:"assetPegWallet"`
}

// NewGenesisAccount : returns GenesisAccount
func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:        acc.Address,
		Coins:          acc.Coins,
		FiatPegWallet:  acc.FiatPegWallet,
		AssetPegWallet: acc.AssetPegWallet,
	}
}

// NewGenesisAccountI : returns GenesisAccount
func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address:        acc.GetAddress(),
		Coins:          acc.GetCoins(),
		FiatPegWallet:  acc.GetFiatPegWallet(),
		AssetPegWallet: acc.GetAssetPegWallet(),
	}
}

// ToAccount : convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address:        ga.Address,
		Coins:          ga.Coins.Sort(),
		FiatPegWallet:  ga.FiatPegWallet,
		AssetPegWallet: ga.AssetPegWallet,
	}
}

// MainAppInit : get app init parameters for server init command
func MainAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)
	
	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagName, "", "validator moniker, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, false, "overwrite the accounts created")
	
	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         MainAppGenTx,
		AppGenState:      MainAppGenStateJSON,
	}
}

// MainGenTx : simple genesis tx
type MainGenTx struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
	PubKey  string         `json:"pub_key"`
}

// MainAppGenTx generates a main genesis transaction.
func MainAppGenTx(
	cdc *wire.Codec, pk crypto.PubKey, genTxConfig config.GenTx,
) (appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	if genTxConfig.Name == "" {
		return nil, nil, tmtypes.GenesisValidator{}, errors.New("Must specify --name (validator moniker)")
	}
	
	var addr sdk.AccAddress
	var kb keys.Keybase
	secret := viper.GetString(client.FlagSeed)
	kb, err = keyss.GetKeyBaseFromDir(genTxConfig.CliRoot)
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
	appGenTx, _, validator, err = MainAppGenTxNF(cdc, pk, addr, genTxConfig.Name)
	
	return appGenTx, cliPrint, validator, err
}

// MainAppGenTxNF : Generate a main genesis transaction without flags
func MainAppGenTxNF(cdc *wire.Codec, pk crypto.PubKey, addr sdk.AccAddress, name string) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	
	var bz []byte
	mainGenTx := MainGenTx{
		Name:    name,
		Address: addr,
		PubKey:  sdk.MustBech32ifyConsPub(pk),
	}
	bz, err = wire.MarshalJSONIndent(cdc, mainGenTx)
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

// MainAppGenState : Create the core parameters for genesis initialization for main
// note that the pubkey input is this machines pubkey
func MainAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {
	
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}
	
	// start with the default staking genesis state
	stakeData := stake.DefaultGenesisState()
	
	// get genesis flag account information
	genaccs := make([]GenesisAccount, len(appGenTxs))
	for i, appGenTx := range appGenTxs {
		
		var genTx MainGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}
		
		// create the genesis account, give'm few steaks and a buncha token with there name
		accAuth := auth.NewBaseAccountWithAddress(genTx.Address)
		accAuth.Coins = sdk.Coins{
			{"comdex", sdk.NewInt(1000)},
			{"steak", freeFermionsAcc},
		}
		
		// //TODO: comdex Use IBC to get all asset hashes
		// for i := 0; i < 1024; i++ {
		// 	pegHash, err := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
		// 	if err == nil {
		// 		accAuth.AssetPegWallet = append(accAuth.AssetPegWallet, sdk.NewBaseAssetPegWithPegHash(pegHash))
		// 		accAuth.FiatPegWallet = append(accAuth.FiatPegWallet, sdk.NewBaseFiatPegWithPegHash(pegHash))
		// 	}
		// }
		acc := NewGenesisAccount(&accAuth)
		genaccs[i] = acc
		stakeData.Pool.LooseTokens = stakeData.Pool.LooseTokens.Add(sdk.NewDecFromInt(freeFermionsAcc)) // increase the supply
		
		// add the validator
		if len(genTx.Name) > 0 {
			desc := stake.NewDescription(genTx.Name, "", "", "")
			validator := stake.NewValidator(
				sdk.ValAddress(genTx.Address), sdk.MustGetConsPubKeyBech32(genTx.PubKey), desc,
			)
			
			stakeData.Pool.LooseTokens = stakeData.Pool.LooseTokens.Add(sdk.NewDec(freeFermionVal)) // increase the supply
			
			// add some new shares to the validator
			var issuedDelShares sdk.Dec
			validator, stakeData.Pool, issuedDelShares = validator.AddTokensFromDel(stakeData.Pool, sdk.NewInt(freeFermionVal))
			stakeData.Validators = append(stakeData.Validators, validator)
			
			// create the self-delegation from the issuedDelShares
			delegation := stake.Delegation{
				DelegatorAddr: sdk.AccAddress(validator.Operator),
				ValidatorAddr: validator.Operator,
				Shares:        issuedDelShares,
				Height:        0,
			}
			
			stakeData.Bonds = append(stakeData.Bonds, delegation)
		}
	}
	
	// create the final app state
	genesisState = GenesisState{
		Accounts:  genaccs,
		StakeData: stakeData,
		GovData:   gov.DefaultGenesisState(),
	}
	return
}

// MainAppGenStateJSON but with JSON
func MainAppGenStateJSON(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	
	// create the final app state
	genesisState, err := MainAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}
