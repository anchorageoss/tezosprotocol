package tezosprotocol

import "golang.org/x/xerrors"

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
