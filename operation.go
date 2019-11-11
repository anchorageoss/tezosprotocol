package tezosprotocol

import (
	"bytes"
	"encoding"
	"fmt"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/xerrors"
)

// OperationContents models one of multiple contents of a tezos operation.
// Reference: http://tezos.gitlab.io/mainnet/api/p2p.html#operation-alpha-contents-determined-from-data-8-bit-tag
type OperationContents interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	fmt.Stringer
	GetTag() ContentsTag
}

// Operation models a tezos operation with variable length contents.
type Operation struct {
	Branch   BranchID
	Contents []OperationContents
}

func (o *Operation) String() string {
	return fmt.Sprintf("Branch: %s, Contents: %s", o.Branch, o.Contents)
}

// MarshalBinary implements encoding.BinaryMarshaler. It encodes the operation
// unsigned, in the format suitable for signing and transmission.
func (o *Operation) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	branchIDBytes, err := o.Branch.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to write branch: %w", err)
	}
	buf.Write(branchIDBytes)

	if len(o.Contents) == 0 {
		return nil, xerrors.New("expected non-zero list of contents in an operation")
	}
	for _, content := range o.Contents {
		contentBytes, err := content.MarshalBinary()
		if err != nil {
			return nil, xerrors.Errorf("failed to marshal operation contents: %#v: %w", content, err)
		}
		buf.Write(contentBytes)
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Operation) UnmarshalBinary(data []byte) (err error) {
	// cleanly recover from out of bounds exceptions
	defer func() {
		if err == nil {
			if r := recover(); r != nil {
				err = catchOutOfRangeExceptions(r)
			}
		}
	}()

	*o = Operation{}
	dataPtr := data
	err = o.Branch.UnmarshalBinary(dataPtr[:BlockHashLen])
	if err != nil {
		return err
	}
	dataPtr = dataPtr[BlockHashLen:]
	for len(dataPtr) > 0 {
		tag := ContentsTag(dataPtr[0])
		var content OperationContents
		switch tag {
		case ContentsTagRevelation:
			content = &Revelation{}
			err = content.UnmarshalBinary(dataPtr)
			if err != nil {
				return xerrors.Errorf("failed to unmarshal revelation: %w", err)
			}
		case ContentsTagTransaction:
			content = &Transaction{}
			err = content.UnmarshalBinary(dataPtr)
			if err != nil {
				return xerrors.Errorf("failed to unmarshal transaction: %w", err)
			}
		case ContentsTagOrigination:
			content = &Origination{}
			err = content.UnmarshalBinary(dataPtr)
			if err != nil {
				return xerrors.Errorf("failed to unmarshal origination: %w", err)
			}
		case ContentsTagDelegation:
			content = &Delegation{}
			err = content.UnmarshalBinary(dataPtr)
			if err != nil {
				return xerrors.Errorf("failed to unmarshal delegation: %w", err)
			}
		default:
			return xerrors.Errorf("unexpected content tag %d", tag)
		}
		o.Contents = append(o.Contents, content)
		marshaled, err := content.MarshalBinary()
		if err != nil {
			return err
		}
		dataPtr = dataPtr[len(marshaled):]
	}

	return nil
}

// SignatureHash returns the hash of the operation to be signed, including watermark
func (o *Operation) SignatureHash() ([]byte, error) {
	operationBytes, err := o.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal operation: %s: %w", o, err)
	}
	bytesWithWatermark := append([]byte{byte(OperationWatermark)}, operationBytes...)
	sigHash := blake2b.Sum256(bytesWithWatermark)
	return sigHash[:], nil
}
