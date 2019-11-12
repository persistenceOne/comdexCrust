package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodeSpace                cTypes.CodespaceType = "types"
	CodeSendDisabled                cTypes.CodeType      = 101
	CodeInvalidInputsOutputs        cTypes.CodeType      = 102
	CodeNegativeAmount              cTypes.CodeType      = 103
	CodeGoValidation                cTypes.CodeType      = 104
	CodeFromName                    cTypes.CodeType      = 105
	CodeZoneIDFromString            cTypes.CodeType      = 106
	CodeAccAddressFromBech32        cTypes.CodeType      = 107
	CodeOrganizationIDFromString    cTypes.CodeType      = 108
	CodeInvalidOrganization         cTypes.CodeType      = 109
	CodeInvalidOrganizationWithZone cTypes.CodeType      = 110
	CodeInvalidACLAccount           cTypes.CodeType      = 111
	CodeUnAuthorizedTransaction     cTypes.CodeType      = 112
	CodeInvalidQuery                cTypes.CodeType      = 113
	CodeInvalidFields               cTypes.CodeType      = 114
	CodeNotEqual                    cTypes.CodeType      = 115
	CodeResponseDataLengthZero      cTypes.CodeType      = 116
	CodePegHashHex                  cTypes.CodeType      = 117
	CodeMarshal                     cTypes.CodeType      = 118
	CodeUnmarshal                   cTypes.CodeType      = 119
	CodeKeyBase                     cTypes.CodeType      = 120
	CodeSign                        cTypes.CodeType      = 121
	CodeNegotiationIDFromString     cTypes.CodeType      = 122
	CodeZoneIDExists                cTypes.CodeType      = 123
	CodeOrganizationIDExists        cTypes.CodeType      = 124
	CodeFeedbackCannotRegister      cTypes.CodeType      = 125
)

func ErrNoInputs(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

func ErrNoOutputs(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidInputsOutputs, "no outputs to send transaction")
}

func ErrInputOutputMismatch(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidInputsOutputs, "sum inputs != sum outputs")
}

func ErrSendDisabled(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeSendDisabled, "send transactions are currently disabled")
}

func ErrNegativeAmount(codeSpace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codeSpace, CodeNegativeAmount, msg)
	}
	return cTypes.NewError(codeSpace, CodeNegativeAmount, "Amount should not be zero")
}

func ErrGoValidator(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeGoValidation, "Error occurred while go validation")
}

func ErrFromName(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeFromName, "Error occurred while fetching from name")
}

func ErrZoneIDFromString(codeSpace cTypes.CodespaceType, ID string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeZoneIDFromString, "Error occurred while converting zoneID from string "+ID)
}

func ErrAccAddressFromBech32(codeSpace cTypes.CodespaceType, address string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeAccAddressFromBech32, "Error occurred while converting Account Address from bech32 "+address)
}

func ErrOrganizationIDFromString(codeSpace cTypes.CodespaceType, ID string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeOrganizationIDFromString, "Error occurred while converting organizationID from string "+ID)
}

func ErrInvalidOrganization(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidOrganization, "Organization is not defined")
}

func ErrInvalidOrganizationWithZone(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidOrganizationWithZone, "Organization is not belongs to respected zone")
}

func ErrInvalidACLAccount(codeSpace cTypes.CodespaceType, address string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidACLAccount, "Could not query account "+address)
}

func ErrUnAuthorizedTransaction(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeUnAuthorizedTransaction, "Unauthorized transaction")
}

func ErrQueryResponseLengthZero(codeSpace cTypes.CodespaceType, query string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeResponseDataLengthZero, "Could not query zone data "+query)
}

func ErrQuery(codeSpace cTypes.CodespaceType, query string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidQuery, "Error occurred while querying "+query+" data")
}

func ErrEmptyRequestFields(codeSpace cTypes.CodespaceType, fields string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidFields, "Request Fields are invalid "+fields)
}

func ErrNotEqual(codeSpace cTypes.CodespaceType, address1 string, address2 string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeNotEqual, address1+"should be equals to "+address2)
}

func ErrPegHashHex(codeSpace cTypes.CodespaceType, _type string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodePegHashHex, "Error converting "+_type+"to hex")
}

func ErrMarshal(codeSpace cTypes.CodespaceType, param string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeMarshal, "Error marshal "+param)
}

func ErrUnmarshal(codeSpace cTypes.CodespaceType, param string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeUnmarshal, "Error unmarshal "+param)
}

func ErrKeyBase(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeKeyBase, "Error occurred while fetching the keybase")
}

func ErrSign(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeSign, "Error occurred while signing")
}

func ErrNegotiationIDFromString(codeSpace cTypes.CodespaceType, ID string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeNegotiationIDFromString, "Error occurred while converting negotiationID from string "+ID)
}

func ErrZoneIDExists(codeSpace cTypes.CodespaceType, ID string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeZoneIDExists, "Zone ID "+ID+" already exists")
}

func ErrOrganizationIDExists(codeSpace cTypes.CodespaceType, ID string) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeOrganizationIDExists, "Organization ID "+ID+" already exists")
}

func ErrFeedbackCannotRegister(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeFeedbackCannotRegister, "You have already given a  traderFeedback for this transaction")
}
