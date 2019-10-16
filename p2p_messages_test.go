package tezosprotocol_test

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"math/big"
	"testing"

	tezosprotocol "github.com/anchorageoss/tezosprotocol/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
)

type encodeDecodeTestCase struct {
	Input    string
	Expected string
}

func TestEncodeContractID(t *testing.T) {
	require := require.New(t)
	testCases := []encodeDecodeTestCase{{
		Input:    "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
		Expected: "000002298c03ed7d454a101eb7022bc95f7e5f41ac78",
	}, {
		Input:    "tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN",
		Expected: "0000e7670f32038107a59a2b9cfefae36ea21f5aa63c",
	}, {
		Input:    "tz29nEixktH9p9XTFX7p8hATUyeLxXEz96KR",
		Expected: "0001101368afffeb1dc3c089facbbe23f5c30b787ce9",
	}, {
		Input:    "tz3Mo3gHekQhCmykfnC58ecqJLXrjMKzkF2Q",
		Expected: "0002101368afffeb1dc3c089facbbe23f5c30b787ce9",
	}, {
		Input:    "KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82",
		Expected: "01aa3358e4da03d38825f1eb133ca823b676c748e000",
	}}

	for _, testCase := range testCases {
		contractID := tezosprotocol.ContractID(testCase.Input)
		observedBytes, err := contractID.MarshalBinary()
		require.NoError(err)
		observed := hex.EncodeToString(observedBytes)
		require.Equal(testCase.Expected, observed, "mismatch for input %s", testCase.Input)
	}
}

func TestDecodeContractID(t *testing.T) {
	require := require.New(t)
	testCases := []encodeDecodeTestCase{{
		Input:    "000002298c03ed7d454a101eb7022bc95f7e5f41ac78",
		Expected: "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	}, {
		Input:    "0000e7670f32038107a59a2b9cfefae36ea21f5aa63c",
		Expected: "tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN",
	}, {
		Input:    "0001101368afffeb1dc3c089facbbe23f5c30b787ce9",
		Expected: "tz29nEixktH9p9XTFX7p8hATUyeLxXEz96KR",
	}, {
		Input:    "0002101368afffeb1dc3c089facbbe23f5c30b787ce9",
		Expected: "tz3Mo3gHekQhCmykfnC58ecqJLXrjMKzkF2Q",
	}, {
		Input:    "01aa3358e4da03d38825f1eb133ca823b676c748e000",
		Expected: "KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82",
	}, {
		// 21 byte $public_key_hash
		Input:    "02101368afffeb1dc3c089facbbe23f5c30b787ce9",
		Expected: "tz3Mo3gHekQhCmykfnC58ecqJLXrjMKzkF2Q",
	}}

	for _, testCase := range testCases {
		var contractID tezosprotocol.ContractID
		inputBytes, err := hex.DecodeString(testCase.Input)
		require.NoError(err)
		require.NoError(contractID.UnmarshalBinary(inputBytes))
		require.Equal(tezosprotocol.ContractID(testCase.Expected), contractID)
	}
}

func TestEncodeRevelation(t *testing.T) {
	require := require.New(t)
	revelation := &tezosprotocol.Revelation{
		Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
		Fee:          big.NewInt(1257),
		Counter:      big.NewInt(1),
		GasLimit:     big.NewInt(10000),
		StorageLimit: big.NewInt(0),
		PublicKey:    tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"),
	}
	encodedBytes, err := revelation.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f"
	require.Equal(expected, encoded)
}

