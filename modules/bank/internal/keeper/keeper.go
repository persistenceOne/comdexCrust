package keeper

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth/exported"
	bankTypes "github.com/commitHub/commitBlockchain/modules/bank/internal/types"
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	"github.com/commitHub/commitBlockchain/modules/orders"
	"github.com/commitHub/commitBlockchain/modules/params"
	"github.com/commitHub/commitBlockchain/modules/reputation"
	"github.com/commitHub/commitBlockchain/types"
)

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper

	DelegateCoins(ctx cTypes.Context, delegatorAddr, moduleAccAddr cTypes.AccAddress, amt cTypes.Coins) cTypes.Error
	UndelegateCoins(ctx cTypes.Context, moduleAccAddr, delegatorAddr cTypes.AccAddress, amt cTypes.Coins) cTypes.Error
}

// BaseKeeper manages transfers between accounts. It implements the Keeper interface.
type BaseKeeper struct {
	BaseSendKeeper

	ak         bankTypes.AccountKeeper
	paramSpace params.Subspace
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(ak bankTypes.AccountKeeper, nk negotiation.Keeper, aclK acl.Keeper, orderKeeper orders.Keeper,
	rk reputation.Keeper, paramSpace params.Subspace, codespace cTypes.CodespaceType) BaseKeeper {

	ps := paramSpace.WithKeyTable(bankTypes.ParamKeyTable())
	return BaseKeeper{
		BaseSendKeeper: NewBaseSendKeeper(ak, nk, aclK, orderKeeper, rk, ps, codespace),
		ak:             ak,
		paramSpace:     ps,
	}
}

// DelegateCoins performs delegation by deducting amt coins from an account with
// address addr. For vesting accounts, delegations amounts are tracked for both
// vesting and vested coins.
// The coins are then transferred from the delegator address to a ModuleAccount address.
// If any of the delegation amounts are negative, an error is returned.
func (keeper BaseKeeper) DelegateCoins(ctx cTypes.Context, delegatorAddr, moduleAccAddr cTypes.AccAddress, amt cTypes.Coins) cTypes.Error {

	delegatorAcc := keeper.ak.GetAccount(ctx, delegatorAddr)
	if delegatorAcc == nil {
		return cTypes.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", delegatorAddr))
	}

	moduleAcc := keeper.ak.GetAccount(ctx, moduleAccAddr)
	if moduleAcc == nil {
		return cTypes.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleAccAddr))
	}

	if !amt.IsValid() {
		return cTypes.ErrInvalidCoins(amt.String())
	}

	oldCoins := delegatorAcc.GetCoins()

	_, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return cTypes.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	if err := trackDelegation(delegatorAcc, ctx.BlockHeader().Time, amt); err != nil {
		return cTypes.ErrInternal(fmt.Sprintf("failed to track delegation: %v", err))
	}

	keeper.ak.SetAccount(ctx, delegatorAcc)

	_, err := keeper.AddCoins(ctx, moduleAccAddr, amt)
	if err != nil {
		return err
	}

	return nil
}

