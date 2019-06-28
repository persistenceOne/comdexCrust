package types

import (
	"encoding/hex"
	"errors"
	"sort"
)

// FiatPeg : peg issued against each fiat transaction
type FiatPeg interface {
	GetPegHash() PegHash
	SetPegHash(PegHash) error
	
	GetTransactionID() string
	SetTransactionID(string) error
	
	GetTransactionAmount() int64
	SetTransactionAmount(int64) error
	
	GetRedeemedAmount() int64
	SetRedeemedAmount(int64) error
	
	GetOwners() []Owner
	SetOwners([]Owner) error
}

// Owner : partial or full owner of a transaction
type Owner struct {
	OwnerAddress AccAddress `json:"ownerAddress"`
	Amount       int64      `json:"amount"`
}

// BaseFiatPeg : fiat peg basic implementation
type BaseFiatPeg struct {
	PegHash           PegHash `json:"pegHash" `
	TransactionID     string  `json:"transactionID" valid:"required~TxID is mandatory,matches(^[A-Z0-9]+$)~Invalid TransactionId,length(2|40)~TransactionId length between 2-40"`
	TransactionAmount int64   `json:"transactionAmount" valid:"required~TransactionAmount is mandatory,matches(^[1-9]{1}[0-9]*$)~Invalid TransactionAmount"`
	RedeemedAmount    int64   `json:"redeemedAmount"`
	Owners            []Owner `json:"owners"`
}

var _ FiatPeg = (*BaseFiatPeg)(nil)

// GetPegHash : getter
func (baseFiatPeg BaseFiatPeg) GetPegHash() PegHash { return baseFiatPeg.PegHash }

// SetPegHash : setter
func (baseFiatPeg *BaseFiatPeg) SetPegHash(pegHash PegHash) error {
	baseFiatPeg.PegHash = pegHash
	return nil
}

// GetTransactionID : getter
func (baseFiatPeg BaseFiatPeg) GetTransactionID() string { return baseFiatPeg.TransactionID }

// SetTransactionID : setter
func (baseFiatPeg *BaseFiatPeg) SetTransactionID(transactionID string) error {
	baseFiatPeg.TransactionID = transactionID
	return nil
}

// GetTransactionAmount : getter
func (baseFiatPeg BaseFiatPeg) GetTransactionAmount() int64 { return baseFiatPeg.TransactionAmount }

// SetTransactionAmount : setter
func (baseFiatPeg *BaseFiatPeg) SetTransactionAmount(transactionAmount int64) error {
	baseFiatPeg.TransactionAmount = transactionAmount
	return nil
}

// GetRedeemedAmount : getter
func (baseFiatPeg BaseFiatPeg) GetRedeemedAmount() int64 { return baseFiatPeg.RedeemedAmount }

// SetRedeemedAmount : setter
func (baseFiatPeg *BaseFiatPeg) SetRedeemedAmount(redeemedAmount int64) error {
	baseFiatPeg.RedeemedAmount = redeemedAmount
	return nil
}

// GetOwners : getter
func (baseFiatPeg BaseFiatPeg) GetOwners() []Owner { return baseFiatPeg.Owners }

// SetOwners : setter
func (baseFiatPeg *BaseFiatPeg) SetOwners(owners []Owner) error {
	baseFiatPeg.Owners = owners
	return nil
}

// FiatPegDecoder : decoder function for fiat peg
type FiatPegDecoder func(fiatPegBytes []byte) (FiatPeg, error)

// GetFiatPegHashHex : convert string to hex peg hash
func GetFiatPegHashHex(pegHashStr string) (pegHash PegHash, err error) {
	if len(pegHashStr) == 0 {
		return pegHash, errors.New("must use provide pegHash")
	}
	bz, err := hex.DecodeString(pegHashStr)
	if err != nil {
		return nil, err
	}
	return PegHash(bz), nil
}

// NewBaseFiatPegWithPegHash a base fiat peg with peg hash
func NewBaseFiatPegWithPegHash(pegHash PegHash) BaseFiatPeg {
	return BaseFiatPeg{
		PegHash: pegHash,
	}
}

// FiatPegWallet : A wallet of fiat peg tokens
type FiatPegWallet []BaseFiatPeg

// ProtoBaseFiatPeg : create a new interface prototype
func ProtoBaseFiatPeg() FiatPeg {
	return &BaseFiatPeg{}
}

// ToFiatPeg : convert base fiat peg to interface fiat peg
func ToFiatPeg(baseFiatPeg BaseFiatPeg) FiatPeg {
	fiatI := ProtoBaseFiatPeg()
	fiatI.SetOwners(baseFiatPeg.Owners)
	fiatI.SetPegHash(baseFiatPeg.PegHash)
	fiatI.SetRedeemedAmount(baseFiatPeg.RedeemedAmount)
	fiatI.SetTransactionAmount(baseFiatPeg.TransactionAmount)
	fiatI.SetTransactionID(baseFiatPeg.TransactionID)
	return fiatI
}

