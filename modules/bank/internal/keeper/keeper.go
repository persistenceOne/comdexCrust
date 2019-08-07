package keeper

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
	
	"github.com/tendermint/tendermint/libs/log"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth/exported"
	"github.com/commitHub/commitBlockchain/modules/negotiation"
	"github.com/commitHub/commitBlockchain/modules/orders"
	"github.com/commitHub/commitBlockchain/modules/reputation"
	
	"github.com/commitHub/commitBlockchain/modules/params"
	
	cmTypes "github.com/commitHub/commitBlockchain/types"
	
	"github.com/commitHub/commitBlockchain/modules/bank/internal/types"
)

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper
	
	DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
}

// BaseKeeper manages transfers between accounts. It implements the Keeper interface.
type BaseKeeper struct {
	BaseSendKeeper
	
	ak         types.AccountKeeper
	paramSpace params.Subspace
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(ak types.AccountKeeper, nk negotiation.Keeper, aclK acl.Keeper, orderKeeper orders.Keeper,
	rk reputation.Keeper, paramSpace params.Subspace, codespace sdk.CodespaceType) BaseKeeper {
	
	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
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
func (keeper BaseKeeper) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	
	delegatorAcc := keeper.ak.GetAccount(ctx, delegatorAddr)
	if delegatorAcc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", delegatorAddr))
	}
	
	moduleAcc := keeper.ak.GetAccount(ctx, moduleAccAddr)
	if moduleAcc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleAccAddr))
	}
	
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	
	oldCoins := delegatorAcc.GetCoins()
	
	_, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}
	
	if err := trackDelegation(delegatorAcc, ctx.BlockHeader().Time, amt); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to track delegation: %v", err))
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
func (keeper BaseKeeper) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	
	delegatorAcc := keeper.ak.GetAccount(ctx, delegatorAddr)
	if delegatorAcc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", delegatorAddr))
	}
	
	moduleAcc := keeper.ak.GetAccount(ctx, moduleAccAddr)
	if moduleAcc == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleAccAddr))
	}
	
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	
	oldCoins := moduleAcc.GetCoins()
	
	newCoins, hasNeg := oldCoins.SafeSub(amt)
	if hasNeg {
		return sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}
	
	err := keeper.SetCoins(ctx, moduleAccAddr, newCoins)
	if err != nil {
		return err
	}
	
	if err := trackUndelegation(delegatorAcc, amt); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("failed to track undelegation: %v", err))
	}
	
	keeper.ak.SetAccount(ctx, delegatorAcc)
	return nil
}

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper
	
	InputOutputCoins(ctx sdk.Context, inputs []types.Input, outputs []types.Output) sdk.Error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error
	
	GetSendEnabled(ctx sdk.Context) bool
	SetSendEnabled(ctx sdk.Context, enabled bool)
	
	IssueAssetsToWallets(ctx sdk.Context, issueAsset types.IssueAsset) sdk.Error
	IssueFiatsToWallets(ctx sdk.Context, issueFiat types.IssueFiat) sdk.Error
	
	RedeemAssetsFromWallets(ctx sdk.Context, redeemAsset types.RedeemAsset) sdk.Error
	RedeemFiatsFromWallets(ctx sdk.Context, redeemFiat types.RedeemFiat) sdk.Error
	
	SendAssetsToWallets(ctx sdk.Context, sendAsset types.SendAsset) sdk.Error
	SendFiatsToWallets(ctx sdk.Context, sendFiat types.SendFiat) sdk.Error
	
	BuyerExecuteTradeOrder(ctx sdk.Context, buyerExecuteOrder types.BuyerExecuteOrder) (sdk.Error, []cmTypes.FiatPegWallet)
	SellerExecuteTradeOrder(ctx sdk.Context, sellerExecuteOrder types.SellerExecuteOrder) (sdk.Error, []cmTypes.AssetPegWallet)
	
	ReleaseLockedAssets(ctx sdk.Context, releaseAsset types.ReleaseAsset) sdk.Error
	DefineZones(ctx sdk.Context, defineZone types.DefineZone) sdk.Error
	DefineOrganizations(ctx sdk.Context, defineOrganization types.DefineOrganization) sdk.Error
	DefineACLs(ctx sdk.Context, defineACL types.DefineACL) sdk.Error
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	BaseViewKeeper
	
	ak               types.AccountKeeper
	nk               negotiation.Keeper
	aclKeeper        acl.Keeper
	orderKeeper      orders.Keeper
	reputationKeeper reputation.Keeper
	paramSpace       params.Subspace
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(ak types.AccountKeeper, nk negotiation.Keeper, aclK acl.Keeper, orderKeeper orders.Keeper,
	rk reputation.Keeper, paramSpace params.Subspace, codespace sdk.CodespaceType) BaseSendKeeper {
	
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
func (keeper BaseSendKeeper) InputOutputCoins(ctx sdk.Context, inputs []types.Input, outputs []types.Output) sdk.Error {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// Check supply invariant and validity of Coins.
	if err := types.ValidateInputsOutputs(inputs, outputs); err != nil {
		return err
	}
	
	for _, in := range inputs {
		_, err := keeper.SubtractCoins(ctx, in.Address, in.Coins)
		if err != nil {
			return err
		}
		
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(types.AttributeKeySender, in.Address.String()),
			),
		)
	}
	
	for _, out := range outputs {
		_, err := keeper.AddCoins(ctx, out.Address, out.Coins)
		if err != nil {
			return err
		}
		
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeTransfer,
				sdk.NewAttribute(types.AttributeKeyRecipient, out.Address.String()),
			),
		)
	}
	
	return nil
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	
	_, err := keeper.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err
	}
	
	_, err = keeper.AddCoins(ctx, toAddr, amt)
	if err != nil {
		return err
	}
	
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, toAddr.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeySender, fromAddr.String()),
		),
	})
	
	return nil
}

