package bank

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/commitHub/commitBlockchain/x/reputation"

	sdk "github.com/commitHub/commitBlockchain/types"
	"github.com/commitHub/commitBlockchain/x/acl"
	"github.com/commitHub/commitBlockchain/x/auth"
	"github.com/commitHub/commitBlockchain/x/negotiation"
	"github.com/commitHub/commitBlockchain/x/order"
)

const (
	costGetCoins      sdk.Gas = 10
	costHasCoins      sdk.Gas = 10
	costSetCoins      sdk.Gas = 100
	costSubtractCoins sdk.Gas = 10
	costAddCoins      sdk.Gas = 10
)

// Keeper manages transfers between accounts
type Keeper struct {
	am auth.AccountMapper
}

// NewKeeper returns a new Keeper
func NewKeeper(am auth.AccountMapper) Keeper {
	return Keeper{am: am}
}

//CheckGenesisAccount : check if the address belongs to the genesis account
func (keeper Keeper) CheckGenesisAccount(ctx sdk.Context, accAddress sdk.AccAddress) bool {
	account := keeper.am.GetAccount(ctx, accAddress)
	if account.GetAccountNumber() == 0 {
		return true
	}
	return false
}

// GetCoins returns the coins at the addr.
func (keeper Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// SetCoins sets the coins at the addr.
func (keeper Keeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	return setCoins(ctx, keeper.am, addr, amt)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper Keeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	return subtractCoins(ctx, keeper.am, addr, amt)
}

// AddCoins adds amt to the coins at the addr.
func (keeper Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	return addCoins(ctx, keeper.am, addr, amt)
}

// SendCoins moves coins from one account to another
func (keeper Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error) {
	return sendCoins(ctx, keeper.am, fromAddr, toAddr, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper Keeper) InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error) {
	return inputOutputCoins(ctx, keeper.am, inputs, outputs)
}

//______________________________________________________________________________________________

// SendKeeper only allows transfers between accounts, without the possibility of creating coins
type SendKeeper struct {
	am auth.AccountMapper
}

// NewSendKeeper returns a new Keeper
func NewSendKeeper(am auth.AccountMapper) SendKeeper {
	return SendKeeper{am: am}
}

// GetCoins returns the coins at the addr.
func (keeper SendKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper SendKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

// SendCoins moves coins from one account to another
func (keeper SendKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error) {
	return sendCoins(ctx, keeper.am, fromAddr, toAddr, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper SendKeeper) InputOutputCoins(ctx sdk.Context, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error) {
	return inputOutputCoins(ctx, keeper.am, inputs, outputs)
}

//______________________________________________________________________________________________

// ViewKeeper only allows reading of balances
type ViewKeeper struct {
	am auth.AccountMapper
}

// NewViewKeeper returns a new Keeper
func NewViewKeeper(am auth.AccountMapper) ViewKeeper {
	return ViewKeeper{am: am}
}

// GetCoins returns the coins at the addr.
func (keeper ViewKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return getCoins(ctx, keeper.am, addr)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper ViewKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return hasCoins(ctx, keeper.am, addr, amt)
}

//______________________________________________________________________________________________

func getCoins(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress) sdk.Coins {
	ctx.GasMeter().ConsumeGas(costGetCoins, "getCoins")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.Coins{}
	}
	return acc.GetCoins()
}

func setCoins(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	ctx.GasMeter().ConsumeGas(costSetCoins, "setCoins")
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress, amt sdk.Coins) bool {
	ctx.GasMeter().ConsumeGas(costHasCoins, "hasCoins")
	return getCoins(ctx, am, addr).IsGTE(amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func subtractCoins(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costSubtractCoins, "subtractCoins")
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Minus(amt)
	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags("sender", []byte(addr.String()))
	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error) {
	ctx.GasMeter().ConsumeGas(costAddCoins, "addCoins")
	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Plus(amt)
	if !newCoins.IsNotNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("%s < %s", oldCoins, amt))
	}
	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags("recipient", []byte(addr.String()))
	return newCoins, tags, err
}

// SendCoins moves coins from one account to another
// NOTE: Make sure to revert state changes from tx on error
func sendCoins(ctx sdk.Context, am auth.AccountMapper, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error) {
	_, subTags, err := subtractCoins(ctx, am, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, am, toAddr, amt)
	if err != nil {
		return nil, err
	}

	return subTags.AppendTags(addTags), nil
}

// InputOutputCoins handles a list of inputs and outputs
// NOTE: Make sure to revert state changes from tx on error
func inputOutputCoins(ctx sdk.Context, am auth.AccountMapper, inputs []Input, outputs []Output) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, am, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {
		_, tags, err := addCoins(ctx, am, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil
}

//*****comdex
func getAssetWallet(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress) sdk.AssetPegWallet {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.AssetPegWallet{}
	}
	return acc.GetAssetPegWallet()
}

func setAssetWallet(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress, asset sdk.AssetPegWallet) sdk.Error {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	acc.SetAssetPegWallet(asset)
	am.SetAccount(ctx, acc)
	return nil
}

func instantiateAndAssignAsset(ctx sdk.Context, am auth.AccountMapper, issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg sdk.AssetPeg) (sdk.AssetPeg, sdk.Tags, sdk.Error) {

	pegHash, _ := sdk.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(am.GetNextAssetPegHash(ctx))))
	assetPeg.SetPegHash(pegHash)
	assetPeg.SetLocked(assetPeg.GetModerated())
	receiverAssetPegWallet := getAssetWallet(ctx, am, toAddress)
	receiverAssetPegWallet = sdk.AddAssetPegToWallet(assetPeg, receiverAssetPegWallet)
	err := setAssetWallet(ctx, am, toAddress, receiverAssetPegWallet)
	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("issuer", []byte(issuerAddress.String()))
	tags = tags.AppendTag("asset", []byte(assetPeg.GetPegHash().String()))
	return assetPeg, tags, err
}

