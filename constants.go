package tezosprotocol

// Field lengths
const (
	// PubKeyHashLen is the length in bytes of a serialized public key hash
	PubKeyHashLen = 20
	// TaggedPubKeyHashLen is the length in bytes of a serialized, tagged public key hash
	TaggedPubKeyHashLen = PubKeyHashLen + 1
	// PubKeyLenEd25519 is the length in bytes of a serialized Ed25519 public key
	PubKeyLenEd25519 = 32
	// PubKeyLenSecp256k1 is the length in bytes of a serialized Secp256k1 public key
	PubKeyLenSecp256k1 = 33
	// PubKeyLenP256 is the length in bytes of a serialized P256 public key
	PubKeyLenP256 = 33
	// ContractHashLen is the length in bytes of a serialized contract hash
	ContractHashLen = 20
	// ContractIDLen is the length in bytes of a serialized contract ID
	ContractIDLen = 22
	// BlockHashLen is the length in bytes of a serialized block hash
	BlockHashLen = 32
	// OperationHashLen is the length in bytes of a serialized operation hash
	OperationHashLen = 32
	// OperationSignatureLen is the length in bytes of a serialized operation signature
	OperationSignatureLen = 64
)

// ContentsTag captures the possible tag values for operation contents
type ContentsTag byte

const (
	// ContentsTagRevelation is the tag for revelations
	ContentsTagRevelation ContentsTag = 107
	// ContentsTagTransaction is the tag for transactions
	ContentsTagTransaction ContentsTag = 108
	// ContentsTagOrigination is the tag for originations
	ContentsTagOrigination ContentsTag = 109
	// ContentsTagDelegation is the tag for delegations
	ContentsTagDelegation ContentsTag = 110
)

// ContractIDTag captures the possible tag values for $contract_id
type ContractIDTag byte

const (
	// ContractIDTagImplicit is the tag for implicit accounts
	ContractIDTagImplicit ContractIDTag = 0
	// ContractIDTagOriginated is the tag for originated accounts
	ContractIDTagOriginated ContractIDTag = 1
)

// PubKeyHashTag captures the possible tag values for $public_key_hash
type PubKeyHashTag byte

const (
	// PubKeyHashTagEd25519 is the tag for Ed25519 pubkey hashes
	PubKeyHashTagEd25519 PubKeyHashTag = 0
	// PubKeyHashTagSecp256k1 is the tag for Secp256k1 pubkey hashes
	PubKeyHashTagSecp256k1 PubKeyHashTag = 1
	// PubKeyHashTagP256 is the tag for P256 pubkey hashes
	PubKeyHashTagP256 PubKeyHashTag = 2
)

// PubKeyTag captures the possible tag values for $public_key
type PubKeyTag byte

const (
	// PubKeyTagEd25519 is the tag for Ed25519 pubkeys
	PubKeyTagEd25519 PubKeyTag = 0
	// PubKeyTagSecp256k1 is the tag for Secp256k1 pubkeys
	PubKeyTagSecp256k1 PubKeyTag = 1
	// PubKeyTagP256 is the tag for P256 pubkeys
	PubKeyTagP256 PubKeyTag = 2
)

// Watermark is the first byte of a signable payload that indicates
// the type of data represented.
type Watermark byte

// References: https://gitlab.com/tezos/tezos/blob/master/src/lib_crypto/signature.ml#L43
const (
	// BlockHeaderWatermark is the special byte prepended to serialized block headers before signing
	BlockHeaderWatermark Watermark = 1
	// EndorsementWatermark is the special byte prepended to serialized endorsements before signing
	EndorsementWatermark Watermark = 2
	// OperationWatermark is the special byte prepended to serialized operations before signing
	OperationWatermark Watermark = 3
	// CustomWatermark is for custom purposes
	CustomWatermark Watermark = 4
	// TextWatermark is the special byte prepended to plaintext messages before signing. It is not
	// yet part of the standard but has some precedent here:
	// https://tezos.stackexchange.com/questions/1177/whats-the-easiest-way-for-an-account-holder-to-verify-sign-that-they-are-the-ri/1178#1178
	TextWatermark Watermark = 5
)

// AccountType is either an implicit account or an originated account
type AccountType string

const (
	// AccountTypeImplicit indicates an implicit account
	AccountTypeImplicit AccountType = "implicit"
	// AccountTypeOriginated indicates an originated account
	AccountTypeOriginated AccountType = "originated"
)