// SubtractCoins subtracts amt from the coins at the addr.
//
// CONTRACT: If the account is a vesting account, the amount has to be spendable.
func (keeper BaseSendKeeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error) {
	
	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}
	
	oldCoins, spendableCoins := sdk.NewCoins(), sdk.NewCoins()
	
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}
	
	// For non-vesting accounts, spendable coins will simply be the original coins.
	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, amt),
		)
	}
	
	newCoins := oldCoins.Sub(amt) // should not panic as spendable coins was already checked
	err := keeper.SetCoins(ctx, addr, newCoins)
	
	return newCoins, err
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseSendKeeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error) {
	
	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}
	
	oldCoins := keeper.GetCoins(ctx, addr)
	newCoins := oldCoins.Add(amt)
	
	if newCoins.IsAnyNegative() {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}
	
	err := keeper.SetCoins(ctx, addr, newCoins)
	return newCoins, err
}

// SetCoins sets the coins at the addr.
func (keeper BaseSendKeeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) sdk.Error {
	
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
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
func (keeper BaseSendKeeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, types.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper BaseSendKeeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeySendEnabled, &enabled)
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool
	
	Codespace() sdk.CodespaceType
}

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	ak        types.AccountKeeper
	codespace sdk.CodespaceType
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(ak types.AccountKeeper, codespace sdk.CodespaceType) BaseViewKeeper {
	return BaseViewKeeper{ak: ak, codespace: codespace}
}

// Logger returns a module-specific logger.
func (keeper BaseViewKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoins()
	}
	return acc.GetCoins()
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool {
	return keeper.GetCoins(ctx, addr).IsAllGTE(amt)
}

// Codespace returns the keeper's codespace.
func (keeper BaseViewKeeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}

// CONTRACT: assumes that amt is valid.
func trackDelegation(acc exported.Account, blockTime time.Time, amt sdk.Coins) error {
	vacc, ok := acc.(exported.VestingAccount)
	if ok {
		// TODO: return error on account.TrackDelegation
		vacc.TrackDelegation(blockTime, amt)
		return nil
	}
	
	return acc.SetCoins(acc.GetCoins().Sub(amt))
}

// CONTRACT: assumes that amt is valid.
func trackUndelegation(acc exported.Account, amt sdk.Coins) error {
	vacc, ok := acc.(exported.VestingAccount)
	if ok {
		// TODO: return error on account.TrackUndelegation
		vacc.TrackUndelegation(amt)
		return nil
	}
	
	return acc.SetCoins(acc.GetCoins().Add(amt))
}

// ######################## commit ##########################
func getAssetWallet(ctx sdk.Context, keeper BaseSendKeeper, addr sdk.AccAddress) cmTypes.AssetPegWallet {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return cmTypes.AssetPegWallet{}
	}
	return acc.GetAssetPegWallet()
}

