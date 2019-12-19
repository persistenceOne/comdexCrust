package host

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)


// SubModuleName defines the ICS 24 host
const SubModuleName = "host"

// // IBCCodeSpace is the codespace for all errors defined in the ibc module
// const IBCCodeSpace = "ibc"

// Error codes specific to the ibc host submodule
const (
	DefaultCodespace sdk.CodespaceType = SubModuleName

	CodeInvalidID     sdk.CodeType = 231
	CodeInvalidPath   sdk.CodeType = 232
	CodeInvalidPacket sdk.CodeType = 233
)

// ErrInvalidID returns a typed ABCI error for an invalid identifier
func ErrInvalidID(codespace sdk.CodespaceType, id string) error {
	return sdk.NewError(
		codespace,
		CodeInvalidID,
		fmt.Sprintf("invalid identifier '%s'", id),
	)
}

// ErrInvalidPath returns a typed ABCI error for an invalid path
func ErrInvalidPath(codespace sdk.CodespaceType, path string) error {
	return sdk.NewError(
		codespace,
		CodeInvalidPath,
		fmt.Sprintf("invalid path '%s'", path),
	)
}

// ErrInvalidPacket returns a typed ABCI error for an invalid identifier
func ErrInvalidPacket(codespace sdk.CodespaceType, msg string) error {
	return sdk.NewError(
		codespace,
		CodeInvalidPacket,
		fmt.Sprintf("invalid packet: '%s'", msg),
	)
}