// ToBaseFiatPeg : convert interface fiat peg to base fiat peg
func ToBaseFiatPeg(fiatPeg FiatPeg) BaseFiatPeg {
	var baseFiatPeg BaseFiatPeg
	baseFiatPeg.Owners = fiatPeg.GetOwners()
	baseFiatPeg.PegHash = fiatPeg.GetPegHash()
	baseFiatPeg.RedeemedAmount = fiatPeg.GetRedeemedAmount()
	baseFiatPeg.TransactionAmount = fiatPeg.GetTransactionAmount()
	baseFiatPeg.TransactionID = fiatPeg.GetTransactionID()
	return baseFiatPeg
}

// Sort interface

func (fiatPegWallet FiatPegWallet) Len() int { return len(fiatPegWallet) }
func (fiatPegWallet FiatPegWallet) Less(i, j int) bool {
	if (fiatPegWallet[i].TransactionAmount - fiatPegWallet[j].TransactionAmount) < 0 {
		return true
	}
	return false
}
func (fiatPegWallet FiatPegWallet) Swap(i, j int) {
	fiatPegWallet[i], fiatPegWallet[j] = fiatPegWallet[j], fiatPegWallet[i]
}

var _ sort.Interface = FiatPegWallet{}

// Sort is a helper function to sort the set of fiat pegs inplace
func (fiatPegWallet FiatPegWallet) Sort() FiatPegWallet {
	sort.Sort(fiatPegWallet)
	return fiatPegWallet
}

// SubtractFiatPegWalletFromWallet : subtract fiat peg  wallet from wallet
func SubtractFiatPegWalletFromWallet(inFiatPegWallet FiatPegWallet, fiatPegWallet FiatPegWallet) FiatPegWallet {
	for _, inFiatPeg := range inFiatPegWallet {
		for i, fiatPeg := range fiatPegWallet {
			if fiatPeg.GetPegHash().String() == inFiatPeg.GetPegHash().String() {
				fiatPegWallet = append(fiatPegWallet[:i], fiatPegWallet[i+1:]...)
				fiatPegWallet = fiatPegWallet.Sort()
				break
			}
		}
	}
	return fiatPegWallet
}

// SubtractAmountFromWallet : subtract fiat peg from wallet
func SubtractAmountFromWallet(amount int64, fiatPegWallet FiatPegWallet) (sentFiatPegWallet FiatPegWallet, oldFiatPegWallet FiatPegWallet) {
	var collected int64
	fiatPegWallet = fiatPegWallet.Sort()
	for _, fiatPeg := range fiatPegWallet {
		if collected < amount {
			if fiatPeg.TransactionAmount <= amount-collected {
				collected += fiatPeg.TransactionAmount
				sentFiatPegWallet = append(sentFiatPegWallet, fiatPeg)
			} else if fiatPeg.TransactionAmount > amount-collected {
				splitFiatPeg := fiatPeg
				splitFiatPeg.TransactionAmount = amount - collected
				fiatPeg.TransactionAmount -= amount - collected
				oldFiatPegWallet = append(oldFiatPegWallet, fiatPeg)
				sentFiatPegWallet = append(sentFiatPegWallet, splitFiatPeg)
				collected += amount - collected
			}
		} else {
			oldFiatPegWallet = append(oldFiatPegWallet, fiatPeg)
		}
	}
	if collected == amount {
		oldFiatPegWallet = oldFiatPegWallet.Sort()
		sentFiatPegWallet = sentFiatPegWallet.Sort()
		return
	}
	return FiatPegWallet{}, FiatPegWallet{}
	
}

// RedeemAmountFromWallet : subtract fiat peg from wallet
func RedeemAmountFromWallet(amount int64, fiatPegWallet FiatPegWallet) (emptiedFiatPegWallet FiatPegWallet, redeemerFiatPegWallet FiatPegWallet) {
	var collected int64
	for _, fiatPeg := range fiatPegWallet {
		if collected < amount {
			if fiatPeg.TransactionAmount <= amount-collected {
				collected += fiatPeg.TransactionAmount
				emptiedFiatPegWallet = append(emptiedFiatPegWallet, fiatPeg)
			} else if fiatPeg.TransactionAmount > amount-collected {
				fiatPeg.TransactionAmount -= amount - collected
				fiatPeg.RedeemedAmount = amount - collected
				redeemerFiatPegWallet = append(redeemerFiatPegWallet, fiatPeg)
				collected += amount - collected
			}
		} else {
			redeemerFiatPegWallet = append(redeemerFiatPegWallet, fiatPeg)
		}
	}
	if collected == amount {
		redeemerFiatPegWallet = redeemerFiatPegWallet.Sort()
		emptiedFiatPegWallet = emptiedFiatPegWallet.Sort()
		return
	}
	return FiatPegWallet{}, FiatPegWallet{}
	
}

