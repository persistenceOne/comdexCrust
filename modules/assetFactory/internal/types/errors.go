package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

type CodeType cTypes.CodeType

const (
	DefaultCodeSpace         cTypes.CodespaceType = ModuleName
	CodeInvalidAmount        cTypes.CodeType      = 201
	CodeInvalidString        cTypes.CodeType      = 203
	CodeInvalidInputsOutputs cTypes.CodeType      = 204
)

func ErrInvalidAmount(codeSpace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codeSpace, CodeInvalidAmount, msg)
	}
	return cTypes.NewError(codeSpace, CodeInvalidAmount, "invalid Amount")
}

func ErrInvalidString(codeSpace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codeSpace, CodeInvalidString, msg)
	}
	return cTypes.NewError(codeSpace, CodeInvalidString, "Invalid string")
}

func ErrNoInputs(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}
