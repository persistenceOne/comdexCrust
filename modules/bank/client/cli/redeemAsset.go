package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	cTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/commitHub/commitBlockchain/types"
	
	"github.com/commitHub/commitBlockchain/codec"
	"github.com/commitHub/commitBlockchain/modules/acl"
	"github.com/commitHub/commitBlockchain/modules/auth"
	"github.com/commitHub/commitBlockchain/modules/auth/client/utils"
	"github.com/commitHub/commitBlockchain/modules/bank/client"
)

func RedeemAssetCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "redeemAsset",
		Short: "Redeem asset from redeemerAddress to the given issuerAddress",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			
			issuerAddress, err := cTypes.AccAddressFromBech32(viper.GetString(FlagTo))
			if err != nil {
				return err
			}
			
			pegHash, err := types.GetAssetPegHashHex(viper.GetString(FlagPegHash))
			if err != nil {
				return err
			}
			
			res, _, err := cliCtx.QueryStore(acl.GetACLAccountKey(issuerAddress), acl.ModuleName)
			if err != nil {
				return err
			}
			
			if len(res) == 0 {
				return cTypes.ErrInternal("Unauthorized transaction")
			}
			
			var account acl.ACLAccount
			err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &account)
			if err != nil {
				return cTypes.ErrInternal("Unmarshal to ACL account faild")
			}
			
			if !account.GetACL().RedeemAsset {
				return cTypes.ErrInternal("Unauthorized transaction")
			}
			
			msg := client.BuildRedeemAssetMsg(issuerAddress, cliCtx.GetFromAddress(), pegHash)
			
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []cTypes.Msg{msg})
		},
	}
	cmd.Flags().AddFlagSet(fsTo)
	cmd.Flags().AddFlagSet(fsPegHash)
	
	return cmd
}