// UndelegateCoins performs undelegation by crediting amt coins to an account with
// address addr. For vesting accounts, undelegation amounts are tracked for both
// vesting and vested coins.
// The coins are then transferred from a ModuleAccount address to the delegator address.
// If any of the undelegation amounts are negative, an error is returned.
func (keeper BaseKeeper) UndelegateCoins(ctx cTypes.Context, moduleAccAddr, delegatorAddr cTypes.AccAddress, amt cTypes.Coins) cTypes.Error {

	delegatorAcc := keeper.ak.GetAccount(ctx, delegatorAddr)
	if delegatorAcc == nil {
		return cTypes.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", delegatorAddr))
	}

	moduleAcc := keeper.ak.GetAccount(ctx, moduleAccAddr)
	if moduleAcc == nil {
		return cTypes.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleAccAddr))
	}

	if !amt.IsValid() {
		return cTypes.ErrInvalidCoins(amt.String())
	}

	oldCoins := moduleAcc.GetCoins()

	newCoins, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return cTypes.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err := keeper.SetCoins(ctx, moduleAccAddr, newCoins)
	if err != nil {
		return err
	}

	if err := trackUndelegation(delegatorAcc, amt); err != nil {
		return cTypes.ErrInternal(fmt.Sprintf("failed to track undelegation: %v", err))
	}

	keeper.ak.SetAccount(ctx, delegatorAcc)
	return nil
}

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper

	InputOutputCoins(ctx cTypes.Context, inputs []bankTypes.Input, outputs []bankTypes.Output) cTypes.Error
	SendCoins(ctx cTypes.Context, fromAddr cTypes.AccAddress, toAddr cTypes.AccAddress, amt cTypes.Coins) cTypes.Error

	SubtractCoins(ctx cTypes.Context, addr cTypes.AccAddress, amt cTypes.Coins) (cTypes.Coins, cTypes.Error)
	AddCoins(ctx cTypes.Context, addr cTypes.AccAddress, amt cTypes.Coins) (cTypes.Coins, cTypes.Error)
	SetCoins(ctx cTypes.Context, addr cTypes.AccAddress, amt cTypes.Coins) cTypes.Error

	GetSendEnabled(ctx cTypes.Context) bool
	SetSendEnabled(ctx cTypes.Context, enabled bool)

	IssueAssetsToWallets(ctx cTypes.Context, issueAsset bankTypes.IssueAsset) (types.AssetPeg, cTypes.Error)
	IssueFiatsToWallets(ctx cTypes.Context, issueFiat bankTypes.IssueFiat) cTypes.Error

	RedeemAssetsFromWallets(ctx cTypes.Context, redeemAsset bankTypes.RedeemAsset) cTypes.Error
	RedeemFiatsFromWallets(ctx cTypes.Context, redeemFiat bankTypes.RedeemFiat) cTypes.Error

	SendAssetsToWallets(ctx cTypes.Context, sendAsset bankTypes.SendAsset) cTypes.Error
	SendFiatsToWallets(ctx cTypes.Context, sendFiat bankTypes.SendFiat) cTypes.Error

	BuyerExecuteTradeOrder(ctx cTypes.Context, buyerExecuteOrder bankTypes.BuyerExecuteOrder) (cTypes.Error, []types.FiatPegWallet)
	SellerExecuteTradeOrder(ctx cTypes.Context, sellerExecuteOrder bankTypes.SellerExecuteOrder) (cTypes.Error, []types.AssetPegWallet)

	ReleaseLockedAssets(ctx cTypes.Context, releaseAsset bankTypes.ReleaseAsset) cTypes.Error
	DefineZones(ctx cTypes.Context, defineZone bankTypes.DefineZone) cTypes.Error
	DefineOrganizations(ctx cTypes.Context, defineOrganization bankTypes.DefineOrganization) cTypes.Error
	DefineACLs(ctx cTypes.Context, defineACL bankTypes.DefineACL) cTypes.Error
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	BaseViewKeeper

	ak               bankTypes.AccountKeeper
	nk               negotiation.Keeper
	aclKeeper        acl.Keeper
	orderKeeper      orders.Keeper
	reputationKeeper reputation.Keeper
	paramSpace       params.Subspace
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(ak bankTypes.AccountKeeper, nk negotiation.Keeper, aclK acl.Keeper, orderKeeper orders.Keeper,
	rk reputation.Keeper, paramSpace params.Subspace, codespace cTypes.CodespaceType) BaseSendKeeper {

	return BaseSendKeeper{
		BaseViewKeeper:   NewBaseViewKeeper(ak, codespace),
		ak:               ak,
		paramSpace:       paramSpace,
		nk:               nk,
		orderKeeper:      orderKeeper,
		reputationKeeper: rk,
		aclKeeper:        aclK,
	}
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseSendKeeper) InputOutputCoins(ctx cTypes.Context, inputs []bankTypes.Input, outputs []bankTypes.Output) cTypes.Error {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// Check supply invariant and validity of Coins.
	if err := bankTypes.ValidateInputsOutputs(inputs, outputs); err != nil {
		return err
	}

	for _, in := range inputs {
		_, err := keeper.SubtractCoins(ctx, in.Address, in.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			cTypes.NewEvent(
				cTypes.EventTypeMessage,
				cTypes.NewAttribute(bankTypes.AttributeKeySender, in.Address.String()),
			),
		)
	}

	for _, out := range outputs {
		_, err := keeper.AddCoins(ctx, out.Address, out.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			cTypes.NewEvent(
				bankTypes.EventTypeTransfer,
				cTypes.NewAttribute(bankTypes.AttributeKeyRecipient, out.Address.String()),
			),
		)
	}

	return nil
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(ctx cTypes.Context, fromAddr cTypes.AccAddress, toAddr cTypes.AccAddress, amt cTypes.Coins) cTypes.Error {

	_, err := keeper.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err
	}

	_, err = keeper.AddCoins(ctx, toAddr, amt)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(cTypes.Events{
		cTypes.NewEvent(
			bankTypes.EventTypeTransfer,
			cTypes.NewAttribute(bankTypes.AttributeKeyRecipient, toAddr.String()),
		),
		cTypes.NewEvent(
			cTypes.EventTypeMessage,
			cTypes.NewAttribute(bankTypes.AttributeKeySender, fromAddr.String()),
		),
	})

	return nil
}

// SubtractCoins subtracts amt from the coins at the addr.
//
// CONTRACT: If the account is a vesting account, the amount has to be spendable.
func (keeper BaseSendKeeper) SubtractCoins(ctx cTypes.Context, addr cTypes.AccAddress, amt cTypes.Coins) (cTypes.Coins, cTypes.Error) {

	if !amt.IsValid() {
		return nil, cTypes.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := cTypes.NewCoins(), cTypes.NewCoins()

	acc := keeper.ak.GetAccount(ctx, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}

	// For non-vesting accounts, spendable coins will simply be the original coins.
	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, cTypes.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoins := oldCoins.Sub(amt) // should not panic as spendable coins was already checked
	err := keeper.SetCoins(ctx, addr, newCoins)

	return newCoins, err
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseSendKeeper) AddCoins(ctx cTypes.Context, addr cTypes.AccAddress, amt cTypes.Coins) (cTypes.Coins, cTypes.Error) {

	if !amt.IsValid() {
		return nil, cTypes.ErrInvalidCoins(amt.String())
	}

	oldCoins := keeper.GetCoins(ctx, addr)
	newCoins := oldCoins.Add(amt)

	if newCoins.IsAnyNegative() {
		return amt, cTypes.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err := keeper.SetCoins(ctx, addr, newCoins)
	return newCoins, err
}

// SetCoins sets the coins at the addr.
func (keeper BaseSendKeeper) SetCoins(ctx cTypes.Context, addr cTypes.AccAddress, amt cTypes.Coins) cTypes.Error {

	if !amt.IsValid() {
		return cTypes.ErrInvalidCoins(amt.String())
	}

	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}

	err := acc.SetCoins(amt)
	if err != nil {
		panic(err)
	}

	keeper.ak.SetAccount(ctx, acc)
	return nil
}

// GetSendEnabled returns the current SendEnabled
// nolint: errcheck
func (keeper BaseSendKeeper) GetSendEnabled(ctx cTypes.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, bankTypes.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper BaseSendKeeper) SetSendEnabled(ctx cTypes.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, bankTypes.ParamStoreKeySendEnabled, &enabled)
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx cTypes.Context, addr cTypes.AccAddress) cTypes.Coins
	HasCoins(ctx cTypes.Context, addr cTypes.AccAddress, amt cTypes.Coins) bool

	Codespace() cTypes.CodespaceType
}

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	ak        bankTypes.AccountKeeper
	codespace cTypes.CodespaceType
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(ak bankTypes.AccountKeeper, codespace cTypes.CodespaceType) BaseViewKeeper {
	return BaseViewKeeper{ak: ak, codespace: codespace}
}

