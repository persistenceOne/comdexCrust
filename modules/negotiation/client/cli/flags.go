package cli

import (
	flag "github.com/spf13/pflag"
)

// noLint
const (
	FlagTo                 = "to"
	FlagFrom               = "from"
	FlagPegHash            = "peg-hash"
	FlagBid                = "bid"
	FlagTime               = "time"
	FlagNegotiationID      = "negotiation-id"
	FlagBuyerContractHash  = "buyer-contract-hash"
	FlagSellerContractHash = "seller-contract-hash"
)

var (
	fsTo                 = flag.NewFlagSet("", flag.ContinueOnError)
	fsPegHash            = flag.NewFlagSet("", flag.ContinueOnError)
	fsBid                = flag.NewFlagSet("", flag.ContinueOnError)
	fsTime               = flag.NewFlagSet("", flag.ContinueOnError)
	fsFrom               = flag.NewFlagSet("", flag.ContinueOnError)
	fsNegotiationID      = flag.NewFlagSet("", flag.ContinueOnError)
	fsBuyerContractHash  = flag.NewFlagSet("", flag.ContinueOnError)
	fsSellerContractHash = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsTo.String(FlagTo, "", "Address to send coins")
	fsPegHash.String(FlagPegHash, "", "Peg Hash to be negotiated ")
	fsBid.String(FlagBid, "", "Amount of fiat to bid against asset")
	fsTime.String(FlagTime, "", "Time to be assumed for contract confirmation")
	fsFrom.String(FlagFrom, "", "address of buyer account")
	fsBuyerContractHash.String(FlagBuyerContractHash, "", "buyer contract hash")
	fsSellerContractHash.String(FlagSellerContractHash, "", "seller contract hash")
	fsNegotiationID.String(FlagNegotiationID, "", "NegotiationID")
}
