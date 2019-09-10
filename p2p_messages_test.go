package tezosprotocol_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol"
	"github.com/stretchr/testify/require"
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
		Source:       tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"),
		Fee:          big.NewInt(1257),
		Counter:      big.NewInt(1),
		GasLimit:     big.NewInt(10000),
		StorageLimit: big.NewInt(0),
		PublicKey:    tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"),
	}
	encodedBytes, err := revelation.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "0701aa3358e4da03d38825f1eb133ca823b676c748e000e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f"
	require.Equal(expected, encoded)
}

func TestDecodeRevelation(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("0701aa3358e4da03d38825f1eb133ca823b676c748e000e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f")
	require.NoError(err)
	revelation := tezosprotocol.Revelation{}
	require.NoError(revelation.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"), revelation.Source)
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
		Counter:      big.NewInt(0),
		GasLimit:     big.NewInt(200),
		StorageLimit: big.NewInt(0),
		Amount:       big.NewInt(100000000),
		Destination:  tezosprotocol.ContractID("tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN"),
	}
	encodedBytes, err := transaction.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "08000002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860300c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c00"
	require.Equal(expected, encoded)
}

func TestDecodeTransaction(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("08000002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860300c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c00")
	require.NoError(err)
	transaction := tezosprotocol.Transaction{}
	require.NoError(transaction.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"), transaction.Source)
	require.Equal("50000", transaction.Fee.String())
	require.Equal("0", transaction.Counter.String())
	require.Equal("200", transaction.GasLimit.String())
	require.Equal("0", transaction.StorageLimit.String())
	require.Equal("100000000", transaction.Amount.String())
	require.Equal(tezosprotocol.ContractID("tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN"), transaction.Destination)
}

func TestEncodeOrigination(t *testing.T) {
	require := require.New(t)
	origination := &tezosprotocol.Origination{
		Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
		Fee:          big.NewInt(1266),
		Counter:      big.NewInt(1),
		GasLimit:     big.NewInt(10100),
		StorageLimit: big.NewInt(277),
		Manager:      tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
		Balance:      big.NewInt(12000000),
		Spendable:    true,
		Delegatable:  true,
		Delegate:     nil,
	}
	encodedBytes, err := origination.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "09000002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e95020002298c03ed7d454a101eb7022bc95f7e5f41ac7880b6dc05ffff0000"
	require.Equal(expected, encoded)
}

func TestDecodeOrigination(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("09000002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e95020002298c03ed7d454a101eb7022bc95f7e5f41ac7880b6dc05ffff0000")
	require.NoError(err)
	origination := tezosprotocol.Origination{}
	require.NoError(origination.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"), origination.Source)
	require.Equal("1266", origination.Fee.String())
	require.Equal("1", origination.Counter.String())
	require.Equal("10100", origination.GasLimit.String())
	require.Equal("277", origination.StorageLimit.String())
	require.Equal(tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"), origination.Manager)
	require.Equal("12000000", origination.Balance.String())
	require.True(origination.Spendable)
	require.True(origination.Delegatable)
	require.Nil(origination.Delegate)
}

func TestEncodeDelegation(t *testing.T) {
	require := require.New(t)
	delegate := tezosprotocol.ContractID("tz1ddb9NMYHZi5UzPdzTZMYQQZoMub195zgv")
	origination := &tezosprotocol.Delegation{
		Source:       tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"),
		Fee:          big.NewInt(1161),
		Counter:      big.NewInt(2),
		GasLimit:     big.NewInt(10100),
		StorageLimit: big.NewInt(0),
		Delegate:     &delegate,
	}
	encodedBytes, err := origination.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "0a01aa3358e4da03d38825f1eb133ca823b676c748e000890902f44e00ff00c55cf02dbeecc978d9c84625dcae72bb77ea4fbd"
	require.Equal(expected, encoded)
}

func TestDecodeDelegation(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("0a01aa3358e4da03d38825f1eb133ca823b676c748e000890902f44e00ff00c55cf02dbeecc978d9c84625dcae72bb77ea4fbd")
	require.NoError(err)
	delegation := tezosprotocol.Delegation{}
	require.NoError(delegation.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"), delegation.Source)
	require.Equal("1161", delegation.Fee.String())
	require.Equal("2", delegation.Counter.String())
	require.Equal("10100", delegation.GasLimit.String())
	require.Equal("0", delegation.StorageLimit.String())
	require.NotNil(delegation.Delegate)
	require.Equal(tezosprotocol.ContractID("tz1ddb9NMYHZi5UzPdzTZMYQQZoMub195zgv"), *delegation.Delegate)
}

//nolint:dupl
func TestEncodeOperation(t *testing.T) {
	require := require.New(t)
	operation := &tezosprotocol.Operation{
		Branch: tezosprotocol.BranchID("BLs171sHn4FoYxrKCdQxs5seHDZL7e1KRfTwh6ZWejrgZtJwPrL"),
		Contents: []tezosprotocol.OperationContents{
			&tezosprotocol.Revelation{
				Source:       tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"),
				Fee:          big.NewInt(1257),
				Counter:      big.NewInt(1),
				GasLimit:     big.NewInt(10000),
				StorageLimit: big.NewInt(0),
				PublicKey:    tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"),
			},
			&tezosprotocol.Transaction{
				Source:       tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"),
				Fee:          big.NewInt(1178),
				Counter:      big.NewInt(2),
				GasLimit:     big.NewInt(10200),
				StorageLimit: big.NewInt(0),
				Amount:       big.NewInt(9870000),
				Destination:  tezosprotocol.ContractID("tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN"),
			},
		},
	}
	encodedBytes, err := operation.MarshalBinary()
	require.NoError(err)
	encoded := hex.EncodeToString(encodedBytes)
	expected := "977f0b9ea521e630bb9f03a02b99fd76c4554bfb39be02d79bc1502e779817cd0701aa3358e4da03d38825f1eb133ca823b676c748e000e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f0801aa3358e4da03d38825f1eb133ca823b676c748e0009a0902d84f00b0b5da040000e7670f32038107a59a2b9cfefae36ea21f5aa63c00"
	require.Equal(expected, encoded)
}

func TestDecodeOperation(t *testing.T) {
	require := require.New(t)
	encoded, err := hex.DecodeString("977f0b9ea521e630bb9f03a02b99fd76c4554bfb39be02d79bc1502e779817cd0701aa3358e4da03d38825f1eb133ca823b676c748e000e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f0801aa3358e4da03d38825f1eb133ca823b676c748e0009a0902d84f00b0b5da040000e7670f32038107a59a2b9cfefae36ea21f5aa63c00")
	require.NoError(err)
	operation := &tezosprotocol.Operation{}
	require.NoError(operation.UnmarshalBinary(encoded))
	require.Equal(tezosprotocol.BranchID("BLs171sHn4FoYxrKCdQxs5seHDZL7e1KRfTwh6ZWejrgZtJwPrL"), operation.Branch)
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
		Branch: tezosprotocol.BranchID("BLs171sHn4FoYxrKCdQxs5seHDZL7e1KRfTwh6ZWejrgZtJwPrL"),
		Contents: []tezosprotocol.OperationContents{
			&tezosprotocol.Revelation{
				Source:       tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"),
				Fee:          big.NewInt(1257),
				Counter:      big.NewInt(1),
				GasLimit:     big.NewInt(10000),
				StorageLimit: big.NewInt(0),
				PublicKey:    tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"),
			},
			&tezosprotocol.Transaction{
				Source:       tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"),
				Fee:          big.NewInt(1178),
				Counter:      big.NewInt(2),
				GasLimit:     big.NewInt(10200),
				StorageLimit: big.NewInt(0),
				Amount:       big.NewInt(9870000),
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
	expected := "977f0b9ea521e630bb9f03a02b99fd76c4554bfb39be02d79bc1502e779817cd0701aa3358e4da03d38825f1eb133ca823b676c748e000e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f0801aa3358e4da03d38825f1eb133ca823b676c748e0009a0902d84f00b0b5da040000e7670f32038107a59a2b9cfefae36ea21f5aa63c00e2538dc042f4919f27ee6d58a8e4155699e03d421fb920b5ff79c81d7547bff9215040bbf783dbe66550cfabcaeb22262786af7506415054f06c8cf18daf9807"
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
	signedOperationBytes, err := hex.DecodeString("977f0b9ea521e630bb9f03a02b99fd76c4554bfb39be02d79bc1502e779817cd09000002298c03ed7d454a101eb7022bc95f7e5f41ac78f10902f44e95020002298c03ed7d454a101eb7022bc95f7e5f41ac7880897aff00000062cd37a350627ddfa683f1df24c8a35f4a9d1ae4288059b0d80e629c003d12ce73ea0543942da6b1aa3e55e7107876e9bc83a41ff7ea948cd618b4425d382808")
	require.NoError(err)
	signedOperation := tezosprotocol.SignedOperation{}
	require.NoError(signedOperation.UnmarshalBinary(signedOperationBytes))
	operationHash, err := signedOperation.GetHash()
	require.NoError(err)
	require.Equal(tezosprotocol.OperationHash("op4jc9sEhzV1MHCD8np6ZcF6qAYQtEJtveNh9LvDtpET5zdDmx8"), operationHash)
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