func setAssetWallet(ctx sdk.Context, keeper BaseSendKeeper, addr sdk.AccAddress, asset cmTypes.AssetPegWallet) sdk.Error {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}
	_ = acc.SetAssetPegWallet(asset)
	keeper.ak.SetAccount(ctx, acc)
	return nil
}

func getFiatWallet(ctx sdk.Context, keeper BaseSendKeeper, addr sdk.AccAddress) cmTypes.FiatPegWallet {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return cmTypes.FiatPegWallet{}
	}
	return acc.GetFiatPegWallet()
}

func setFiatWallet(ctx sdk.Context, keeper BaseSendKeeper, addr sdk.AccAddress, fiat cmTypes.FiatPegWallet) sdk.Error {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}
	_ = acc.SetFiatPegWallet(fiat)
	keeper.ak.SetAccount(ctx, acc)
	return nil
}

// IssueAssetsToWallets handles a list of IssueAsset messages
func (keeper BaseSendKeeper) IssueAssetsToWallets(ctx sdk.Context, issueAsset types.IssueAsset) sdk.Error {
	var _acl acl.ACL
	var err sdk.Error
	
	moderated := issueAsset.AssetPeg.GetModerated()
	if moderated {
		_acl, err = keeper.aclKeeper.CheckZoneAndGetACL(ctx, issueAsset.IssuerAddress, issueAsset.ToAddress)
		if err != nil {
			return err
		}
	} else {
		aclAccount, err := keeper.aclKeeper.GetAccountACLDetails(ctx, issueAsset.IssuerAddress)
		if err != nil {
			return err
		}
		_acl = aclAccount.GetACL()
	}
	if !_acl.IssueAsset {
		return sdk.ErrInternal(fmt.Sprintf("Assets cant be issued to account %v.", issueAsset.ToAddress.String()))
	}
	err = instantiateAndAssignAsset(ctx, issueAsset.IssuerAddress, issueAsset.ToAddress, issueAsset.AssetPeg, keeper)
	if err != nil {
		return nil
	}
	
	return nil
	
}

func instantiateAndAssignAsset(ctx sdk.Context, issuerAddress sdk.AccAddress, toAddress sdk.AccAddress, assetPeg cmTypes.AssetPeg, keeper BaseSendKeeper) sdk.Error {
	pegHash, _ := cmTypes.GetAssetPegHashHex(fmt.Sprintf("%x", strconv.Itoa(keeper.ak.GetNextAssetPegHash(ctx))))
	_ = assetPeg.SetPegHash(pegHash)
	_ = assetPeg.SetLocked(assetPeg.GetModerated())
	receiverAssetPegWallet := getAssetWallet(ctx, keeper, toAddress)
	receiverAssetPegWallet = cmTypes.AddAssetPegToWallet(assetPeg, receiverAssetPegWallet)
	_ = setAssetWallet(ctx, keeper, toAddress, receiverAssetPegWallet)
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeIssueAsset,
			sdk.NewAttribute("recipient", toAddress.String()),
			sdk.NewAttribute("issuer", issuerAddress.String()),
			sdk.NewAttribute("asset", assetPeg.GetPegHash().String()),
		))
	return nil
}

func (keeper BaseSendKeeper) IssueFiatsToWallets(ctx sdk.Context, issueFiat types.IssueFiat) sdk.Error {
	_acl, err := keeper.aclKeeper.CheckZoneAndGetACL(ctx, issueFiat.IssuerAddress, issueFiat.ToAddress)
	if err != nil {
		return err
	}
	
	if !_acl.IssueFiat {
		return sdk.ErrInternal(fmt.Sprintf("Fiats can't be issued to account %v.", issueFiat.ToAddress.String()))
	}
	
	err = instantiateAndAssignFiat(ctx, keeper, issueFiat.IssuerAddress, issueFiat.ToAddress, issueFiat.FiatPeg)
	if err != nil {
		return err
	}
	return nil
}

func instantiateAndAssignFiat(ctx sdk.Context, keeper BaseSendKeeper, issuerAddress sdk.AccAddress,
	toAddress sdk.AccAddress, fiatPeg cmTypes.FiatPeg) sdk.Error {
	
	pegHash, _ := cmTypes.GetFiatPegHashHex(fmt.Sprintf("%x", strconv.Itoa(keeper.ak.GetNextFiatPegHash(ctx))))
	_ = fiatPeg.SetPegHash(pegHash)
	receiverFiatPegWallet := getFiatWallet(ctx, keeper, toAddress)
	receiverFiatPegWallet = cmTypes.AddFiatPegToWallet(receiverFiatPegWallet, []cmTypes.BaseFiatPeg{cmTypes.ToBaseFiatPeg(fiatPeg)})
	
	_ = setFiatWallet(ctx, keeper, toAddress, receiverFiatPegWallet)
	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIssueFiat,
		sdk.NewAttribute("recipient", toAddress.String()),
		sdk.NewAttribute("issuer", issuerAddress.String()),
		sdk.NewAttribute("fiat", fiatPeg.GetPegHash().String()),
	))
	return nil
}

