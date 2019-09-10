package tezosprotocol_test

import (
	"encoding/hex"
	"testing"

	"github.com/anchorageoss/tezosprotocol"
	"github.com/stretchr/testify/require"
)

type b58CheckTestCase struct {
	PayloadHex         string
	Prefix             tezosprotocol.Base58CheckPrefix
	Base58CheckEncoded string
	String             string
}

var testCases = []b58CheckTestCase{{
	PayloadHex:         "02298c03ed7d454a101eb7022bc95f7e5f41ac78",
	Prefix:             tezosprotocol.PrefixEd25519PublicKeyHash,
	Base58CheckEncoded: "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	String:             "tz1(36)",
}, {
	PayloadHex:         "101368afffeb1dc3c089facbbe23f5c30b787ce9",
	Prefix:             tezosprotocol.PrefixSecp256k1PublicKeyHash,
	Base58CheckEncoded: "tz29nEixktH9p9XTFX7p8hATUyeLxXEz96KR",
	String:             "tz2(36)",
}, {
	PayloadHex:         "101368afffeb1dc3c089facbbe23f5c30b787ce9",
	Prefix:             tezosprotocol.PrefixP256PublicKeyHash,
	Base58CheckEncoded: "tz3Mo3gHekQhCmykfnC58ecqJLXrjMKzkF2Q",
	String:             "tz3(36)",
}, {
	PayloadHex:         "866052fe3b476d7abf950530b7db7fefc58a279e6a7fe02b59911445054a45ab",
	Prefix:             tezosprotocol.PrefixBlockHash,
	Base58CheckEncoded: "BLjToBJ9Y8CHdzCbdfZ8famZ6Yk9c3yG1uP7p99dkYgPhozrvvj",
	String:             "B(51)",
}, {
	PayloadHex:         "f2342b8bc076c65f83a286152634e9c172ad08de",
	Prefix:             tezosprotocol.PrefixContractHash,
	Base58CheckEncoded: "KT1WfRb2j1YPot5PR1CRPKowiteVmKGaA5NA",
	String:             "KT1(36)",
}, {
	PayloadHex:         "6a5c3d425cfb5c4e2f8a4033098acdb732868950a73777316dcd499d5304b4391bc367618ad8005290f866a9776a1ad564b1eea429a9a3080d2297d4e4b28a0e",
	Prefix:             tezosprotocol.PrefixEd25519Signature,
	Base58CheckEncoded: "edsigtmiq6NN7djPAXTQbyztgaLgbojoCdr2hUkZU2qsevHSL8vq7ZfQYC7cvPRb6sudzjKzy4DDJb1f4aFFpL7KNidaMaztevk",
	String:             "edsig(99)",
}, {
	PayloadHex:         "7a06a770",
	Prefix:             tezosprotocol.PrefixChainID,
	Base58CheckEncoded: "NetXdQprcVkpaWU",
	String:             "Net(15)",
}}

func TestBase58CheckEncode(t *testing.T) {
	require := require.New(t)
	for _, testCase := range testCases {
		payloadBytes, err := hex.DecodeString(testCase.PayloadHex)
		require.NoError(err)
		observed, err := tezosprotocol.Base58CheckEncode(testCase.Prefix, payloadBytes)
		require.NoError(err)
		require.Equal(testCase.Base58CheckEncoded, observed)
		require.Equal(testCase.String, testCase.Prefix.String())
	}
}

func TestBase58CheckDecode(t *testing.T) {
	require := require.New(t)
	for _, testCase := range testCases {
		prefix, payloadBytes, err := tezosprotocol.Base58CheckDecode(testCase.Base58CheckEncoded)
		require.NoError(err)
		payloadHex := hex.EncodeToString(payloadBytes)
		require.Equal(testCase.Prefix, prefix)
		require.Equal(testCase.PayloadHex, payloadHex)
	}
}

func TestBase58CheckDecodeNegativeCases(t *testing.T) {
	require := require.New(t)

	// empty string
	_, _, err := tezosprotocol.Base58CheckDecode("")
	require.Error(err)

	// bad checksum
	_, _, err = tezosprotocol.Base58CheckDecode("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSR")
	require.Error(err)
	require.Contains(err.Error(), "checksum")

	// unknown prefix
	_, _, err = tezosprotocol.Base58CheckDecode("zz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LoDpVc2")
	require.Error(err)
	require.Contains(err.Error(), "prefix")

	// incorrect length
	_, _, err = tezosprotocol.Base58CheckDecode("8Fy8oBr77jCfuUas")
	require.Error(err)
	require.Contains(err.Error(), "unexpected length")
}
