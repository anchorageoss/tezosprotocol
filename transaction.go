package tezosprotocol

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/anchorageoss/tezosprotocol/v2/zarith"
	"golang.org/x/xerrors"
)

// Transaction models the tezos transaction type
type Transaction struct {
	Source       ContractID
	Fee          *big.Int
	Counter      *big.Int
	GasLimit     *big.Int
	StorageLimit *big.Int
	Amount       *big.Int
	Destination  ContractID
	Parameters   *TransactionParameters
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

// MarshalBinary implements encoding.BinaryMarshaler
func (t *Transaction) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	// tag
	buf.WriteByte(byte(t.GetTag()))

	// source
	sourceBytes, err := t.Source.EncodePubKeyHash()
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

	// parameters
	paramsFollow := t.Parameters != nil
	buf.WriteByte(serializeBoolean(paramsFollow))
	if paramsFollow {
		paramsBytes, err := t.Parameters.MarshalBinary()
		if err != nil {
			return nil, xerrors.Errorf("failed to write transaction parameters: %w", err)
		}
		buf.Write(paramsBytes)
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (t *Transaction) UnmarshalBinary(data []byte) (err error) {
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
	if tag != ContentsTagTransaction {
		return xerrors.Errorf("invalid tag for transaction. Expected %d, saw %d", ContentsTagTransaction, tag)
	}
	dataPtr = dataPtr[1:]

	// source
	err = t.Source.UnmarshalBinary(dataPtr[:TaggedPubKeyHashLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal source: %w", err)
	}
	dataPtr = dataPtr[TaggedPubKeyHashLen:]

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
	dataPtr = dataPtr[1:]
	if err != nil {
		return xerrors.Errorf("failed to deserialialize presence of field \"parameters\": %w", err)
	}
	if hasParameters {
		t.Parameters = &TransactionParameters{Value: &TransactionParametersValueRawBytes{}}
		err = t.Parameters.UnmarshalBinary(dataPtr)
		if err != nil {
			return xerrors.Errorf("failed to deserialize transaction parameters: %w", err)
		}
	}

	return nil
}
