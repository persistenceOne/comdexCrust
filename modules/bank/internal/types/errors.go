package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace cTypes.CodespaceType = ModuleName

	CodeSendDisabled         cTypes.CodeType = 101
	CodeInvalidInputsOutputs cTypes.CodeType = 102
	CodeNegativeAmount       cTypes.CodeType = 103
)

// ErrNoInputs is an error
func ErrNoInputs(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

// ErrNoOutputs is an error
func ErrNoOutputs(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidInputsOutputs, "no outputs to send transaction")
}

// ErrInputOutputMismatch is an error
func ErrInputOutputMismatch(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeInvalidInputsOutputs, "sum inputs != sum outputs")
}

// ErrSendDisabled is an error
func ErrSendDisabled(codespace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codespace, CodeSendDisabled, "send transactions are currently disabled")
}

func ErrNegativeAmount(codeSpace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codeSpace, CodeNegativeAmount, msg)
	}
	return cTypes.NewError(codeSpace, CodeNegativeAmount, "Amount should not be zero")
}
