# tezosprotocol

[![godoc](https://godoc.org/github.com/google/wire?status.svg)][godoc] [![CircleCI](https://circleci.com/gh/anchorageoss/tezosprotocol.svg?style=svg)](https://circleci.com/gh/anchorageoss/tezosprotocol) [![codecov](https://codecov.io/gh/anchorageoss/tezosprotocol/branch/master/graph/badge.svg)](https://codecov.io/gh/anchorageoss/tezosprotocol)

An implementation of the Tezos peer-to-peer communications protocol.

Supports offline fee calculation and the parsing, signing, and encoding of Tezos operations.

[godoc]: https://godoc.org/github.com/anchorageoss/tezosprotocol

## Installation
```bash
$ go get github.com/anchorageoss/tezosprotocol
```

## Examples

### Operation Parsing
Parsing a Tezos wire-format operation:
```go
signedOperationBytes, _ := hex.DecodeString("977f0b9ea521e630bb9f03a02b99fd76c4554bfb39be02d79bc1502e779817cd09000002298c03ed7d454a101eb7022bc95f7e5f41ac78f10902f44e95020002298c03ed7d454a101eb7022bc95f7e5f41ac7880897aff00000062cd37a350627ddfa683f1df24c8a35f4a9d1ae4288059b0d80e629c003d12ce73ea0543942da6b1aa3e55e7107876e9bc83a41ff7ea948cd618b4425d382808")
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
	Branch: tezosprotocol.BranchID("BLs171sHn4FoYxrKCdQxs5seHDZL7e1KRfTwh6ZWejrgZtJwPrL"),
	Contents: []tezosprotocol.OperationContents{
		&tezosprotocol.Revelation{
			Source:       tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"),
			Fee:          big.NewInt(1257),
			Counter:      big.NewInt(1),
			GasLimit:     big.NewInt(10000),
			StorageLimit: big.NewInt(0),
			PublicKey:    tezosprotocol.PublicKey("edpkuBknW28nW72KG6RoHtYW7p12T6GKc7nAbwYX5m8Wd9sDVC9yav"),
		},
		&tezosprotocol.Transaction{
			Source:       tezosprotocol.ContractID("KT1Q6hx3bJayhQYfMDL1z2ugd7GXGckVAV82"),
			Fee:          big.NewInt(1178),
			Counter:      big.NewInt(2),
			GasLimit:     big.NewInt(10200),
			StorageLimit: big.NewInt(0),
			Amount:       big.NewInt(9870000),
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
