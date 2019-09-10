package tezosprotocol

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"encoding"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/xerrors"

	"github.com/anchorageoss/tezosprotocol/zarith"
	"github.com/btcsuite/btcd/btcec"
)

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

// PubKeyTag captures the possible tag values for $public_key
type PubKeyTag byte

// ContentsTag captures the possible tag values for operation contents
type ContentsTag byte

const (
	// ContentsTagRevelation is the tag for revelations
	ContentsTagRevelation ContentsTag = 7
	// ContentsTagTransaction is the tag for transactions
	ContentsTagTransaction ContentsTag = 8
	// ContentsTagOrigination is the tag for originations
	ContentsTagOrigination ContentsTag = 9
	// ContentsTagDelegation is the tag for delegations
	ContentsTagDelegation ContentsTag = 10
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

func serializeBoolean(b bool) byte {
	if b {
		return byte(255)
	}
	return byte(0)
}

func deserializeBoolean(b byte) (bool, error) {
	switch b {
	case 0:
		return false, nil
	case 255:
		return true, nil
	default:
		return false, xerrors.Errorf("byte value %d not a valid boolean encoding", b)
	}
}

func catchOutOfRangeExceptions() error {
	if r := recover(); r != nil {
		if strings.Contains(fmt.Sprintf("%s", r), "out of range") {
			return xerrors.New("out of bounds exception while parsing operation")
		}
		panic(r)
	}
	return nil
}

// ContractID encodes a tezos contract ID (either implicit or originated) in
// base58check encoding.
type ContractID string

// NewContractIDFromPublicKey creates a new contract ID from a public key.
// AccountType is "implicit."
func NewContractIDFromPublicKey(pubKey PublicKey) (ContractID, error) {
	// pubkey bytes
	cryptoPubKey, err := pubKey.CryptoPublicKey()
	if err != nil {
		return "", err
	}
	var pubKeyBytes []byte
	switch key := cryptoPubKey.(type) {
	case ed25519.PublicKey:
		pubKeyBytes = []byte(key)
	default:
		return "", xerrors.Errorf("unknown public key type %T", cryptoPubKey)
	}

	// pubkey hash
	pubKeyHash, err := blake2b.New(PubKeyHashLen, nil)
	if err != nil {
		panic(xerrors.Errorf("failed to create blake2b hash: %w", err))
	}
	_, err = pubKeyHash.Write(pubKeyBytes)
	if err != nil {
		panic(xerrors.Errorf("failed to write pubkey to hash: %w", err))
	}
	pubKeyHashBytes := pubKeyHash.Sum([]byte{})

	// base58check
	tz1Addr, err := Base58CheckEncode(PrefixEd25519PublicKeyHash, pubKeyHashBytes)
	if err != nil {
		return "", xerrors.Errorf("failed to base58check encode hash: %w", err)
	}

	return ContractID(tz1Addr), nil
}

// NewContractIDFromOrigination returns the address (contract ID) of an account that
// would be originated by this operation. Nonce disambiguates which account in
// the case that multiple accounts would be originated by this same operation.
// Nonce starts at 0 for the first account.
// AccountType is "originated."
func NewContractIDFromOrigination(operationHash OperationHash, nonce uint32) (ContractID, error) {
	contractHash, err := blake2b.New(ContractHashLen, nil)
	if err != nil {
		return "", err
	}

	// operation hash
	operationHashBytes, err := operationHash.MarshalBinary()
	if err != nil {
		return "", err
	}
	_, err = contractHash.Write(operationHashBytes)
	if err != nil {
		return "", err
	}

	// nonce
	nonceBytesBuf := new(bytes.Buffer)
	err = binary.Write(nonceBytesBuf, binary.BigEndian, nonce)
	if err != nil {
		return "", err
	}
	_, err = contractHash.Write(nonceBytesBuf.Bytes())
	if err != nil {
		return "", err
	}

	// encode the hash
	contractHashBytes := contractHash.Sum([]byte{})
	contractHashBytes = append(contractHashBytes, 0) // one byte of padding
	contractIDBytes := append([]byte{byte(ContractIDTagOriginated)}, contractHashBytes...)
	var contractID ContractID
	err = contractID.UnmarshalBinary(contractIDBytes)
	return contractID, err
}

// MarshalBinary implements encoding.BinaryMarshaler. Reference:
// http://tezos.gitlab.io/mainnet/api/p2p.html#contract-id-22-bytes-8-bit-tag
func (c ContractID) MarshalBinary() ([]byte, error) {
	b58prefix, b58decoded, err := Base58CheckDecode(string(c))
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}

	switch b58prefix {
	case PrefixEd25519PublicKeyHash, PrefixSecp256k1PublicKeyHash, PrefixP256PublicKeyHash:
		buf.WriteByte(byte(ContractIDTagImplicit))
		switch b58prefix {
		case PrefixEd25519PublicKeyHash:
			buf.WriteByte(byte(PubKeyHashTagEd25519))
		case PrefixSecp256k1PublicKeyHash:
			buf.WriteByte(byte(PubKeyHashTagSecp256k1))
		case PrefixP256PublicKeyHash:
			buf.WriteByte(byte(PubKeyHashTagP256))
		}
		// public key hash
		if len(b58decoded) != PubKeyHashLen {
			return nil, xerrors.Errorf("expected address %s to have %d bytes for PKH. Saw %d bytes", c, PubKeyHashLen, len(b58decoded))
		}
		buf.Write(b58decoded)

	case PrefixContractHash:
		buf.WriteByte(byte(ContractIDTagOriginated))
		// contract hash
		if len(b58decoded) != ContractHashLen {
			return nil, xerrors.Errorf("saw %d byte contract hash for address %s instead of %d bytes", len(b58decoded), c, ContractHashLen)
		}
		buf.Write(b58decoded)
		// padding
		buf.WriteByte(0)

	default:
		return nil, xerrors.Errorf("unexpected base58check prefix %s in %s", b58prefix, c)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (c *ContractID) UnmarshalBinary(data []byte) error {
	if len(data) < ContractIDLen {
		return xerrors.Errorf("expected %d bytes for contract ID, received %d", ContractIDLen, len(data))
	}
	contractIDTag := ContractIDTag(data[0])
	switch contractIDTag {
	case ContractIDTagImplicit:
		pubKeyHashTag := PubKeyHashTag(data[1])
		pubKeyHash := data[2:]
		switch pubKeyHashTag {
		case PubKeyHashTagEd25519:
			encoded, err := Base58CheckEncode(PrefixEd25519PublicKeyHash, pubKeyHash)
			*c = ContractID(encoded)
			return err
		case PubKeyHashTagSecp256k1:
			encoded, err := Base58CheckEncode(PrefixSecp256k1PublicKeyHash, pubKeyHash)
			*c = ContractID(encoded)
			return err
		case PubKeyHashTagP256:
			encoded, err := Base58CheckEncode(PrefixP256PublicKeyHash, pubKeyHash)
			*c = ContractID(encoded)
			return err
		default:
			return xerrors.Errorf("unexpected pub_key_hash tag %d", pubKeyHashTag)
		}
	case ContractIDTagOriginated:
		contractHash := data[1 : 1+ContractHashLen]
		encoded, err := Base58CheckEncode(PrefixContractHash, contractHash)
		*c = ContractID(encoded)
		return err
	default:
		return xerrors.Errorf("unexpected contract ID tag %d", contractIDTag)
	}
}

// EncodePubKeyHash returns the public key hash corresponding to this contract
// ID. This is only possible for implicit addresses, which are themselves just
// a base58check encoding of a public key hash. Method returns an error for
// originated addresses, whose public key hashes are not inferrable from their
// contract ID.
func (c ContractID) EncodePubKeyHash() ([]byte, error) {
	b58prefix, _, err := Base58CheckDecode(string(c))
	if err != nil {
		return nil, err
	}

	switch b58prefix {
	case PrefixEd25519PublicKeyHash, PrefixSecp256k1PublicKeyHash, PrefixP256PublicKeyHash:
		binaryEncoded, err := c.MarshalBinary()
		if err != nil {
			return nil, err
		}
		// implicit address encoding is a 0 byte plus the PKH encoding
		return binaryEncoded[1:], nil
	default:
		return nil, xerrors.Errorf("can't infer pubkeyhash for %s", c)
	}
}

// AccountType returns the account type represented by this contract ID
func (c ContractID) AccountType() (AccountType, error) {
	b58prefix, _, err := Base58CheckDecode(string(c))
	if err != nil {
		return "", xerrors.Errorf("invalid base58check: %q: %w", c, err)
	}

	switch b58prefix {
	case PrefixEd25519PublicKeyHash, PrefixSecp256k1PublicKeyHash, PrefixP256PublicKeyHash:
		return AccountTypeImplicit, nil
	case PrefixContractHash:
		return AccountTypeOriginated, nil
	default:
		return "", xerrors.Errorf("unknown contract type for %q", c)
	}
}

// BranchID encodes a tezos branch ID in base58check encoding
type BranchID string

// MarshalBinary implements encoding.BinaryMarshaler.
func (b BranchID) MarshalBinary() ([]byte, error) {
	b58prefix, b58decoded, err := Base58CheckDecode(string(b))
	if err != nil {
		return nil, err
	}
	if b58prefix != PrefixBlockHash {
		return nil, xerrors.Errorf("unexpected base58check prefix for branch ID %s", b)
	}
	return b58decoded, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (b *BranchID) UnmarshalBinary(data []byte) error {
	if len(data) != BlockHashLen {
		return xerrors.Errorf("expect branch ID to be %d bytes but received %d", BlockHashLen, len(data))
	}
	b58checkEncoded, err := Base58CheckEncode(PrefixBlockHash, data)
	if err != nil {
		return err
	}
	*b = BranchID(b58checkEncoded)
	return nil
}

// OperationHash encodes an operation hash in base58check encoding
type OperationHash string

// MarshalBinary implements encoding.BinaryMarshaler.
func (o OperationHash) MarshalBinary() ([]byte, error) {
	b58prefix, b58decoded, err := Base58CheckDecode(string(o))
	if err != nil {
		return nil, err
	}
	if b58prefix != PrefixOperationHash {
		return nil, xerrors.Errorf("unexpected base58check prefix for operation hash %s", o)
	}
	return b58decoded, nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (o *OperationHash) UnmarshalBinary(data []byte) error {
	if len(data) != OperationHashLen {
		return xerrors.Errorf("expect operation hash to be %d bytes but received %d", OperationHashLen, len(data))
	}
	b58checkEncoded, err := Base58CheckEncode(PrefixOperationHash, data)
	if err != nil {
		return err
	}
	*o = OperationHash(b58checkEncoded)
	return nil
}

// Operation models a tezos operation with variable length contents.
type Operation struct {
	Branch   BranchID
	Contents []OperationContents
}

func (o *Operation) String() string {
	return fmt.Sprintf("Branch: %s, Contents: %s", o.Branch, o.Contents)
}

// MarshalBinary implements encoding.BinaryMarshaler. It encodes the operation
// unsigned, in the format suitable for signing and transmission.
func (o *Operation) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	branchIDBytes, err := o.Branch.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write branch: %w", err)
	}
	buf.Write(branchIDBytes)

	if len(o.Contents) == 0 {
		return nil, xerrors.New("expected non-zero list of contents in an operation")
	}
	for _, content := range o.Contents {
		contentBytes, err := content.MarshalBinary()
		if err != nil {
			return nil, xerrors.Errorf("failed to marshal operation contents: %#v: %w", content, err)
		}
		buf.Write(contentBytes)
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Operation) UnmarshalBinary(data []byte) (err error) {
	// cleanly recover from out of bounds exceptions
	defer func() {
		if err == nil {
			err = catchOutOfRangeExceptions()
		}
	}()

	*o = Operation{}
	dataPtr := data
	err = o.Branch.UnmarshalBinary(dataPtr[:BlockHashLen])
	if err != nil {
		return err
	}
	dataPtr = dataPtr[BlockHashLen:]
	for len(dataPtr) > 0 {
		tag := ContentsTag(dataPtr[0])
		var content OperationContents
		switch tag {
		case ContentsTagRevelation:
			content = &Revelation{}
			err = content.UnmarshalBinary(dataPtr)
			if err != nil {
				return xerrors.Errorf("failed to unmarshal revelation: %w", err)
			}
		case ContentsTagTransaction:
			content = &Transaction{}
			err = content.UnmarshalBinary(dataPtr)
			if err != nil {
				return xerrors.Errorf("failed to unmarshal transaction: %w", err)
			}
		case ContentsTagOrigination:
			content = &Origination{}
			err = content.UnmarshalBinary(dataPtr)
			if err != nil {
				return xerrors.Errorf("failed to unmarshal origination: %w", err)
			}
		case ContentsTagDelegation:
			content = &Delegation{}
			err = content.UnmarshalBinary(dataPtr)
			if err != nil {
				return xerrors.Errorf("failed to unmarshal delegation: %w", err)
			}
		default:
			return xerrors.Errorf("unexpected content tag %d", tag)
		}
		o.Contents = append(o.Contents, content)
		marshaled, err := content.MarshalBinary()
		if err != nil {
			return err
		}
		dataPtr = dataPtr[len(marshaled):]
	}

	return nil
}

// SignatureHash returns the hash of the operation to be signed, including watermark
func (o *Operation) SignatureHash() ([]byte, error) {
	operationBytes, err := o.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal operation: %s: %w", o, err)
	}
	bytesWithWatermark := append([]byte{byte(OperationWatermark)}, operationBytes...)
	sigHash := blake2b.Sum256(bytesWithWatermark)
	return sigHash[:], nil
}

// OperationContents models one of multiple contents of a tezos operation.
// Reference: http://tezos.gitlab.io/mainnet/api/p2p.html#operation-alpha-contents-determined-from-data-8-bit-tag
type OperationContents interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	fmt.Stringer
	GetTag() ContentsTag
}

// Revelation models the revelation operation type
type Revelation struct {
	Source       ContractID
	Fee          *big.Int
	Counter      *big.Int
	GasLimit     *big.Int
	StorageLimit *big.Int
	PublicKey    PublicKey
}

func (r *Revelation) String() string {
	return fmt.Sprintf("%#v", r)
}

// GetTag implements OperationContents
func (r *Revelation) GetTag() ContentsTag {
	return ContentsTagRevelation
}

// GetSource returns the operation's source
func (r *Revelation) GetSource() ContractID {
	return r.Source
}

// MarshalBinary implements encoding.BinaryMarshaler
func (r *Revelation) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	// tag
	buf.WriteByte(byte(r.GetTag()))

	// source
	sourceBytes, err := r.Source.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write source: %w", err)
	}
	buf.Write(sourceBytes)

	// fee
	fee, err := zarith.Encode(r.Fee)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Fee: %w", err)
	}
	buf.Write(fee)

	// counter
	counter, err := zarith.Encode(r.Counter)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Counter: %w", err)
	}
	buf.Write(counter)

	// gas limit
	gasLimit, err := zarith.Encode(r.GasLimit)
	if err != nil {
		return nil, xerrors.Errorf("failed to write GasLimit: %w", err)
	}
	buf.Write(gasLimit)

	// storage limit
	storageLimit, err := zarith.Encode(r.StorageLimit)
	if err != nil {
		return nil, xerrors.Errorf("failed to write StorageLimit: %w", err)
	}
	buf.Write(storageLimit)

	// public key
	pubKeyBytes, err := r.PublicKey.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write pubKey: %w", err)
	}
	buf.Write(pubKeyBytes)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (r *Revelation) UnmarshalBinary(data []byte) (err error) {
	// cleanly recover from out of bounds exceptions
	defer func() {
		if err == nil {
			err = catchOutOfRangeExceptions()
		}
	}()

	dataPtr := data

	// tag
	tag := ContentsTag(dataPtr[0])
	if tag != ContentsTagRevelation {
		return xerrors.Errorf("invalid tag for revelation. Expected %d, saw %d", ContentsTagRevelation, tag)
	}
	dataPtr = dataPtr[1:]

	// source
	err = r.Source.UnmarshalBinary(dataPtr[:ContractIDLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal source: %w", err)
	}
	dataPtr = dataPtr[ContractIDLen:]

	// fee
	var bytesRead int
	r.Fee, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal fee: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// counter
	r.Counter, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal counter: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// gas limit
	r.GasLimit, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal gas limit: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// storage limit
	r.StorageLimit, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal storage limit: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// public key
	err = r.PublicKey.UnmarshalBinary(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal public key: %w", err)
	}

	return nil
}

// Transaction models the tezos transaction type
type Transaction struct {
	Source       ContractID
	Fee          *big.Int
	Counter      *big.Int
	GasLimit     *big.Int
	StorageLimit *big.Int
	Amount       *big.Int
	Destination  ContractID
	// TODO: parameters
}

func (t *Transaction) String() string {
	return fmt.Sprintf("%#v", t)
}

// GetTag implements OperationContents
func (t *Transaction) GetTag() ContentsTag {
	return ContentsTagTransaction
}

// GetSource returns the operation's source
func (t *Transaction) GetSource() ContractID {
	return t.Source
}

// MmarshalBinary implements encoding.BinaryMarshaler
func (t *Transaction) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	// tag
	buf.WriteByte(byte(t.GetTag()))

	// source
	sourceBytes, err := t.Source.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write source: %w", err)
	}
	buf.Write(sourceBytes)

	// fee
	fee, err := zarith.Encode(t.Fee)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Fee: %w", err)
	}
	buf.Write(fee)

	// counter
	counter, err := zarith.Encode(t.Counter)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Counter: %w", err)
	}
	buf.Write(counter)

	// gas limit
	gasLimit, err := zarith.Encode(t.GasLimit)
	if err != nil {
		return nil, xerrors.Errorf("failed to write GasLimit: %w", err)
	}
	buf.Write(gasLimit)

	// storage limit
	storageLimit, err := zarith.Encode(t.StorageLimit)
	if err != nil {
		return nil, xerrors.Errorf("failed to write StorageLimit: %w", err)
	}
	buf.Write(storageLimit)

	// amount
	amount, err := zarith.Encode(t.Amount)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Amount: %w", err)
	}
	buf.Write(amount)

	// destination
	destinationBytes, err := t.Destination.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write destination: %w", err)
	}
	buf.Write(destinationBytes)

	// no parameters follow
	buf.WriteByte(0)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (t *Transaction) UnmarshalBinary(data []byte) (err error) {
	// cleanly recover from out of bounds exceptions
	defer func() {
		if err == nil {
			err = catchOutOfRangeExceptions()
		}
	}()

	dataPtr := data

	// tag
	tag := ContentsTag(dataPtr[0])
	if tag != ContentsTagTransaction {
		return xerrors.Errorf("invalid tag for transaction. Expected %d, saw %d", ContentsTagTransaction, tag)
	}
	dataPtr = dataPtr[1:]

	// source
	err = t.Source.UnmarshalBinary(dataPtr[:ContractIDLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal source: %w", err)
	}
	dataPtr = dataPtr[ContractIDLen:]

	// fee
	var bytesRead int
	t.Fee, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal fee: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// counter
	t.Counter, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal counter: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// gas limit
	t.GasLimit, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal gas limit: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// storage limit
	t.StorageLimit, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal storage limit: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// amount
	t.Amount, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal counter: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// destination
	err = t.Destination.UnmarshalBinary(dataPtr[:ContractIDLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal destination: %w", err)
	}
	dataPtr = dataPtr[ContractIDLen:]

	// parameters
	hasParameters, err := deserializeBoolean(dataPtr[0])
	if err != nil {
		return xerrors.Errorf("failed to deserialialize presence of field \"parameters\": %w", err)
	}
	if hasParameters {
		return xerrors.Errorf("deserializing parameters not supported")
	}

	return nil
}

// Origination models the tezos origination operation type.
type Origination struct {
	Source       ContractID
	Fee          *big.Int
	Counter      *big.Int
	GasLimit     *big.Int
	StorageLimit *big.Int
	Manager      ContractID
	Balance      *big.Int
	Spendable    bool
	Delegatable  bool
	Delegate     *ContractID
	// TODO: script
}

func (o *Origination) String() string {
	return fmt.Sprintf("%#v", o)
}

// GetTag implements OperationContents
func (o *Origination) GetTag() ContentsTag {
	return ContentsTagOrigination
}

// GetSource returns the operation's source
func (o *Origination) GetSource() ContractID {
	return o.Source
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Origination) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	// tag
	buf.WriteByte(byte(o.GetTag()))

	// source
	sourceBytes, err := o.Source.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write source: %w", err)
	}
	buf.Write(sourceBytes)

	// fee
	fee, err := zarith.Encode(o.Fee)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Fee: %w", err)
	}
	buf.Write(fee)

	// counter
	counter, err := zarith.Encode(o.Counter)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Counter: %w", err)
	}
	buf.Write(counter)

	// gas limit
	gasLimit, err := zarith.Encode(o.GasLimit)
	if err != nil {
		return nil, xerrors.Errorf("failed to write GasLimit: %w", err)
	}
	buf.Write(gasLimit)

	// storage limit
	storageLimit, err := zarith.Encode(o.StorageLimit)
	if err != nil {
		return nil, xerrors.Errorf("failed to write StorageLimit: %w", err)
	}
	buf.Write(storageLimit)

	// manager pub key hash
	managerPubKeyBytes, err := o.Manager.EncodePubKeyHash()
	if err != nil {
		return nil, xerrors.Errorf("failed to write managerPubKey: %w", err)
	}
	buf.Write(managerPubKeyBytes)

	// balance
	balance, err := zarith.Encode(o.Balance)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Balance: %w", err)
	}
	buf.Write(balance)

	// spendable
	buf.WriteByte(serializeBoolean(o.Spendable))

	// delegatable
	buf.WriteByte(serializeBoolean(o.Delegatable))

	// delegate
	hasDelegate := o.Delegate != nil
	buf.WriteByte(serializeBoolean(hasDelegate))
	if hasDelegate {
		delegatePubKeyHashBytes, err := o.Delegate.EncodePubKeyHash()
		if err != nil {
			return nil, xerrors.Errorf("failed to write delegate: %w", err)
		}
		buf.Write(delegatePubKeyHashBytes)
	}

	// script
	hasScript := false
	buf.WriteByte(serializeBoolean(hasScript))

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Origination) UnmarshalBinary(data []byte) (err error) {
	// cleanly recover from out of bounds exceptions
	defer func() {
		if err == nil {
			err = catchOutOfRangeExceptions()
		}
	}()

	dataPtr := data

	// tag
	tag := ContentsTag(dataPtr[0])
	if tag != ContentsTagOrigination {
		return xerrors.Errorf("invalid tag for origination. Expected %d, saw %d", ContentsTagOrigination, tag)
	}
	dataPtr = dataPtr[1:]

	// source
	err = o.Source.UnmarshalBinary(dataPtr[:ContractIDLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal source: %w", err)
	}
	dataPtr = dataPtr[ContractIDLen:]

	// fee
	var bytesRead int
	o.Fee, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal fee: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// counter
	o.Counter, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal counter: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// gas limit
	o.GasLimit, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal gas limit: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// storage limit
	o.StorageLimit, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal storage limit: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// manager (from pub key hash)
	taggedPubKeyHash := dataPtr[:TaggedPubKeyHashLen]
	managerContractID := append([]byte{byte(ContractIDTagImplicit)}, taggedPubKeyHash...)
	err = o.Manager.UnmarshalBinary(managerContractID)
	dataPtr = dataPtr[TaggedPubKeyHashLen:]

	// balance
	o.Balance, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal balance: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// spendable
	o.Spendable, err = deserializeBoolean(dataPtr[0])
	if err != nil {
		return xerrors.Errorf("failed to deserialize spendable: %w", err)
	}
	dataPtr = dataPtr[1:]

	// delegatable
	o.Delegatable, err = deserializeBoolean(dataPtr[0])
	if err != nil {
		return xerrors.Errorf("failed to deserialize delegatable: %w", err)
	}
	dataPtr = dataPtr[1:]

	// delegate
	hasDelegate, err := deserializeBoolean(dataPtr[0])
	if err != nil {
		return xerrors.Errorf("failed to deserialize presence of field \"delegate\": %w", err)
	}
	dataPtr = dataPtr[1:]
	if hasDelegate {
		taggedPubKeyHash = dataPtr[:TaggedPubKeyHashLen]
		delegateContractIDBytes := append([]byte{byte(ContractIDTagImplicit)}, taggedPubKeyHash...)
		var delegate ContractID
		err = delegate.UnmarshalBinary(delegateContractIDBytes)
		if err != nil {
			return xerrors.Errorf("failed to deserialize delegate: %w", err)
		}
		o.Delegate = &delegate
		dataPtr = dataPtr[TaggedPubKeyHashLen:]
	}

	// script
	hasScript, err := deserializeBoolean(dataPtr[0])
	if err != nil {
		return xerrors.Errorf("failed to deserialize presence of field \"script\": %w", err)
	}
	if hasScript {
		return xerrors.New("deserializing scripts not yet supported")
	}

	return nil
}

// Delegation models the tezos delegation operation type
type Delegation struct {
	Source       ContractID
	Fee          *big.Int
	Counter      *big.Int
	GasLimit     *big.Int
	StorageLimit *big.Int
	Delegate     *ContractID
}

func (d *Delegation) String() string {
	return fmt.Sprintf("%#v", d)
}

// GetTag implements OperationContents
func (d *Delegation) GetTag() ContentsTag {
	return ContentsTagDelegation
}

// GetSource returns the operation's source
func (d *Delegation) GetSource() ContractID {
	return d.Source
}

// MarshalBinary implements encoding.BinaryMarshaler
func (d *Delegation) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	// tag
	buf.WriteByte(byte(d.GetTag()))

	// source
	sourceBytes, err := d.Source.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write source: %w", err)
	}
	buf.Write(sourceBytes)

	// fee
	fee, err := zarith.Encode(d.Fee)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Fee: %w", err)
	}
	buf.Write(fee)

	// counter
	counter, err := zarith.Encode(d.Counter)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Counter: %w", err)
	}
	buf.Write(counter)

	// gas limit
	gasLimit, err := zarith.Encode(d.GasLimit)
	if err != nil {
		return nil, xerrors.Errorf("failed to write GasLimit: %w", err)
	}
	buf.Write(gasLimit)

	// storage limit
	storageLimit, err := zarith.Encode(d.StorageLimit)
	if err != nil {
		return nil, xerrors.Errorf("failed to write StorageLimit: %w", err)
	}
	buf.Write(storageLimit)

	// delegate
	hasDelegate := d.Delegate != nil
	buf.WriteByte(serializeBoolean(hasDelegate))
	if hasDelegate {
		delegatePubKeyHashBytes, err := d.Delegate.EncodePubKeyHash()
		if err != nil {
			return nil, xerrors.Errorf("failed to write delegate: %w", err)
		}
		buf.Write(delegatePubKeyHashBytes)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (d *Delegation) UnmarshalBinary(data []byte) (err error) {
	// cleanly recover from out of bounds exceptions
	defer func() {
		if err == nil {
			err = catchOutOfRangeExceptions()
		}
	}()

	dataPtr := data

	// tag
	tag := ContentsTag(dataPtr[0])
	if tag != ContentsTagDelegation {
		return xerrors.Errorf("invalid tag for delegation. Expected %d, saw %d", ContentsTagDelegation, tag)
	}
	dataPtr = dataPtr[1:]

	// source
	err = d.Source.UnmarshalBinary(dataPtr[:ContractIDLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal source: %w", err)
	}
	dataPtr = dataPtr[ContractIDLen:]

	// fee
	var bytesRead int
	d.Fee, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal fee: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// counter
	d.Counter, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal counter: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// gas limit
	d.GasLimit, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal gas limit: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// storage limit
	d.StorageLimit, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal storage limit: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	// delegate
	hasDelegate, err := deserializeBoolean(dataPtr[0])
	if err != nil {
		return xerrors.Errorf("failed to deserialize presence of field \"delegate\": %w", err)
	}
	dataPtr = dataPtr[1:]
	if hasDelegate {
		taggedPubKeyHash := dataPtr[:TaggedPubKeyHashLen]
		delegateContractIDBytes := append([]byte{byte(ContractIDTagImplicit)}, taggedPubKeyHash...)
		var delegate ContractID
		err = delegate.UnmarshalBinary(delegateContractIDBytes)
		if err != nil {
			return xerrors.Errorf("failed to deserialize delegate: %w", err)
		}
		d.Delegate = &delegate
	}

	return nil
}

// Signature is a tezos base58check encoded signature. It may be in either the generic or non-generic format.
type Signature string

// MarshalBinary implements encoding.BinaryMarshaler
func (s Signature) MarshalBinary() ([]byte, error) {
	prefix, payload, err := Base58CheckDecode(string(s))
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal signature: %s: %w", s, err)
	}
	switch prefix {
	case PrefixEd25519Signature, PrefixP256Signature, PrefixSecp256k1Signature, PrefixGenericSignature:
		return payload, nil
	default:
		return nil, xerrors.Errorf("unexpected base58check prefix (%s) for signature %s", prefix.String(), s)
	}
}

// SignedOperation represents a signed operation
type SignedOperation struct {
	Operation *Operation
	Signature Signature
}

// SignOperation signs the given tezos operation using the provided
// signing key. The returned bytes are the signed operation, encoded as
// (operation bytes || signature bytes).
func SignOperation(operation *Operation, privateKey PrivateKey) (SignedOperation, error) {
	// serialize operation
	operationBytes, err := operation.MarshalBinary()
	if err != nil {
		return SignedOperation{}, xerrors.Errorf("failed to marshal operation: %s: %w", operation, err)
	}

	// sign
	signature, err := signGeneric(OperationWatermark, operationBytes, privateKey)
	return SignedOperation{Operation: operation, Signature: signature}, err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (s SignedOperation) MarshalBinary() ([]byte, error) {
	opBytes, err := s.Operation.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal operation: %w", err)
	}
	sigBytes, err := s.Signature.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal signature: %w", err)
	}
	return append(opBytes, sigBytes...), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler. In cases where
// the signature type cannot be inferred, PrefixGenericSignature is used instead.
func (s *SignedOperation) UnmarshalBinary(data []byte) error {
	if len(data) < OperationSignatureLen {
		return xerrors.Errorf("signed operation too short, probably not a signed operation: %d", len(data))
	}

	// operation
	operationLen := len(data) - OperationSignatureLen
	s.Operation = &Operation{}
	err := s.Operation.UnmarshalBinary(data[:operationLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal operation in signed operation: %w", err)
	}

	// signature
	signatureBytes := data[operationLen:]
	for _, content := range s.Operation.Contents {
		sourceableContent, ok := content.(interface{ GetSource() ContractID })
		if ok {
			sourceContract := sourceableContent.GetSource()
			var sourceContractType Base58CheckPrefix
			sourceContractType, _, err = Base58CheckDecode(string(sourceContract))
			if err != nil {
				return err
			}
			var signature string
			switch sourceContractType {
			case PrefixEd25519PublicKeyHash:
				signature, err = Base58CheckEncode(PrefixEd25519Signature, signatureBytes)
				s.Signature = Signature(signature)
				return err
			case PrefixP256PublicKeyHash:
				signature, err = Base58CheckEncode(PrefixP256Signature, signatureBytes)
				s.Signature = Signature(signature)
				return err
			case PrefixSecp256k1PublicKeyHash:
				signature, err = Base58CheckEncode(PrefixSecp256k1Signature, signatureBytes)
				s.Signature = Signature(signature)
				return err
			case PrefixContractHash:
				// manager (signer) not known -- continue searching operation contents
			}
		}
	}
	// could not determine signature type -- most likely because the source is an originated account
	signature, err := Base58CheckEncode(PrefixGenericSignature, signatureBytes)
	s.Signature = Signature(signature)
	return err
}

// GetHash returns the hash of a signed operation.
func (s SignedOperation) GetHash() (OperationHash, error) {
	signedOpBytes, err := s.MarshalBinary()
	if err != nil {
		return "", err
	}
	hashBytes := blake2b.Sum256(signedOpBytes)
	var hashEncoded OperationHash
	err = hashEncoded.UnmarshalBinary(hashBytes[:])
	return hashEncoded, err
}

// SignMessage signs the given text based message using the provided
// signing key. It returns the base58check-encoded signature which does not include the message.
// It uses the 0x04 non-standard watermark.
func SignMessage(message string, privateKey PrivateKey) (Signature, error) {
	return signGeneric(TextWatermark, []byte(message), privateKey)
}

func signGeneric(watermark Watermark, message []byte, privateKey PrivateKey) (Signature, error) {
	// prepend the tezos operation watermark
	bytesWithWatermark := append([]byte{byte(watermark)}, message...)

	// hash unsigned operation
	payloadHash := blake2b.Sum256(bytesWithWatermark)

	// sign the hash
	cryptoPrivateKey, err := privateKey.CryptoPrivateKey()
	if err != nil {
		return "", err
	}
	switch key := cryptoPrivateKey.(type) {
	case ed25519.PrivateKey:
		signatureBytes := ed25519.Sign(key, payloadHash[:])
		signature, err := Base58CheckEncode(PrefixEd25519Signature, signatureBytes)
		return Signature(signature), err
	case ecdsa.PrivateKey:
		btcecPrivKey := btcec.PrivateKey(key)
		btcecSignature, err := btcecPrivKey.Sign(payloadHash[:])
		if err != nil {
			return "", err
		}
		signature, err := Base58CheckEncode(PrefixGenericSignature, btcecSignature.Serialize())
		return Signature(signature), err
	default:
		return "", xerrors.Errorf("unsupported private key type: %T", cryptoPrivateKey)
	}
}

// VerifyMessage verifies the signature on a human readable message
func VerifyMessage(message string, signature Signature, publicKey crypto.PublicKey) error {
	return verifyGeneric(TextWatermark, []byte(message), signature, publicKey)
}

func verifyGeneric(watermark Watermark, message []byte, signature Signature, publicKey crypto.PublicKey) error {
	// prepend the tezos operation watermark
	bytesWithWatermark := append([]byte{byte(watermark)}, message...)

	// hash
	payloadHash := blake2b.Sum256(bytesWithWatermark)

	// verify signature over hash
	sigPrefix, sigBytes, err := Base58CheckDecode(string(signature))
	if err != nil {
		return xerrors.Errorf("failed to decode signature: %s: %w", signature, err)
	}
	var ok bool
	switch key := publicKey.(type) {
	case ed25519.PublicKey:
		if sigPrefix != PrefixEd25519Signature && sigPrefix != PrefixGenericSignature {
			return xerrors.Errorf("signature type %s does not match public key type %T", sigPrefix, publicKey)
		}
		ok = ed25519.Verify(key, payloadHash[:], sigBytes)
	default:
		return xerrors.Errorf("unsupported public key type: %T", publicKey)
	}
	if !ok {
		return xerrors.Errorf("invalid signature %s for public key %s", signature, publicKey)
	}
	return nil
}
