package tezosprotocol_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v3"
	"github.com/stretchr/testify/require"
)

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
