package types

import cTypes "github.com/cosmos/cosmos-sdk/types"

type CodeType cTypes.CodeType

const (
	DefaultCodeSpace cTypes.CodespaceType = ModuleName

	CodeInvalidNegotiation   cTypes.CodeType = 600
	CodeInvalidSignature     cTypes.CodeType = 601
	CodeInvalidBid           cTypes.CodeType = 602
	CodeUnauthorized         cTypes.CodeType = 603
	CodeInvalidInputsOutputs cTypes.CodeType = 604
	CodeNegativeAmount       cTypes.CodeType = 605
)

func ErrInvalidNegotiationID(codespace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codespace, CodeInvalidNegotiation, msg)
	}
	return cTypes.NewError(codespace, CodeInvalidNegotiation, "negotiation doesn't exist")
}

func ErrVerifySignature(codespace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codespace, CodeInvalidSignature, msg)
	}
	return cTypes.NewError(codespace, CodeInvalidSignature, "signature verification failed")
}

func ErrInvalidBid(codespace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codespace, CodeInvalidBid, msg)
	}
	return cTypes.NewError(codespace, CodeInvalidBid, "")
}

func ErrUnauthorized(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeUnauthorized, "Unauthorized transaction")
}

// ErrNoInputs is an error
func ErrNoInputs(codeSpace cTypes.CodespaceType) cTypes.Error {
	return cTypes.NewError(codeSpace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

func ErrNegativeAmount(codeSpace cTypes.CodespaceType, msg string) cTypes.Error {
	if msg != "" {
		return cTypes.NewError(codeSpace, CodeNegativeAmount, msg)
	}
	return cTypes.NewError(codeSpace, CodeNegativeAmount, "Amount should not be zero")
}
