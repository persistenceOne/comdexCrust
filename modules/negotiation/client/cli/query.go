package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/comdexCrust/codec"
	negotiationTypes "github.com/persistenceOne/comdexCrust/modules/negotiation/internal/types"
	"github.com/persistenceOne/comdexCrust/types"
)

func GetNegotiationCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "[negotiation-id]",
		Short: "Query negotiation details",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext()

			negotiationID, err := types.GetNegotiationIDFromString(args[0])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryStore(negotiationID.Bytes(), negotiationTypes.ModuleName)
			if err != nil {
				return err
			}

			if res == nil {
				return cTypes.ErrUnknownAddress("No negotiation with negotiationID " + args[0] +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			var _negotiation types.Negotiation
			err = cdc.UnmarshalBinaryBare(res, &_negotiation)
			if err != nil {
				return err
			}

			output, err := cdc.MarshalJSON(_negotiation)
			if err != nil {
				return err
			}

			fmt.Println(string(output))
			return nil
		},
	}

	cmd.Flags().AddFlagSet(fsNegotiationID)
	return cmd
}
