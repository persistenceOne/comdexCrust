package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName
	
	CodeSendDisabled         sdk.CodeType = 101
	CodeInvalidInputsOutputs sdk.CodeType = 102
	CodeNegativeAmount       sdk.CodeType = 103
)

// ErrNoInputs is an error
func ErrNoInputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

// ErrNoOutputs is an error
func ErrNoOutputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "no outputs to send transaction")
}

// ErrInputOutputMismatch is an error
func ErrInputOutputMismatch(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "sum inputs != sum outputs")
}

// ErrSendDisabled is an error
func ErrSendDisabled(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSendDisabled, "send transactions are currently disabled")
}

func ErrNegativeAmount(codeSpace sdk.CodespaceType, msg string) sdk.Error {
	if msg != "" {
		return sdk.NewError(codeSpace, CodeNegativeAmount, msg)
	}
	return sdk.NewError(codeSpace, CodeNegativeAmount, "Amount should not be zero")
}
