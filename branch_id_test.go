package tezosprotocol_test

import (
	"reflect"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v2"
)

func TestBranchID_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		b       tezosprotocol.BranchID
		want    []byte
		wantErr bool
	}{
		{
			name: "Off by one bit",
			//Valid Branch is BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB
			//Let's try an invalid one
			b:       tezosprotocol.BranchID("BNTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB"),
			want:    nil,
			wantErr: true,
		},
		{
			name: "Contracts are not branches",
			//Valid Branch is BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB
			//Let's try an invalid one
			b:       tezosprotocol.BranchID("KT19ZKrg4XVKV9z5zbYav8SonZrGVmxKuRHB"),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("BranchID.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BranchID.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBranchID_UnmarshalBinary(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		b       *tezosprotocol.BranchID
		args    args
		wantErr bool
	}{
		{
			name: "Very short",
			b:    nil,
			args: args{
				data: []byte{1, 52},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("BranchID.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
