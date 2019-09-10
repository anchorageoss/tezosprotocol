package tezosprotocol_test

import (
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol"
	"github.com/stretchr/testify/require"
)

// Check that the computed values of the minimum fees match what we've seen from
// the tezos client
func TestMinimumFees(t *testing.T) {
	require := require.New(t)

	require.Equal(int64(1252), tezosprotocol.OriginationMinimumFee)
	require.Equal("1252", tezosprotocol.ComputeMinimumFee(big.NewInt(tezosprotocol.OriginationGasLimit), big.NewInt(tezosprotocol.MinimumOriginationSizeBytes)).String())

	require.Equal(int64(257000), tezosprotocol.OriginationStorageBurn)
}
