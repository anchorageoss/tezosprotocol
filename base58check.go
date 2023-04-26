package tezosprotocol

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/xerrors"
)

// Base58CheckPrefix in an enum that models a base58check prefix used specifically by tezos
type Base58CheckPrefix int

type base58CheckPrefixInfo struct {
	id            int
	payloadLength int
	prefixBytes   []byte
}

var base58CheckPrefixInfos = map[Base58CheckPrefix]base58CheckPrefixInfo{}

func registerBase58CheckPrefix(info base58CheckPrefixInfo) Base58CheckPrefix {
	if info.payloadLength == 0 {
		panic("no payload length set")
	}
	info.id = len(base58CheckPrefixInfos)
	base58CheckPrefix := Base58CheckPrefix(info.id)
	AllBase58CheckPrefixes = append(AllBase58CheckPrefixes, base58CheckPrefix)
	base58CheckPrefixInfos[base58CheckPrefix] = info
	return base58CheckPrefix
}

// PayloadLength is the number of bytes expected to be in the base58 encoded payload
func (b Base58CheckPrefix) PayloadLength() int {
	return base58CheckPrefixInfos[b].payloadLength
}

// PrefixBytes are the bytes to append as a prefix before base58 encoding
func (b Base58CheckPrefix) PrefixBytes() []byte {
	return base58CheckPrefixInfos[b].prefixBytes
}

// String prints a human regodnizable string representation of this prefix
func (b Base58CheckPrefix) String() string {
	// Try to guess the prefix as a string
	zeros := make([]byte, base58CheckPrefixInfos[b].payloadLength)
	zerosStr, err := Base58CheckEncode(b, zeros)
	if err != nil {
		panic(err)
	}
	ones := bytes.Repeat([]byte{255}, base58CheckPrefixInfos[b].payloadLength)
	onesStr, err := Base58CheckEncode(b, ones)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s(%d)", commonPrefix(zerosStr, onesStr), len(zerosStr))
}

func commonPrefix(a string, bs ...string) string {
	prefix := []byte{}
	for i := 0; i < len(a); i++ {
		c := a[i]
		eq := true
		for _, b := range bs {
			if b[i] != c {
				eq = false
				break
			}
		}
		if eq {
			prefix = append(prefix, c)
		} else {
			return string(prefix)
		}
	}
	return string(prefix)
}

// Base58Check prefixes
var (
	// AllBase58CheckPrefixes is the list of all defined base58check prefixes
	AllBase58CheckPrefixes = []Base58CheckPrefix{}

	PrefixBlockHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{1, 52},
	})
	PrefixOperationHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{5, 116},
	})
	PrefixOperationListHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{133, 233},
	})
	PrefixOperationListListHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{29, 159, 109},
	})
	PrefixProtocolHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{2, 170},
	})
	PrefixContextHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{79, 199},
	})
	PrefixEd25519PublicKeyHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 20,
		prefixBytes:   []byte{6, 161, 159},
	})
	PrefixSecp256k1PublicKeyHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 20,
		prefixBytes:   []byte{6, 161, 161},
	})
	PrefixP256PublicKeyHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 20,
		prefixBytes:   []byte{6, 161, 164},
	})
	PrefixCryptoboxPublicKeyHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 16,
		prefixBytes:   []byte{153, 103},
	})
	PrefixEd25519Seed = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{13, 15, 58, 7},
	})
	PrefixEd25519PublicKey = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{13, 15, 37, 217},
	})
	PrefixSecp256k1SecretKey = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{17, 162, 224, 201},
	})
	PrefixP256SecretKey = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 32,
		prefixBytes:   []byte{16, 81, 238, 189},
	})
	PrefixEd25519EncryptedSeed = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 56,
		prefixBytes:   []byte{7, 90, 60, 179, 41},
	})
	PrefixSecp256k1EncryptedSecretKey = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 56,
		prefixBytes:   []byte{9, 237, 241, 174, 150},
	})
	PrefixP256EncryptedSecretKey = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 56,
		prefixBytes:   []byte{9, 48, 57, 115, 171},
	})
	PrefixSecp256k1PublicKey = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 33,
		prefixBytes:   []byte{3, 254, 226, 86},
	})
	PrefixP256PublicKey = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 33,
		prefixBytes:   []byte{3, 178, 139, 127},
	})
	PrefixSecp256k1Scalar = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 33,
		prefixBytes:   []byte{38, 248, 136},
	})
	PrefixSecp256k1Element = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 33,
		prefixBytes:   []byte{5, 92, 0},
	})
	PrefixEd25519SecretKey = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 64,
		prefixBytes:   []byte{43, 246, 78, 7},
	})
	PrefixEd25519Signature = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 64,
		prefixBytes:   []byte{9, 245, 205, 134, 18},
	})
	PrefixSecp256k1Signature = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 64,
		prefixBytes:   []byte{13, 115, 101, 19, 63},
	})
	PrefixP256Signature = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 64,
		prefixBytes:   []byte{54, 240, 44, 52},
	})
	PrefixGenericSignature = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 64,
		prefixBytes:   []byte{4, 130, 43},
	})
	PrefixChainID = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 4,
		prefixBytes:   []byte{87, 82, 0},
	})
	// PrefixContractHash is referenced from https://gitlab.com/tezos/tezos/blob/master/src/proto_alpha/lib_protocol/contract_hash.ml#L26
	PrefixContractHash = registerBase58CheckPrefix(base58CheckPrefixInfo{
		payloadLength: 20,
		prefixBytes:   []byte{2, 90, 121},
	})
)