// Logger returns a module-specific logger.
func (keeper BaseViewKeeper) Logger(ctx cTypes.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", bankTypes.ModuleName))
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx cTypes.Context, addr cTypes.AccAddress) cTypes.Coins {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return cTypes.NewCoins()
	}
	return acc.GetCoins()
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx cTypes.Context, addr cTypes.AccAddress, amt cTypes.Coins) bool {
	return keeper.GetCoins(ctx, addr).IsAllGTE(amt)
}

// Codespace returns the keeper's codespace.
func (keeper BaseViewKeeper) Codespace() cTypes.CodespaceType {
	return keeper.codespace
}

// CONTRACT: assumes that amt is valid.
func trackDelegation(acc exported.Account, blockTime time.Time, amt cTypes.Coins) error {
	vacc, ok := acc.(exported.VestingAccount)
	if ok {
		// TODO: return error on account.TrackDelegation
		vacc.TrackDelegation(blockTime, amt)
		return nil
	}

	return acc.SetCoins(acc.GetCoins().Sub(amt))
}

// CONTRACT: assumes that amt is valid.
func trackUndelegation(acc exported.Account, amt cTypes.Coins) error {
	vacc, ok := acc.(exported.VestingAccount)
	if ok {
		// TODO: return error on account.TrackUndelegation
		vacc.TrackUndelegation(amt)
		return nil
	}

	return acc.SetCoins(acc.GetCoins().Add(amt))
}

// ######################## commit ##########################
func getAssetWallet(ctx cTypes.Context, keeper BaseSendKeeper, addr cTypes.AccAddress) types.AssetPegWallet {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return types.AssetPegWallet{}
	}
	return acc.GetAssetPegWallet()
}

func setAssetWallet(ctx cTypes.Context, keeper BaseSendKeeper, addr cTypes.AccAddress, asset types.AssetPegWallet) cTypes.Error {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}
	_ = acc.SetAssetPegWallet(asset)
	keeper.ak.SetAccount(ctx, acc)
	return nil
}

func getFiatWallet(ctx cTypes.Context, keeper BaseSendKeeper, addr cTypes.AccAddress) types.FiatPegWallet {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return types.FiatPegWallet{}
	}
	return acc.GetFiatPegWallet()
}

func setFiatWallet(ctx cTypes.Context, keeper BaseSendKeeper, addr cTypes.AccAddress, fiat types.FiatPegWallet) cTypes.Error {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}
	_ = acc.SetFiatPegWallet(fiat)
	keeper.ak.SetAccount(ctx, acc)
	return nil
}

// IssueAssetsToWallets handles a list of IssueAsset messages
func (keeper BaseSendKeeper) IssueAssetsToWallets(ctx cTypes.Context, issueAsset bankTypes.IssueAsset) (types.AssetPeg, cTypes.Error) {
	var _acl acl.ACL
	var err cTypes.Error

	moderated := issueAsset.AssetPeg.GetModerated()
	if moderated {
		_acl, err = keeper.aclKeeper.CheckZoneAndGetACL(ctx, issueAsset.IssuerAddress, issueAsset.ToAddress)
		if err != nil {
			return nil, err
		}
	} else {
		aclAccount, err := keeper.aclKeeper.GetAccountACLDetails(ctx, issueAsset.IssuerAddress)
		if err != nil {
			return nil, err
		}
		_acl = aclAccount.GetACL()
	}
	if !_acl.IssueAsset {
		return nil, cTypes.ErrInternal(fmt.Sprintf("Assets cant be issued to account %v.", issueAsset.ToAddress.String()))
	}
	issuedAssetPeg, err := instantiateAndAssignAsset(ctx, issueAsset.IssuerAddress, issueAsset.ToAddress, issueAsset.AssetPeg, keeper)
	if err != nil {
		return nil, err
	}

	return issuedAssetPeg, nil

}

func instantiateAndAssignAsset(ctx cTypes.Context, issuerAddress cTypes.AccAddress, toAddress cTypes.AccAddress, assetPeg types.AssetPeg, keeper BaseSendKeeper) (types.AssetPeg, cTypes.Error) {
	pegHash, _ := types.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(keeper.ak.GetNextAssetPegHash(ctx))))
	_ = assetPeg.SetPegHash(pegHash)
	_ = assetPeg.SetLocked(assetPeg.GetModerated())
	receiverAssetPegWallet := getAssetWallet(ctx, keeper, toAddress)
	receiverAssetPegWallet = types.AddAssetPegToWallet(assetPeg, receiverAssetPegWallet)
	_ = setAssetWallet(ctx, keeper, toAddress, receiverAssetPegWallet)

	ctx.EventManager().EmitEvent(
		cTypes.NewEvent(
			bankTypes.EventTypeIssueAsset,
			cTypes.NewAttribute("recipient", toAddress.String()),
			cTypes.NewAttribute("issuer", issuerAddress.String()),
			cTypes.NewAttribute("asset", assetPeg.GetPegHash().String()),
		))
	return assetPeg, nil
}

func (keeper BaseSendKeeper) IssueFiatsToWallets(ctx cTypes.Context, issueFiat bankTypes.IssueFiat) cTypes.Error {
	_acl, err := keeper.aclKeeper.CheckZoneAndGetACL(ctx, issueFiat.IssuerAddress, issueFiat.ToAddress)
	if err != nil {
		return err
	}

	if !_acl.IssueFiat {
		return cTypes.ErrInternal(fmt.Sprintf("Fiats can't be issued to account %v.", issueFiat.ToAddress.String()))
	}

	err = instantiateAndAssignFiat(ctx, keeper, issueFiat.IssuerAddress, issueFiat.ToAddress, issueFiat.FiatPeg)
	if err != nil {
		return err
	}
	return nil
}

