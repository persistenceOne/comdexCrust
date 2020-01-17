package cli

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	commitContext "github.com/persistenceOne/persistenceSDK/client/context"
	commitFlags "github.com/persistenceOne/persistenceSDK/client/flags"
	channel "github.com/persistenceOne/persistenceSDK/modules/ibc/04-channel"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	abci "github.com/tendermint/tendermint/abci/types"
)

// GetTxCmd returns the transaction commands for IBC fungible token transfer
func GetQueryCmd(cdc *codec.Codec, storeKey string) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:   "transfer",
		Short: "IBC fungible token transfer query subcommands",
	}

	queryCmd.AddCommand(
		GetCmdQueryNextSequence(cdc, storeKey),
	)

	return queryCmd
}

// GetCmdQueryNextSequence defines the command to query a next receive sequence
func GetCmdQueryNextSequence(cdc *codec.Codec, queryRoute string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next-recv [port-id] [channel-id]",
		Short: "Query a next receive sequence",
		Long: strings.TrimSpace(fmt.Sprintf(`Query an IBC channel end
		
Example:
$ %s query ibc channel next-recv [port-id] [channel-id]
		`, version.ClientName),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			portID := args[0]
			channelID := args[1]

			req := abci.RequestQuery{
				Path:  "store/ibc/key",
				Data:  channel.KeyNextSequenceRecv(portID, channelID),
				Prove: true,
			}

			res, err := commitContext.QueryABCI(cliCtx, req)
			if err != nil {
				return err
			}

			sequence := binary.BigEndian.Uint64(res.Value)

			return commitContext.PrintOutput(cliCtx, sequence)
		},
	}
	cmd.Flags().Bool(commitFlags.FlagProve, true, "show proofs for the query results")

	return cmd
}