func TestDecodeRevelation(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f")
	require.NoError(err)
	revelation := tezosprotocol.Revelation{}
	require.NoError(revelation.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"), revelation.Source)
	require.Equal("1257", revelation.Fee.String())
	require.Equal("1", revelation.Counter.String())
	require.Equal("10000", revelation.GasLimit.String())
	require.Equal("0", revelation.StorageLimit.String())
	require.Equal(tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"), revelation.PublicKey)
}

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
	paramsValueBytes, err := hex.DecodeString("000000050200000000")
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

func TestEncodeOrigination(t *testing.T) {
	require := require.New(t)
	primUnit, err := hex.DecodeString("036c") // 03 <prim0> 6c <unit>
	require.NoError(err)
	dummyScript := tezosprotocol.ContractScript{
		Code:    primUnit,
		Storage: primUnit,
	}
	delegate := tezosprotocol.ContractID("tz1ddb9NMYHZi5UzPdzTZMYQQZoMub195zgv")
	origination := &tezosprotocol.Origination{
		Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
		Fee:          big.NewInt(1266),
		Counter:      big.NewInt(1),
		GasLimit:     big.NewInt(10100),
		StorageLimit: big.NewInt(277),
		Balance:      big.NewInt(12000000),
		Delegate:     &delegate,
		Script:       dummyScript,
	}
	encodedBytes, err := origination.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	// source:
	// tezos-client rpc post /chains/main/blocks/head/helpers/forge/operations with '{
	// "branch": "BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB",
	// "contents":
	// 	[ { "kind": "origination",
	// 		"source": "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	// 		"fee": "1266", "counter": "1", "gas_limit": "10100", "delegate": "tz1ddb9NMYHZi5UzPdzTZMYQQZoMub195zgv",
	// 		"storage_limit": "277",  "balance": "12000000", "script": { "code": {"prim": "unit"}, "storage": {"prim": "unit"} } } ]
	// }'
	expected := "6d0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e950280b6dc05ff00c55cf02dbeecc978d9c84625dcae72bb77ea4fbd00000002036c00000002036c"
	require.Equal(expected, encoded)
}

func TestDecodeOrigination(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("6d0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e950280b6dc05ff00c55cf02dbeecc978d9c84625dcae72bb77ea4fbd00000002036c00000002036c")
	require.NoError(err)
	origination := tezosprotocol.Origination{}
	require.NoError(origination.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"), origination.Source)
	require.Equal("1266", origination.Fee.String())
	require.Equal("1", origination.Counter.String())
	require.Equal("10100", origination.GasLimit.String())
	require.Equal("277", origination.StorageLimit.String())
	require.Equal("12000000", origination.Balance.String())
	require.NotNil(origination.Delegate)
	require.Equal(tezosprotocol.ContractID("tz1ddb9NMYHZi5UzPdzTZMYQQZoMub195zgv"), *origination.Delegate)

	// check the script
	primUnit, err := hex.DecodeString("036c") // 03 <prim0> 6c <unit>
	require.NoError(err)
	require.Equal(primUnit, origination.Script.Code)
	require.Equal(primUnit, origination.Script.Storage)
}

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

//nolint:dupl
func TestEncodeOperation(t *testing.T) {
	require := require.New(t)
	operation := &tezosprotocol.Operation{
		Branch: tezosprotocol.BranchID("BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB"),
		Contents: []tezosprotocol.OperationContents{
			&tezosprotocol.Revelation{
				Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
				Fee:          big.NewInt(1257),
				Counter:      big.NewInt(1),
				GasLimit:     big.NewInt(10000),
				StorageLimit: big.NewInt(0),
				PublicKey:    tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"),
			},
			&tezosprotocol.Transaction{
				Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
				Fee:          big.NewInt(50000),
				Counter:      big.NewInt(2),
				GasLimit:     big.NewInt(200),
				StorageLimit: big.NewInt(0),
				Amount:       big.NewInt(100000000),
				Destination:  tezosprotocol.ContractID("tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN"),
			},
		},
	}
	encodedBytes, err := operation.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860302c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c00"
	require.Equal(expected, encoded)
}

func TestDecodeOperation(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860302c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c00")
	require.NoError(err)
	operation := &tezosprotocol.Operation{}
	require.NoError(operation.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.BranchID("BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB"), operation.Branch)
	require.Len(operation.Contents, 2)
	require.IsType(&tezosprotocol.Revelation{}, operation.Contents[0])
	require.IsType(&tezosprotocol.Transaction{}, operation.Contents[1])
}

// checks the SignOperation function against a known operation, private key, and
// signature. Note that this is possible because Ed25519 signatures are deterministic.
//nolint:dupl
func TestSignOperation(t *testing.T) {
	require := require.New(t)
	operation := &tezosprotocol.Operation{
		Branch: tezosprotocol.BranchID("BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB"),
		Contents: []tezosprotocol.OperationContents{
			&tezosprotocol.Revelation{
				Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
				Fee:          big.NewInt(1257),
				Counter:      big.NewInt(1),
				GasLimit:     big.NewInt(10000),
				StorageLimit: big.NewInt(0),
				PublicKey:    tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"),
			},
			&tezosprotocol.Transaction{
				Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
				Fee:          big.NewInt(50000),
				Counter:      big.NewInt(2),
				GasLimit:     big.NewInt(200),
				StorageLimit: big.NewInt(0),
				Amount:       big.NewInt(100000000),
				Destination:  tezosprotocol.ContractID("tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN"),
			},
		},
	}
	privateKey := tezosprotocol.PrivateKey("edskRwAubEVzMEsaPYnTx3DCttC8zYrGjzPMzTfDr7jfDaihYuh95CFrrYj6kyJoqYhycQPXMZHsZR5mPQRtDgjY6KHJxpeKnZ")
	signedOperation, err := tezosprotocol.SignOperation(operation, privateKey)
	require.NoError(err)
	signedOperationBytes, err := signedOperation.MarshalBinary()
	require.NoError(err)
	signedOperationHex := hex.EncodeToString(signedOperationBytes)
	expected := "e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860302c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c0065667ade71f0c28dcd8c6f443be8b2ff9ebe9f3d2bd8a95d8a29df74319ef24e46bb8abe3e2553dec2a81353f059093861229869ad3c468ade4d9366be3e1308"
	require.Equal(expected, signedOperationHex)

	// rehydrate serialized SignedOperation
	deserialized := tezosprotocol.SignedOperation{}
	require.NoError(deserialized.UnmarshalBinary(signedOperationBytes))
	require.Equal(signedOperation.Operation, deserialized.Operation)
	originalSignatureBytes, err := signedOperation.Signature.MarshalBinary()
	require.NoError(err)
	deserializedSignatureBytes, err := deserialized.Signature.MarshalBinary()
	require.NoError(err)
	require.Equal(originalSignatureBytes, deserializedSignatureBytes)
}

func TestGetOperationHash(t *testing.T) {
	require := require.New(t)
	signedOperationBytes, err := hex.DecodeString("e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860302c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c0065667ade71f0c28dcd8c6f443be8b2ff9ebe9f3d2bd8a95d8a29df74319ef24e46bb8abe3e2553dec2a81353f059093861229869ad3c468ade4d9366be3e1308")
	require.NoError(err)
	signedOperation := tezosprotocol.SignedOperation{}
	require.NoError(signedOperation.UnmarshalBinary(signedOperationBytes))
	operationHash, err := signedOperation.GetHash()
	require.NoError(err)
	require.Equal(tezosprotocol.OperationHash("onvk5LwVA1AXnUEvcz17HE2jt2DLkYbqxkbboX53utEJQ56sThr"), operationHash)
}

func TestDeriveOriginatedAddress(t *testing.T) {
	require := require.New(t)
	// reference operation: e805ceeaec0942f1e9fd30f901f102758c027e7da96968c54ed1319608e9674209000002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860304f44e95020002298c03ed7d454a101eb7022bc95f7e5f41ac787bff00000009000002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860305f44e95020002298c03ed7d454a101eb7022bc95f7e5f41ac787bff00000073f8327fc2ed94037230d2c1c88b55001d65371ff0d2f53fc9c60f5ace9024c2839deba6c94e0f68e5a52aee506f5d486a1a2f99e41d8acb2db7349ea9319203
	operationHash := tezosprotocol.OperationHash("onwZr5efqY6eT8r7sUf8WAvDKAPQ2qYkyvqP1UAbSoWWeq45Ut5")
	originatedAddr0, err := tezosprotocol.NewContractIDFromOrigination(operationHash, 0)
	require.NoError(err)
	require.Equal(tezosprotocol.ContractID("KT19ZKrg4XVKV9z5zbYav8SonZrGVmxKuRHB"), originatedAddr0)
	originatedAddr1, err := tezosprotocol.NewContractIDFromOrigination(operationHash, 1)
	require.NoError(err)
	require.Equal(tezosprotocol.ContractID("KT1MXc7s1ZtoVZvbws7vrmz1oLeVGPFoBqpL"), originatedAddr1)
}

func TestNewContractIDFromPublicKey(t *testing.T) {
	require := require.New(t)
	publicKey := tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav")
	expected := tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx")
	observed, err := tezosprotocol.NewContractIDFromPublicKey(publicKey)
	require.NoError(err)
	require.Equal(expected, observed)
}

func TestUnmarshalingIndexOutOfBoundsException(t *testing.T) {
	require := require.New(t)
	emptyBytes := []byte{}
	unmarshalers := []encoding.BinaryUnmarshaler{
		&tezosprotocol.Operation{},
		&tezosprotocol.Revelation{},
		&tezosprotocol.Transaction{},
		&tezosprotocol.Delegation{},
		&tezosprotocol.Origination{},
	}
	for _, unmarshaler := range unmarshalers {
		err := unmarshaler.UnmarshalBinary(emptyBytes)
		require.Error(err, "%T", unmarshaler)
		require.Contains(err.Error(), "out of bounds exception", "%T", unmarshaler)
	}
}

func TestAccountType(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		Input    string
		Expected tezosprotocol.AccountType
	}{{
		Input:    "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
		Expected: tezosprotocol.AccountTypeImplicit,
	}, {
		Input:    "tz29nEixktH9p9XTFX7p8hATUyeLxXEz96KR",
		Expected: tezosprotocol.AccountTypeImplicit,
	}, {
		Input:    "tz3Mo3gHekQhCmykfnC58ecqJLXrjMKzkF2Q",
		Expected: tezosprotocol.AccountTypeImplicit,
	}, {
		Input:    "KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82",
		Expected: tezosprotocol.AccountTypeOriginated,
	}}

	for _, testCase := range testCases {
		contractID := tezosprotocol.ContractID(testCase.Input)
		observedAccountType, err := contractID.AccountType()
		require.NoError(err, contractID)
		require.Equal(testCase.Expected, observedAccountType, "mismatch for input %s", testCase.Input)
	}
}

func TestNewContractIDGeneration(t *testing.T) {
	require := require.New(t)
	cryptoPublicKey, _, err := ed25519.GenerateKey(bytes.NewReader(randSeed))
	require.NoError(err)
	publicKey, err := tezosprotocol.NewPublicKeyFromCryptoPublicKey(cryptoPublicKey)
	require.NoError(err)
	_, err = tezosprotocol.NewContractIDFromPublicKey(publicKey)
	require.NoError(err)
}

func TestMessageSignatureVerification(t *testing.T) {
	require := require.New(t)
	msg := "Hi, my name is Werner Brandes. My voice is my passport. Verify Me."
	cryptoPublicKey, cryptoPrivateKey, err := ed25519.GenerateKey(bytes.NewReader(randSeed))
	require.NoError(err)
	privateKey, err := tezosprotocol.NewPrivateKeyFromCryptoPrivateKey(cryptoPrivateKey)
	require.NoError(err)
	sig, err := tezosprotocol.SignMessage(msg, privateKey)
	require.NoError(err)
	err = tezosprotocol.VerifyMessage(msg, sig, cryptoPublicKey)
	require.NoError(err)
}