func instantiateAndAssignFiat(ctx cTypes.Context, keeper BaseSendKeeper, issuerAddress cTypes.AccAddress,
	toAddress cTypes.AccAddress, fiatPeg types.FiatPeg) cTypes.Error {

	pegHash, _ := types.GetFiatPegHashHex(fmt.Sprintf("%x", strconv.Itoa(keeper.ak.GetNextFiatPegHash(ctx))))
	_ = fiatPeg.SetPegHash(pegHash)
	receiverFiatPegWallet := getFiatWallet(ctx, keeper, toAddress)
	receiverFiatPegWallet = types.AddFiatPegToWallet(receiverFiatPegWallet, []types.BaseFiatPeg{types.ToBaseFiatPeg(fiatPeg)})

	_ = setFiatWallet(ctx, keeper, toAddress, receiverFiatPegWallet)

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		bankTypes.EventTypeIssueFiat,
		cTypes.NewAttribute("recipient", toAddress.String()),
		cTypes.NewAttribute("issuer", issuerAddress.String()),
		cTypes.NewAttribute("fiat", fiatPeg.GetPegHash().String()),
	))
	return nil
}

func (keeper BaseSendKeeper) RedeemAssetsFromWallets(ctx cTypes.Context, redeemAsset bankTypes.RedeemAsset) cTypes.Error {

	_acl, err := keeper.aclKeeper.CheckZoneAndGetACL(ctx, redeemAsset.IssuerAddress, redeemAsset.RedeemerAddress)
	if err != nil {
		return err
	}
	if !_acl.RedeemAsset {
		return cTypes.ErrInternal(fmt.Sprintf("Assets can't be redeemed from account %v.", redeemAsset.RedeemerAddress.String()))
	}
	err = instantiateAndRedeemAsset(ctx, keeper, redeemAsset.IssuerAddress, redeemAsset.RedeemerAddress, redeemAsset.PegHash)
	if err != nil {
		return err
	}
	return nil
}

func instantiateAndRedeemAsset(ctx cTypes.Context, keeper BaseSendKeeper, issuerAddress cTypes.AccAddress,
	redeemerAddress cTypes.AccAddress, pegHash types.PegHash) cTypes.Error {

	redeemerPegHashWallet := getAssetWallet(ctx, keeper, redeemerAddress)
	issuerPegHashWallet := getAssetWallet(ctx, keeper, issuerAddress)
	var assetPeg types.AssetPeg
	length := len(redeemerPegHashWallet)
	if length == 0 {
		return cTypes.ErrInternal("No Assets Found!") // Codespace and CodeType needs to be defined
	}
	i := redeemerPegHashWallet.SearchAssetPeg(pegHash)
	if i == length {
		return cTypes.ErrInternal("No Assets With Given PegHash Found!") // Codespace and CodeType needs to be defined
	}
	assetPeg, redeemerPegHashWallet = types.SubtractAssetPegFromWallet(pegHash, redeemerPegHashWallet)
	unSetAssetPeg := types.NewBaseAssetPegWithPegHash(assetPeg.GetPegHash())
	issuerPegHashWallet = types.AddAssetPegToWallet(&unSetAssetPeg, issuerPegHashWallet)
	err := setAssetWallet(ctx, keeper, redeemerAddress, redeemerPegHashWallet)
	if err == nil {
		err = setAssetWallet(ctx, keeper, issuerAddress, issuerPegHashWallet)
	}

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		bankTypes.EventTypeRedeemAsset,
		cTypes.NewAttribute("redeemer", redeemerAddress.String()),
		cTypes.NewAttribute("issuer", issuerAddress.String()),
		cTypes.NewAttribute("fiat", assetPeg.GetPegHash().String()),
	))
	return nil
}

func (keeper BaseSendKeeper) RedeemFiatsFromWallets(ctx cTypes.Context, redeemFiat bankTypes.RedeemFiat) cTypes.Error {
	_acl, err := keeper.aclKeeper.CheckZoneAndGetACL(ctx, redeemFiat.IssuerAddress, redeemFiat.RedeemerAddress)
	if err != nil {
		return err
	}
	if !_acl.RedeemFiat {
		return cTypes.ErrInternal(fmt.Sprintf("Fiats can't be redeemed from account %v.", redeemFiat.RedeemerAddress.String()))
	}
	err = instantiateAndRedeemFiat(ctx, keeper, redeemFiat.IssuerAddress, redeemFiat.RedeemerAddress, redeemFiat.Amount)
	if err != nil {
		return err
	}
	return nil
}
func instantiateAndRedeemFiat(ctx cTypes.Context, keeper BaseSendKeeper, issuerAddress cTypes.AccAddress,
	redeemerAddress cTypes.AccAddress, amount int64) cTypes.Error {

	fromOldFiatWallet := getFiatWallet(ctx, keeper, redeemerAddress)

	emptiedFiatPegWallet, redeemerFiatPegWallet := types.RedeemAmountFromWallet(amount, fromOldFiatWallet)
	if len(redeemerFiatPegWallet) == 0 && len(emptiedFiatPegWallet) == 0 {
		return cTypes.ErrInsufficientCoins(fmt.Sprintf("Redeemed amount higher than the account balance"))
	}

	err := setFiatWallet(ctx, keeper, redeemerAddress, redeemerFiatPegWallet)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		bankTypes.EventTypeRedeemFiat,
		cTypes.NewAttribute("redeemer", redeemerAddress.String()),
	))
	return nil
}

