package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/comdexCrust/modules/acl"
	"github.com/persistenceOne/comdexCrust/modules/assetFactory"
	client "github.com/persistenceOne/comdexCrust/modules/ibc/02-client"
	connection "github.com/persistenceOne/comdexCrust/modules/ibc/03-connection"
	channel "github.com/persistenceOne/comdexCrust/modules/ibc/04-channel"
	port "github.com/persistenceOne/comdexCrust/modules/ibc/05-port"
	transfer "github.com/persistenceOne/comdexCrust/modules/ibc/20-transfer"
)

// Keeper defines each ICS keeper for IBC
type Keeper struct {
	ClientKeeper     client.Keeper
	ConnectionKeeper connection.Keeper
	ChannelKeeper    channel.Keeper
	PortKeeper       port.Keeper
	TransferKeeper   transfer.Keeper
}

// NewKeeper creates a new ibc Keeper
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType,
	bk transfer.BankKeeper, sk transfer.SupplyKeeper, assetKeeper assetFactory.Keeper, aclKeeper acl.Keeper,
) Keeper {
	clientKeeper := client.NewKeeper(cdc, key, codespace)
	connectionKeeper := connection.NewKeeper(cdc, key, codespace, clientKeeper)
	portKeeper := port.NewKeeper(cdc, key, codespace)
	channelKeeper := channel.NewKeeper(cdc, key, codespace, clientKeeper, connectionKeeper, portKeeper)
	transferKeeper := transfer.NewKeeper(cdc, key, codespace, clientKeeper, connectionKeeper, channelKeeper, bk, sk, assetKeeper, aclKeeper)

	return Keeper{
		ClientKeeper:     clientKeeper,
		ConnectionKeeper: connectionKeeper,
		ChannelKeeper:    channelKeeper,
		PortKeeper:       portKeeper,
		TransferKeeper:   transferKeeper,
	}
}
