package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	
	"github.com/comdex-blockchain/client"
	keyss "github.com/comdex-blockchain/client/keys"
	"github.com/comdex-blockchain/server"
	"github.com/comdex-blockchain/server/config"
	sdk "github.com/comdex-blockchain/types"
	"github.com/comdex-blockchain/wire"
	"github.com/comdex-blockchain/x/auth"
	"github.com/comdex-blockchain/x/stake"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
)

// DefaultKeyPass : default password for genesis account
const DefaultKeyPass = "1234567890"

var (
	flagName       = "name"
	flagClientHome = "home-client"
	flagOWK        = "owk"
	
	// bonded tokens given to genesis validators/accounts
	freeFermionVal  = int64(100)
	freeFermionsAcc = sdk.NewInt(50)
)

// GenesisState : State to Unmarshal
type GenesisState struct {
	Accounts  []GenesisAccount   `json:"accounts"`
	Fiats     []GenesisFiatPeg   `json:"fiats"`
	StakeData stake.GenesisState `json:"stake"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.AccAddress `json:"address"`
	Coins   sdk.Coins      `json:"coins"`
}

// NewGenesisAccount : returns a new genesis state account
func NewGenesisAccount(acc *auth.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
	}
}

// NewGenesisAccountI : new genesis account from already existing account
func NewGenesisAccountI(acc auth.Account) GenesisAccount {
	return GenesisAccount{
		Address: acc.GetAddress(),
		Coins:   acc.GetCoins(),
	}
}

// ToAccount : convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount) {
	return &auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
}

// GenesisFiatPeg : genesis state of fiats
type GenesisFiatPeg struct {
	PegHash           sdk.PegHash `json:"pegHash"`
	TransactionID     string      `json:"transactionID"`
	TransactionAmount int64       `json:"transactionAmount"`
	RedeemedAmount    int64       `json:"redeemedAmount"`
	Owners            []sdk.Owner `json:"owner"`
}

// NewGenesisFiatPegI : returns a new genesis state account
func NewGenesisFiatPegI(fiatPeg sdk.FiatPeg) GenesisFiatPeg {
	return GenesisFiatPeg{
		PegHash:           fiatPeg.GetPegHash(),
		TransactionID:     fiatPeg.GetTransactionID(),
		TransactionAmount: fiatPeg.GetTransactionAmount(),
		RedeemedAmount:    fiatPeg.GetRedeemedAmount(),
		Owners:            fiatPeg.GetOwners(),
	}
}

// NewGenesisFiatPeg : returns a new genesis state account
func NewGenesisFiatPeg(fiatPeg *sdk.BaseFiatPeg) GenesisFiatPeg {
	return GenesisFiatPeg{
		PegHash:           fiatPeg.PegHash,
		TransactionID:     fiatPeg.TransactionID,
		TransactionAmount: fiatPeg.TransactionAmount,
		RedeemedAmount:    fiatPeg.RedeemedAmount,
		Owners:            fiatPeg.Owners,
	}
}

// ToFiatPeg : convert GenesisFiatPeg to sdk.BaseFiatPeg
func (ga *GenesisFiatPeg) ToFiatPeg() (fiat *sdk.BaseFiatPeg) {
	return &sdk.BaseFiatPeg{
		PegHash:           ga.PegHash,
		TransactionID:     ga.TransactionID,
		TransactionAmount: ga.TransactionAmount,
		RedeemedAmount:    ga.RedeemedAmount,
		Owners:            ga.Owners,
	}
}

// FiatAppInit : get app init parameters for server init command
func FiatAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)
	
	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagName, "", "validator moniker, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, false, "overwrite the accounts created")
	
	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         FiatAppGenTx,
		AppGenState:      FiatAppGenStateJSON,
	}
}

// FiatGenTx : simple genesis tx
type FiatGenTx struct {
	Name    string         `json:"name"`
	Address sdk.AccAddress `json:"address"`
	PubKey  string         `json:"pub_key"`
}

// FiatAppGenTx generates a Gaia genesis transaction.
func FiatAppGenTx(
	cdc *wire.Codec, pk crypto.PubKey, genTxConfig config.GenTx,
) (appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	
	if genTxConfig.Name == "" {
		return nil, nil, tmtypes.GenesisValidator{}, errors.New("Must specify --name (validator moniker)")
	}
	
	secret := viper.GetString(client.FlagSeed)
	var addr sdk.AccAddress
	kb, err := keyss.GetKeyBaseFromDir(genTxConfig.CliRoot)
	if err != nil {
		return appGenTx, cliPrint, validator, err
	}
	
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
	appGenTx, _, validator, err = FiatAppGenTxNF(cdc, pk, addr, genTxConfig.Name)
	
	return appGenTx, cliPrint, validator, err
}

// FiatAppGenTxNF : Generate a  fiat genesis transaction without flags
func FiatAppGenTxNF(cdc *wire.Codec, pk crypto.PubKey, addr sdk.AccAddress, name string) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	
	var bz []byte
	fiatGenTx := FiatGenTx{
		Name:    name,
		Address: addr,
		PubKey:  sdk.MustBech32ifyConsPub(pk),
	}
	bz, err = wire.MarshalJSONIndent(cdc, fiatGenTx)
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

// FiatAppGenState : Create the core parameters for genesis initialization for  fiat
// note that the pubkey input is this machines pubkey
func FiatAppGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {
	
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}
	
	stakeData := stake.DefaultGenesisState()
	
	genaccs := make([]GenesisAccount, len(appGenTxs))
	for i, appGenTx := range appGenTxs {
		
		var genTx FiatGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}
		
		accAuth := auth.NewBaseAccountWithAddress(genTx.Address)
		
		accAuth.Coins = sdk.Coins{
			{
				Denom:  "comdex",
				Amount: sdk.NewInt(1000),
			},
			{
				Denom:  "steak",
				Amount: freeFermionsAcc,
			},
		}
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
	
	// Generate empty fiat tokens
	genesisFiatPegs := make([]GenesisFiatPeg, 10000)
	for i := 0; i < 10000; i++ {
		pegHash, err := sdk.GetFiatPegHashHex(fmt.Sprintf("%x", strconv.Itoa(i)))
		if err == nil {
			genesisFiatPegs[i] = NewGenesisFiatPeg(&sdk.BaseFiatPeg{
				PegHash: pegHash,
			})
		}
	}
	
	// create the final app state
	genesisState = GenesisState{
		Accounts:  genaccs,
		Fiats:     genesisFiatPegs,
		StakeData: stakeData,
	}
	return
}

// FiatAppGenStateJSON : FiatAppGenState but with JSON
func FiatAppGenStateJSON(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	
	// create the final app state
	genesisState, err := FiatAppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}
