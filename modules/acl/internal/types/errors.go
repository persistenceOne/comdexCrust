package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

type CodeType cTypes.CodeType

const (
	DefaultCodeSpace cTypes.CodespaceType = ModuleName
	
	CodeInvalidInputsOutputs cTypes.CodeType = 101
	CodeInvalidID            cTypes.CodeType = 102
	CodeInvalidAddress       cTypes.CodeType = 103
)

func ErrNoInputs(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

func ErrInvalidID(codespace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codespace, CodeInvalidID, msg)
	}
	return cTypes.NewError(codespace, CodeInvalidID, "invalid ID")
}

func ErrInvalidAddress(codespace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codespace, CodeInvalidAddress, msg)
	}
	return cTypes.NewError(codespace, CodeInvalidAddress, "")
}
