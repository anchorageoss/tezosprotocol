package tezosprotocol_test

import (
	"encoding/hex"
	"math"
	"strings"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v2"
	"github.com/stretchr/testify/require"
)

func TestContractScriptUnmarshalBinary(t *testing.T) {
	require := require.New(t)

	// invalid code length
	err := (&tezosprotocol.ContractScript{}).UnmarshalBinary([]byte{})
	require.Error(err)
	require.Contains(err.Error(), "failed to read code length")

	// invalid code
	badCode, err := hex.DecodeString("00000002")
	require.NoError(err)
	err = (&tezosprotocol.ContractScript{}).UnmarshalBinary(badCode)
	require.Error(err)
	require.Contains(err.Error(), "failed to read code")

	// invalid storage length
	badStorageLength, err := hex.DecodeString("00000002C0DE00")
	require.NoError(err)
	err = (&tezosprotocol.ContractScript{}).UnmarshalBinary(badStorageLength)
	require.Error(err)
	require.Contains(err.Error(), "failed to read storage length")

	// invalid storage
	badStorage, err := hex.DecodeString("00000002C0DE00000007")
	require.NoError(err)
	err = (&tezosprotocol.ContractScript{}).UnmarshalBinary(badStorage)
	require.Error(err)
	require.Contains(err.Error(), "failed to read storage")
}

func TestSerializeTransactionParameters(t *testing.T) {
	require := require.New(t)

	// "do" entrypoint
	// ---------------
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
	params := tezosprotocol.TransactionParameters{
		Entrypoint: tezosprotocol.EntrypointDo,
		Value:      &paramsValue,
	}
	expectedBytes := "02000000050200000000"
	observedBytes, err := params.MarshalBinary()
	require.NoError(err)
	require.Equal(expectedBytes, hex.EncodeToString(observedBytes))
	reserialized := tezosprotocol.TransactionParameters{}
	require.NoError(reserialized.UnmarshalBinary(observedBytes))
	require.Equal(params, reserialized)
}

func TestSerializeNamedEntrypoint(t *testing.T) {
	require := require.New(t)

	// misc named entrypoint
	// ---------------------
	// tezos-client rpc post /chains/main/blocks/head/helpers/forge/operations with '{
	// 	"branch": "BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB",
	// 	"contents":
	// 		[ { "kind": "transaction",
	// 			"source": "tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx",
	// 			"fee": "1266", "counter": "1", "gas_limit": "10100",
	// 			"storage_limit": "277",  "amount": "0",
	// 			"destination": "KT1GrStTuhgMMpzbNWKTt7NoXGrYiufrHDYq",
	// 			"parameters": {"entrypoint": "dummy", "value": {}} } ]
	// }'
	// e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78f20901f44e950200015ab81204ccd229281b9c462edaf0a43e78075f4600ffff0564756d6d79000000050200000000
	paramsValueBytes, err := hex.DecodeString("0200000000")
	require.NoError(err)
	entrypoint, err := tezosprotocol.NewNamedEntrypoint("dummy")
	require.NoError(err)
	paramsValue := tezosprotocol.TransactionParametersValueRawBytes(paramsValueBytes)
	expectedBytes := "ff0564756d6d79000000050200000000"
	params := tezosprotocol.TransactionParameters{
		Entrypoint: entrypoint,
		Value:      &paramsValue,
	}
	observedBytes, err := params.MarshalBinary()
	require.NoError(err)
	require.Equal(expectedBytes, hex.EncodeToString(observedBytes))
	reserialized := tezosprotocol.TransactionParameters{}
	require.NoError(reserialized.UnmarshalBinary(observedBytes))
	require.Equal(params, reserialized)
}

func TestEndpointNameTooLong(t *testing.T) {
	_, err := tezosprotocol.NewNamedEntrypoint(strings.Repeat("a", math.MaxUint8+1))
	require.Error(t, err)
}

