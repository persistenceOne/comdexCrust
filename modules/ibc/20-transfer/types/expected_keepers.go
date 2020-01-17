package types

import (
	"github.com/persistenceOne/persistenceSDK/modules/bank"
	clientexported "github.com/persistenceOne/persistenceSDK/modules/ibc/02-client/exported"
	connection "github.com/persistenceOne/persistenceSDK/modules/ibc/03-connection"
	channel "github.com/persistenceOne/persistenceSDK/modules/ibc/04-channel"
	channelexported "github.com/persistenceOne/persistenceSDK/modules/ibc/04-channel/exported"
	commitment "github.com/persistenceOne/persistenceSDK/modules/ibc/23-commitment"
	supplyexported "github.com/persistenceOne/persistenceSDK/modules/supply/exported"
	persistenceSDKTypes "github.com/persistenceOne/persistenceSDK/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	IssueAssetsToWallets(ctx sdk.Context, issueAsset bank.IssueAsset) (persistenceSDKTypes.AssetPeg, sdk.Error)
}

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channel.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, packet channelexported.PacketI, portCapability persistenceSDKTypes.CapabilityKey) error
	RecvPacket(ctx sdk.Context, packet channelexported.PacketI, proof commitment.ProofI, proofHeight uint64, acknowledgement []byte, portCapability persistenceSDKTypes.CapabilityKey) (channelexported.PacketI, error)
}

// ClientKeeper defines the expected IBC client keeper
type ClientKeeper interface {
	GetConsensusState(ctx sdk.Context, clientID string) (connection clientexported.ConsensusState, found bool)
}

// ConnectionKeeper defines the expected IBC connection keeper
type ConnectionKeeper interface {
	GetConnection(ctx sdk.Context, connectionID string) (connection connection.ConnectionEnd, found bool)
}

// SupplyKeeper expected supply keeper
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
}
