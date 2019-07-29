package types

import (
	"bytes"
	"encoding/hex"
	"sort"

	"github.com/tendermint/tendermint/libs/common"
)

//PegHash : reference address of asset peg
type PegHash = common.HexBytes

//AssetPeg : comdex asset interface
type AssetPeg interface {
	GetPegHash() PegHash
	SetPegHash(PegHash) error

	GetDocumentHash() string
	SetDocumentHash(string) error

	GetAssetType() string
	SetAssetType(string) error

	GetAssetQuantity() int64
	SetAssetQuantity(int64) error

	GetAssetPrice() int64
	SetAssetPrice(int64) error

	GetQuantityUnit() string
	SetQuantityUnit(string) error

	GetOwnerAddress() AccAddress
	SetOwnerAddress(AccAddress) error

	GetLocked() bool
	SetLocked(bool) error

	GetModerated() bool
	SetModerated(bool) error

	GetTakerAddress() AccAddress
	SetTakerAddress(AccAddress) error
}

//BaseAssetPeg : base asset type
type BaseAssetPeg struct {
	PegHash       PegHash    `json:"pegHash"`
	DocumentHash  string     `json:"documentHash" valid:"required~Mandatory parameter documentHash missing"`
	AssetType     string     `json:"assetType" valid:"required~Mandatory parameter AssetType missing,matches(^[A-Za-z ]+$)~Parameter AssetType is Invalid"`
	AssetQuantity int64      `json:"assetQuantity" valid:"required~Mandatory parameter AssetQuantity missing,matches(^[1-9]{1}[0-9]*$)~Parameter AssetQuantity is Invalid"`
	AssetPrice    int64      `json:"assetPrice" valid:"required~Mandatory parameter assetPrice missing,matches(^[1-9]{1}[0-9]*$)~Parameter assetPrice is Invalid"`
	QuantityUnit  string     `json:"quantityUnit" valid:"required~Mandatory parameter QuantityUnit missing,matches(^[A-Za-z]*$)~Parameter QuantityUnit is Invalid"`
	OwnerAddress  AccAddress `json:"ownerAddress"`
	Locked        bool       `json:"locked"`
	Moderated     bool       `json:"moderated"`
	TakerAddress  AccAddress `json:"takerAddress"`
}

//NewBaseAssetPegWithPegHash a base asset peg with peg hash
func NewBaseAssetPegWithPegHash(pegHash PegHash) BaseAssetPeg {
	return BaseAssetPeg{
		PegHash: pegHash,
	}
}

// ProtoBaseAssetPeg : converts concrete to assetPeg
func ProtoBaseAssetPeg() AssetPeg {
	return &BaseAssetPeg{}
}

var _ AssetPeg = (*BaseAssetPeg)(nil)

//GetPegHash : getter
func (baseAssetPeg BaseAssetPeg) GetPegHash() PegHash { return baseAssetPeg.PegHash }

//SetPegHash : setter
func (baseAssetPeg *BaseAssetPeg) SetPegHash(pegHash PegHash) error {
	baseAssetPeg.PegHash = pegHash
	return nil
}

//GetDocumentHash : getter
func (baseAssetPeg BaseAssetPeg) GetDocumentHash() string { return baseAssetPeg.DocumentHash }

//SetDocumentHash : setter
func (baseAssetPeg *BaseAssetPeg) SetDocumentHash(documentHash string) error {
	baseAssetPeg.DocumentHash = documentHash
	return nil
}

//GetAssetType : getter
func (baseAssetPeg BaseAssetPeg) GetAssetType() string { return baseAssetPeg.AssetType }

//SetAssetType : setter
func (baseAssetPeg *BaseAssetPeg) SetAssetType(assetType string) error {
	baseAssetPeg.AssetType = assetType
	return nil
}

//GetAssetPrice : getter
func (baseAssetPeg BaseAssetPeg) GetAssetPrice() int64 { return baseAssetPeg.AssetPrice }

//SetAssetPrice : setter
func (baseAssetPeg *BaseAssetPeg) SetAssetPrice(assetPrice int64) error {
	baseAssetPeg.AssetPrice = assetPrice
	return nil
}

//GetAssetQuantity : getter
func (baseAssetPeg BaseAssetPeg) GetAssetQuantity() int64 { return baseAssetPeg.AssetQuantity }

//SetAssetQuantity : setter
func (baseAssetPeg *BaseAssetPeg) SetAssetQuantity(assetQuantity int64) error {
	baseAssetPeg.AssetQuantity = assetQuantity
	return nil
}

//GetQuantityUnit : getter
func (baseAssetPeg BaseAssetPeg) GetQuantityUnit() string { return baseAssetPeg.QuantityUnit }

//SetQuantityUnit : setter
func (baseAssetPeg *BaseAssetPeg) SetQuantityUnit(quantityUnit string) error {
	baseAssetPeg.QuantityUnit = quantityUnit
	return nil
}

//GetOwnerAddress : getter
func (baseAssetPeg BaseAssetPeg) GetOwnerAddress() AccAddress { return baseAssetPeg.OwnerAddress }

//SetOwnerAddress : setter
func (baseAssetPeg *BaseAssetPeg) SetOwnerAddress(ownerAddress AccAddress) error {
	baseAssetPeg.OwnerAddress = ownerAddress
	return nil
}

//GetLocked : getter
func (baseAssetPeg BaseAssetPeg) GetLocked() bool { return baseAssetPeg.Locked }

//SetLocked : setter
func (baseAssetPeg *BaseAssetPeg) SetLocked(locked bool) error {
	baseAssetPeg.Locked = locked
	return nil
}

//GetModerated : getter
func (baseAssetPeg *BaseAssetPeg) GetModerated() bool { return baseAssetPeg.Moderated }

