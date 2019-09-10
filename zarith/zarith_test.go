package zarith_test

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol/zarith"
	"github.com/stretchr/testify/require"
)

type zarithTestCase struct {
	input    string
	expected string
}

func TestEncode(t *testing.T) {
	require := require.New(t)
	testCases := []zarithTestCase{{
		input:    "0",
		expected: "00",
	}, {
		input:    "50000",
		expected: "d08603",
	}, {
		input:    "200",
		expected: "c801",
	}, {
		input:    "100000000",
		expected: "80c2d72f",
	}, {
		input:    "1",
		expected: "01",
	}, {
		input:    "10100",
		expected: "f44e",
	}}

	for _, testCase := range testCases {
		input := new(big.Int)
		_, ok := input.SetString(testCase.input, 10)
		require.True(ok)
		observed, err := zarith.EncodeToHex(input)
		require.NoError(err)
		require.Equal(testCase.expected, observed, "mismatch for input %s", testCase.input)
	}
}

func TestDecode(t *testing.T) {
	require := require.New(t)
	testCases := []zarithTestCase{{
		input:    "00",
		expected: "0",
	}, {
		input:    "d08603",
		expected: "50000",
	}, {
		input:    "c801",
		expected: "200",
	}, {
		input:    "80c2d72f",
		expected: "100000000",
	}, {
		input:    "01",
		expected: "1",
	}, {
		input:    "f44e",
		expected: "10100",
	}}

	for _, testCase := range testCases {
		observedDecimal, err := zarith.DecodeHex(testCase.input)
		require.NoError(err)
		observed := observedDecimal.String()
		require.Equal(testCase.expected, observed, "mismatch for input %s", testCase.input)
	}
}

func TestReadNext(t *testing.T) {
	require := require.New(t)

	inputNoExtraBytes, err := hex.DecodeString("d08603")
	require.NoError(err)
	decoded, bytesRead, err := zarith.ReadNext(inputNoExtraBytes)
	require.NoError(err)
	require.Equal(len(inputNoExtraBytes), bytesRead)
	require.Equal("50000", decoded.String())

	inputExtraBytes := append(inputNoExtraBytes, byte(128))
	decoded, bytesRead, err = zarith.ReadNext(inputExtraBytes)
	require.NoError(err)
	require.Equal(len(inputNoExtraBytes), bytesRead)
	require.Equal("50000", decoded.String())

	inputSingleByte := []byte{5}
	decoded, bytesRead, err = zarith.ReadNext(inputSingleByte)
	require.NoError(err)
	require.Equal(1, bytesRead)
	require.Equal("5", decoded.String())

	// 4 bytes of 11111111 -- never terminates because there is never a leading
	// zero.
	inputNonTerminatingZarithNumber := bytes.Repeat([]byte{255}, 4)
	_, _, err = zarith.ReadNext(inputNonTerminatingZarithNumber)
	require.Error(err)
}

func TestNegative(t *testing.T) {
	require := require.New(t)
	input := big.NewInt(-10)
	_, err := zarith.Encode(input)
	require.Error(err)
}