func (keeper BaseSendKeeper) RedeemAssetsFromWallets(ctx sdk.Context, redeemAsset types.RedeemAsset) sdk.Error {
	
	_acl, err := keeper.aclKeeper.CheckZoneAndGetACL(ctx, redeemAsset.IssuerAddress, redeemAsset.RedeemerAddress)
	if err != nil {
		return err
	}
	if !_acl.RedeemAsset {
		return sdk.ErrInternal(fmt.Sprintf("Assets can't be redeemed from account %v.", redeemAsset.RedeemerAddress.String()))
	}
	err = instantiateAndRedeemAsset(ctx, keeper, redeemAsset.IssuerAddress, redeemAsset.RedeemerAddress, redeemAsset.PegHash)
	if err != nil {
		return err
	}
	return nil
}

func instantiateAndRedeemAsset(ctx sdk.Context, keeper BaseSendKeeper, issuerAddress sdk.AccAddress,
	redeemerAddress sdk.AccAddress, pegHash cmTypes.PegHash) sdk.Error {
	
	redeemerPegHashWallet := getAssetWallet(ctx, keeper, redeemerAddress)
	issuerPegHashWallet := getAssetWallet(ctx, keeper, issuerAddress)
	var assetPeg cmTypes.AssetPeg
	length := len(redeemerPegHashWallet)
	if length == 0 {
		return sdk.ErrInternal("No Assets Found!") // Codespace and CodeType needs to be defined
	}
	i := redeemerPegHashWallet.SearchAssetPeg(pegHash)
	if i == length {
		return sdk.ErrInternal("No Assets With Given PegHash Found!") // Codespace and CodeType needs to be defined
	}
	assetPeg, redeemerPegHashWallet = cmTypes.SubtractAssetPegFromWallet(pegHash, redeemerPegHashWallet)
	unSetAssetPeg := cmTypes.NewBaseAssetPegWithPegHash(assetPeg.GetPegHash())
	issuerPegHashWallet = cmTypes.AddAssetPegToWallet(&unSetAssetPeg, issuerPegHashWallet)
	err := setAssetWallet(ctx, keeper, redeemerAddress, redeemerPegHashWallet)
	if err == nil {
		err = setAssetWallet(ctx, keeper, issuerAddress, issuerPegHashWallet)
	}
	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRedeemAsset,
		sdk.NewAttribute("redeemer", redeemerAddress.String()),
		sdk.NewAttribute("issuer", issuerAddress.String()),
		sdk.NewAttribute("fiat", assetPeg.GetPegHash().String()),
	))
	return nil
}

func (keeper BaseSendKeeper) RedeemFiatsFromWallets(ctx sdk.Context, redeemFiat types.RedeemFiat) sdk.Error {
	_acl, err := keeper.aclKeeper.CheckZoneAndGetACL(ctx, redeemFiat.IssuerAddress, redeemFiat.RedeemerAddress)
	if err != nil {
		return err
	}
	if !_acl.RedeemFiat {
		return sdk.ErrInternal(fmt.Sprintf("Fiats can't be redeemed from account %v.", redeemFiat.RedeemerAddress.String()))
	}
	err = instantiateAndRedeemFiat(ctx, keeper, redeemFiat.IssuerAddress, redeemFiat.RedeemerAddress, redeemFiat.Amount)
	if err != nil {
		return err
	}
	return nil
}
func instantiateAndRedeemFiat(ctx sdk.Context, keeper BaseSendKeeper, issuerAddress sdk.AccAddress,
	redeemerAddress sdk.AccAddress, amount int64) sdk.Error {
	
	fromOldFiatWallet := getFiatWallet(ctx, keeper, redeemerAddress)
	
	emptiedFiatPegWallet, redeemerFiatPegWallet := cmTypes.RedeemAmountFromWallet(amount, fromOldFiatWallet)
	if len(redeemerFiatPegWallet) == 0 && len(emptiedFiatPegWallet) == 0 {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("Redeemed amount higher than the account balance"))
	}
	
	err := setFiatWallet(ctx, keeper, redeemerAddress, redeemerFiatPegWallet)
	if err != nil {
		return err
	}
	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRedeemFiat,
		sdk.NewAttribute("redeemer", redeemerAddress.String()),
	))
	return nil
}