// AddFiatPegToWallet : add fiat peg to wallet
func AddFiatPegToWallet(fiatPegWallet FiatPegWallet, inFiatPegWallet FiatPegWallet) FiatPegWallet {
	for _, inFiatPeg := range inFiatPegWallet {
		added := false
		for i, fiatPeg := range fiatPegWallet {
			if fiatPeg.PegHash.String() == inFiatPeg.PegHash.String() {
				inFiatPeg.TransactionAmount += fiatPeg.TransactionAmount
				fiatPegWallet[i] = inFiatPeg
				added = true
				break
			}
		}
		if !added {
			fiatPegWallet = append(fiatPegWallet, inFiatPeg)
		}
	}
	fiatPegWallet = fiatPegWallet.Sort()
	return fiatPegWallet
}

// IssueFiatPeg : issues fiat peg from the zones wallet to the provided wallet
func IssueFiatPeg(issuerFiatPegWallet FiatPegWallet, receiverFiatPegWallet FiatPegWallet, fiatPeg FiatPeg) (FiatPegWallet, FiatPegWallet, FiatPeg) {
	issuedFiatPegHash := issuerFiatPegWallet[len(issuerFiatPegWallet)-1].PegHash
	issuerFiatPegWallet = issuerFiatPegWallet[:len(issuerFiatPegWallet)-1]
	fiatPeg.SetPegHash(issuedFiatPegHash)
	receiverFiatPegWallet = AddFiatPegToWallet(receiverFiatPegWallet, []BaseFiatPeg{ToBaseFiatPeg(fiatPeg)})
	return issuerFiatPegWallet, receiverFiatPegWallet, fiatPeg
}

// GetFiatPegWalletBalance :  gets the total sum of all fiat pegs in a wallet
func GetFiatPegWalletBalance(fiatPegWallet FiatPegWallet) int64 {
	var balance int64
	for _, fiatPeg := range fiatPegWallet {
		balance += fiatPeg.TransactionAmount
	}
	return balance
}

// TransferFiatPegsToWallet : subtracts and changes owners of fiat peg in fiat chain
func TransferFiatPegsToWallet(fiatPegWallet FiatPegWallet, oldFiatPegWallet FiatPegWallet, fromAddress AccAddress, toAddress AccAddress) FiatPegWallet {
	for _, fiatPeg := range fiatPegWallet {
		transfered := false
		for j, oldFiatPeg := range oldFiatPegWallet {
			if fiatPeg.GetPegHash().String() == oldFiatPeg.GetPegHash().String() {
				subtracted := 0
				added := 0
				for i, owner := range oldFiatPeg.Owners {
					if owner.OwnerAddress.String() == fromAddress.String() && owner.Amount >= fiatPeg.TransactionAmount {
						owner.Amount -= fiatPeg.TransactionAmount
						oldFiatPeg.Owners[i] = owner
						subtracted++
					} else if owner.OwnerAddress.String() == toAddress.String() {
						owner.Amount += fiatPeg.TransactionAmount
						oldFiatPeg.Owners[i] = owner
						added++
					}
				}
				if added == 0 {
					owner := Owner{toAddress, fiatPeg.TransactionAmount}
					oldFiatPeg.Owners = append(oldFiatPeg.Owners, owner)
					added++
				}
				if subtracted != 1 || added != 1 {
					return nil
				}
				transfered = true
				oldFiatPegWallet[j] = oldFiatPeg
				break
			}
		}
		if !transfered {
			return nil
		}
	}
	return oldFiatPegWallet
}

// RedeemFiatPegsFromWallet : subtracts and changes owners of fiat peg in fiat chain
func RedeemFiatPegsFromWallet(fiatPegWallet FiatPegWallet, oldFiatPegWallet FiatPegWallet, fromAddress AccAddress) FiatPegWallet {
	for _, fiatPeg := range fiatPegWallet {
		transfered := false
		for j, oldFiatPeg := range oldFiatPegWallet {
			if fiatPeg.GetPegHash().String() == oldFiatPeg.GetPegHash().String() {
				subtracted := 0
				
				for i, owner := range oldFiatPeg.Owners {
					if owner.OwnerAddress.String() == fromAddress.String() && owner.Amount > fiatPeg.RedeemedAmount {
						owner.Amount -= fiatPeg.RedeemedAmount
						oldFiatPeg.Owners[i] = owner
						subtracted++
					} else if owner.OwnerAddress.String() == fromAddress.String() && owner.Amount == fiatPeg.RedeemedAmount {
						oldFiatPeg.Owners = append(oldFiatPeg.Owners[:i], oldFiatPeg.Owners[i+1:]...)
						subtracted++
					}
				}
				
				if subtracted != 1 {
					return nil
				}
				oldFiatPeg.TransactionAmount -= fiatPeg.RedeemedAmount
				oldFiatPeg.RedeemedAmount += fiatPeg.RedeemedAmount
				
				transfered = true
				oldFiatPegWallet[j] = oldFiatPeg
				break
			}
		}
		if !transfered {
			return nil
		}
	}
	return oldFiatPegWallet
}
