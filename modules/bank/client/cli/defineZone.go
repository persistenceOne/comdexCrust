package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/persistenceSDK/codec"
	"github.com/persistenceOne/persistenceSDK/modules/acl"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/auth/client/utils"
	bankTypes "github.com/persistenceOne/persistenceSDK/modules/bank/internal/types"
)

func DefineZoneCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defineZone",
		Short: "define a zone address in acl",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)
			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			strZoneID := viper.GetString(FlagZoneID)
			zoneID, err := acl.GetZoneIDFromString(strZoneID)
			if err != nil {
				return nil
			}

			msg := bankTypes.BuildMsgDefineZone(cliCtx.GetFromAddress(), to, zoneID)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsZoneID)
	cmd.Flags().AddFlagSet(fsTo)
	return cmd
}
