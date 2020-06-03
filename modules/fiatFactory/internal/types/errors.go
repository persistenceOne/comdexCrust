package types

import (
	ctypes "github.com/cosmos/cosmos-sdk/types"
)

type CodeType ctypes.CodeType

const (
	DefaultCodeSpace         ctypes.CodespaceType = ModuleName
	CodeInvalidAmount        ctypes.CodeType      = 601
	CodeInvalidString        ctypes.CodeType      = 602
	CodeInvalidInputsOutputs ctypes.CodeType      = 603
	CodeInvalidPegHash       ctypes.CodeType      = 604
)

func ErrInvalidAmount(codespace ctypes.CodespaceType, msg string) ctypes.Error {
	if msg != "" {
		return ctypes.NewError(codespace, CodeInvalidAmount, msg)
	}
	return ctypes.NewError(codespace, CodeInvalidAmount, "invalid Amount")
}

func ErrInvalidString(codespace ctypes.CodespaceType, msg string) ctypes.Error {
	if msg != "" {
		return ctypes.NewError(codespace, CodeInvalidString, msg)
	}
	return ctypes.NewError(codespace, CodeInvalidString, "Invalid string")
}

func ErrNoInputs(codespace ctypes.CodespaceType) ctypes.Error {
	return ctypes.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

func ErrInvalidPegHash(codespace ctypes.CodespaceType) ctypes.Error {
	return ctypes.NewError(codespace, CodeInvalidPegHash, "invalid peg hash")
}
