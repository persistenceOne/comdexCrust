package types

import (
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

type CodeType cTypes.CodeType

const (
	DefaultCodeSpace cTypes.CodespaceType = ModuleName

	CodeInvalidInputsOutputs cTypes.CodeType = 701
	CodeUnauthorized         cTypes.CodeType = 702
)

func ErrNoInputsOutputs(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

func ErrUnauthorized(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeUnauthorized, "Unauthorized transaction")
}