func checksum(input []byte) [4]byte {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	cksum := [4]byte{}
	copy(cksum[:], h2[:4])
	return cksum
}

// Base58CheckEncode encodes the given binary payload to base58check. Prefix
// must be a valid tezos base58check prefix.
func Base58CheckEncode(b58Prefix Base58CheckPrefix, input []byte) (string, error) {
	lengthExpected := b58Prefix.PayloadLength()
	if len(input) != lengthExpected {
		return "", xerrors.Errorf("unexpected length when encoding base58 input: %d != %d", len(input), lengthExpected)
	}

	prefixBytes := b58Prefix.PrefixBytes()
	payload := append(prefixBytes, input...)
	cksum := checksum(payload)
	payload = append(payload, cksum[:]...)
	return base58.Encode(payload), nil
}

// Base58CheckDecode decodes the given base58check string and returns the
// payload and prefix. Errors if the given string does not include a tezos
// prefix, or if the checksum does not match.
func Base58CheckDecode(input string) (Base58CheckPrefix, []byte, error) {
	decoded := base58.Decode(input)

	// checksum
	if len(decoded) < 5 {
		return 0, nil, xerrors.Errorf("%s not valid base58check", input)
	}
	var cksum [4]byte
	copy(cksum[:], decoded[len(decoded)-4:])
	if checksum(decoded[:len(decoded)-4]) != cksum {
		return 0, nil, xerrors.Errorf("b58check checksum failed: %s", input)
	}
	decoded = decoded[:len(decoded)-4]

	// prefix
	var b58prefix Base58CheckPrefix
	found := false
	for _, candidateB58Prefix := range AllBase58CheckPrefixes {
		binaryPrefix := candidateB58Prefix.PrefixBytes()
		if bytes.HasPrefix(decoded, binaryPrefix) {
			b58prefix = candidateB58Prefix
			decoded = decoded[len(binaryPrefix):]
			found = true
			break
		}
	}
	if !found {
		return 0, nil, xerrors.Errorf("unknown base58check prefix: %s", input)
	}

	lengthExpected := b58prefix.PayloadLength()
	if len(decoded) != lengthExpected {
		return 0, nil, xerrors.Errorf("unexpected length when decoding base58 input with prefix %s: %d != %d", b58prefix, len(decoded), lengthExpected)
	}

	return b58prefix, decoded, nil
}
