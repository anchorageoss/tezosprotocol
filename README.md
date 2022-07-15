# tezosprotocol

[![godoc](https://godoc.org/github.com/google/wire?status.svg)][godoc] [![CircleCI](https://circleci.com/gh/anchorageoss/tezosprotocol.svg?style=svg)](https://circleci.com/gh/anchorageoss/tezosprotocol) [![codecov](https://codecov.io/gh/anchorageoss/tezosprotocol/branch/master/graph/badge.svg)](https://codecov.io/gh/anchorageoss/tezosprotocol)

An implementation of the Tezos peer-to-peer communications protocol.

Supports offline fee calculation and the parsing, signing, and encoding of Tezos operations.

[godoc]: https://godoc.org/github.com/anchorageoss/tezosprotocol

## Installation

```bash
$ go get github.com/anchorageoss/tezosprotocol/v3
```

```go
import "github.com/anchorageoss/tezosprotocol/v3"
```

## Examples

### Operation Parsing

Parsing a Tezos wire-format operation:

```go
signedOperationBytes, _ := hex.DecodeString("e655948a282fcfc31b98abe9b37a82038c4c0e9b8e11f60ea0c7b33e6ecc625f6b0002298c03ed7d454a101eb7022bc95f7e5f41ac78e90901904e00004798d2cc98473d7e250c898885718afd2e4efbcb1a1595ab9730761ed830de0f6c0002298c03ed7d454a101eb7022bc95f7e5f41ac78d0860302c8010080c2d72f0000e7670f32038107a59a2b9cfefae36ea21f5aa63c0065667ade71f0c28dcd8c6f443be8b2ff9ebe9f3d2bd8a95d8a29df74319ef24e46bb8abe3e2553dec2a81353f059093861229869ad3c468ade4d9366be3e1308")
var signedOperation tezosprotocol.SignedOperation
_ = signedOperation.UnmarshalBinary(signedOperationBytes)

fmt.Printf("%v\n", signedOperation)
operationHash, _ := signedOperation.GetHash()
fmt.Println(operationHash)
```

### Operation Signing and Encoding

Creating, signing, and encoding a Tezos wire-format operation:

```go
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
signedOperation, _ := tezosprotocol.SignOperation(operation, privateKey)
signedOperationBytes, _ := signedOperation.MarshalBinary()
fmt.Printf("%x\n", signedOperationBytes)
```

## Protocol Upgrades

All backwards-incompatible changes will be handled via standard semantic versioning conventions. Separate tags for the latest stable release for each Tezos protocol will be maintained in parallel.

A separate go module will be published for each Tezos protocol upgrade to support applications migrating between the two.
```go
import (
	babylon "github.com/anchorageoss/tezosprotocol/v3"
	athens "github.com/anchorageoss/tezosprotocol"
)
```