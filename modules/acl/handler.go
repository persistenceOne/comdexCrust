package acl

import (
	"fmt"
	
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) cTypes.Handler {
	return func(ctx cTypes.Context, msg cTypes.Msg) cTypes.Result {
		ctx = ctx.WithEventManager(cTypes.NewEventManager())
		switch msg := msg.(type) {
		case MsgDefineZones:
			return handleMsgDefineZones(ctx, k, msg)
		case MsgDefineOrganizations:
			return handleMsgDefineOrganizations(ctx, k, msg)
		case MsgDefineACLs:
			return handleMsgDefineACLAccounts(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return cTypes.ErrUnknownRequest(errMsg).Result()
			
		}
	}
}

func handleMsgDefineZones(ctx cTypes.Context, k Keeper, msg MsgDefineZones) cTypes.Result {
	events := cTypes.EmptyEvents()
	
	for _, zone := range msg.DefineZones {
		events, err := defineZone(ctx, k, zone)
		if err != nil {
			return err.Result()
		}
		events.AppendEvents(events)
	}
	
	ctx.EventManager().EmitEvents(events)
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func defineZone(ctx cTypes.Context, k Keeper, zone DefineZone) (cTypes.Events, cTypes.Error) {
	events := cTypes.EmptyEvents()
	
	if !(k.CheckValidGenesisAddress(ctx, zone.From)) {
		return nil, ErrInvalidAddress(DefaultCodeSpace,
			fmt.Sprintf("Account %v is not the genesis account. Zones can only be defined by the genesis account.", zone.From.String()))
	}
	
	err := k.SetZoneAddress(ctx, zone.ZoneID, zone.To)
	if err != nil {
		return nil, err
	}
	
	event := cTypes.NewEvent(
		EventTypeDefineZone,
		cTypes.NewAttribute(AttributeKeyZoneID, zone.ZoneID.String()),
		cTypes.NewAttribute(AttributeKeyZoneAddress, zone.To.String()),
	)
	events.AppendEvent(event)
	
	return events, nil
}

func handleMsgDefineOrganizations(ctx cTypes.Context, k Keeper, msg MsgDefineOrganizations) cTypes.Result {
	events := cTypes.EmptyEvents()
	
	for _, organization := range msg.DefineOrganizations {
		events, err := defineOrganization(ctx, k, organization)
		if err != nil {
			return err.Result()
		}
		events.AppendEvents(events)
	}
	
	ctx.EventManager().EmitEvents(events)
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func defineOrganization(ctx cTypes.Context, k Keeper, organization DefineOrganization) (cTypes.Events, cTypes.Error) {
	events := cTypes.EmptyEvents()
	
	if !(k.CheckValidZoneAddress(ctx, organization.ZoneID, organization.From)) {
		return nil, ErrInvalidAddress(DefaultCodeSpace, fmt.Sprintf("Account %v is not the zone account. Organizations can only be defined by the zone account.", organization.From.String()))
	}
	
	newOrganization := NewOrganization(organization.To, organization.ZoneID)
	
	err := k.SetOrganization(ctx, organization.OrganizationID, newOrganization)
	if err != nil {
		return nil, err
	}
	
	event := cTypes.NewEvent(
		EventTypeDefineOrganization,
		cTypes.NewAttribute(AttributeKeyOrganizationID, organization.OrganizationID.String()),
		cTypes.NewAttribute(AttributeKeyOrganizationAddress, organization.To.String()),
	)
	events.AppendEvent(event)
	
	return events, nil
}

func handleMsgDefineACLAccounts(ctx cTypes.Context, k Keeper, msg MsgDefineACLs) cTypes.Result {
	events := cTypes.EmptyEvents()
	
	for _, acl := range msg.DefineACLs {
		events, err := defineACLAccount(ctx, k, acl)
		if err != nil {
			return err.Result()
		}
		events.AppendEvents(events)
	}
	
	ctx.EventManager().EmitEvents(events)
	return cTypes.Result{
		Events: ctx.EventManager().Events(),
	}
}

func defineACLAccount(ctx cTypes.Context, k Keeper, acl DefineACL) (cTypes.Events, cTypes.Error) {
	events := cTypes.EmptyEvents()
	
	if !(k.CheckValidGenesisAddress(ctx, acl.From)) {
		if !(k.CheckValidZoneAddress(ctx, acl.ACLAccount.GetZoneID(), acl.From)) {
			if !(k.CheckValidOrganizationAddress(ctx, acl.ACLAccount.GetZoneID(), acl.ACLAccount.GetOrganizationID(), acl.From)) {
				return nil, ErrInvalidAddress(DefaultCodeSpace, fmt.Sprintf("Account %v does not have access to define acl for account %v.", acl.From.String(), acl.To.String()))
			}
		}
	}
	
	err := k.SetACLAccount(ctx, acl.ACLAccount)
	if err != nil {
		return nil, err
	}
	
	event := cTypes.NewEvent(
		EventTypeDefineACL,
		cTypes.NewAttribute(AttributeACLAccountAddress, acl.To.String()),
	)
	
	events.AppendEvent(event)
	
	return events, nil
}
