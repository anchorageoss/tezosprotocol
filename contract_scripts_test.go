package tezosprotocol_test

import (
	"encoding/hex"
	"math"
	"strings"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v2"
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

func TestSerializeTransactionParameters(t *testing.T) {
	require := require.New(t)

	// "do" entrypoint
	// ---------------
	// tezos-client rpc post /chains/main/blocks/head/helpers/forge/operations with '{
	// 	"branch": "BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB",
	// 	"contents":
	// 		[ { "kind": "transaction",
	// 			"source": "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	// 			"fee": "1266", "counter": "1", "gas_limit": "10100",
	// 			"storage_limit": "277",  "amount": "0",
	// 			"destination": "KT1GrStTuhgMMpzbNWKTt7NoXGrYiufrHDYq",
	// 			"parameters": {"entrypoint": "do", "value": {}} } ]
	// }'
	// e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e950200015ab81204ccd229281b9c462edaf0a43e78075f4600ff02000000050200000000
	paramsValueBytes, err := hex.DecodeString("0200000000")
	require.NoError(err)
	paramsValue := tezosprotocol.TransactionParametersValueRawBytes(paramsValueBytes)
	params := tezosprotocol.TransactionParameters{
		Entrypoint: tezosprotocol.EntrypointDo,
		Value:      &paramsValue,
	}
	expectedBytes := "02000000050200000000"
	observedBytes, err := params.MarshalBinary()
	require.NoError(err)
	require.Equal(expectedBytes, hex.EncodeToString(observedBytes))
	reserialized := tezosprotocol.TransactionParameters{}
	require.NoError(reserialized.UnmarshalBinary(observedBytes))
	require.Equal(params, reserialized)
}

func TestSerializeNamedEntrypoint(t *testing.T) {
	require := require.New(t)

	// misc named entrypoint
	// ---------------------
	// tezos-client rpc post /chains/main/blocks/head/helpers/forge/operations with '{
	// 	"branch": "BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB",
	// 	"contents":
	// 		[ { "kind": "transaction",
	// 			"source": "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	// 			"fee": "1266", "counter": "1", "gas_limit": "10100",
	// 			"storage_limit": "277",  "amount": "0",
	// 			"destination": "KT1GrStTuhgMMpzbNWKTt7NoXGrYiufrHDYq",
	// 			"parameters": {"entrypoint": "dummy", "value": {}} } ]
	// }'
	// e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e950200015ab81204ccd229281b9c462edaf0a43e78075f4600ffff0564756d6d79000000050200000000
	paramsValueBytes, err := hex.DecodeString("0200000000")
	require.NoError(err)
	entrypoint, err := tezosprotocol.NewNamedEntrypoint("dummy")
	require.NoError(err)
	paramsValue := tezosprotocol.TransactionParametersValueRawBytes(paramsValueBytes)
	expectedBytes := "ff0564756d6d79000000050200000000"
	params := tezosprotocol.TransactionParameters{
		Entrypoint: entrypoint,
		Value:      &paramsValue,
	}
	observedBytes, err := params.MarshalBinary()
	require.NoError(err)
	require.Equal(expectedBytes, hex.EncodeToString(observedBytes))
	reserialized := tezosprotocol.TransactionParameters{}
	require.NoError(reserialized.UnmarshalBinary(observedBytes))
	require.Equal(params, reserialized)
}

func TestEndpointNameTooLong(t *testing.T) {
	_, err := tezosprotocol.NewNamedEntrypoint(strings.Repeat("a", math.MaxUint8+1))
	require.Error(t, err)
}
