package tezosprotocol_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v3"
	"github.com/stretchr/testify/require"
)

func TestEncodeOrigination(t *testing.T) {
	require := require.New(t)
	micheline := tezosprotocol.MichelinePrim{Prim: tezosprotocol.PrimT_unit}
	michelineBytes, err := micheline.MarshalBinary()
	require.NoError(err)
	dummyScript := tezosprotocol.ContractScript{
		Code:    michelineBytes,
		Storage: michelineBytes,
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
