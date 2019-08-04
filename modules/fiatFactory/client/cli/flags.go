package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagTo                = "to"
	FlagPegHash           = "peg-hash"
	FlagFiatPegHash       = "fiat-peg-hash"
	FlagAssetPegHash      = "asset-peg-hash"
	FlagTransactionID     = "transaction-id"
	FlagTransactionAmount = "transaction-amount"
	FlagOwnerAddress      = "owner"
	FlagAmount            = "amount"
)

var (
	fsTo                = flag.NewFlagSet("", flag.ContinueOnError)
	fsPegHash           = flag.NewFlagSet("", flag.ContinueOnError)
	fsFiatPegHash       = flag.NewFlagSet("", flag.ContinueOnError)
	fsAssetPegHash      = flag.NewFlagSet("", flag.ContinueOnError)
	fsTransactionID     = flag.NewFlagSet("", flag.ContinueOnError)
	fsTransactionAmount = flag.NewFlagSet("", flag.ContinueOnError)
	fsOwnerAddress      = flag.NewFlagSet("", flag.ContinueOnError)
	fsAmount            = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsTo.String(FlagTo, "", "Address to send coins")
	fsPegHash.String(FlagPegHash, "", "Peg Hash to be transferred")
	fsFiatPegHash.String(FlagFiatPegHash, "", "Peg Hash to be transferred")
	fsAssetPegHash.String(FlagAssetPegHash, "", "Peg Hash to be transferred")
	fsTransactionID.String(FlagTransactionID, "", "Fiat deposit transaction ID")
	fsTransactionAmount.String(FlagTransactionAmount, "", "Fiat deposit transaction amount")
	fsOwnerAddress.String(FlagOwnerAddress, "", "Address of current owner of asset peg")
	fsAmount.String(FlagAmount, "", "Amount to be transfered")
}
