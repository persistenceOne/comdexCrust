package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/persistenceOne/comdexCrust/codec"
	"github.com/persistenceOne/comdexCrust/modules/acl"
	"github.com/persistenceOne/comdexCrust/modules/auth"
	"github.com/persistenceOne/comdexCrust/modules/auth/client/utils"
	bankTypes "github.com/persistenceOne/comdexCrust/modules/bank/internal/types"
)

func DefineOrganizationCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "defineOrganization",
		Short: "define an organization address in acl",
		RunE: func(cmd *cobra.Command, args []string) error {

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			toStr := viper.GetString(FlagTo)
			to, err := cTypes.AccAddressFromBech32(toStr)
			if err != nil {
				return nil
			}

			strOrganizationID := viper.GetString(FlagOrganizationID)
			organizationID, err := acl.GetOrganizationIDFromString(strOrganizationID)
			if err != nil {
				return nil
			}

			strZoneID := viper.GetString(FlagZoneID)
			zoneID, err := acl.GetZoneIDFromString(strZoneID)
			if err != nil {
				return nil
			}

			msg := bankTypes.BuildMsgDefineOrganization(cliCtx.GetFromAddress(), to, organizationID, zoneID)

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsOrganizationID)
	cmd.Flags().AddFlagSet(fsZoneID)
	cmd.Flags().AddFlagSet(fsTo)
	return cmd
}