func (keeper BaseSendKeeper) SendAssetsToWallets(ctx cTypes.Context, sendAsset bankTypes.SendAsset) cTypes.Error {
	aclStore, err := keeper.aclKeeper.GetAccountACLDetails(ctx, sendAsset.FromAddress)
	if err != nil {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	account := aclStore.GetACL()
	if !account.SendAsset {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	err = sendAssetToOrder(ctx, keeper, sendAsset.FromAddress, sendAsset.ToAddress, sendAsset.PegHash)
	if err != nil {
		return err
	}
	keeper.reputationKeeper.SetSendAssetsPositiveTx(ctx, sendAsset.FromAddress)
	return nil
}

func sendAssetToOrder(ctx cTypes.Context, keeper BaseSendKeeper, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress,
	pegHash types.PegHash) cTypes.Error {

	_negotiation, err := keeper.nk.GetNegotiationDetails(ctx, toAddress, fromAddress, pegHash)
	if err != nil {
		return err
	}
	_time := _negotiation.GetTime() + _negotiation.GetSellerBlockHeight()
	if ctx.BlockHeight() > _time {
		return cTypes.ErrInvalidSequence("Negotiation time expired.")
	}
	if _negotiation.GetSellerSignature() == nil || _negotiation.GetBuyerSignature() == nil {
		return cTypes.ErrInternal("Signatures are not present")
	}

	fromOldAssetWallet := getAssetWallet(ctx, keeper, fromAddress)
	sentAsset, fromNewAssetPegWallet := types.SubtractAssetPegFromWallet(pegHash, fromOldAssetWallet)
	if sentAsset == nil {
		return cTypes.ErrInsufficientCoins("Asset not found.")
	}
	if sentAsset.GetLocked() {
		return cTypes.ErrInsufficientCoins("Asset locked.")
	}
	err = keeper.orderKeeper.SendAssetsToOrder(ctx, fromAddress, toAddress, sentAsset)
	if err == nil {
		err = setAssetWallet(ctx, keeper, fromAddress, fromNewAssetPegWallet)
	}

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		bankTypes.EventTypeSendAsset,
		cTypes.NewAttribute("recipient", toAddress.String()),
		cTypes.NewAttribute("sender", fromAddress.String()),
		cTypes.NewAttribute("asset", sentAsset.GetPegHash().String()),
	))

	return err
}

func (keeper BaseSendKeeper) SendFiatsToWallets(ctx cTypes.Context, sendFiat bankTypes.SendFiat) cTypes.Error {

	_acl, err := keeper.aclKeeper.GetAccountACLDetails(ctx, sendFiat.FromAddress)
	if err != nil {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	account := _acl.GetACL()
	if !account.SendFiat {
		return cTypes.ErrInternal("Unauthorized transaction")
	}
	err = sendFiatToOrder(ctx, keeper, sendFiat.FromAddress, sendFiat.ToAddress, sendFiat.PegHash, sendFiat.Amount)
	if err != nil {
		return err
	}

	keeper.reputationKeeper.SetSendFiatsPositiveTx(ctx, sendFiat.FromAddress)
	return nil
}

func sendFiatToOrder(ctx cTypes.Context, keeper BaseSendKeeper, fromAddress cTypes.AccAddress, toAddress cTypes.AccAddress,
	pegHash types.PegHash, amount int64) cTypes.Error {

	_negotiation, err := keeper.nk.GetNegotiationDetails(ctx, fromAddress, toAddress, pegHash)
	if err != nil {
		return err
	}
	_time := _negotiation.GetTime() + _negotiation.GetBuyerBlockHeight()
	if ctx.BlockHeight() > _time {
		return cTypes.ErrInvalidSequence("Negotiation time expired.")
	}
	if _negotiation.GetSellerSignature() == nil || _negotiation.GetBuyerSignature() == nil {
		return cTypes.ErrInternal("Signatures are not present")
	}

	fromOldFiatWallet := getFiatWallet(ctx, keeper, fromAddress)
	sentFiatPegWallet, oldFiatPegWallet := types.SubtractAmountFromWallet(amount, fromOldFiatWallet)
	if len(sentFiatPegWallet) == 0 && len(oldFiatPegWallet) == 0 {
		return cTypes.ErrInsufficientCoins(fmt.Sprintf("Insufficient funds"))
	}

	err = keeper.orderKeeper.SendFiatsToOrder(ctx, fromAddress, toAddress, pegHash, sentFiatPegWallet)
	if err == nil {
		err = setFiatWallet(ctx, keeper, fromAddress, oldFiatPegWallet)
	}

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		bankTypes.EventTypeSendFiat,
		cTypes.NewAttribute("recipient", toAddress.String()),
		cTypes.NewAttribute("sender", fromAddress.String()),
	))

	return nil
}

func (keeper BaseSendKeeper) BuyerExecuteTradeOrder(ctx cTypes.Context, buyerExecuteOrder bankTypes.BuyerExecuteOrder) (
	cTypes.Error, []types.FiatPegWallet) {

	var fiatPegWallets []types.FiatPegWallet
	var _acl acl.ACL
	var err cTypes.Error

	_, assetWallet, _, _, _ := keeper.orderKeeper.GetOrderDetails(ctx, buyerExecuteOrder.BuyerAddress,
		buyerExecuteOrder.SellerAddress, buyerExecuteOrder.PegHash)

	if len(assetWallet) == 0 {
		return cTypes.ErrInsufficientCoins("Asset token not found!"), fiatPegWallets
	}
	i := assetWallet.SearchAssetPeg(buyerExecuteOrder.PegHash)
	if i < len(assetWallet) && assetWallet[i].GetPegHash().String() == buyerExecuteOrder.PegHash.String() {
		assetPeg := assetWallet[i]
		moderated := assetPeg.GetModerated()
		if !moderated {
			aclAccount, err := keeper.aclKeeper.GetAccountACLDetails(ctx, buyerExecuteOrder.BuyerAddress)
			if err != nil {
				return err, fiatPegWallets
			}
			if !reflect.DeepEqual(buyerExecuteOrder.MediatorAddress, buyerExecuteOrder.BuyerAddress) {
				return cTypes.ErrUnauthorized(fmt.Sprintf("Trade cannot be executed for account %v."+
					" Access Denied.", buyerExecuteOrder.MediatorAddress.String())), fiatPegWallets
			}
			_acl = aclAccount.GetACL()
			if !_acl.BuyerExecuteOrder {
				return cTypes.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v."+
					" Access Denied.", buyerExecuteOrder.BuyerAddress.String())), fiatPegWallets
			}
			err, fiatPegWallet, _ := privateExchangeOrderTokens(ctx, keeper, buyerExecuteOrder.MediatorAddress,
				buyerExecuteOrder.BuyerAddress, buyerExecuteOrder.SellerAddress, buyerExecuteOrder.PegHash,
				buyerExecuteOrder.FiatProofHash, "")

			if err != nil {
				return err, fiatPegWallets
			}
			fiatPegWallets = append(fiatPegWallets, fiatPegWallet)
			return nil, fiatPegWallets
		}
		_acl, err = keeper.aclKeeper.CheckZoneAndGetACL(ctx, buyerExecuteOrder.MediatorAddress, buyerExecuteOrder.BuyerAddress)
		if err != nil {
			return err, fiatPegWallets

		}
	}
	if !_acl.BuyerExecuteOrder {
		return cTypes.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.",
			buyerExecuteOrder.BuyerAddress.String())), fiatPegWallets
	}

	err, fiatPegWallet, _ := exchangeOrderTokens(ctx, keeper, buyerExecuteOrder.MediatorAddress,
		buyerExecuteOrder.BuyerAddress, buyerExecuteOrder.SellerAddress, buyerExecuteOrder.PegHash, buyerExecuteOrder.FiatProofHash, "")

	if err != nil {
		return err, fiatPegWallets
	}
	fiatPegWallets = append(fiatPegWallets, fiatPegWallet)
	return nil, fiatPegWallets
}

