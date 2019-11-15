package tezosprotocol

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/anchorageoss/tezosprotocol/v2/zarith"
	"golang.org/x/xerrors"
)

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
	sourceBytes, err := d.Source.EncodePubKeyHash()
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
			if r := recover(); r != nil {
				err = catchOutOfRangeExceptions(r)
			}
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
	err = d.Source.UnmarshalBinary(dataPtr[:TaggedPubKeyHashLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal source: %w", err)
	}
	dataPtr = dataPtr[TaggedPubKeyHashLen:]

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
		var delegate ContractID
		err = delegate.UnmarshalBinary(taggedPubKeyHash)
		if err != nil {
			return xerrors.Errorf("failed to deserialize delegate: %w", err)
		}
		d.Delegate = &delegate
	}

	return nil
}
