package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func IsOK(res sdk.TxResponse) bool {
	for _, lg := range res.Logs {
		if !lg.Success {
			return false
		}
	}
	return true
}