func (keeper BaseSendKeeper) SendAssetsToWallets(ctx sdk.Context, sendAsset types.SendAsset) sdk.Error {
	aclStore, err := keeper.aclKeeper.GetAccountACLDetails(ctx, sendAsset.FromAddress)
	if err != nil {
		return sdk.ErrInternal("Unauthorized transaction")
	}
	account := aclStore.GetACL()
	if !account.SendAsset {
		return sdk.ErrInternal("Unauthorized transaction")
	}
	err = sendAssetToOrder(ctx, keeper, sendAsset.FromAddress, sendAsset.ToAddress, sendAsset.PegHash)
	if err != nil {
		return err
	}
	keeper.reputationKeeper.SetSendAssetsPositiveTx(ctx, sendAsset.FromAddress)
	return nil
}

func sendAssetToOrder(ctx sdk.Context, keeper BaseSendKeeper, fromAddress sdk.AccAddress, toAddress sdk.AccAddress,
	pegHash cmTypes.PegHash) sdk.Error {
	
	_negotiation, err := keeper.nk.GetNegotiationDetails(ctx, toAddress, fromAddress, pegHash)
	if err != nil {
		return err
	}
	_time := _negotiation.GetTime() + _negotiation.GetSellerBlockHeight()
	if ctx.BlockHeight() > _time {
		return sdk.ErrInvalidSequence("Negotiation time expired.")
	}
	if _negotiation.GetSellerSignature() == nil || _negotiation.GetBuyerSignature() == nil {
		return sdk.ErrInternal("Signatures are not present")
	}
	
	fromOldAssetWallet := getAssetWallet(ctx, keeper, fromAddress)
	sentAsset, fromNewAssetPegWallet := cmTypes.SubtractAssetPegFromWallet(pegHash, fromOldAssetWallet)
	if sentAsset == nil {
		return sdk.ErrInsufficientCoins("Asset not found.")
	}
	if sentAsset.GetLocked() {
		return sdk.ErrInsufficientCoins("Asset locked.")
	}
	err = keeper.orderKeeper.SendAssetsToOrder(ctx, fromAddress, toAddress, sentAsset)
	if err == nil {
		err = setAssetWallet(ctx, keeper, fromAddress, fromNewAssetPegWallet)
	}
	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendAsset,
		sdk.NewAttribute("recipient", toAddress.String()),
		sdk.NewAttribute("sender", fromAddress.String()),
		sdk.NewAttribute("asset", sentAsset.GetPegHash().String()),
	))
	
	return err
}

func (keeper BaseSendKeeper) SendFiatsToWallets(ctx sdk.Context, sendFiat types.SendFiat) sdk.Error {
	
	_acl, err := keeper.aclKeeper.GetAccountACLDetails(ctx, sendFiat.FromAddress)
	if err != nil {
		return sdk.ErrInternal("Unauthorized transaction")
	}
	account := _acl.GetACL()
	if !account.SendFiat {
		return sdk.ErrInternal("Unauthorized transaction")
	}
	err = sendFiatToOrder(ctx, keeper, sendFiat.FromAddress, sendFiat.ToAddress, sendFiat.PegHash, sendFiat.Amount)
	if err != nil {
		return err
	}
	
	keeper.reputationKeeper.SetSendFiatsPositiveTx(ctx, sendFiat.FromAddress)
	return nil
}

