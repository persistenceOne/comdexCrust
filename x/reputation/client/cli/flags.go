package cli

import (
	flag "github.com/spf13/pflag"
)

// noLint
const (
	FlagTo      = "to"
	FlagPegHash = "pegHash"
	FlagRating  = "rating"
)

var (
	fsTo      = flag.NewFlagSet("", flag.ContinueOnError)
	fsPeghash = flag.NewFlagSet("", flag.ContinueOnError)
	fsRating  = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsTo.String(FlagTo, "", "Address to rate")
	fsPeghash.String(FlagPegHash, "", "Peg Hash to be negotiated ")
	fsRating.String(FlagRating, "", "Ratings from feedback")
}