func issueAssetsToWallets(ctx sdk.Context, am auth.AccountMapper, issueAsset []IssueAsset, ak acl.Keeper) (sdk.Tags, sdk.Error, []sdk.AssetPeg) {
	allTags := sdk.EmptyTags()
	var issuedAssetPegs []sdk.AssetPeg
	var acl sdk.ACL
	var err sdk.Error
	for _, req := range issueAsset {
		moderated := req.AssetPeg.GetModerated()
		if moderated {
			acl, err = ak.CheckZoneAndGetACL(ctx, req.IssuerAddress, req.ToAddress)
			if err != nil {
				return nil, err, issuedAssetPegs
			}
		} else {
			aclAccount, err1 := ak.GetAccountACLDetails(ctx, req.IssuerAddress)
			if err1 != nil {
				return nil, err1, issuedAssetPegs
			}
			acl = aclAccount.GetACL()
		}

		if !acl.IssueAsset {
			return nil, sdk.ErrInternal(fmt.Sprintf("Assets cant be issued to account %v.", req.ToAddress.String())), issuedAssetPegs
		}

		issuedAsset, tags, err := instantiateAndAssignAsset(ctx, am, req.IssuerAddress, req.ToAddress, req.AssetPeg)
		if err != nil {
			return nil, err, issuedAssetPegs
		}
		issuedAssetPegs = append(issuedAssetPegs, issuedAsset)
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil, issuedAssetPegs
}

//IssueAssetsToWallets handles a list of IssueAsset messages
func (keeper Keeper) IssueAssetsToWallets(ctx sdk.Context, issueAssets []IssueAsset, ak acl.Keeper) (sdk.Tags, sdk.Error, []sdk.AssetPeg) {
	return issueAssetsToWallets(ctx, keeper.am, issueAssets, ak)
}

func instantiateAndRedeemAsset(ctx sdk.Context, am auth.AccountMapper, issuerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, pegHash sdk.PegHash) (sdk.PegHash, sdk.Tags, sdk.Error) {
	redeemerPegHashWallet := getAssetWallet(ctx, am, redeemerAddress)
	issuerPegHashWallet := getAssetWallet(ctx, am, issuerAddress)
	var assetPeg sdk.AssetPeg
	length := len(redeemerPegHashWallet)
	if length == 0 {
		return nil, nil, sdk.ErrInternal("No Assets Found!") // Codespace and CodeType needs to be defined
	}
	i := redeemerPegHashWallet.SearchAssetPeg(pegHash)
	if i == length {
		return nil, nil, sdk.ErrInternal("No Assets With Given PegHash Found!") // Codespace and CodeType needs to be defined
	}
	assetPeg, redeemerPegHashWallet = sdk.SubtractAssetPegFromWallet(pegHash, redeemerPegHashWallet)
	unSetAssetPeg := sdk.NewBaseAssetPegWithPegHash(assetPeg.GetPegHash())
	issuerPegHashWallet = sdk.AddAssetPegToWallet(&unSetAssetPeg, issuerPegHashWallet)
	err := setAssetWallet(ctx, am, redeemerAddress, redeemerPegHashWallet)
	if err == nil {
		err = setAssetWallet(ctx, am, issuerAddress, issuerPegHashWallet)
	}
	tags := sdk.NewTags("redeemer", []byte(redeemerAddress.String()))
	tags = tags.AppendTag("issuer", []byte(issuerAddress.String()))
	tags = tags.AppendTag("asset", []byte(assetPeg.GetPegHash().String()))

	return assetPeg.GetPegHash(), tags, nil
}

func redeemAssetsFromWallets(ctx sdk.Context, am auth.AccountMapper, redeemAssets []RedeemAsset, ak acl.Keeper) (sdk.Tags, sdk.Error, []sdk.PegHash) {
	allTags := sdk.EmptyTags()
	var redeemedPegHashs []sdk.PegHash
	for _, req := range redeemAssets {
		acl, err := ak.CheckZoneAndGetACL(ctx, req.IssuerAddress, req.RedeemerAddress)
		if err != nil {
			return nil, err, redeemedPegHashs
		}
		if !acl.RedeemAsset {
			return nil, sdk.ErrInternal(fmt.Sprintf("Assets can't be redeemed from account %v.", req.RedeemerAddress.String())), redeemedPegHashs
		}
		redeemedPegHash, tags, err := instantiateAndRedeemAsset(ctx, am, req.IssuerAddress, req.RedeemerAddress, req.PegHash)
		if err != nil {
			return nil, err, redeemedPegHashs
		}
		redeemedPegHashs = append(redeemedPegHashs, redeemedPegHash)
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil, redeemedPegHashs
}

//RedeemAssetsFromWallets handles a list of RedeemAsset messages
func (keeper Keeper) RedeemAssetsFromWallets(ctx sdk.Context, redeemAssets []RedeemAsset, ak acl.Keeper) (sdk.Tags, sdk.Error, []sdk.PegHash) {
	return redeemAssetsFromWallets(ctx, keeper.am, redeemAssets, ak)
}

//***** Issue Fiat
func getFiatWallet(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress) sdk.FiatPegWallet {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.FiatPegWallet{}
	}
	return acc.GetFiatPegWallet()
}

func setFiatWallet(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress, fiat sdk.FiatPegWallet) sdk.Error {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	acc.SetFiatPegWallet(fiat)
	am.SetAccount(ctx, acc)
	return nil
}

func instantiateAndAssignFiat(ctx sdk.Context, am auth.AccountMapper, issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, fiatPeg sdk.FiatPeg) (sdk.FiatPeg, sdk.Tags, sdk.Error) {

	pegHash, _ := sdk.GetFiatPegHashHex(fmt.Sprintf("%x", strconv.Itoa(am.GetNextFiatPegHash(ctx))))
	fiatPeg.SetPegHash(pegHash)
	receiverFiatPegWallet := getFiatWallet(ctx, am, toAddress)
	receiverFiatPegWallet = sdk.AddFiatPegToWallet(receiverFiatPegWallet, []sdk.BaseFiatPeg{sdk.ToBaseFiatPeg(fiatPeg)})

	err := setFiatWallet(ctx, am, toAddress, receiverFiatPegWallet)

	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("issuer", []byte(issuerAddress.String()))
	tags = tags.AppendTag("fiat", []byte(fiatPeg.GetPegHash().String()))
	return fiatPeg, tags, err
}

func issueFiatsToWallets(ctx sdk.Context, am auth.AccountMapper, issueFiat []IssueFiat, ak acl.Keeper) (sdk.Tags, sdk.Error, []sdk.FiatPeg) {
	allTags := sdk.EmptyTags()
	var issuedFiatPegs []sdk.FiatPeg
	for _, req := range issueFiat {

		acl, err := ak.CheckZoneAndGetACL(ctx, req.IssuerAddress, req.ToAddress)
		if err != nil {
			return nil, err, issuedFiatPegs
		}
		if !acl.IssueFiat {
			return nil, sdk.ErrInternal(fmt.Sprintf("Fiats can't be issued to account %v.", req.ToAddress.String())), issuedFiatPegs
		}

		issuedFiat, tags, err := instantiateAndAssignFiat(ctx, am, req.IssuerAddress, req.ToAddress, req.FiatPeg)
		if err != nil {
			return nil, err, issuedFiatPegs
		}
		issuedFiatPegs = append(issuedFiatPegs, issuedFiat)
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil, issuedFiatPegs
}

//IssueFiatsToWallets handles a list of IssueFiat messages
func (keeper Keeper) IssueFiatsToWallets(ctx sdk.Context, issueFiats []IssueFiat, ak acl.Keeper) (sdk.Tags, sdk.Error, []sdk.FiatPeg) {
	return issueFiatsToWallets(ctx, keeper.am, issueFiats, ak)
}

//##### Issue Fiat

//##### Redeem Fiat

func instantiateAndRedeemFiat(ctx sdk.Context, am auth.AccountMapper, issuerAddress sdk.AccAddress, redeemerAddress sdk.AccAddress, amount int64) (sdk.FiatPegWallet, sdk.Tags, sdk.Error) {
	var redeemerFiatPegWallet sdk.FiatPegWallet

	fromOldFiatWallet := getFiatWallet(ctx, am, redeemerAddress)

	emptiedFiatPegWallet, redeemerFiatPegWallet := sdk.RedeemAmountFromWallet(amount, fromOldFiatWallet)
	if len(redeemerFiatPegWallet) == 0 && len(emptiedFiatPegWallet) == 0 {
		return fromOldFiatWallet, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("Redeemed amount higher than the account balance"))
	}

	err := setFiatWallet(ctx, am, redeemerAddress, redeemerFiatPegWallet)
	if err != nil {
		return nil, nil, err
	}

	tags := sdk.NewTags("redeemer", []byte(redeemerAddress.String()))
	return redeemerFiatPegWallet, tags, err
}

func redeemFiatsFromWallets(ctx sdk.Context, am auth.AccountMapper, redeemFiats []RedeemFiat, ak acl.Keeper) (sdk.Tags, sdk.Error, []sdk.FiatPegWallet) {
	allTags := sdk.EmptyTags()
	var redeemerFiatPegWallets []sdk.FiatPegWallet

	for _, req := range redeemFiats {
		acl, err := ak.CheckZoneAndGetACL(ctx, req.IssuerAddress, req.RedeemerAddress)
		if err != nil {
			return nil, err, redeemerFiatPegWallets
		}
		if !acl.RedeemFiat {
			return nil, sdk.ErrInternal(fmt.Sprintf("Fiats can't be redeemed from account %v.", req.RedeemerAddress.String())), redeemerFiatPegWallets
		}
		redeemerFiatPegWallet, tags, err := instantiateAndRedeemFiat(ctx, am, req.IssuerAddress, req.RedeemerAddress, req.Amount)
		if err != nil {
			return nil, sdk.ErrInternal(""), redeemerFiatPegWallets
		}
		redeemerFiatPegWallets = append(redeemerFiatPegWallets, redeemerFiatPegWallet)
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil, redeemerFiatPegWallets
}

//RedeemFiatsFromWallets handles a list of Redeem Fiat messages
func (keeper Keeper) RedeemFiatsFromWallets(ctx sdk.Context, redeemFiats []RedeemFiat, ak acl.Keeper) (sdk.Tags, sdk.Error, []sdk.FiatPegWallet) {
	return redeemFiatsFromWallets(ctx, keeper.am, redeemFiats, ak)
}

//##### Send Asset

func sendAssetToOrder(ctx sdk.Context, am auth.AccountMapper, orderKeeper order.Keeper, negotiationKeeper negotiation.Keeper, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash) (sdk.AssetPeg, sdk.Tags, sdk.Error) {
	err, negotiation := negotiationKeeper.GetNegotiation(ctx, toAddress, fromAddress, pegHash)
	if err != nil {
		return nil, nil, err
	}
	time := negotiation.GetTime() + negotiation.GetSellerBlockHeight()
	if ctx.BlockHeight() > time {
		return nil, nil, sdk.ErrInvalidSequence("Negotiation time expired.")
	}
	if negotiation.GetSellerSignature() == nil || negotiation.GetBuyerSignature() == nil {
		return nil, nil, sdk.ErrInternal("Signatures are not present")
	}

	fromOldAssetWallet := getAssetWallet(ctx, am, fromAddress)
	sentAsset, fromNewAssetPegWallet := sdk.SubtractAssetPegFromWallet(pegHash, fromOldAssetWallet)
	if sentAsset == nil {
		return nil, nil, sdk.ErrInsufficientCoins("Asset not found.")
	}
	if sentAsset.GetLocked() {
		return nil, nil, sdk.ErrInsufficientCoins("Asset locked.")
	}
	err = orderKeeper.SendAssetsToOrder(ctx, fromAddress, toAddress, sentAsset)
	if err == nil {
		err = setAssetWallet(ctx, am, fromAddress, fromNewAssetPegWallet)
	}
	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("sender", []byte(fromAddress.String()))
	tags = tags.AppendTag("asset", []byte(sentAsset.GetPegHash().String()))
	return sentAsset, tags, err
}

func sendAssetsToWallets(ctx sdk.Context, am auth.AccountMapper, orderKeeper order.Keeper, negotiationKeeper negotiation.Keeper, sendAsset []SendAsset, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, []sdk.AssetPeg) {
	allTags := sdk.EmptyTags()
	var sentAssetPegs []sdk.AssetPeg

	for _, req := range sendAsset {
		aclStore, err := ak.GetAccountACLDetails(ctx, req.FromAddress)
		if err != nil {
			return nil, sdk.ErrInternal("Unauthorized transaction"), nil
		}
		account := aclStore.GetACL()
		if !account.SendAsset {
			return nil, sdk.ErrInternal("Unauthorized transaction"), nil
		}
		sentAsset, tags, err := sendAssetToOrder(ctx, am, orderKeeper, negotiationKeeper, req.FromAddress, req.ToAddress, req.PegHash)
		if err != nil {
			return nil, err, sentAssetPegs
		}
		sentAssetPegs = append(sentAssetPegs, sentAsset)
		reputationKeeper.SetSendAssetsPositiveTx(ctx, req.FromAddress)
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil, sentAssetPegs
}

//SendAssetsToWallets handles a list of SendAsset messages
func (keeper Keeper) SendAssetsToWallets(ctx sdk.Context, orderKeeper order.Keeper, negotiationKeeper negotiation.Keeper, sendAssets []SendAsset, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, []sdk.AssetPeg) {
	return sendAssetsToWallets(ctx, keeper.am, orderKeeper, negotiationKeeper, sendAssets, ak, reputationKeeper)
}

//***** Send Asset

//##### Send Fiat

func sendFiatToOrder(ctx sdk.Context, am auth.AccountMapper, orderKeeper order.Keeper, negotiationKeeper negotiation.Keeper, fromAddress sdk.AccAddress, toAddress sdk.AccAddress, pegHash sdk.PegHash, amount int64) (sdk.FiatPegWallet, sdk.Tags, sdk.Error) {
	var sentFiatPegWallet, oldFiatPegWallet sdk.FiatPegWallet
	err, negotiation := negotiationKeeper.GetNegotiation(ctx, fromAddress, toAddress, pegHash)
	if err != nil {
		return sentFiatPegWallet, nil, err
	}
	time := negotiation.GetTime() + negotiation.GetBuyerBlockHeight()
	if ctx.BlockHeight() > time {
		return nil, nil, sdk.ErrInvalidSequence("Negotiation time expired.")
	}
	if negotiation.GetSellerSignature() == nil || negotiation.GetBuyerSignature() == nil {
		return nil, nil, sdk.ErrInternal("Signatures are not present")
	}

	fromOldFiatWallet := getFiatWallet(ctx, am, fromAddress)
	sentFiatPegWallet, oldFiatPegWallet = sdk.SubtractAmountFromWallet(amount, fromOldFiatWallet)
	if len(sentFiatPegWallet) == 0 && len(oldFiatPegWallet) == 0 {
		return fromOldFiatWallet, nil, sdk.ErrInsufficientCoins(fmt.Sprintf("Insufficient funds"))
	}

	err = orderKeeper.SendFiatsToOrder(ctx, fromAddress, toAddress, pegHash, sentFiatPegWallet)
	if err == nil {
		err = setFiatWallet(ctx, am, fromAddress, oldFiatPegWallet)
	}
	tags := sdk.NewTags("recepient", []byte(toAddress.String()))
	tags = tags.AppendTag("sender", []byte(fromAddress.String()))
	return sentFiatPegWallet, tags, err
}

func sendFiatsToWallets(ctx sdk.Context, am auth.AccountMapper, orderKeeper order.Keeper, negotiationKeeper negotiation.Keeper, sendFiat []SendFiat, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, []sdk.FiatPegWallet) {
	allTags := sdk.EmptyTags()
	var sentFiatPegWallets []sdk.FiatPegWallet

	for _, req := range sendFiat {
		store, err := ak.GetAccountACLDetails(ctx, req.FromAddress)
		if err != nil {
			return nil, sdk.ErrInternal("Unauthorized transaction"), nil
		}
		account := store.GetACL()
		if !account.SendFiat {
			return nil, sdk.ErrInternal("Unauthorized transaction"), nil
		}
		sentFiatPegWallet, tags, err := sendFiatToOrder(ctx, am, orderKeeper, negotiationKeeper, req.FromAddress, req.ToAddress, req.PegHash, req.Amount)
		if err != nil {
			return nil, err, sentFiatPegWallets
		}
		sentFiatPegWallets = append(sentFiatPegWallets, sentFiatPegWallet)
		reputationKeeper.SetSendFiatsPositiveTx(ctx, req.FromAddress)
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil, sentFiatPegWallets
}

//SendFiatsToWallets handles a list of SendFiat messages
func (keeper Keeper) SendFiatsToWallets(ctx sdk.Context, orderKeeper order.Keeper, negotiationKeeper negotiation.Keeper, sendFiats []SendFiat, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, []sdk.FiatPegWallet) {
	return sendFiatsToWallets(ctx, keeper.am, orderKeeper, negotiationKeeper, sendFiats, ak, reputationKeeper)
}

//***** Send Fiat

//##### Execute Trade Order

func exchangeOrderTokens(ctx sdk.Context, am auth.AccountMapper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, fiatProofHash string, awbProofHash string, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, sdk.FiatPegWallet, sdk.AssetPegWallet) {
	err, assetPegWallet, fiatPegWallet, orderFiatProofHash, orderAWBProofHash := orderKeeper.GetOrderDetails(ctx, buyerAddress, sellerAddress, pegHash)
	if err != nil {
		return nil, err, fiatPegWallet, assetPegWallet
	}
	err, negotiation := negotiationKeeper.GetNegotiation(ctx, buyerAddress, sellerAddress, pegHash)
	if err != nil {
		return nil, err, fiatPegWallet, assetPegWallet
	}

	var reverseOrder bool
	var oldFiatPegWallet sdk.FiatPegWallet
	if len(fiatPegWallet) == 0 || negotiation.GetBid() > sdk.GetFiatPegWalletBalance(fiatPegWallet) {
		if negotiation.GetTime() < ctx.BlockHeight() {
			return nil, sdk.ErrInsufficientCoins("Fiat tokens not found!"), fiatPegWallet, assetPegWallet
		}
		reverseOrder = true
		reputationKeeper.SetBuyerExecuteOrderNegativeTx(ctx, buyerAddress)
	}
	if len(assetPegWallet) != 1 || assetPegWallet[0].GetPegHash().String() != pegHash.String() {
		if negotiation.GetTime() < ctx.BlockHeight() {
			return nil, sdk.ErrInsufficientCoins("Asset token not found!"), fiatPegWallet, assetPegWallet
		}
		reverseOrder = true
		reputationKeeper.SetSellerExecuteOrderNegativeTx(ctx, sellerAddress)

	}
	buyerTime := negotiation.GetTime() + negotiation.GetBuyerBlockHeight()
	sellerTime := negotiation.GetTime() + negotiation.GetSellerBlockHeight()
	time := ctx.BlockHeight()
	if time > buyerTime || time > sellerTime {
		reverseOrder = true
	}

	if negotiation.GetBid() < sdk.GetFiatPegWalletBalance(fiatPegWallet) {
		fiatPegWallet, oldFiatPegWallet = sdk.SubtractAmountFromWallet(negotiation.GetBid(), fiatPegWallet)
	}
	var executed bool
	if !reverseOrder {
		sellerFiatWallet := getFiatWallet(ctx, am, sellerAddress)
		buyerAssetWallet := getAssetWallet(ctx, am, buyerAddress)

		if orderFiatProofHash == "" && fiatProofHash != "" {
			orderKeeper.SetOrderFiatProofHash(ctx, buyerAddress, sellerAddress, pegHash, fiatProofHash)
		}
		if orderAWBProofHash == "" && awbProofHash != "" {
			orderKeeper.SetOrderAWBProofHash(ctx, buyerAddress, sellerAddress, pegHash, awbProofHash)
		}
		err, _, _, orderFiatProofHash, orderAWBProofHash = orderKeeper.GetOrderDetails(ctx, buyerAddress, sellerAddress, pegHash)
		if err != nil {
			return nil, err, nil, nil
		}
		if orderFiatProofHash != "" && orderAWBProofHash != "" {
			executed = true
			reputationKeeper.SetSellerExecuteOrderPositiveTx(ctx, sellerAddress)
			reputationKeeper.SetBuyerExecuteOrderPositiveTx(ctx, buyerAddress)

			buyerAssetWallet = sdk.AddAssetPegToWallet(&assetPegWallet[0], buyerAssetWallet)
			sellerFiatWallet = sdk.AddFiatPegToWallet(sellerFiatWallet, fiatPegWallet)

			fiatPegWallet = orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
			assetPegWallet = orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])
		}

		setFiatWallet(ctx, am, sellerAddress, sellerFiatWallet)
		setAssetWallet(ctx, am, buyerAddress, buyerAssetWallet)

	}

	if executed == true || reverseOrder == true {
		buyerFiatWallet := getFiatWallet(ctx, am, buyerAddress)
		sellerAssetWallet := getAssetWallet(ctx, am, sellerAddress)
		if len(oldFiatPegWallet) != 0 {
			buyerFiatWallet = sdk.AddFiatPegToWallet(buyerFiatWallet, oldFiatPegWallet)
		}

		if len(fiatPegWallet) != 0 {
			buyerFiatWallet = sdk.AddFiatPegToWallet(buyerFiatWallet, fiatPegWallet)
			orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
		}
		if len(assetPegWallet) != 0 {
			sellerAssetWallet = sdk.AddAssetPegToWallet(&assetPegWallet[0], sellerAssetWallet)
			orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])
		}

		setFiatWallet(ctx, am, buyerAddress, buyerFiatWallet)
		setAssetWallet(ctx, am, sellerAddress, sellerAssetWallet)
	}
	tags := sdk.NewTags("buyer", []byte(buyerAddress.String()))
	tags = tags.AppendTag("seller", []byte(sellerAddress.String()))
	tags = tags.AppendTag("assetPegHash", []byte(pegHash.String()))
	tags = tags.AppendTag("executed", []byte(strconv.FormatBool(executed)))
	tags = tags.AppendTag("assetPrice", []byte(strconv.FormatInt(negotiation.GetBid(), 10)))
	tags = tags.AppendTag("reversed", []byte(strconv.FormatBool(reverseOrder)))
	return tags, nil, fiatPegWallet, assetPegWallet
}

