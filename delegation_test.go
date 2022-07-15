package tezosprotocol_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v3"
	"github.com/stretchr/testify/require"
)

func TestEncodeDelegation(t *testing.T) {
	require := require.New(t)
	delegate := tezosprotocol.ContractID("tz1ddb9NMYHZi5UzPdzTZMYQQZoMub195zgv")
	origination := &tezosprotocol.Delegation{
		Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
		Fee:          big.NewInt(1266),
		Counter:      big.NewInt(1),
		GasLimit:     big.NewInt(10100),
		StorageLimit: big.NewInt(277),
		Delegate:     &delegate,
	}
	encodedBytes, err := origination.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "6e0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e9502ff00c55cf02dbeecc978d9c84625dcae72bb77ea4fbd"
	require.Equal(expected, encoded)
}

func TestDecodeDelegation(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("6e0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e9502ff00c55cf02dbeecc978d9c84625dcae72bb77ea4fbd")
	require.NoError(err)
	delegation := tezosprotocol.Delegation{}
	require.NoError(delegation.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"), delegation.Source)
	require.Equal("1266", delegation.Fee.String())
	require.Equal("1", delegation.Counter.String())
	require.Equal("10100", delegation.GasLimit.String())
	require.Equal("277", delegation.StorageLimit.String())
	require.NotNil(delegation.Delegate)
	require.Equal(tezosprotocol.ContractID("tz1ddb9NMYHZi5UzPdzTZMYQQZoMub195zgv"), *delegation.Delegate)
}