func privateExchangeOrderTokens(ctx cTypes.Context, keeper BaseSendKeeper, mediatorAddress cTypes.AccAddress,
	buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash, fiatProofHash string,
	awbProofHash string) (cTypes.Error, types.FiatPegWallet, types.AssetPegWallet) {

	var reverseOrder bool
	var executed bool

	err, assetPegWallet, fiatPegWallet, orderFiatProofHash, orderAWBProofHash := keeper.orderKeeper.GetOrderDetails(ctx,
		buyerAddress, sellerAddress, pegHash)

	if err != nil {
		return err, fiatPegWallet, assetPegWallet
	}

	negotiationID := types.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	_negotiation, err := keeper.nk.GetNegotiation(ctx, negotiationID)
	if err != nil {
		return err, fiatPegWallet, assetPegWallet
	}

	if len(assetPegWallet) != 1 || assetPegWallet[0].GetPegHash().String() != pegHash.String() {
		if _negotiation.GetTime() < ctx.BlockHeight() {
			return cTypes.ErrInsufficientCoins("Asset token not found!"), fiatPegWallet, assetPegWallet
		}
		reverseOrder = true
		keeper.reputationKeeper.SetSellerExecuteOrderNegativeTx(ctx, sellerAddress)

	}

	buyerAssetWallet := getAssetWallet(ctx, keeper, buyerAddress)

	if orderFiatProofHash == "" && fiatProofHash != "" {
		keeper.orderKeeper.SetOrderFiatProofHash(ctx, buyerAddress, sellerAddress, pegHash, fiatProofHash)
	}
	if orderAWBProofHash == "" && awbProofHash != "" {
		keeper.orderKeeper.SetOrderAWBProofHash(ctx, buyerAddress, sellerAddress, pegHash, awbProofHash)
	}
	err, _, _, orderFiatProofHash, orderAWBProofHash = keeper.orderKeeper.GetOrderDetails(ctx, buyerAddress, sellerAddress, pegHash)
	if err != nil {
		return err, nil, nil
	}
	if orderFiatProofHash != "" && orderAWBProofHash != "" {
		executed = true
		keeper.reputationKeeper.SetSellerExecuteOrderPositiveTx(ctx, sellerAddress)
		keeper.reputationKeeper.SetBuyerExecuteOrderPositiveTx(ctx, buyerAddress)

		buyerAssetWallet = types.AddAssetPegToWallet(&assetPegWallet[0], buyerAssetWallet)
		assetPegWallet = keeper.orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])

		_ = setAssetWallet(ctx, keeper, buyerAddress, buyerAssetWallet)
	}
	if orderAWBProofHash == "" && orderFiatProofHash == "" {
		reverseOrder = true
	}
	if executed == true || reverseOrder == true {
		buyerFiatWallet := getFiatWallet(ctx, keeper, buyerAddress)
		sellerAssetWallet := getAssetWallet(ctx, keeper, sellerAddress)

		if len(fiatPegWallet) != 0 {
			buyerFiatWallet = types.AddFiatPegToWallet(buyerFiatWallet, fiatPegWallet)
			keeper.orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
		}
		if len(assetPegWallet) != 0 {
			sellerAssetWallet = types.AddAssetPegToWallet(&assetPegWallet[0], sellerAssetWallet)
			keeper.orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])
		}

		_ = setFiatWallet(ctx, keeper, buyerAddress, buyerFiatWallet)
		_ = setAssetWallet(ctx, keeper, sellerAddress, sellerAssetWallet)
	}

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		bankTypes.EventTypeExecuteOrder,
		cTypes.NewAttribute("buyer", buyerAddress.String()),
		cTypes.NewAttribute("seller", sellerAddress.String()),
		cTypes.NewAttribute("assetPegHash", pegHash.String()),
		cTypes.NewAttribute("executed", strconv.FormatBool(executed)),
		cTypes.NewAttribute("assetPrice", strconv.FormatInt(_negotiation.GetBid(), 10)),
		cTypes.NewAttribute("reversed", strconv.FormatBool(reverseOrder)),
	))

	return nil, fiatPegWallet, assetPegWallet

}