func sendFiatToOrder(ctx sdk.Context, keeper BaseSendKeeper, fromAddress sdk.AccAddress, toAddress sdk.AccAddress,
	pegHash cmTypes.PegHash, amount int64) sdk.Error {
	
	_negotiation, err := keeper.nk.GetNegotiationDetails(ctx, fromAddress, toAddress, pegHash)
	if err != nil {
		return err
	}
	_time := _negotiation.GetTime() + _negotiation.GetBuyerBlockHeight()
	if ctx.BlockHeight() > _time {
		return sdk.ErrInvalidSequence("Negotiation time expired.")
	}
	if _negotiation.GetSellerSignature() == nil || _negotiation.GetBuyerSignature() == nil {
		return sdk.ErrInternal("Signatures are not present")
	}
	
	fromOldFiatWallet := getFiatWallet(ctx, keeper, fromAddress)
	sentFiatPegWallet, oldFiatPegWallet := cmTypes.SubtractAmountFromWallet(amount, fromOldFiatWallet)
	if len(sentFiatPegWallet) == 0 && len(oldFiatPegWallet) == 0 {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("Insufficient funds"))
	}
	
	err = keeper.orderKeeper.SendFiatsToOrder(ctx, fromAddress, toAddress, pegHash, sentFiatPegWallet)
	if err == nil {
		err = setFiatWallet(ctx, keeper, fromAddress, oldFiatPegWallet)
	}
	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSendFiat,
		sdk.NewAttribute("recipient", toAddress.String()),
		sdk.NewAttribute("sender", fromAddress.String()),
	))
	
	return nil
}

func (keeper BaseSendKeeper) BuyerExecuteTradeOrder(ctx sdk.Context, buyerExecuteOrder types.BuyerExecuteOrder) (
	sdk.Error, []cmTypes.FiatPegWallet) {
	
	var fiatPegWallets []cmTypes.FiatPegWallet
	var _acl acl.ACL
	var err sdk.Error
	
	_, assetWallet, _, _, _ := keeper.orderKeeper.GetOrderDetails(ctx, buyerExecuteOrder.BuyerAddress,
		buyerExecuteOrder.SellerAddress, buyerExecuteOrder.PegHash)
	
	if len(assetWallet) == 0 {
		return sdk.ErrInsufficientCoins("Asset token not found!"), fiatPegWallets
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
				return sdk.ErrUnauthorized(fmt.Sprintf("Trade cannot be executed for account %v."+
					" Access Denied.", buyerExecuteOrder.MediatorAddress.String())), fiatPegWallets
			}
			_acl = aclAccount.GetACL()
			if !_acl.BuyerExecuteOrder {
				return sdk.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v."+
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
		return sdk.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.",
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

func privateExchangeOrderTokens(ctx sdk.Context, keeper BaseSendKeeper, mediatorAddress sdk.AccAddress,
	buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash cmTypes.PegHash, fiatProofHash string,
	awbProofHash string) (sdk.Error, cmTypes.FiatPegWallet, cmTypes.AssetPegWallet) {
	
	err, assetPegWallet, fiatPegWallet, orderFiatProofHash, orderAWBProofHash := keeper.orderKeeper.GetOrderDetails(ctx,
		buyerAddress, sellerAddress, pegHash)
	
	if err != nil {
		return err, fiatPegWallet, assetPegWallet
	}
	negotiationID := negotiation.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	_negotiation, err := keeper.nk.GetNegotiation(ctx, negotiationID)
	if err != nil {
		return err, fiatPegWallet, assetPegWallet
	}
	
	var reverseOrder bool
	var executed bool
	
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
		
		buyerAssetWallet = cmTypes.AddAssetPegToWallet(&assetPegWallet[0], buyerAssetWallet)
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
			buyerFiatWallet = cmTypes.AddFiatPegToWallet(buyerFiatWallet, fiatPegWallet)
			keeper.orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
		}
		if len(assetPegWallet) != 0 {
			sellerAssetWallet = cmTypes.AddAssetPegToWallet(&assetPegWallet[0], sellerAssetWallet)
			keeper.orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])
		}
		
		_ = setFiatWallet(ctx, keeper, buyerAddress, buyerFiatWallet)
		_ = setAssetWallet(ctx, keeper, sellerAddress, sellerAssetWallet)
	}
	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeExecuteOrder,
		sdk.NewAttribute("buyer", buyerAddress.String()),
		sdk.NewAttribute("seller", sellerAddress.String()),
		sdk.NewAttribute("assetPegHash", pegHash.String()),
		sdk.NewAttribute("executed", strconv.FormatBool(executed)),
		sdk.NewAttribute("assetPrice", strconv.FormatInt(_negotiation.GetBid(), 10)),
		sdk.NewAttribute("reversed", strconv.FormatBool(reverseOrder)),
	))
	
	return nil, fiatPegWallet, assetPegWallet
	
}

