package tezosprotocol_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v3"
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

func TestNewContractIDGeneration(t *testing.T) {
	require := require.New(t)
	cryptoPublicKey, _, err := ed25519.GenerateKey(bytes.NewReader(randSeed))
	require.NoError(err)
	publicKey, err := tezosprotocol.NewPublicKeyFromCryptoPublicKey(cryptoPublicKey)
	require.NoError(err)
	_, err = tezosprotocol.NewContractIDFromPublicKey(publicKey)
	require.NoError(err)
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
