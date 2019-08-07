package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodeSpace cTypes.CodespaceType = ModuleName
	
	CodeFeedbackCannotRegister cTypes.CodeType = 800
	CodeInvalidInputsOutputs   cTypes.CodeType = 801
)

func ErrFeedbackCannotRegister(msg string) cTypes.Error {
	return cTypes.NewError(DefaultCodeSpace, CodeFeedbackCannotRegister, msg)
}

// ErrNoInputs is an error
func ErrNoInputs(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}