func exchangeOrderTokens(ctx sdk.Context, keeper BaseSendKeeper, mediatorAddress sdk.AccAddress,
	buyerAddress sdk.AccAddress, sellerAddress sdk.AccAddress, pegHash cmTypes.PegHash, fiatProofHash string,
	awbProofHash string) (sdk.Error, cmTypes.FiatPegWallet, cmTypes.AssetPegWallet) {
	
	err, assetPegWallet, fiatPegWallet, orderFiatProofHash, orderAWBProofHash := keeper.orderKeeper.GetOrderDetails(ctx, buyerAddress, sellerAddress, pegHash)
	if err != nil {
		return err, fiatPegWallet, assetPegWallet
	}
	
	negotiationID := negotiation.NegotiationID(append(append(buyerAddress.Bytes(), sellerAddress.Bytes()...), pegHash.Bytes()...))
	_negotiation, err := keeper.nk.GetNegotiation(ctx, negotiationID)
	if err != nil {
		return err, fiatPegWallet, assetPegWallet
	}
	
	var reverseOrder bool
	var oldFiatPegWallet cmTypes.FiatPegWallet
	if len(fiatPegWallet) == 0 || _negotiation.GetBid() > cmTypes.GetFiatPegWalletBalance(fiatPegWallet) {
		if _negotiation.GetTime() < ctx.BlockHeight() {
			return sdk.ErrInsufficientCoins("Fiat tokens not found!"), fiatPegWallet, assetPegWallet
		}
		reverseOrder = true
		keeper.reputationKeeper.SetBuyerExecuteOrderNegativeTx(ctx, buyerAddress)
	}
	if len(assetPegWallet) != 1 || assetPegWallet[0].GetPegHash().String() != pegHash.String() {
		if _negotiation.GetTime() < ctx.BlockHeight() {
			return sdk.ErrInsufficientCoins("Asset token not found!"), fiatPegWallet, assetPegWallet
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
	
	if _negotiation.GetBid() < cmTypes.GetFiatPegWalletBalance(fiatPegWallet) {
		fiatPegWallet, oldFiatPegWallet = cmTypes.SubtractAmountFromWallet(_negotiation.GetBid(), fiatPegWallet)
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
			
			buyerAssetWallet = cmTypes.AddAssetPegToWallet(&assetPegWallet[0], buyerAssetWallet)
			sellerFiatWallet = cmTypes.AddFiatPegToWallet(sellerFiatWallet, fiatPegWallet)
			
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
			buyerFiatWallet = cmTypes.AddFiatPegToWallet(buyerFiatWallet, oldFiatPegWallet)
		}
		
		if len(fiatPegWallet) != 0 {
			buyerFiatWallet = cmTypes.AddFiatPegToWallet(buyerFiatWallet, fiatPegWallet)
			keeper.orderKeeper.SendFiatsFromOrder(ctx, buyerAddress, sellerAddress, pegHash, fiatPegWallet)
		}
		if len(assetPegWallet) != 0 {
			sellerAssetWallet = cmTypes.AddAssetPegToWallet(&assetPegWallet[0], sellerAssetWallet)
			keeper.orderKeeper.SendAssetFromOrder(ctx, buyerAddress, sellerAddress, &assetPegWallet[0])
		}
		
		_ = setFiatWallet(ctx, keeper, buyerAddress, buyerFiatWallet)
		_ = setAssetWallet(ctx, keeper, sellerAddress, sellerAssetWallet)
	}
	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeExecuteOrder,
		sdk.NewAttribute("buyer", buyerAddress.String()),
		sdk.NewAttribute("seller", sellerAddress.String()),
		sdk.NewAttribute("assetPegHash", pegHash.String()),
		sdk.NewAttribute("executed", strconv.FormatBool(executed)),
		sdk.NewAttribute("assetPrice", strconv.FormatInt(_negotiation.GetBid(), 10)),
		sdk.NewAttribute("reversed", strconv.FormatBool(reverseOrder)),
	))
	
	return nil, fiatPegWallet, assetPegWallet
}

