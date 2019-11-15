package tezosprotocol_test

import (
	"encoding"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v2"
	"github.com/stretchr/testify/require"
)

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
