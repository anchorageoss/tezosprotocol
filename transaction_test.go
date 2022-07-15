package tezosprotocol_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v3"
	"github.com/stretchr/testify/require"
)

func TestEncodeTransaction(t *testing.T) {
	require := require.New(t)
	transaction := &tezosprotocol.Transaction{
		Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
		Fee:          big.NewInt(50000),
		Counter:      big.NewInt(1),
		GasLimit:     big.NewInt(200),
		StorageLimit: big.NewInt(0),
		Amount:       big.NewInt(100000000),
		Destination:  tezosprotocol.ContractID("tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN"),
	}
	encodedBytes, err := transaction.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860301c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c00"
	require.Equal(expected, encoded)
}

func TestDecodeTransaction(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860301c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c00")
	require.NoError(err)
	transaction := tezosprotocol.Transaction{}
	require.NoError(transaction.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"), transaction.Source)
	require.Equal("50000", transaction.Fee.String())
	require.Equal("1", transaction.Counter.String())
	require.Equal("200", transaction.GasLimit.String())
	require.Equal("0", transaction.StorageLimit.String())
	require.Equal("100000000", transaction.Amount.String())
	require.Equal(tezosprotocol.ContractID("tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN"), transaction.Destination)
	require.Nil(transaction.Parameters)
}

func TestEncodeTransactionWithParameters(t *testing.T) {
	require := require.New(t)
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
	transaction := &tezosprotocol.Transaction{
		Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
		Fee:          big.NewInt(1266),
		Counter:      big.NewInt(1),
		GasLimit:     big.NewInt(10100),
		StorageLimit: big.NewInt(277),
		Amount:       big.NewInt(0),
		Destination:  tezosprotocol.ContractID("KT1GrStTuhgMMpzbNWKTt7NoXGrYiufrHDYq"),
		Parameters: &tezosprotocol.TransactionParameters{
			Entrypoint: tezosprotocol.EntrypointDo,
			Value:      &paramsValue,
		},
	}
	encodedBytes, err := transaction.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e950200015ab81204ccd229281b9c462edaf0a43e78075f4600ff02000000050200000000"
	require.Equal(expected, encoded)
}

func TestDecodeTransactionWithParameters(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e950200015ab81204ccd229281b9c462edaf0a43e78075f4600ff02000000050200000000")
	require.NoError(err)
	transaction := tezosprotocol.Transaction{}
	require.NoError(transaction.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"), transaction.Source)
	require.Equal("1266", transaction.Fee.String())
	require.Equal("1", transaction.Counter.String())
	require.Equal("10100", transaction.GasLimit.String())
	require.Equal("277", transaction.StorageLimit.String())
	require.Equal("0", transaction.Amount.String())
	require.Equal(tezosprotocol.ContractID("KT1GrStTuhgMMpzbNWKTt7NoXGrYiufrHDYq"), transaction.Destination)
	require.NotNil(transaction.Parameters)
	require.Equal(tezosprotocol.EntrypointDo, transaction.Parameters.Entrypoint)
	expectedParamsValue, err := hex.DecodeString("000000050200000000")
	require.NoError(err)
	observedParamsValue, err := transaction.Parameters.Value.MarshalBinary()
	require.NoError(err)
	require.Equal(expectedParamsValue, observedParamsValue)
}