func (keeper BaseSendKeeper) SellerExecuteTradeOrder(ctx sdk.Context, sellerExecuteOrder types.SellerExecuteOrder) (
	sdk.Error, []cmTypes.AssetPegWallet) {
	
	var assetPegWallets []cmTypes.AssetPegWallet
	var _acl acl.ACL
	var err sdk.Error
	var assetPegWallet cmTypes.AssetPegWallet
	
	_, assetWallet, _, _, _ := keeper.orderKeeper.GetOrderDetails(ctx, sellerExecuteOrder.BuyerAddress,
		sellerExecuteOrder.SellerAddress, sellerExecuteOrder.PegHash)
	
	if len(assetWallet) == 0 {
		return sdk.ErrInsufficientCoins("Asset token not found!"), assetPegWallets
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
				return sdk.ErrUnauthorized(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.",
					sellerExecuteOrder.MediatorAddress.String())), assetPegWallets
			}
			_acl = aclAccount.GetACL()
			if !_acl.SellerExecuteOrder {
				return sdk.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.",
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
		return sdk.ErrInternal(fmt.Sprintf("Trade cannot be executed for account %v. Access Denied.",
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

func (keeper BaseSendKeeper) ReleaseLockedAssets(ctx sdk.Context, releaseAsset types.ReleaseAsset) sdk.Error {
	
	_acl, err := keeper.aclKeeper.CheckZoneAndGetACL(ctx, releaseAsset.ZoneAddress, releaseAsset.OwnerAddress)
	if err != nil {
		return err
	}
	if !_acl.ReleaseAsset {
		return sdk.ErrInternal(fmt.Sprintf("Assets cannot be released for account %v. Access Denied.",
			releaseAsset.OwnerAddress.String()))
	}
	
	err = releaseAssets(ctx, keeper, releaseAsset.ZoneAddress, releaseAsset.OwnerAddress, releaseAsset.PegHash)
	if err != nil {
		return err
	}
	return nil
}

func releaseAssets(ctx sdk.Context, keeper BaseSendKeeper, zoneAddress sdk.AccAddress, ownerAddress sdk.AccAddress,
	pegHash cmTypes.PegHash) sdk.Error {
	
	ownerAssetWallet := getAssetWallet(ctx, keeper, ownerAddress)
	if !cmTypes.ReleaseAssetPegInWallet(ownerAssetWallet, pegHash) {
		return sdk.ErrInternal("Asset peg not found.")
	}
	_ = setAssetWallet(ctx, keeper, ownerAddress, ownerAssetWallet)
	
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeReleaseAsset,
		sdk.NewAttribute("zone", zoneAddress.String()),
		sdk.NewAttribute("owner", ownerAddress.String()),
		sdk.NewAttribute("asset", pegHash.String()),
	))
	
	return nil
}

func (keeper BaseSendKeeper) DefineZones(ctx sdk.Context, defineZone types.DefineZone) sdk.Error {
	if !keeper.aclKeeper.CheckValidGenesisAddress(ctx, defineZone.From) {
		return sdk.ErrInternal(fmt.Sprintf("Account %v is not the genesis account. Zones can only be"+
			" defined by the genesis account.", defineZone.From.String()))
	}
	
	err := keeper.aclKeeper.DefineZoneAddress(ctx, defineZone.To, defineZone.ZoneID)
	if err != nil {
		return err
	}
	return nil
}

func (keeper BaseSendKeeper) DefineOrganizations(ctx sdk.Context, defineOrganization types.DefineOrganization) sdk.Error {
	if !keeper.aclKeeper.CheckValidZoneAddress(ctx, defineOrganization.ZoneID, defineOrganization.From) {
		return sdk.ErrInternal(fmt.Sprintf("Account %v is not the zone account. Organizations can only "+
			"be defined by the zone account.", defineOrganization.From.String()))
	}
	err := keeper.aclKeeper.DefineOrganizationAddress(ctx, defineOrganization.To,
		defineOrganization.OrganizationID, defineOrganization.ZoneID)
	
	if err != nil {
		return err
	}
	return nil
}

func (keeper BaseSendKeeper) DefineACLs(ctx sdk.Context, defineACL types.DefineACL) sdk.Error {
	if !keeper.aclKeeper.CheckValidGenesisAddress(ctx, defineACL.From) {
		if !keeper.aclKeeper.CheckValidZoneAddress(ctx, defineACL.ACLAccount.GetZoneID(), defineACL.From) {
			if !keeper.aclKeeper.CheckValidOrganizationAddress(ctx, defineACL.ACLAccount.GetZoneID(),
				defineACL.ACLAccount.GetOrganizationID(), defineACL.From) {
				
				return sdk.ErrInternal(fmt.Sprintf("Account %v does not have access to define acl "+
					"for account %v.", defineACL.From.String(), defineACL.To.String()))
			}
		}
	}
	
	err := keeper.aclKeeper.DefineACLAccount(ctx, defineACL.To, defineACL.ACLAccount)
	if err != nil {
		return err
	}
	return nil
}