func exchangeOrderTokens(ctx cTypes.Context, keeper BaseSendKeeper, mediatorAddress cTypes.AccAddress,
	buyerAddress cTypes.AccAddress, sellerAddress cTypes.AccAddress, pegHash types.PegHash, fiatProofHash string,
	awbProofHash string) (cTypes.Error, types.FiatPegWallet, types.AssetPegWallet) {

	err, assetPegWallet, fiatPegWallet, orderFiatProofHash, orderAWBProofHash := keeper.orderKeeper.GetOrderDetails(ctx, buyerAddress, sellerAddress, pegHash)
	if err != nil {
		return err, fiatPegWallet, assetPegWallet
	}

	negotiationID := types.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	_negotiation, err := keeper.nk.GetNegotiation(ctx, negotiationID)
	if err != nil {
		return err, fiatPegWallet, assetPegWallet
	}

	var reverseOrder bool
	var oldFiatPegWallet types.FiatPegWallet
	if len(fiatPegWallet) == 0 || _negotiation.GetBid() > types.GetFiatPegWalletBalance(fiatPegWallet) {
		if _negotiation.GetTime() < ctx.BlockHeight() {
			return cTypes.ErrInsufficientCoins("Fiat tokens not found!"), fiatPegWallet, assetPegWallet
		}
		reverseOrder = true
		keeper.reputationKeeper.SetBuyerExecuteOrderNegativeTx(ctx, buyerAddress)
	}
	if len(assetPegWallet) != 1 || assetPegWallet[0].GetPegHash().String() != pegHash.String() {
		if _negotiation.GetTime() < ctx.BlockHeight() {
			return cTypes.ErrInsufficientCoins("Asset token not found!"), fiatPegWallet, assetPegWallet
		}
		reverseOrder = true
		keeper.reputationKeeper.SetSellerExecuteOrderNegativeTx(ctx, sellerAddress)

	}
	buyerTime := _negotiation.GetTime() + _negotiation.GetBuyerBlockHeight()
	sellerTime := _negotiation.GetTime() + _negotiation.GetSellerBlockHeight()
	_time := ctx.BlockHeight()
	if _time > buyerTime || _time > sellerTime {
		reverseOrder = true
	}

	if _negotiation.GetBid() < types.GetFiatPegWalletBalance(fiatPegWallet) {
		fiatPegWallet, oldFiatPegWallet = types.SubtractAmountFromWallet(_negotiation.GetBid(), fiatPegWallet)
	}
	var executed bool
	if !reverseOrder {
		sellerFiatWallet := getFiatWallet(ctx, keeper, sellerAddress)
		buyerAssetWallet := getAssetWallet(ctx, keeper, buyerAddress)

		if orderFiatProofHash == "" && fiatProofHash != "" {
			keeper.orderKeeper.SetOrderFiatProofHash(ctx, buyerAddress, sellerAddress, pegHash, fiatProofHash)
		}
		if orderAWBProofHash == "" && awbProofHash != "" {
			keeper.orderKeeper.SetOrderAWBProofHash(ctx, buyerAddress, sellerAddress, pegHash, awbProofHash)
		}
		err, _, _, orderFiatProofHash, orderAWBProofHash = keeper.orderKeeper.GetOrderDetails(ctx, buyerAddress,
			sellerAddress, pegHash)

		if err != nil {
			return err, nil, nil
		}
		if orderFiatProofHash != "" && orderAWBProofHash != "" {
			executed = true
			keeper.reputationKeeper.SetSellerExecuteOrderPositiveTx(ctx, sellerAddress)
			keeper.reputationKeeper.SetBuyerExecuteOrderPositiveTx(ctx, buyerAddress)

			buyerAssetWallet = types.AddAssetPegToWallet(&assetPegWallet[0], buyerAssetWallet)
			sellerFiatWallet = types.AddFiatPegToWallet(sellerFiatWallet, fiatPegWallet)

			fiatPegWallet = keeper.orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
			assetPegWallet = keeper.orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])
		}

		_ = setFiatWallet(ctx, keeper, sellerAddress, sellerFiatWallet)
		_ = setAssetWallet(ctx, keeper, buyerAddress, buyerAssetWallet)

	}

	if executed == true || reverseOrder == true {
		buyerFiatWallet := getFiatWallet(ctx, keeper, buyerAddress)
		sellerAssetWallet := getAssetWallet(ctx, keeper, sellerAddress)
		if len(oldFiatPegWallet) != 0 {
			buyerFiatWallet = types.AddFiatPegToWallet(buyerFiatWallet, oldFiatPegWallet)
		}

		if len(fiatPegWallet) != 0 {
			buyerFiatWallet = types.AddFiatPegToWallet(buyerFiatWallet, fiatPegWallet)
			keeper.orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
		}
		if len(assetPegWallet) != 0 {
			sellerAssetWallet = types.AddAssetPegToWallet(&assetPegWallet[0], sellerAssetWallet)
			keeper.orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])
		}

		_ = setFiatWallet(ctx, keeper, buyerAddress, buyerFiatWallet)
		_ = setAssetWallet(ctx, keeper, sellerAddress, sellerAssetWallet)
	}

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		bankTypes.EventTypeExecuteOrder,
		cTypes.NewAttribute("buyer", buyerAddress.String()),
		cTypes.NewAttribute("seller", sellerAddress.String()),
		cTypes.NewAttribute("assetPegHash", pegHash.String()),
		cTypes.NewAttribute("executed", strconv.FormatBool(executed)),
		cTypes.NewAttribute("assetPrice", strconv.FormatInt(_negotiation.GetBid(), 10)),
		cTypes.NewAttribute("reversed", strconv.FormatBool(reverseOrder)),
	))

	return nil, fiatPegWallet, assetPegWallet
}

