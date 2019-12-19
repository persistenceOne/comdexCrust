package context

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	yaml "gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/tendermint/tendermint/libs/log"
	tmlite "github.com/tendermint/tendermint/lite"
	tmliteproxy "github.com/tendermint/tendermint/lite/proxy"
)

const (
	cacheSize = 10
)

// NewCLIContextIBC takes additional arguements
func NewCLIContextIBC(from string, chainID string, nodeURI string) context.CLIContext {
	var rpc rpcclient.Client

	genOnly := viper.GetBool(flags.FlagGenerateOnly)
	fromAddress, fromName, err := context.GetFromFields(from, genOnly)
	if err != nil {
		fmt.Printf("failed to get from fields: %v", err)
		os.Exit(1)
	}

	if !genOnly {
		if nodeURI != "" {
			rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
		}
	}

	ctx := context.CLIContext{
		Client:        rpc,
		Output:        os.Stdout,
		NodeURI:       nodeURI,
		From:          from,
		OutputFormat:  viper.GetString(cli.OutputFlag),
		Height:        viper.GetInt64(flags.FlagHeight),
		TrustNode:     viper.GetBool(flags.FlagTrustNode),
		UseLedger:     viper.GetBool(flags.FlagUseLedger),
		BroadcastMode: viper.GetString(flags.FlagBroadcastMode),
		Simulate:      viper.GetBool(flags.FlagDryRun),
		GenerateOnly:  genOnly,
		FromAddress:   fromAddress,
		FromName:      fromName,
		Indent:        viper.GetBool(flags.FlagIndentResponse),
		SkipConfirm:   viper.GetBool(flags.FlagSkipConfirmation),
	}

	// create a verifier for the specific chain ID and RPC client
	verifier, err := IBCCreateVerifier(ctx, chainID, cacheSize)
	if err != nil && viper.IsSet(flags.FlagTrustNode) {
		fmt.Printf("failed to create verifier: %s\n", err)
		os.Exit(1)
	}

	return ctx.WithVerifier(verifier)
}

// CreateVerifier returns a Tendermint verifier from a CLIContext object and
// cache size. An error is returned if the CLIContext is missing required values
// or if the verifier could not be created. A CLIContext must at the very least
// have the chain ID and home directory set. If the CLIContext has TrustNode
// enabled, no verifier will be created.
func IBCCreateVerifier(ctx context.CLIContext, chainID string, cacheSize int) (tmlite.Verifier, error) {
	if ctx.TrustNode {
		return nil, nil
	}

	homeDir := viper.GetString(flags.FlagHome)

	switch {
	case chainID == "":
		return nil, errors.New("must provide a valid chain ID to create verifier")

	case homeDir == "":
		return nil, errors.New("must provide a valid home directory to create verifier")

	case ctx.Client == nil && ctx.NodeURI == "":
		return nil, errors.New("must provide a valid RPC client or RPC URI to create verifier")
	}

	// create an RPC client based off of the RPC URI if no RPC client exists
	client := ctx.Client
	if client == nil {
		client = rpcclient.NewHTTP(ctx.NodeURI, "/websocket")
	}

	return tmliteproxy.NewVerifier(
		chainID, filepath.Join(homeDir, chainID, ".lite_verifier"),
		client, log.NewNopLogger(), cacheSize,
	)
}

// NOTE: pass in marshalled structs that have been unmarshaled
// because this function will panic on marshaling errors
func PrintOutput(ctx context.CLIContext, toPrint interface{}) error {
	var (
		out []byte
		err error
	)

	switch ctx.OutputFormat {
	case "text":
		out, err = yaml.Marshal(&toPrint)

	case "json":
		if ctx.Indent {
			out, err = ctx.Codec.MarshalJSONIndent(toPrint, "", "  ")
		} else {
			out, err = ctx.Codec.MarshalJSON(toPrint)
		}
	}

	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}
