package cli

import (
	flag "github.com/spf13/pflag"
)

// noLint
const (
	FlagTo            = "to"
	FlagAmount        = "amount"
	FlagDocumentHash  = "documentHash"
	FlagAssetType     = "assetType"
	FlagAssetPrice    = "assetPrice"
	FlagAssetQuantity = "assetQuantity"
	FlagQuantityUnit  = "quantityUnit"
	FlagPegHash       = "pegHash"
	FlagOwnerAddress  = "owner"
)

var (
	fsTo            = flag.NewFlagSet("", flag.ContinueOnError)
	fsAmount        = flag.NewFlagSet("", flag.ContinueOnError)
	fsDocumentHash  = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetType     = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetPrice    = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetQuantity = flag.NewFlagSet("", flag.ContinueOnError)
	fsQuantityUnit  = flag.NewFlagSet("", flag.ContinueOnError)
	fsPegHash       = flag.NewFlagSet("", flag.ContinueOnError)
	fsOwnerAddress  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsTo.String(FlagTo, "", "Address to send coins")
	fsAmount.String(FlagAmount, "", "Amount of coins to send")
	fsDocumentHash.String(FlagDocumentHash, "", "Hash of the asset doccuments of the asset")
	fsAssetType.String(FlagAssetType, "", "Type of the asset")
	fsAssetPrice.String(FlagAssetPrice, "", "Type of the asset")
	fsAssetQuantity.String(FlagAssetQuantity, "", "Quantity of the assent in integer")
	fsQuantityUnit.String(FlagQuantityUnit, "", "The unit of the qunatity")
	fsPegHash.String(FlagPegHash, "", "Peg Hash to be transferred ")
	fsOwnerAddress.String(FlagOwnerAddress, "", "Address of current owner of asset peg")
}