func (keeper BaseSendKeeper) SellerExecuteTradeOrder(ctx cTypes.Context, sellerExecuteOrder bankTypes.SellerExecuteOrder) (
	cTypes.Error, []types.AssetPegWallet) {

	var assetPegWallets []types.AssetPegWallet
	var _acl acl.ACL
	var err cTypes.Error
	var assetPegWallet types.AssetPegWallet

	_, assetWallet, _, _, _ := keeper.orderKeeper.GetOrderDetails(ctx, sellerExecuteOrder.BuyerAddress,
		sellerExecuteOrder.SellerAddress, sellerExecuteOrder.PegHash)

	if len(assetWallet) == 0 {
		return cTypes.ErrInsufficientCoins("Asset token not found!"), assetPegWallets
	}
	i := assetWallet.SearchAssetPeg(sellerExecuteOrder.PegHash)
	if i < len(assetWallet) && assetWallet[i].GetPegHash().String() == sellerExecuteOrder.PegHash.String() {
		assetPeg := assetWallet[i]
		moderated := assetPeg.GetModerated()
		if !moderated {
			aclAccount, err := keeper.aclKeeper.GetAccountACLDetails(ctx, sellerExecuteOrder.SellerAddress)
			if err != nil {
				return err, assetPegWallets
			}
			if !reflect.DeepEqual(sellerExecuteOrder.MediatorAddress, sellerExecuteOrder.SellerAddress) {
				return cTypes.ErrUnauthorized(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.",
					sellerExecuteOrder.MediatorAddress.String())), assetPegWallets
			}
			_acl = aclAccount.GetACL()
			if !_acl.SellerExecuteOrder {
				return cTypes.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.",
					sellerExecuteOrder.SellerAddress.String())), assetPegWallets
			}
			err, _, assetPegWallet := privateExchangeOrderTokens(ctx, keeper,
				sellerExecuteOrder.MediatorAddress, sellerExecuteOrder.BuyerAddress, sellerExecuteOrder.SellerAddress,
				sellerExecuteOrder.PegHash, "", sellerExecuteOrder.AWBProofHash)

			if err != nil {
				return err, assetPegWallets
			}
			assetPegWallets = append(assetPegWallets, assetPegWallet)
			return nil, assetPegWallets
		}
		_acl, err = keeper.aclKeeper.CheckZoneAndGetACL(ctx, sellerExecuteOrder.MediatorAddress, sellerExecuteOrder.SellerAddress)
		if err != nil {
			return err, assetPegWallets

		}
	}

	if !_acl.SellerExecuteOrder {
		return cTypes.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.",
			sellerExecuteOrder.SellerAddress.String())), assetPegWallets
	}

	err, _, assetPegWallet = exchangeOrderTokens(ctx, keeper, sellerExecuteOrder.MediatorAddress,
		sellerExecuteOrder.BuyerAddress, sellerExecuteOrder.SellerAddress, sellerExecuteOrder.PegHash,
		"", sellerExecuteOrder.AWBProofHash)

	if err != nil {
		return err, assetPegWallets
	}
	assetPegWallets = append(assetPegWallets, assetPegWallet)
	return nil, assetPegWallets
}

func (keeper BaseSendKeeper) ReleaseLockedAssets(ctx cTypes.Context, releaseAsset bankTypes.ReleaseAsset) cTypes.Error {

	_acl, err := keeper.aclKeeper.CheckZoneAndGetACL(ctx, releaseAsset.ZoneAddress, releaseAsset.OwnerAddress)
	if err != nil {
		return err
	}
	if !_acl.ReleaseAsset {
		return cTypes.ErrInternal(fmt.Sprintf("Assets cannot be released for account %v. Access Denied.",
			releaseAsset.OwnerAddress.String()))
	}

	err = releaseAssets(ctx, keeper, releaseAsset.ZoneAddress, releaseAsset.OwnerAddress, releaseAsset.PegHash)
	if err != nil {
		return err
	}
	return nil
}

func releaseAssets(ctx cTypes.Context, keeper BaseSendKeeper, zoneAddress cTypes.AccAddress, ownerAddress cTypes.AccAddress,
	pegHash types.PegHash) cTypes.Error {

	ownerAssetWallet := getAssetWallet(ctx, keeper, ownerAddress)
	if !types.ReleaseAssetPegInWallet(ownerAssetWallet, pegHash) {
		return cTypes.ErrInternal("Asset peg not found.")
	}
	_ = setAssetWallet(ctx, keeper, ownerAddress, ownerAssetWallet)

	ctx.EventManager().EmitEvent(cTypes.NewEvent(
		bankTypes.EventTypeReleaseAsset,
		cTypes.NewAttribute("zone", zoneAddress.String()),
		cTypes.NewAttribute("owner", ownerAddress.String()),
		cTypes.NewAttribute("asset", pegHash.String()),
	))

	return nil
}

func (keeper BaseSendKeeper) DefineZones(ctx cTypes.Context, defineZone bankTypes.DefineZone) cTypes.Error {
	if !keeper.aclKeeper.CheckValidGenesisAddress(ctx, defineZone.From) {
		return cTypes.ErrInternal(fmt.Sprintf("Account %v is not the genesis account. Zones can only be"+
			" defined by the genesis account.", defineZone.From.String()))
	}

	acc := keeper.ak.GetAccount(ctx, defineZone.To)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, defineZone.To)
		keeper.ak.SetAccount(ctx, acc)
	}

	err := keeper.aclKeeper.DefineZoneAddress(ctx, defineZone.To, defineZone.ZoneID)
	if err != nil {
		return err
	}
	return nil
}

func (keeper BaseSendKeeper) DefineOrganizations(ctx cTypes.Context, defineOrganization bankTypes.DefineOrganization) cTypes.Error {
	if !keeper.aclKeeper.CheckValidZoneAddress(ctx, defineOrganization.ZoneID, defineOrganization.From) {
		return cTypes.ErrInternal(fmt.Sprintf("Account %v is not the zone account. Organizations can only "+
			"be defined by the zone account.", defineOrganization.From.String()))
	}

	acc := keeper.ak.GetAccount(ctx, defineOrganization.To)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, defineOrganization.To)
		keeper.ak.SetAccount(ctx, acc)
	}

	err := keeper.aclKeeper.DefineOrganizationAddress(ctx, defineOrganization.To, defineOrganization.OrganizationID, defineOrganization.ZoneID)

	if err != nil {
		return err
	}
	return nil
}

func (keeper BaseSendKeeper) DefineACLs(ctx cTypes.Context, defineACL bankTypes.DefineACL) cTypes.Error {
	if !keeper.aclKeeper.CheckValidGenesisAddress(ctx, defineACL.From) {
		if !keeper.aclKeeper.CheckValidZoneAddress(ctx, defineACL.ACLAccount.GetZoneID(), defineACL.From) {
			if !keeper.aclKeeper.CheckValidOrganizationAddress(ctx, defineACL.ACLAccount.GetZoneID(),
				defineACL.ACLAccount.GetOrganizationID(), defineACL.From) {

				return cTypes.ErrInternal(fmt.Sprintf("Account %v does not have access to define acl "+
					"for account %v.", defineACL.From.String(), defineACL.To.String()))
			}
		}
	}

	acc := keeper.ak.GetAccount(ctx, defineACL.To)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, defineACL.To)
		keeper.ak.SetAccount(ctx, acc)
	}

	err := keeper.aclKeeper.DefineACLAccount(ctx, defineACL.To, defineACL.ACLAccount)
	if err != nil {
		return err
	}
	return nil
}
