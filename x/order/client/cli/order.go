package cli

import (
	"fmt"

	"github.com/commitHub/commitBlockchain/client/context"
	sdk "github.com/commitHub/commitBlockchain/types"

	wire "github.com/commitHub/commitBlockchain/wire"
	"github.com/commitHub/commitBlockchain/x/order"
	"github.com/spf13/cobra"
)

//GetOrderCmd : command to get order details
func GetOrderCmd(storeName string, cdc *wire.Codec, decoder sdk.OrderDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "order [negotiationID]",
		Short: "Query order details",
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx := context.NewCLIContext()

			// negotiationID := sdk.NegotiationID(viper.GetString(FlagNegotiationID))
			nego := args[0]

			negotiationID, err := sdk.GetNegotiationIDHex(nego)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(order.StoreKey([]byte(negotiationID)), storeName)
			if err != nil {
				return err
			}

			// Check if account was found
			if res == nil {
				return sdk.ErrUnknownAddress("No order with negotiationID " + nego +
					" was found in the state.\nAre you sure there has been a transaction involving it?")
			}

			// decode the value
			orderResponse, err := decoder(res)
			if err != nil {
				return err
			}

			// print out whole account
			output, err := wire.MarshalJSONIndent(cdc, orderResponse)
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

//GetOrderDecoder : get order decoder
func GetOrderDecoder(cdc *wire.Codec) sdk.OrderDecoder {
	return func(orderBytes []byte) (order sdk.Order, err error) {
		err = cdc.UnmarshalBinaryBare(orderBytes, &order)
		if err != nil {
			panic(err)
		}
		return order, err
	}
}
