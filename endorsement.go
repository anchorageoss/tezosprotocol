package tezosprotocol

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/anchorageoss/tezosprotocol/v2/zarith"
	"golang.org/x/xerrors"
)

// Endorsement models the tezos endorsement operation type
type Endorsement struct {
	Level *big.Int
}

func (e *Endorsement) String() string {
	return fmt.Sprintf("%#v", e)
}

// GetTag implements OperationContents
func (e *Endorsement) GetTag() ContentsTag {
	return ContentsTagEndorsement
}

// MarshalBinary implements encoding.BinaryMarshaler
func (e *Endorsement) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	// tag
	buf.WriteByte(byte(e.GetTag()))

	// Level
	level, err := zarith.Encode(e.Level)
	if err != nil {
		return nil, xerrors.Errorf("failed to write Level: %w", err)
	}
	buf.Write(level)

	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (e *Endorsement) UnmarshalBinary(data []byte) (err error) {
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
	if tag != ContentsTagEndorsement {
		return xerrors.Errorf("invalid tag for endorsement. Expected %d, saw %d", ContentsTagEndorsement, tag)
	}
	dataPtr = dataPtr[1:]

	// counter
	var bytesRead int
	e.Level, bytesRead, err = zarith.ReadNext(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal level: %w", err)
	}
	dataPtr = dataPtr[bytesRead:]

	return nil
}
