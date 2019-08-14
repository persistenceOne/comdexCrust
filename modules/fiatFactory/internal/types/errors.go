package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

type CodeType cTypes.CodeType

const (
	DefaultCodeSpace         cTypes.CodespaceType = ModuleName
	CodeInvalidAmount        cTypes.CodeType      = 601
	CodeInvalidString        cTypes.CodeType      = 602
	CodeInvalidInputsOutputs cTypes.CodeType      = 603
	CodeInvalidPegHash       cTypes.CodeType      = 604
)

func ErrInvalidAmount(codespace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codespace, CodeInvalidAmount, msg)
	}
	return cTypes.NewError(codespace, CodeInvalidAmount, "invalid Amount")
}

func ErrInvalidString(codespace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codespace, CodeInvalidString, msg)
	}
	return cTypes.NewError(codespace, CodeInvalidString, "Invalid string")
}

func ErrNoInputs(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

func ErrInvalidPegHash(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidPegHash, "invalid peg hash")
}
