package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/store/errors"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/orders/internal/types"
	types2 "github.com/persistenceOne/comdexCrust/types"
)

func GetOrderCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "[negotiation-id]",
		Short: "Query order details",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext()

			negotiationID, err := types2.GetNegotiationIDFromString(args[0])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryStore(negotiationID, types.ModuleName)
			if err != nil {
				return err
			}

			if res == nil {
				return errors.ErrInternal("No order with negotiationID " + args[0] +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			var order types2.Order
			err = cdc.UnmarshalBinaryBare(res, &order)
			if err != nil {
				return err
			}

			output, err := cdc.MarshalJSONIndent(order, "", " ")
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
