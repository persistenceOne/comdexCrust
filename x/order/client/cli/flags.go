package cli

import (
	flag "github.com/spf13/pflag"
)

//noLint
const (
	FlagNegotiationID = "negotiationID"
)

var (
	fsNegotiationID = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsNegotiationID.String(FlagNegotiationID, "", "NegotiationID")
}
