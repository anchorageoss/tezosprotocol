package tezosprotocol

import "math/big"

// ComputeMinimumFee returns the minimum fee required according to the constraint:
//   fees >= (minimal_fees + minimal_nanotez_per_byte * size + minimal_nanotez_per_gas_unit * gas)
// Amount returned is in units of mutez.
// Reference: http://tezos.gitlab.io/mainnet/protocols/003_PsddFKi3.html#baker
func ComputeMinimumFee(gasLimit, operationSizeBytes *big.Int) *big.Int {
	storageFee := new(big.Int).Mul(operationSizeBytes, big.NewInt(DefaultMinimalNanotezPerByte))
	storageFee = new(big.Int).Div(storageFee, big.NewInt(1000))

	gasFee := new(big.Int).Mul(gasLimit, big.NewInt(DefaultMinimalNanotezPerGasUnit))
	gasFee = new(big.Int).Div(gasFee, big.NewInt(1000))

	totalFee := new(big.Int).Add(storageFee, gasFee)
	totalFee = new(big.Int).Add(totalFee, big.NewInt(DefaultMinimalFees))

	return totalFee
}

// Common values for fees
const (
	// StorageCostPerByte is the amount of mutez burned per byte of storage used.
	// Reference: https://gitlab.com/tezos/tezos/blob/f5c50c8ba1670b7a2ee58bed8a7806f00c43340c/src/proto_alpha/lib_protocol/constants_repr.ml#L126
	StorageCostPerByte = int64(1000)

	// NewAccountStorageLimitBytes is the storage needed to create a new
	// account, either implicit or originated.
	NewAccountStorageLimitBytes = int64(257)

	// NewAccountCreationBurn is the cost in mutez burned from an account that signs
	// an operation creating a new account, either by a transferring to a new implicit address
	// or by originating a KT1 address. The value is equal to ꜩ0.257
	NewAccountCreationBurn = NewAccountStorageLimitBytes * StorageCostPerByte

	// DefaultMinimalFees is a flat fee that represents the cost of broadcasting
	// an operation to the network. This flat fee is added to the variable minimal
	// fees for gas spent and storage used.
	// Reference: https://gitlab.com/tezos/tezos/blob/f5c50c8ba1670b7a2ee58bed8a7806f00c43340c/src/proto_alpha/lib_client/client_proto_args.ml#L251
	DefaultMinimalFees = int64(100)

	// DefaultMinimalMutezPerGasUnit is the default fee rate in mutez that nodes expect
	// per unit gas spent by an operation (and all its contents).
	// Reference: https://gitlab.com/tezos/tezos/blob/f5c50c8ba1670b7a2ee58bed8a7806f00c43340c/src/proto_alpha/lib_client/client_proto_args.ml#L252
	DefaultMinimalNanotezPerGasUnit = int64(100)

	// DefaultMinimalMutezPerByte is the default fee rate in mutez that nodes expect per
	// byte of a serialized, signed operation -- including header and all contents.
	// Reference: https://gitlab.com/tezos/tezos/blob/f5c50c8ba1670b7a2ee58bed8a7806f00c43340c/src/proto_alpha/lib_client/client_proto_args.ml#L253
	DefaultMinimalNanotezPerByte = int64(1000)

	// OriginationGasLimit is the gas consumed by a simple origination.
	// reference: http://tezos.gitlab.io/mainnet/protocols/003_PsddFKi3.html#more-details-on-fees-and-cost-model
	OriginationGasLimit = int64(10000)

	// MinimumOriginationSizeBytes is the smallest size in bytes of a serialized,
	// signed origination operation
	MinimumOriginationSizeBytes = int64(152)

	// OriginationMinimumFee is the minimum amount to be paid to a baker for an
	// operation with one origination
	OriginationMinimumFee = DefaultMinimalFees +
		DefaultMinimalNanotezPerByte*MinimumOriginationSizeBytes/int64(1000) +
		DefaultMinimalNanotezPerGasUnit*OriginationGasLimit/int64(1000)

	// OriginationStorageLimitBytes is the storage limit required for originations
	OriginationStorageLimitBytes = NewAccountStorageLimitBytes

	// OriginationStorageBurn is the amount of mutez burned by an account as a consequence
	// of signing an origination.
	OriginationStorageBurn = OriginationStorageLimitBytes * StorageCostPerByte

	// reference: http://tezos.gitlab.io/mainnet/protocols/003_PsddFKi3.html#more-details-on-fees-and-cost-model
	MinimumOriginatedAccountTransferGasLimit  = int64(10100)
	MinimumOriginatedAccountTransferSizeBytes = int64(215)

	// OriginatedAccountTransferMinimumFee is the minimum amount to be paid to a baker
	// for a transfer from an originated account
	OriginatedAccountTransferMinimumFee = DefaultMinimalFees +
		DefaultMinimalNanotezPerByte*MinimumOriginatedAccountTransferSizeBytes/int64(1000) +
		DefaultMinimalNanotezPerGasUnit*MinimumOriginatedAccountTransferGasLimit/int64(1000)

	// RevelationGasLimit is the gas consumed by a revelation
	RevelationGasLimit = int64(10000)

	// RevelationStorageLimitBytes is the storage limit required for revelations. Note that
	// it is zero.
	RevelationStorageLimitBytes = int64(0)

	// RevelationStorageBurn is the amount burned by an account as a consequence
	// of signing a revelation. Note that it is zero.
	RevelationStorageBurn = RevelationStorageLimitBytes * StorageCostPerByte

	// MinimumTransactionGasLimit is the gas consumed by a transaction with no parameters
	// that does not result in any Michelson code execution.
	MinimumTransactionGasLimit = int64(10200)

	// DelegationGasLimit is the gas consumed by a delegation
	DelegationGasLimit = int64(10000)

	// DelegationStorageLimitBytes is the storage limit required for delegations. Note that
	// it is zero.
	DelegationStorageLimitBytes = int64(0)

	// DelegationStorageBurn is the amount burned by an account as a consequence
	// of signing a delegation. Note that it is zero.
	DelegationStorageBurn = DelegationStorageLimitBytes * StorageCostPerByte
)
