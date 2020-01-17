package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	persistenceSDKContext "github.com/persistenceOne/persistenceSDK/client/context"
	"github.com/persistenceOne/persistenceSDK/modules/auth"
	"github.com/persistenceOne/persistenceSDK/modules/auth/client/utils"
	"github.com/persistenceOne/persistenceSDK/modules/bank"
	ibcclient "github.com/persistenceOne/persistenceSDK/modules/ibc/02-client"
	"github.com/persistenceOne/persistenceSDK/modules/ibc/02-client/types/tendermint"
	cutils "github.com/persistenceOne/persistenceSDK/modules/ibc/04-channel/client/utils"
	"github.com/persistenceOne/persistenceSDK/modules/ibc/20-transfer/types"
	commitment "github.com/persistenceOne/persistenceSDK/modules/ibc/23-commitment"
	persistenceSDKTypes "github.com/persistenceOne/persistenceSDK/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IBC transfer flags
var (
	FlagSource   = "source"
	FlagNode1    = "node1"
	FlagNode2    = "node2"
	FlagFrom1    = "from1"
	FlagFrom2    = "from2"
	FlagChainId2 = "chain-id2"
	FlagSequence = "packet-sequence"
	FlagTimeout  = "timeout"
)

// GetTxCmd returns the transaction commands for IBC fungible token transfer
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "transfer",
		Short: "IBC fungible token transfer transaction subcommands",
	}
	txCmd.AddCommand(
		GetTransferTxCmd(cdc),
		GetMsgRecvPacketCmd(cdc),
		GetIssueAssetTransferTxCmd(cdc),
		GetReceiveIssueAssetPacket(cdc),
	)

	return txCmd
}

// GetTransferTxCmd returns the command to create a NewMsgTransfer transaction
func GetTransferTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [src-port] [src-channel] [receiver] [amount]",
		Short: "Transfer fungible token through IBC",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			ctx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)

			sender := ctx.GetFromAddress()
			srcPort := args[0]
			srcChannel := args[1]
			receiver, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			// parse coin trying to be sent
			coins, err := sdk.ParseCoins(args[3])
			if err != nil {
				return err
			}

			source := viper.GetBool(FlagSource)

			msg := types.NewMsgTransfer(srcPort, srcChannel, coins, sender, receiver, source)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Bool(FlagSource, false, "Pass flag for sending token from the source chain")
	cmd.Flags().String(flags.FlagFrom, "", "key in local keystore to send from")
	return cmd
}

// GetMsgRecvPacketCmd returns the command to create a MsgRecvTransferPacket transaction
func GetMsgRecvPacketCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recv-packet [sending-port-id] [sending-channel-id] [client-id]",
		Short: "Creates and sends a SendPacket message",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)

			node2 := viper.GetString(FlagNode2)
			cid1 := viper.GetString(flags.FlagChainID)
			cid2 := viper.GetString(FlagChainId2)
			cliCtx2 := persistenceSDKContext.NewCLIContextIBC(cliCtx.GetFromAddress().String(), cid2, node2).WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)

			header, err := tendermint.GetHeader(cliCtx2)
			if err != nil {
				return err
			}

			sourcePort, sourceChannel, clientid := args[0], args[1], args[2]

			passphrase, err := keys.GetPassphrase(viper.GetString(flags.FlagFrom))
			if err != nil {
				return nil
			}

			viper.Set(flags.FlagChainID, cid1)
			msgUpdateClient := ibcclient.NewMsgUpdateClient(clientid, header, cliCtx.GetFromAddress())
			fmt.Printf("%v <- %-23v", cid2, msgUpdateClient.Type())
			res, err := utils.CompleteAndBroadcastTx(txBldr, cliCtx, []sdk.Msg{msgUpdateClient}, passphrase)
			if err != nil || !persistenceSDKTypes.IsOK(res) {
				fmt.Println(res)
				return err
			}
			fmt.Printf(" [OK] txid(%v) client(%v)\n", res.TxHash, clientid)

			viper.Set(flags.FlagChainID, cid2)
			sequence := uint64(viper.GetInt(FlagSequence))
			packetRes, err := cutils.QueryPacket(cliCtx2.WithHeight(header.Height-1), sourcePort, sourceChannel, sequence, uint64(viper.GetInt(FlagTimeout)), "ibc")
			if err != nil {
				return err
			}
			viper.Set(flags.FlagChainID, cid1)

			msg := types.NewMsgRecvPacket(packetRes.Packet, []commitment.Proof{packetRes.Proof}, packetRes.ProofHeight, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]
	cmd.Flags().Bool(FlagSource, false, "Pass flag for sending token from the source chain")
	cmd.Flags().String(FlagNode2, "tcp://localhost:26657", "RPC port for the second chain")
	cmd.Flags().String(FlagChainId2, "", "chain-id for the second chain")
	cmd.Flags().String(FlagSequence, "", "sequence for the packet")
	cmd.Flags().String(FlagTimeout, "", "timeout for the packet")
	return cmd
}

func GetIssueAssetTransferTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibcIssueAsset [src-port] [src-channel] [to-address]",
		Short: "Transfer command to create asset on ibc and main chain through IBC",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			ctx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)

			issuerAddress := ctx.GetFromAddress()
			srcPort := args[0]
			srcChannel := args[1]

			toAddress, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			assetPeg := &persistenceSDKTypes.BaseAssetPeg{
				AssetQuantity: viper.GetInt64(FlagAssetQuantity),
				AssetType:     viper.GetString(FlagAssetType),
				AssetPrice:    viper.GetInt64(FlagAssetPrice),
				DocumentHash:  viper.GetString(FlagDocumentHash),
				QuantityUnit:  viper.GetString(FlagQuantityUnit),
				Moderated:     true,
			}

			msg := types.NewMsgIssueAssetTransfer(srcPort, srcChannel, bank.NewIssueAsset(issuerAddress, toAddress, assetPeg))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(ctx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flags.FlagFrom, "", "key in local keystore to send from")
	cmd.Flags().AddFlagSet(fsDocumentHash)
	cmd.Flags().AddFlagSet(fsAssetType)
	cmd.Flags().AddFlagSet(fsAssetPrice)
	cmd.Flags().AddFlagSet(fsAssetQuantity)
	cmd.Flags().AddFlagSet(fsQuantityUnit)

	cmd.MarkFlagRequired(FlagAssetQuantity)
	cmd.MarkFlagRequired(FlagAssetType)
	cmd.MarkFlagRequired(FlagAssetPrice)
	cmd.MarkFlagRequired(FlagDocumentHash)
	cmd.MarkFlagRequired(FlagQuantityUnit)
	return cmd
}

func GetReceiveIssueAssetPacket(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receiveIssueAsset [sending-port-id] [sending-channel-id] [client-id]",
		Short: "Receive command to create asset on ibc and main chain through IBC",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)

			node2 := viper.GetString(FlagNode2)
			cid1 := viper.GetString(flags.FlagChainID)
			cid2 := viper.GetString(FlagChainId2)
			cliCtx2 := persistenceSDKContext.NewCLIContextIBC(cliCtx.GetFromAddress().String(), cid2, node2).WithCodec(cdc).WithBroadcastMode(flags.BroadcastBlock)

			header, err := tendermint.GetHeader(cliCtx2)
			if err != nil {
				return err
			}

			sourcePort, sourceChannel, clientid := args[0], args[1], args[2]

			passphrase, err := keys.GetPassphrase(viper.GetString(flags.FlagFrom))
			if err != nil {
				return nil
			}

			viper.Set(flags.FlagChainID, cid1)
			msgUpdateClient := ibcclient.NewMsgUpdateClient(clientid, header, cliCtx.GetFromAddress())
			fmt.Printf("%v <- %-23v", cid2, msgUpdateClient.Type())
			res, err := utils.CompleteAndBroadcastTx(txBldr, cliCtx, []sdk.Msg{msgUpdateClient}, passphrase)
			if err != nil || !persistenceSDKTypes.IsOK(res) {
				fmt.Println(res)
				return err
			}
			fmt.Printf(" [OK] txid(%v) client(%v)\n", res.TxHash, clientid)

			viper.Set(flags.FlagChainID, cid2)
			sequence := uint64(viper.GetInt(FlagSequence))
			packetRes, err := cutils.QueryPacket(cliCtx2.WithHeight(header.Height-1), sourcePort, sourceChannel, sequence, uint64(viper.GetInt(FlagTimeout)), "ibc")
			if err != nil {
				return err
			}
			viper.Set(flags.FlagChainID, cid1)

			msg := types.NewMsgReceiveIssueAssetPacket(packetRes.Packet, []commitment.Proof{packetRes.Proof}, packetRes.ProofHeight, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]
	cmd.Flags().String(FlagNode2, "tcp://localhost:26657", "RPC port for the second chain")
	cmd.Flags().String(FlagChainId2, "", "chain-id for the second chain")
	cmd.Flags().String(FlagSequence, "", "sequence for the packet")
	cmd.Flags().String(FlagTimeout, "", "timeout for the packet")
	return cmd
}
