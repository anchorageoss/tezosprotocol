package tezosprotocol

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/xerrors"
)

// PublicKey encodes a tezos public key in base58check encoding
type PublicKey string

// NewPublicKeyFromCryptoPublicKey creates a new PublicKey from a crypto.PublicKey
func NewPublicKeyFromCryptoPublicKey(cryptoPubKey crypto.PublicKey) (PublicKey, error) {
	switch key := cryptoPubKey.(type) {
	case ed25519.PublicKey:
		ret, err := Base58CheckEncode(PrefixEd25519PublicKey, key)
		return PublicKey(ret), err
	case ecdsa.PublicKey:
		switch key.Curve {
		case btcec.S256():
			btcSuitePublicKey := btcec.PublicKey(key)
			compressedPubKeyBytes := btcSuitePublicKey.SerializeCompressed()
			ret, err := Base58CheckEncode(PrefixSecp256k1PublicKey, compressedPubKeyBytes)
			return PublicKey(ret), err
		case elliptic.P256():
			btcSuitePublicKey := btcec.PublicKey(key)
			compressedPubKeyBytes := btcSuitePublicKey.SerializeCompressed()
			ret, err := Base58CheckEncode(PrefixP256PublicKey, compressedPubKeyBytes)
			return PublicKey(ret), err
		default:
			return "", xerrors.Errorf("unsupported curve %s", key.Curve)
		}
	default:
		return "", xerrors.Errorf("unsupported public key type %T", cryptoPubKey)
	}
}

// CryptoPublicKey returns a crypto.PublicKey
func (p PublicKey) CryptoPublicKey() (crypto.PublicKey, error) {
	b58prefix, b58decoded, err := Base58CheckDecode(string(p))
	if err != nil {
		return nil, err
	}
	switch b58prefix {
	case PrefixEd25519PublicKey:
		return ed25519.PublicKey(b58decoded), nil
	case PrefixSecp256k1PublicKey:
		btcecPublicKey, err := btcec.ParsePubKey(b58decoded, btcec.S256())
		if err != nil {
			return nil, err
		}
		return btcecPublicKey.ToECDSA(), nil
	case PrefixP256PublicKey:
		return nil, xerrors.New("unable to deserialize compressed P256 keys")
	default:
		return nil, xerrors.Errorf("unexpected base58check prefix: %s", p)
	}
}

// MarshalBinary implements encoding.BinaryMarshaler. Reference:
// http://tezos.gitlab.io/mainnet/api/p2p.html#public-key-determined-from-data-8-bit-tag
func (p PublicKey) MarshalBinary() ([]byte, error) {
	b58prefix, b58decoded, err := Base58CheckDecode(string(p))
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}

	// write the tag byte
	var expectedPkLength int
	switch b58prefix {
	case PrefixEd25519PublicKey:
		expectedPkLength = PubKeyLenEd25519
		buf.WriteByte(byte(PubKeyTagEd25519))
	case PrefixSecp256k1PublicKey:
		expectedPkLength = PubKeyLenSecp256k1
		buf.WriteByte(byte(PubKeyTagSecp256k1))
	case PrefixP256PublicKey:
		expectedPkLength = PubKeyLenP256
		buf.WriteByte(byte(PubKeyTagP256))
	default:
		return nil, xerrors.Errorf("unexpected base58check prefix: %s", p)
	}

	// write the public key
	if len(b58decoded) != expectedPkLength {
		return nil, xerrors.Errorf("expected public key for addr %s to be %d bytes long, saw %d", p, expectedPkLength, len(b58decoded))
	}
	buf.Write(b58decoded)
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (p *PublicKey) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return xerrors.Errorf("too few bytes to unmarshal public_key")
	}
	pubKeyTag := PubKeyTag(data[0])
	pubKey := data[1:]
	var expectedLength int
	var base58checkPrefix Base58CheckPrefix

	switch pubKeyTag {
	case PubKeyTagEd25519:
		expectedLength = PubKeyLenEd25519
		base58checkPrefix = PrefixEd25519PublicKey
	case PubKeyTagSecp256k1:
		expectedLength = PubKeyLenSecp256k1
		base58checkPrefix = PrefixSecp256k1PublicKey
	case PubKeyTagP256:
		expectedLength = PubKeyLenP256
		base58checkPrefix = PrefixP256PublicKey
	default:
		return xerrors.Errorf("invalid public_key tag %d", pubKeyTag)
	}

	if len(pubKey) < expectedLength {
		return xerrors.Errorf("too few bytes to unmarshal public_key")
	}
	encoded, err := Base58CheckEncode(base58checkPrefix, pubKey[:expectedLength])
	if err != nil {
		return err
	}
	*p = PublicKey(encoded)
	return nil
}