//##### Buyer Execute Trade Order
func buyerExecuteTradeOrders(ctx sdk.Context, am auth.AccountMapper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, buyerExecuteOrders []BuyerExecuteOrder, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, []sdk.FiatPegWallet) {
	allTags := sdk.EmptyTags()
	var fiatPegWallets []sdk.FiatPegWallet
	var acl sdk.ACL
	var err sdk.Error
	var tags sdk.Tags
	var fiatPegWallet sdk.FiatPegWallet
	for _, req := range buyerExecuteOrders {
		_, assetWallet, _, _, _ := orderKeeper.GetOrderDetails(ctx, req.BuyerAddress, req.SellerAddress, req.PegHash)
		if len(assetWallet) == 0 {
			return nil, sdk.ErrInsufficientCoins("Asset token not found!"), fiatPegWallets
		}
		i := assetWallet.SearchAssetPeg(req.PegHash)
		if i < len(assetWallet) && assetWallet[i].GetPegHash().String() == req.PegHash.String() {
			assetPeg := assetWallet[i]
			moderated := assetPeg.GetModerated()
			if !moderated {
				aclAccount, err := ak.GetAccountACLDetails(ctx, req.BuyerAddress)
				if err != nil {
					return nil, err, fiatPegWallets
				}
				if !reflect.DeepEqual(req.MediatorAddress, req.BuyerAddress) {
					return nil, sdk.ErrUnauthorized(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", req.MediatorAddress.String())), fiatPegWallets
				}
				acl = aclAccount.GetACL()
				if !acl.BuyerExecuteOrder {
					return nil, sdk.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", req.BuyerAddress.String())), fiatPegWallets
				}
				tags, err, fiatPegWallet, _ = privateExchangeOrderTokens(ctx, am, negotiationKeeper, orderKeeper, req.MediatorAddress, req.BuyerAddress, req.SellerAddress, req.PegHash, req.FiatProofHash, "", reputationKeeper)
				if err != nil {
					return nil, err, fiatPegWallets
				}
				fiatPegWallets = append(fiatPegWallets, fiatPegWallet)
				if err != nil {
					return nil, err, fiatPegWallets
				}
				allTags = allTags.AppendTags(tags)
				return allTags, nil, fiatPegWallets
			}
			acl, err = ak.CheckZoneAndGetACL(ctx, req.MediatorAddress, req.BuyerAddress)
			if err != nil {
				return nil, err, fiatPegWallets

			}
		}
		if !acl.BuyerExecuteOrder {
			return nil, sdk.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", req.BuyerAddress.String())), fiatPegWallets
		}
		tags, err, fiatPegWallet, _ = exchangeOrderTokens(ctx, am, negotiationKeeper, orderKeeper, req.MediatorAddress, req.BuyerAddress, req.SellerAddress, req.PegHash, req.FiatProofHash, "", reputationKeeper)
		if err != nil {
			return nil, err, fiatPegWallets
		}
		fiatPegWallets = append(fiatPegWallets, fiatPegWallet)
		if err != nil {
			return nil, err, fiatPegWallets
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil, fiatPegWallets
}

//BuyerExecuteTradeOrders :  haddles a list of Buyer Execute Order messages
func (keeper Keeper) BuyerExecuteTradeOrders(ctx sdk.Context, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, buyerExecuteOrders []BuyerExecuteOrder, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, []sdk.FiatPegWallet) {
	return buyerExecuteTradeOrders(ctx, keeper.am, negotiationKeeper, orderKeeper, buyerExecuteOrders, ak, reputationKeeper)
}

//##### Buyer Execute Trade Order
//##### Seller Execute Trade Order
func sellerExecuteTradeOrders(ctx sdk.Context, am auth.AccountMapper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, sellerExecuteOrders []SellerExecuteOrder, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, []sdk.AssetPegWallet) {
	allTags := sdk.EmptyTags()
	var assetPegWallets []sdk.AssetPegWallet
	var acl sdk.ACL
	var err sdk.Error
	var tags sdk.Tags
	var assetPegWallet sdk.AssetPegWallet
	for _, req := range sellerExecuteOrders {

		_, assetWallet, _, _, _ := orderKeeper.GetOrderDetails(ctx, req.BuyerAddress, req.SellerAddress, req.PegHash)
		if len(assetWallet) == 0 {
			return nil, sdk.ErrInsufficientCoins("Asset token not found!"), assetPegWallets
		}
		i := assetWallet.SearchAssetPeg(req.PegHash)
		if i < len(assetWallet) && assetWallet[i].GetPegHash().String() == req.PegHash.String() {
			assetPeg := assetWallet[i]
			moderated := assetPeg.GetModerated()
			if !moderated {
				aclAccount, err := ak.GetAccountACLDetails(ctx, req.SellerAddress)
				if err != nil {
					return nil, err, assetPegWallets
				}
				if !reflect.DeepEqual(req.MediatorAddress, req.SellerAddress) {
					return nil, sdk.ErrUnauthorized(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", req.MediatorAddress.String())), assetPegWallets
				}
				acl = aclAccount.GetACL()
				if !acl.SellerExecuteOrder {
					return nil, sdk.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", req.SellerAddress.String())), assetPegWallets
				}
				tags, err, _, assetPegWallet = privateExchangeOrderTokens(ctx, am, negotiationKeeper, orderKeeper, req.MediatorAddress, req.BuyerAddress, req.SellerAddress, req.PegHash, "", req.AWBProofHash, reputationKeeper)
				if err != nil {
					return nil, err, assetPegWallets
				}
				assetPegWallets = append(assetPegWallets, assetPegWallet)
				if err != nil {
					return nil, err, assetPegWallets
				}
				allTags = allTags.AppendTags(tags)
				return allTags, nil, assetPegWallets
			}
			acl, err = ak.CheckZoneAndGetACL(ctx, req.MediatorAddress, req.SellerAddress)
			if err != nil {
				return nil, err, assetPegWallets

			}
		}

		if !acl.SellerExecuteOrder {
			return nil, sdk.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.", req.SellerAddress.String())), assetPegWallets
		}

		tags, err, _, assetPegWallet = exchangeOrderTokens(ctx, am, negotiationKeeper, orderKeeper, req.MediatorAddress, req.BuyerAddress, req.SellerAddress, req.PegHash, "", req.AWBProofHash, reputationKeeper)
		if err != nil {
			return nil, err, assetPegWallets
		}
		assetPegWallets = append(assetPegWallets, assetPegWallet)
		if err != nil {
			return nil, err, assetPegWallets
		}
		allTags = allTags.AppendTags(tags)

	}

	return allTags, nil, assetPegWallets
}

//SellerExecuteTradeOrders :  haddles a list of Seller Execute Order messages
func (keeper Keeper) SellerExecuteTradeOrders(ctx sdk.Context, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, sellerExecuteOrders []SellerExecuteOrder, ak acl.Keeper, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, []sdk.AssetPegWallet) {
	return sellerExecuteTradeOrders(ctx, keeper.am, negotiationKeeper, orderKeeper, sellerExecuteOrders, ak, reputationKeeper)
}

//##### Seller Execute Trade Order
//***** Release Assets

func releaseAsset(ctx sdk.Context, am auth.AccountMapper, zoneAddress sdk.AccAddress, ownerAddress sdk.AccAddress, pegHash sdk.PegHash) (sdk.Tags, sdk.Error) {
	ownerAssetWallet := getAssetWallet(ctx, am, ownerAddress)
	if !sdk.ReleaseAssetPegInWallet(ownerAssetWallet, pegHash) {
		return nil, sdk.ErrInternal("Asset peg not found.")
	}
	setAssetWallet(ctx, am, ownerAddress, ownerAssetWallet)
	tags := sdk.NewTags("zone", []byte(zoneAddress.String()))
	tags = tags.AppendTag("owner", []byte(ownerAddress.String()))
	tags = tags.AppendTag("asset", []byte(pegHash.String()))
	return tags, nil
}

func releaseLockedAssets(ctx sdk.Context, am auth.AccountMapper, releaseAssets []ReleaseAsset, ak acl.Keeper) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	for _, req := range releaseAssets {

		acl, err := ak.CheckZoneAndGetACL(ctx, req.ZoneAddress, req.OwnerAddress)
		if err != nil {
			return nil, err
		}
		if !acl.ReleaseAsset {
			return nil, sdk.ErrInternal(fmt.Sprintf("Assets cannot be released for account %v. Access Denied.", req.OwnerAddress.String()))
		}

		tags, err := releaseAsset(ctx, am, req.ZoneAddress, req.OwnerAddress, req.PegHash)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

//ReleaseLockedAssets :  haddles a list of release asset messages
func (keeper Keeper) ReleaseLockedAssets(ctx sdk.Context, releaseAssets []ReleaseAsset, ak acl.Keeper) (sdk.Tags, sdk.Error) {
	return releaseLockedAssets(ctx, keeper.am, releaseAssets, ak)
}

//##### Release Assets

//*****ACL

//DefineZones : handles a list of Define Zone messages
func (keeper Keeper) DefineZones(ctx sdk.Context, aclKeeper acl.Keeper, defineZone []DefineZone) (sdk.Tags, sdk.Error) {
	return defineZones(ctx, keeper, aclKeeper, defineZone)
}

func defineZones(ctx sdk.Context, keeper Keeper, aclKeeper acl.Keeper, defineZone []DefineZone) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	for _, in := range defineZone {

		if !keeper.CheckGenesisAccount(ctx, in.From) {
			return nil, sdk.ErrInternal(fmt.Sprintf("Account %v is not the genesis account. Zones can only be defined by the genesis account.", in.From.String()))
		}

		tags, err := aclKeeper.DefineZoneAddress(ctx, in.To, in.ZoneID)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

//DefineOrganizations : handles a list of Define Organization messages
func (keeper Keeper) DefineOrganizations(ctx sdk.Context, aclKeeper acl.Keeper, defineOrganization []DefineOrganization) (sdk.Tags, sdk.Error) {
	return defineOrganizations(ctx, keeper, aclKeeper, defineOrganization)
}

func defineOrganizations(ctx sdk.Context, keeper Keeper, aclKeeper acl.Keeper, defineOrganization []DefineOrganization) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	for _, in := range defineOrganization {
		if !aclKeeper.CheckIfZoneAccount(ctx, in.ZoneID, in.From) {
			return nil, sdk.ErrInternal(fmt.Sprintf("Account %v is not the zone account. Organizations can only be defined by the zone account.", in.From.String()))
		}
		tags, err := aclKeeper.DefineOrganizationAddress(ctx, in.To, in.OrganizationID, in.ZoneID)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

//DefineACLs : handles a list of Define ACL messages
func (keeper Keeper) DefineACLs(ctx sdk.Context, aclKeeper acl.Keeper, defineACL []DefineACL) (sdk.Tags, sdk.Error) {
	return defineACLs(ctx, keeper, aclKeeper, defineACL)
}

func defineACLs(ctx sdk.Context, keeper Keeper, aclKeeper acl.Keeper, defineACL []DefineACL) (sdk.Tags, sdk.Error) {
	allTags := sdk.EmptyTags()
	for _, in := range defineACL {

		if !keeper.CheckGenesisAccount(ctx, in.From) {
			if !aclKeeper.CheckIfZoneAccount(ctx, in.ACLAccount.GetZoneID(), in.From) {
				if !aclKeeper.CheckIfOrganizationAccount(ctx, in.ACLAccount.GetOrganizationID(), in.From) {
					return nil, sdk.ErrInternal(fmt.Sprintf("Account %v does not have access to define acl for account %v.", in.From.String(), in.To.String()))
				}
			}
		}

		tags, err := aclKeeper.DefineACLAccount(ctx, in.To, in.ACLAccount)
		if err != nil {
			return nil, err
		}

		allTags = allTags.AppendTags(tags)
	}
	return allTags, nil
}

func privateExchangeOrderTokens(ctx sdk.Context, am auth.AccountMapper, negotiationKeeper negotiation.Keeper, orderKeeper order.Keeper, mediatorAddress sdk.AccAddress, buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash sdk.PegHash, fiatProofHash string, awbProofHash string, reputationKeeper reputation.Keeper) (sdk.Tags, sdk.Error, sdk.FiatPegWallet, sdk.AssetPegWallet) {
	err, assetPegWallet, fiatPegWallet, orderFiatProofHash, orderAWBProofHash := orderKeeper.GetOrderDetails(ctx, buyerAddress, sellerAddress, pegHash)
	if err != nil {
		return nil, err, fiatPegWallet, assetPegWallet
	}
	err, negotiation := negotiationKeeper.GetNegotiation(ctx, buyerAddress, sellerAddress, pegHash)
	if err != nil {
		return nil, err, fiatPegWallet, assetPegWallet
	}

	var reverseOrder bool
	var executed bool

	sellerFiatWallet := getFiatWallet(ctx, am, sellerAddress)
	buyerAssetWallet := getAssetWallet(ctx, am, buyerAddress)

	if orderFiatProofHash == "" && fiatProofHash != "" {
		orderKeeper.SetOrderFiatProofHash(ctx, buyerAddress, sellerAddress, pegHash, fiatProofHash)
	}
	if orderAWBProofHash == "" && awbProofHash != "" {
		orderKeeper.SetOrderAWBProofHash(ctx, buyerAddress, sellerAddress, pegHash, awbProofHash)
	}
	err, _, _, orderFiatProofHash, orderAWBProofHash = orderKeeper.GetOrderDetails(ctx, buyerAddress, sellerAddress, pegHash)
	if err != nil {
		return nil, err, nil, nil
	}
	if orderFiatProofHash != "" && orderAWBProofHash != "" {
		executed = true
		reputationKeeper.SetSellerExecuteOrderPositiveTx(ctx, sellerAddress)
		reputationKeeper.SetBuyerExecuteOrderPositiveTx(ctx, buyerAddress)

		buyerAssetWallet = sdk.AddAssetPegToWallet(&assetPegWallet[0], buyerAssetWallet)
		sellerFiatWallet = sdk.AddFiatPegToWallet(sellerFiatWallet, fiatPegWallet)

		fiatPegWallet = orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
		assetPegWallet = orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])

		setFiatWallet(ctx, am, sellerAddress, sellerFiatWallet)
		setAssetWallet(ctx, am, buyerAddress, buyerAssetWallet)
	}
	if orderAWBProofHash == "" && orderFiatProofHash == "" {
		reverseOrder = true
	}
	if executed == true || reverseOrder == true {
		buyerFiatWallet := getFiatWallet(ctx, am, buyerAddress)
		sellerAssetWallet := getAssetWallet(ctx, am, sellerAddress)

		if len(fiatPegWallet) != 0 {
			buyerFiatWallet = sdk.AddFiatPegToWallet(buyerFiatWallet, fiatPegWallet)
			orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
		}
		if len(assetPegWallet) != 0 {
			sellerAssetWallet = sdk.AddAssetPegToWallet(&assetPegWallet[0], sellerAssetWallet)
			orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])
		}

		setFiatWallet(ctx, am, buyerAddress, buyerFiatWallet)
		setAssetWallet(ctx, am, sellerAddress, sellerAssetWallet)
	}
	tags := sdk.NewTags("buyer", []byte(buyerAddress.String()))
	tags = tags.AppendTag("seller", []byte(sellerAddress.String()))
	tags = tags.AppendTag("assetPegHash", []byte(pegHash.String()))
	tags = tags.AppendTag("executed", []byte(strconv.FormatBool(executed)))
	tags = tags.AppendTag("assetPrice", []byte(strconv.FormatInt(negotiation.GetBid(), 10)))
	tags = tags.AppendTag("reversed", []byte(strconv.FormatBool(reverseOrder)))
	return tags, nil, fiatPegWallet, assetPegWallet

}

//#####ACL
