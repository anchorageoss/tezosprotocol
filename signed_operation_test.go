package tezosprotocol_test

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v3"
	"github.com/stretchr/testify/require"
)

// checks the SignOperation function against a known operation, private key, and
// signature. Note that this is possible because Ed25519 signatures are deterministic.
//nolint:dupl
func TestSignOperation(t *testing.T) {
	require := require.New(t)
	operation := &tezosprotocol.Operation{
		Branch: tezosprotocol.BranchID("BMTiv62VhjkVXZJL9Cu5s56qTAJxyciQB2fzA9vd2EiVMsaucWB"),
		Contents: []tezosprotocol.OperationContents{
			&tezosprotocol.Revelation{
				Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
				Fee:          big.NewInt(1257),
				Counter:      big.NewInt(1),
				GasLimit:     big.NewInt(10000),
				StorageLimit: big.NewInt(0),
				PublicKey:    tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"),
			},
			&tezosprotocol.Transaction{
				Source:       tezosprotocol.ContractID("tz1KqTpEZ7Yob7QbPE4Hy4Wo8fHG8LhKxZSx"),
				Fee:          big.NewInt(50000),
				Counter:      big.NewInt(2),
				GasLimit:     big.NewInt(200),
				StorageLimit: big.NewInt(0),
				Amount:       big.NewInt(100000000),
				Destination:  tezosprotocol.ContractID("tz1gjaF81ZRRvdzjobyfVNsAeSC6PScjfQwN"),
			},
		},
	}
	privateKey := tezosprotocol.PrivateKey("edskRwAubEVzMEsaPYnTx3DCttC8zYrGjzPMzTfDr7jfDaihYuh95CFrrYj6kyJoqYhycQPXMZHsZR5mPQRtDgjY6KHJxpeKnZ")
	signedOperation, err := tezosprotocol.SignOperation(operation, privateKey)
	require.NoError(err)
	signedOperationBytes, err := signedOperation.MarshalBinary()
	require.NoError(err)
	signedOperationHex := hex.EncodeToString(signedOperationBytes)
	expected := "e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860302c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c0065667ade71f0c28dcd8c6f443be8b2ff9ebe9f3d2bd8a95d8a29df74319ef24e46bb8abe3e2553dec2a81353f059093861229869ad3c468ade4d9366be3e1308"
	require.Equal(expected, signedOperationHex)

	// rehydrate serialized SignedOperation
	deserialized := tezosprotocol.SignedOperation{}
	require.NoError(deserialized.UnmarshalBinary(signedOperationBytes))
	require.Equal(signedOperation.Operation, deserialized.Operation)
	originalSignatureBytes, err := signedOperation.Signature.MarshalBinary()
	require.NoError(err)
	deserializedSignatureBytes, err := deserialized.Signature.MarshalBinary()
	require.NoError(err)
	require.Equal(originalSignatureBytes, deserializedSignatureBytes)
}

func TestGetSignedOperationHash(t *testing.T) {
	require := require.New(t)
	signedOperationBytes, err := hex.DecodeString("e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860302c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c0065667ade71f0c28dcd8c6f443be8b2ff9ebe9f3d2bd8a95d8a29df74319ef24e46bb8abe3e2553dec2a81353f059093861229869ad3c468ade4d9366be3e1308")
	require.NoError(err)
	signedOperation := tezosprotocol.SignedOperation{}
	require.NoError(signedOperation.UnmarshalBinary(signedOperationBytes))
	operationHash, err := signedOperation.GetHash()
	require.NoError(err)
	require.Equal(tezosprotocol.OperationHash("onvk5LwVA1AXnUEvcz17HE2jt2DLkYbqxkbboX53utEJQ56sThr"), operationHash)
}

func TestMessageSignatureVerification(t *testing.T) {
	require := require.New(t)
	msg := "Hi, my name is Werner Brandes. My voice is my passport. Verify Me."
	cryptoPublicKey, cryptoPrivateKey, err := ed25519.GenerateKey(bytes.NewReader(randSeed))
	require.NoError(err)
	privateKey, err := tezosprotocol.NewPrivateKeyFromCryptoPrivateKey(cryptoPrivateKey)
	require.NoError(err)
	sig, err := tezosprotocol.SignMessage(msg, privateKey)
	require.NoError(err)
	err = tezosprotocol.VerifyMessage(msg, sig, cryptoPublicKey)
	require.NoError(err)
}