// PrivateKey encodes a tezos private key in base58check encoding
type PrivateKey string

// NewPrivateKeyFromCryptoPrivateKey creates a new PrivateKey from a crypto.PrivateKey
func NewPrivateKeyFromCryptoPrivateKey(cryptoPrivateKey crypto.PrivateKey) (PrivateKey, error) {
	switch key := cryptoPrivateKey.(type) {
	case ed25519.PrivateKey:
		ret, err := Base58CheckEncode(PrefixEd25519SecretKey, []byte(key))
		if err != nil {
			return "", xerrors.Errorf("unable to base58check encode private key: %w", err)
		}
		return PrivateKey(ret), nil
	case *ecdsa.PrivateKey:
		switch key.PublicKey.Curve {
		case btcec.S256():
			btcSuitePrivateKey := btcec.PrivateKey(*key)
			privKeyBytes := btcSuitePrivateKey.Serialize()
			ret, err := Base58CheckEncode(PrefixSecp256k1SecretKey, privKeyBytes)
			return PrivateKey(ret), err
		case elliptic.P256():
			btcSuitePrivateKey := btcec.PrivateKey(*key)
			privKeyBytes := btcSuitePrivateKey.Serialize()
			ret, err := Base58CheckEncode(PrefixP256SecretKey, privKeyBytes)
			return PrivateKey(ret), err
		default:
			return "", xerrors.Errorf("unsupported curve %s", key.Curve)
		}
	default:
		return "", xerrors.Errorf("unsupported private key type %T", cryptoPrivateKey)
	}
}

// CryptoPrivateKey returns a crypto.PrivateKey
func (p PrivateKey) CryptoPrivateKey() (crypto.PrivateKey, error) {
	b58prefix, b58decoded, err := Base58CheckDecode(string(p))
	if err != nil {
		return nil, xerrors.Errorf("unable to base58check decode private key: %w", err)
	}
	switch b58prefix {
	case PrefixEd25519SecretKey:
		return ed25519.PrivateKey(b58decoded), nil
	case PrefixSecp256k1SecretKey:
		privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), b58decoded)
		return privateKey.ToECDSA(), nil
	case PrefixP256SecretKey:
		privateKey, _ := btcec.PrivKeyFromBytes(elliptic.P256(), b58decoded)
		return privateKey.ToECDSA(), nil
	default:
		return nil, xerrors.Errorf("unexpected base58check private key prefix %s", b58prefix)
	}
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (p PrivateKey) MarshalBinary() ([]byte, error) {
	b58prefix, b58decoded, err := Base58CheckDecode(string(p))
	if err != nil {
		return nil, xerrors.Errorf("unable to base58check encode private key: %w", err)
	}
	switch b58prefix {
	case PrefixEd25519SecretKey, PrefixSecp256k1SecretKey, PrefixP256SecretKey:
		return b58decoded, nil
	default:
		return nil, xerrors.Errorf("unexpected base58check private key prefix %s", b58prefix)
	}
}

// PrivateKeySeed encodes a tezos private key seed in base58check encoding.
type PrivateKeySeed string

// PrivateKey returns the private key derived from this private key seed.
func (p PrivateKeySeed) PrivateKey() (PrivateKey, error) {
	b58prefix, seedBytes, err := Base58CheckDecode(string(p))
	if err != nil {
		return "", xerrors.Errorf("failed to base58check decode seed: %w", err)
	}
	switch b58prefix {
	case PrefixEd25519Seed:
		cryptoPrivateKey := ed25519.NewKeyFromSeed(seedBytes)
		return NewPrivateKeyFromCryptoPrivateKey(cryptoPrivateKey)
	default:
		return "", xerrors.Errorf("unsupported private key seed prefix %s", b58prefix)
	}
}
