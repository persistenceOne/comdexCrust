package types

import (
	"encoding/json"

	cTypes "github.com/cosmos/cosmos-sdk/types"
)

// DefineZone : singular define zone message
type DefineZone struct {
	From   cTypes.AccAddress `json:"from"`
	To     cTypes.AccAddress `json:"to"`
	ZoneID ZoneID            `json:"zoneID"`
}

// NewDefineZone : new define zone struct
func NewDefineZone(from cTypes.AccAddress, to cTypes.AccAddress, zoneID ZoneID) DefineZone {
	return DefineZone{from, to, zoneID}
}

// GetSignBytes : get bytes to sign
func (in DefineZone) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		From   string `json:"from"`
		To     string `json:"to"`
		ZoneID string `json:"zoneID"`
	}{
		From:   in.From.String(),
		To:     in.To.String(),
		ZoneID: in.ZoneID.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineZone) ValidateBasic() cTypes.Error {
	if len(in.From) == 0 {
		return cTypes.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return cTypes.ErrInvalidAddress(in.To.String())
	} else if len(in.ZoneID) == 0 {
		return cTypes.ErrInvalidAddress(in.ZoneID.String())
	}
	return nil
}

// MsgDefineZones : message define zones
type MsgDefineZones struct {
	DefineZones []DefineZone `json:"defineZones"`
}

// NewMsgDefineZones : new message define zones
func NewMsgDefineZones(defineZones []DefineZone) MsgDefineZones {
	return MsgDefineZones{defineZones}
}

var _ cTypes.Msg = MsgDefineZones{}

// Type : implements msg
func (msg MsgDefineZones) Type() string { return "defineZone" }

