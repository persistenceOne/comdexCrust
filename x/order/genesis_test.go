package order

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_InitOrder(t *testing.T) {
	err := InitOrder(ctx, orderKeeper)
	require.Nil(t, err)
}
