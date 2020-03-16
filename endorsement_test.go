package tezosprotocol_test

import (
	"encoding/hex"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v2"
	"github.com/stretchr/testify/require"
)

func TestEncodeEndorsement(t *testing.T) {
	require := require.New(t)
	origination := &tezosprotocol.Endorsement{
		Level: 999,
	}
	encodedBytes, err := origination.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "00000003e7"
	require.Equal(expected, encoded)
}

func TestDecodeEndorsement(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("00000003e7")
	require.NoError(err)
	endorsement := tezosprotocol.Endorsement{}
	require.NoError(endorsement.UnmarshalBinary(encoded))
	require.Equal(int32(999), endorsement.Level)
}
