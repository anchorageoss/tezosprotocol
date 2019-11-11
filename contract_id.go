package tezosprotocol

import (
	"bytes"
	"encoding/binary"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/xerrors"
)

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
		return "", xerrors.Errorf("unsupported public key type %T", cryptoPubKey)
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

// UnmarshalBinary implements encoding.BinaryUnmarshaler. It accepts a 22 byte $contract_id or
// a 21 byte $public_key_hash
func (c *ContractID) UnmarshalBinary(data []byte) error {
	var contractIDTag ContractIDTag
	switch len(data) {
	case ContractIDLen:
		contractIDTag = ContractIDTag(data[0])
	case TaggedPubKeyHashLen:
		// prepend a byte to data so we can pretend the caller supplied a $contract_id instead of
		// a $public_key_hash and reuse the same parsing code below
		data = append([]byte{byte(ContractIDTagImplicit)}, data...)
		contractIDTag = ContractIDTagImplicit
	default:
		return xerrors.Errorf("expected %d bytes for contract ID or %d bytes for tagged public key hash; received %d", ContractIDLen, TaggedPubKeyHashLen, len(data))
	}
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
	accountType, err := c.AccountType()
	if err != nil {
		return nil, err
	}
	if accountType != AccountTypeImplicit {
		return nil, xerrors.Errorf("contract ID %s does not represent an implicit account", c)
	}

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
