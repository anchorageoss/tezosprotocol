package tezosprotocol

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"golang.org/x/xerrors"
)

// Endorsement models the tezos endorsement operation type
type Endorsement struct {
	Level int32
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
	levelBytesBuf := new(bytes.Buffer)
	err := binary.Write(levelBytesBuf, binary.BigEndian, e.Level)
	if err != nil {
		return []byte(""), err
	}

	_, err = buf.Write(levelBytesBuf.Bytes())
	if err != nil {
		return []byte(""), err
	}

	return buf.Bytes(), nil
}

func readInt32(data []byte) (ret int32, err error) {
	buf := bytes.NewBuffer(data)
	err = binary.Read(buf, binary.BigEndian, &ret)
	return ret, err
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

	// Level
	level, err := readInt32(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal level: %w", err)
	}
	e.Level = level

	return nil
}
