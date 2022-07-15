package tezosprotocol_test

import (
	"testing"

	tezosprotocol "github.com/anchorageoss/tezosprotocol/v3"
	"github.com/stretchr/testify/require"
)

func TestMichelineEncodings(t *testing.T) {
	emptyString := ""
	shortString := "a"
	tests := []struct {
		name    string
		node    tezosprotocol.MichelineNode
		want    []byte
		wantErr bool
	}{
		{
			name: "empty string",
			node: (*tezosprotocol.MichelineString)(&emptyString),
			want: []byte{0x1, 0x0, 0x0, 0x0, 0x0},
		}, {
			name: "short string",
			node: (*tezosprotocol.MichelineString)(&shortString),
			want: []byte{0x1, 0x0, 0x0, 0x0, 0x1, 0x61},
		}, {
			name: "prim0",
			node: &tezosprotocol.MichelinePrim{Prim: tezosprotocol.PrimT_unit},
			want: []byte{0x3, 0x6c},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.node.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("MichelineInt.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