func (msg MsgDefineZones) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgDefineZones) ValidateBasic() cTypes.Error {
	if len(msg.DefineZones) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.DefineZones {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineZones) GetSignBytes() []byte {
	var defineZones []json.RawMessage
	for _, defineZone := range msg.DefineZones {
		defineZones = append(defineZones, defineZone.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		DefineZones []json.RawMessage `json:"defineZones"`
	}{
		DefineZones: defineZones,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineZones) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.DefineZones))
	for i, in := range msg.DefineZones {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineZones : build define zones message
func BuildMsgDefineZones(from cTypes.AccAddress, to cTypes.AccAddress, zoneID ZoneID, msgs []DefineZone) []DefineZone {
	defineZone := NewDefineZone(from, to, zoneID)
	msgs = append(msgs, defineZone)
	return msgs
}

// BuildMsgDefineZoneWithDefineZones : build define zones message
func BuildMsgDefineZoneWithDefineZones(msgs []DefineZone) cTypes.Msg {
	return NewMsgDefineZones(msgs)
}

// BuildMsgDefineZone : build define zones message
func BuildMsgDefineZone(from cTypes.AccAddress, to cTypes.AccAddress, zoneID ZoneID) cTypes.Msg {
	defineZone := NewDefineZone(from, to, zoneID)
	return NewMsgDefineZones([]DefineZone{defineZone})
}

// DefineOrganization : singular define organization message
type DefineOrganization struct {
	From           cTypes.AccAddress `json:"from"`
	To             cTypes.AccAddress `json:"to"`
	OrganizationID OrganizationID    `json:"organizationID"`
	ZoneID         ZoneID            `json:"zoneID"`
}

// NewDefineOrganization : new define organization struct
func NewDefineOrganization(from cTypes.AccAddress, to cTypes.AccAddress, organizationID OrganizationID, zoneID ZoneID) DefineOrganization {
	return DefineOrganization{from, to, organizationID, zoneID}
}

// GetSignBytes : get bytes to sign
func (in DefineOrganization) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		From           string `json:"from"`
		To             string `json:"to"`
		OrganizationID string `json:"organizationID"`
		ZoneID         string `json:"zoneID"`
	}{
		From:           in.From.String(),
		To:             in.To.String(),
		OrganizationID: in.OrganizationID.String(),
		ZoneID:         in.ZoneID.String(),
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineOrganization) ValidateBasic() cTypes.Error {
	if len(in.From) == 0 {
		return cTypes.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return cTypes.ErrInvalidAddress(in.To.String())
	} else if len(in.OrganizationID) == 0 {
		return cTypes.ErrInvalidAddress(in.OrganizationID.String())
	} else if len(in.ZoneID) == 0 {
		return cTypes.ErrInvalidAddress(in.ZoneID.String())
	}
	return nil
}

// MsgDefineOrganizations : message define organizations
type MsgDefineOrganizations struct {
	DefineOrganizations []DefineOrganization `json:"defineOrganizations"`
}

// NewMsgDefineOrganizations : new message define organizations
func NewMsgDefineOrganizations(defineOrganizations []DefineOrganization) MsgDefineOrganizations {
	return MsgDefineOrganizations{defineOrganizations}
}

var _ cTypes.Msg = MsgDefineOrganizations{}

// Type : implements msg
func (msg MsgDefineOrganizations) Type() string  { return "defineOrganizations" }
func (msg MsgDefineOrganizations) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgDefineOrganizations) ValidateBasic() cTypes.Error {
	if len(msg.DefineOrganizations) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.DefineOrganizations {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineOrganizations) GetSignBytes() []byte {
	var defineOrganizations []json.RawMessage
	for _, defineOrganization := range msg.DefineOrganizations {
		defineOrganizations = append(defineOrganizations, defineOrganization.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		DefineOrganizations []json.RawMessage `json:"defineOrganizations"`
	}{
		DefineOrganizations: defineOrganizations,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineOrganizations) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.DefineOrganizations))
	for i, in := range msg.DefineOrganizations {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineOrganizations : build define organization message
func BuildMsgDefineOrganizations(from cTypes.AccAddress, to cTypes.AccAddress, organizationID OrganizationID, zoneID ZoneID, msgs []DefineOrganization) []DefineOrganization {
	defineOrganization := NewDefineOrganization(from, to, organizationID, zoneID)
	msgs = append(msgs, defineOrganization)
	return msgs
}

// BuildMsgDefineOrganizationWithMsgs : build define organization message
func BuildMsgDefineOrganizationWithMsgs(msgs []DefineOrganization) cTypes.Msg {
	return NewMsgDefineOrganizations(msgs)
}

// BuildMsgDefineOrganization : build define organization message
func BuildMsgDefineOrganization(from cTypes.AccAddress, to cTypes.AccAddress, organizationID OrganizationID, zoneID ZoneID) cTypes.Msg {
	defineOrganization := NewDefineOrganization(from, to, organizationID, zoneID)
	return NewMsgDefineOrganizations([]DefineOrganization{defineOrganization})
}

// DefineACL : indular define acl message
type DefineACL struct {
	From       cTypes.AccAddress `json:"from"`
	To         cTypes.AccAddress `json:"to"`
	ACLAccount ACLAccount        `json:"aclAccount"`
}

// NewDefineACL : new define acl struct
func NewDefineACL(from cTypes.AccAddress, to cTypes.AccAddress, aclAccount ACLAccount) DefineACL {
	return DefineACL{from, to, aclAccount}
}

// GetSignBytes : get bytes to sign
func (in DefineACL) GetSignBytes() []byte {
	bin, err := ModuleCdc.MarshalJSON(struct {
		From       string     `json:"from"`
		To         string     `json:"to"`
		ACLAccount ACLAccount `json:"aclAccount"`
	}{
		From:       in.From.String(),
		To:         in.To.String(),
		ACLAccount: in.ACLAccount,
	})
	if err != nil {
		panic(err)
	}
	return bin
}

// ValidateBasic : Validate Basic
func (in DefineACL) ValidateBasic() cTypes.Error {
	if len(in.From) == 0 {
		return cTypes.ErrInvalidAddress(in.From.String())
	} else if len(in.To) == 0 {
		return cTypes.ErrInvalidAddress(in.To.String())
	} else if len(in.ACLAccount.GetAddress()) == 0 {
		return cTypes.ErrInvalidAddress(in.ACLAccount.GetAddress().String())
	} else if len(in.ACLAccount.GetZoneID()) == 0 {
		return cTypes.ErrUnauthorized("ZoneID should not be empty.")
	} else if len(in.ACLAccount.GetOrganizationID()) == 0 {
		return cTypes.ErrUnauthorized("OrganizationID should not be empty.")
	}
	return nil
}

// MsgDefineACLs : message define acls
type MsgDefineACLs struct {
	DefineACLs []DefineACL `json:"defineACLs"`
}

// NewMsgDefineACLs : new message define acls
func NewMsgDefineACLs(defineACLs []DefineACL) MsgDefineACLs {
	return MsgDefineACLs{defineACLs}
}

var _ cTypes.Msg = MsgDefineACLs{}

// Type : implements msg
func (msg MsgDefineACLs) Type() string { return "defineACL" }

func (msg MsgDefineACLs) Route() string { return RouterKey }

// ValidateBasic : implements msg
func (msg MsgDefineACLs) ValidateBasic() cTypes.Error {
	if len(msg.DefineACLs) == 0 {
		return ErrNoInputs(DefaultCodeSpace).TraceSDK("")
	}
	for _, in := range msg.DefineACLs {
		if err := in.ValidateBasic(); err != nil {
			return err.TraceSDK("")
		}
	}
	return nil
}

// GetSignBytes : implements msg
func (msg MsgDefineACLs) GetSignBytes() []byte {
	var defineACLs []json.RawMessage
	for _, defineACL := range msg.DefineACLs {
		defineACLs = append(defineACLs, defineACL.GetSignBytes())
	}

	b, err := ModuleCdc.MarshalJSON(struct {
		DefineACLs []json.RawMessage `json:"defineACLs"`
	}{
		DefineACLs: defineACLs,
	})
	if err != nil {
		panic(err)
	}
	return b
}

// GetSigners : implements msg
func (msg MsgDefineACLs) GetSigners() []cTypes.AccAddress {
	addrs := make([]cTypes.AccAddress, len(msg.DefineACLs))
	for i, in := range msg.DefineACLs {
		addrs[i] = in.From
	}
	return addrs
}

// BuildMsgDefineACLs : build define acls message
func BuildMsgDefineACLs(from cTypes.AccAddress, to cTypes.AccAddress, aclAccount ACLAccount, msgs []DefineACL) []DefineACL {
	defineACL := NewDefineACL(from, to, aclAccount)
	msgs = append(msgs, defineACL)
	return msgs
}

// BuildMsgDefineACLWithACLs : build define acls message
func BuildMsgDefineACLWithACLs(msgs []DefineACL) cTypes.Msg {
	return NewMsgDefineACLs(msgs)
}

// BuildMsgDefineACL : build define acls message
func BuildMsgDefineACL(from cTypes.AccAddress, to cTypes.AccAddress, aclAccount ACLAccount) cTypes.Msg {
	defineACL := NewDefineACL(from, to, aclAccount)
	return NewMsgDefineACLs([]DefineACL{defineACL})
}