func TestEntrypoint_Name(t *testing.T) {
	type fields struct {
		tag  tezosprotocol.EntrypointTag
		name string
	}
	tests := []struct {
		name    string
		bytes   []byte
		want    string
		wantErr bool
	}{
		{
			name:    "default",
			bytes:   []byte{byte(tezosprotocol.EntrypointTagDefault)},
			want:    "default",
			wantErr: false,
		},
		{
			name:    "root",
			bytes:   []byte{byte(tezosprotocol.EntrypointTagRoot)},
			want:    "root",
			wantErr: false,
		},
		{
			name:    "do",
			bytes:   []byte{byte(tezosprotocol.EntrypointTagDo)},
			want:    "do",
			wantErr: false,
		},
		{
			name:    "set_delegate",
			bytes:   []byte{byte(tezosprotocol.EntrypointTagSetDelegate)},
			want:    "set_delegate",
			wantErr: false,
		},
		{
			name:    "remove_delegate",
			bytes:   []byte{byte(tezosprotocol.EntrypointTagRemoveDelegate)},
			want:    "remove_delegate",
			wantErr: false,
		},
		{
			name:    "named",
			bytes:   append([]byte{byte(tezosprotocol.EntrypointTagNamed), 4}, []byte("tada")...),
			want:    "tada",
			wantErr: false,
		},
		{
			name:    "empty named",
			bytes:   []byte{byte(tezosprotocol.EntrypointTagNamed), 0},
			want:    "",
			wantErr: true,
		},
		{
			name:    "potato",
			bytes:   []byte{byte(42)},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e tezosprotocol.Entrypoint
			err := e.UnmarshalBinary(tt.bytes)
			if err != nil {
				t.Errorf("UnmarshalBinary(%v) error = %v", tt.bytes, err)
				return
			}
			got, err := e.Name()
			if (err != nil) != tt.wantErr {
				t.Errorf("Entrypoint.Name() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Entrypoint.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntrypoint_Tag(t *testing.T) {
	type fields struct {
		tag  tezosprotocol.EntrypointTag
		name string
	}
	tests := []struct {
		name  string
		bytes []byte
		want  tezosprotocol.EntrypointTag
	}{
		{
			name:  "default",
			bytes: []byte{byte(tezosprotocol.EntrypointTagDefault)},
			want:  tezosprotocol.EntrypointTagDefault,
		},
		{
			name:  "root",
			bytes: []byte{byte(tezosprotocol.EntrypointTagRoot)},
			want:  tezosprotocol.EntrypointTagRoot,
		},
		{
			name:  "do",
			bytes: []byte{byte(tezosprotocol.EntrypointTagDo)},
			want:  tezosprotocol.EntrypointTagDo,
		},
		{
			name:  "set_delegate",
			bytes: []byte{byte(tezosprotocol.EntrypointTagSetDelegate)},
			want:  tezosprotocol.EntrypointTagSetDelegate,
		},
		{
			name:  "remove_delegate",
			bytes: []byte{byte(tezosprotocol.EntrypointTagRemoveDelegate)},
			want:  tezosprotocol.EntrypointTagRemoveDelegate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e tezosprotocol.Entrypoint
			err := e.UnmarshalBinary(tt.bytes)
			if err != nil {
				t.Errorf("UnmarshalBinary(%v) error = %v", tt.bytes, err)
				return
			}
			if got := e.Tag(); got != tt.want {
				t.Errorf("Entrypoint.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntrypoint_String(t *testing.T) {
	type fields struct {
		tag  tezosprotocol.EntrypointTag
		name string
	}
	tests := []struct {
		name  string
		bytes []byte
		want  string
	}{
		{
			name:  "default",
			bytes: []byte{byte(tezosprotocol.EntrypointTagDefault)},
			want:  "default",
		},
		{
			name:  "root",
			bytes: []byte{byte(tezosprotocol.EntrypointTagRoot)},
			want:  "root",
		},
		{
			name:  "do",
			bytes: []byte{byte(tezosprotocol.EntrypointTagDo)},
			want:  "do",
		},
		{
			name:  "set_delegate",
			bytes: []byte{byte(tezosprotocol.EntrypointTagSetDelegate)},
			want:  "set_delegate",
		},
		{
			name:  "remove_delegate",
			bytes: []byte{byte(tezosprotocol.EntrypointTagRemoveDelegate)},
			want:  "remove_delegate",
		},
		{
			name:  "named",
			bytes: append([]byte{byte(tezosprotocol.EntrypointTagNamed), 4}, []byte("tada")...),
			want:  "tada",
		},
		{
			name:  "empty named",
			bytes: []byte{byte(tezosprotocol.EntrypointTagNamed), 0},
			want:  "<invalid entrypoint>",
		},
		{
			name:  "potato",
			bytes: []byte{byte(42)},
			want:  "<invalid entrypoint>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e tezosprotocol.Entrypoint
			err := e.UnmarshalBinary(tt.bytes)
			if err != nil {
				t.Errorf("UnmarshalBinary(%v) error = %v", tt.bytes, err)
				return
			}
			if tt.want == "<invalid entrypoint>" {
				if got := e.String(); got != tt.want {
					t.Errorf("Entrypoint.String() = %v, want %v", got, tt.want)
				}
			} else {
				if got := e.String(); got != "%"+tt.want {
					t.Errorf("Entrypoint.String() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
