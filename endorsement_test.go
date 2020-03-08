package tezosprotocol_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v2"
	"github.com/stretchr/testify/require"
)

func TestEncodeEndorsement(t *testing.T) {
	require := require.New(t)
	origination := &tezosprotocol.Endorsement{
		Level: big.NewInt(9),
	}
	encodedBytes, err := origination.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "0000000009"
	require.Equal(expected, encoded)
}

func TestDecodeEndorsement(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("0000000009")
	require.NoError(err)
	endorsement := tezosprotocol.Endorsement{}
	require.NoError(endorsement.UnmarshalBinary(encoded))
	require.Equal("9", endorsement.Level.String())
}
