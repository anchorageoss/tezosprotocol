package tezosprotocol

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/anchorageoss/tezosprotocol/v3/zarith"
	"golang.org/x/xerrors"
)

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
	sourceBytes, err := r.Source.EncodePubKeyHash()
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
			if r := recover(); r != nil {
				err = catchOutOfRangeExceptions(r)
			}
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
	err = r.Source.UnmarshalBinary(dataPtr[:TaggedPubKeyHashLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal source: %w", err)
	}
	dataPtr = dataPtr[TaggedPubKeyHashLen:]

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
