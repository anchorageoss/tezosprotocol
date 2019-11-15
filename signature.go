package tezosprotocol

import "golang.org/x/xerrors"

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
