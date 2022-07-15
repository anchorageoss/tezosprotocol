package tezosprotocol

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/anchorageoss/tezosprotocol/v3/zarith"
	"golang.org/x/xerrors"
)

// Origination models the tezos origination operation type.
type Origination struct {
	Source       ContractID
	Fee          *big.Int
	Counter      *big.Int
	GasLimit     *big.Int
	StorageLimit *big.Int
	Balance      *big.Int
	Delegate     *ContractID
	Script       ContractScript
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
	sourceBytes, err := o.Source.EncodePubKeyHash()
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

	// balance
	balance, err := zarith.Encode(o.Balance)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Balance: %w", err)
	}
	buf.Write(balance)

	// delegate
	hasDelegate := o.Delegate != nil
	buf.WriteByte(serializeBoolean(hasDelegate))
	if hasDelegate {
		//nolint:govet
		delegatePubKeyHashBytes, err := o.Delegate.EncodePubKeyHash()
		if err != nil {
			return nil, xerrors.Errorf("failed to write delegate: %w", err)
		}
		buf.Write(delegatePubKeyHashBytes)
	}

	// script
	scriptBytes, err := o.Script.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write Script: %w", err)
	}
	buf.Write(scriptBytes)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Origination) UnmarshalBinary(data []byte) (err error) {
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
	if tag != ContentsTagOrigination {
		return xerrors.Errorf("invalid tag for origination. Expected %d, saw %d", ContentsTagOrigination, tag)
	}
	dataPtr = dataPtr[1:]

	// source
	err = o.Source.UnmarshalBinary(dataPtr[:TaggedPubKeyHashLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal source: %w", err)
	}
	dataPtr = dataPtr[TaggedPubKeyHashLen:]

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

	// balance
	o.Balance, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal balance: %w", err)
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
		o.Delegate = &delegate
		dataPtr = dataPtr[TaggedPubKeyHashLen:]
	}

	// script
	err = o.Script.UnmarshalBinary(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to deserialize script: %w", err)
	}

	return nil
}
