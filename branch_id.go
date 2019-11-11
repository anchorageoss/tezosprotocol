package tezosprotocol

import "golang.org/x/xerrors"

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
