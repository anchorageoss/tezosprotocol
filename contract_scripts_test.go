package tezosprotocol_test

import (
	"encoding/hex"
	"testing"

	tezosprotocol "github.com/anchorageoss/tezosprotocol/v2"
	"github.com/stretchr/testify/require"
)

func TestContractScriptUnmarshalBinary(t *testing.T) {
	require := require.New(t)

	// invalid code length
	err := (&tezosprotocol.ContractScript{}).UnmarshalBinary([]byte{})
	require.Error(err)
	require.Contains(err.Error(), "failed to read code length")

	// invalid code
	badCode, err := hex.DecodeString("00000002")
	require.NoError(err)
	err = (&tezosprotocol.ContractScript{}).UnmarshalBinary(badCode)
	require.Error(err)
	require.Contains(err.Error(), "failed to read code")

	// invalid storage length
	badStorageLength, err := hex.DecodeString("00000002C0DE00")
	require.NoError(err)
	err = (&tezosprotocol.ContractScript{}).UnmarshalBinary(badStorageLength)
	require.Error(err)
	require.Contains(err.Error(), "failed to read storage length")

	// invalid storage
	badStorage, err := hex.DecodeString("00000002C0DE00000007")
	require.NoError(err)
	err = (&tezosprotocol.ContractScript{}).UnmarshalBinary(badStorage)
	require.Error(err)
	require.Contains(err.Error(), "failed to read storage")
}
