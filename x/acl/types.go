package acl

import (
	sdk "github.com/commitHub/commitBlockchain/types"
	wire "github.com/commitHub/commitBlockchain/wire"
)

//GetACLAccountDecoder : return a decode for acl accounts
func GetACLAccountDecoder(cdc *wire.Codec) sdk.ACLAccountDecoder {
	return func(aclBytes []byte) (acl sdk.ACLAccount, err error) {
		// acct := new(auth.BaseAccount)
		err = cdc.UnmarshalBinaryBare(aclBytes, &acl)
		if err != nil {
			panic(err)
		}
		return acl, err
	}
}

//GetOrganizationDecoder : return a decode for organization accounts
func GetOrganizationDecoder(cdc *wire.Codec) sdk.OrgDecoder {
	return func(orgBytes []byte) (org sdk.Organization, err error) {
		// acct := new(auth.BaseAccount)
		err = cdc.UnmarshalBinaryBare(orgBytes, &org)
		if err != nil {
			panic(err)
		}
		return org, err
	}
}
