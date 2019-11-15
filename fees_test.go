package tezosprotocol

import (
	"math/big"
	"reflect"
	"testing"
)

func TestComputeMinimumFee(t *testing.T) {
	type args struct {
		gasLimit           *big.Int
		operationSizeBytes *big.Int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "Default",
			args: args{
				gasLimit:           big.NewInt(1),
				operationSizeBytes: big.NewInt(1173),
			},
			want: big.NewInt(1273),
		},
	}
	for _, tt := range tests {
		//Addresses lint issues: using the variable on range scope `tt` in function literal
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeMinimumFee(tt.args.gasLimit, tt.args.operationSizeBytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeMinimumFee() = %v, want %v", got, tt.want)
			}
		})
	}
}