//SetModerated : setter
func (baseAssetPeg *BaseAssetPeg) SetModerated(moderated bool) error {
	baseAssetPeg.Moderated = moderated
	return nil
}

//GetTakerAddress : getter
func (baseAssetPeg *BaseAssetPeg) GetTakerAddress() AccAddress {
	return baseAssetPeg.TakerAddress
}

//SetTakerAddress : setter
func (baseAssetPeg *BaseAssetPeg) SetTakerAddress(takerAddress AccAddress) error {
	baseAssetPeg.TakerAddress = takerAddress
	return nil
}

//AssetPegDecoder : decoder function for asset peg
type AssetPegDecoder func(assetPegBytes []byte) (AssetPeg, error)

//GetAssetPegHashHex : convert string to hex peg hash
func GetAssetPegHashHex(pegHashStr string) (pegHash PegHash, err error) {
	bz, err := hex.DecodeString(pegHashStr)
	if err != nil {
		return nil, err
	}
	return PegHash(bz), nil
}

//ToBaseAssetPeg : convert interface to concrete
func ToBaseAssetPeg(assetPeg AssetPeg) BaseAssetPeg {
	var baseAssetPeg BaseAssetPeg
	baseAssetPeg.AssetQuantity = assetPeg.GetAssetQuantity()
	baseAssetPeg.AssetPrice = assetPeg.GetAssetPrice()
	baseAssetPeg.AssetType = assetPeg.GetAssetType()
	baseAssetPeg.DocumentHash = assetPeg.GetDocumentHash()
	baseAssetPeg.PegHash = assetPeg.GetPegHash()
	baseAssetPeg.QuantityUnit = assetPeg.GetQuantityUnit()
	baseAssetPeg.OwnerAddress = assetPeg.GetOwnerAddress()
	baseAssetPeg.Locked = assetPeg.GetLocked()
	baseAssetPeg.Moderated = assetPeg.GetModerated()
	baseAssetPeg.TakerAddress = assetPeg.GetTakerAddress()
	return baseAssetPeg
}

//AssetPegWallet : A wallet of AssetPegTokens
type AssetPegWallet []BaseAssetPeg

// Sort interface

func (assetPegWallet AssetPegWallet) Len() int { return len(assetPegWallet) }

func (assetPegWallet AssetPegWallet) Less(i, j int) bool {
	return bytes.Compare(assetPegWallet[i].PegHash, assetPegWallet[j].PegHash) < 0
}

func (assetPegWallet AssetPegWallet) Swap(i, j int) {
	assetPegWallet[i], assetPegWallet[j] = assetPegWallet[j], assetPegWallet[i]
}

var _ sort.Interface = AssetPegWallet{}

// Sort is a helper function to sort the set of asset pegs inplace
func (assetPegWallet AssetPegWallet) Sort() AssetPegWallet {
	sort.Sort(assetPegWallet)
	return assetPegWallet
}

//GetAssetPeg :
func (assetPegWallet AssetPegWallet) SearchAssetPeg(pegHash PegHash) int {
	index := sort.Search(assetPegWallet.Len(), func(i int) bool {
		return bytes.Compare(assetPegWallet[i].GetPegHash(), pegHash) != -1
	})
	return index
}

//SubtractAssetPegFromWallet : subtract asset peg from wallet
func SubtractAssetPegFromWallet(pegHash PegHash, assetPegWallet AssetPegWallet) (AssetPeg, AssetPegWallet) {
	i := assetPegWallet.SearchAssetPeg(pegHash)
	if i < len(assetPegWallet) && assetPegWallet[i].GetPegHash().String() == pegHash.String() {
		assetPeg := assetPegWallet[i]
		assetPegWallet = append(assetPegWallet[:i], assetPegWallet[i+1:]...)
		assetPegWallet = assetPegWallet.Sort()
		return &assetPeg, assetPegWallet
	}
	return nil, assetPegWallet

}

//AddAssetPegToWallet : add asset peg to wallet
func AddAssetPegToWallet(assetPeg AssetPeg, assetPegWallet AssetPegWallet) AssetPegWallet {
	i := assetPegWallet.SearchAssetPeg(assetPeg.GetPegHash())
	if i < len(assetPegWallet) && assetPegWallet[i].GetPegHash().String() == assetPeg.GetPegHash().String() {
		return assetPegWallet
	}
	assetPegWallet = append(assetPegWallet, ToBaseAssetPeg(assetPeg))
	assetPegWallet = assetPegWallet.Sort()
	return assetPegWallet

}

//IssueAssetPeg : issues asset peg from the zones wallet to the provided wallet
func IssueAssetPeg(issuerAssetPegWallet AssetPegWallet, receiverAssetPegWallet AssetPegWallet, assetPeg AssetPeg) (AssetPegWallet, AssetPegWallet, AssetPeg) {
	issuedAssetPegHash := issuerAssetPegWallet[len(issuerAssetPegWallet)-1].PegHash
	issuerAssetPegWallet = issuerAssetPegWallet[:len(issuerAssetPegWallet)-1]
	assetPeg.SetPegHash(issuedAssetPegHash)
	assetPeg.SetLocked(true)
	receiverAssetPegWallet = AddAssetPegToWallet(assetPeg, receiverAssetPegWallet)
	return issuerAssetPegWallet, receiverAssetPegWallet, assetPeg
}

//ReleaseAssetPegInWallet : get an asset peg in wallet and set locked to false
func ReleaseAssetPegInWallet(assetPegWallet AssetPegWallet, pegHash PegHash) bool {
	i := assetPegWallet.SearchAssetPeg(pegHash)
	if i < len(assetPegWallet) && assetPegWallet[i].GetPegHash().String() == pegHash.String() {
		assetPegWallet[i].SetLocked(false)
		return true
	}
	return false
}
