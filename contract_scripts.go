package tezosprotocol

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"math"

	"golang.org/x/xerrors"
)

const maxUint30 = 1<<30 - 1

// ContractScript models $scripted.contracts
type ContractScript struct {
	Code    []byte
	Storage []byte
}

// MarshalBinary implements encoding.BinaryMarshaler. Reference:
// http://tezos.gitlab.io/mainnet/api/p2p.html#contract-id-22-bytes-8-bit-tag
func (c ContractScript) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if len(c.Code) > maxUint30 {
		return nil, xerrors.Errorf("script code cannot exceed %d bytes (uint30_max)", maxUint30)
	}
	if len(c.Storage) > maxUint30 {
		return nil, xerrors.Errorf("script storage cannot exceed %d bytes (uint30_max)", maxUint30)
	}
	err := binary.Write(buf, binary.BigEndian, uint32(len(c.Code)))
	if err != nil {
		return nil, xerrors.Errorf("failed to write code length: %w", err)
	}
	buf.Write(c.Code)
	err = binary.Write(buf, binary.BigEndian, uint32(len(c.Storage)))
	if err != nil {
		return nil, xerrors.Errorf("failed to write storage length: %w", err)
	}
	buf.Write(c.Storage)
	return buf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (c *ContractScript) UnmarshalBinary(data []byte) error {
	var codeLen uint32
	var storageLen uint32
	bytesReader := bytes.NewReader(data)

	// code length
	err := binary.Read(bytesReader, binary.BigEndian, &codeLen)
	if err != nil {
		return xerrors.Errorf("failed to read code length: %w", err)
	}

	// code
	c.Code = make([]byte, codeLen)
	numRead, err := bytesReader.Read(c.Code)
	if err != nil {
		return xerrors.Errorf("failed to read code: %w", err)
	}
	if numRead != int(codeLen) {
		return xerrors.Errorf("failed to read code")
	}

	// storage length
	err = binary.Read(bytesReader, binary.BigEndian, &storageLen)
	if err != nil {
		return xerrors.Errorf("failed to read storage length: %w", err)
	}

	// storage
	c.Storage = make([]byte, storageLen)
	numRead, err = bytesReader.Read(c.Storage)
	if err != nil {
		return xerrors.Errorf("failed to read storage: %w", err)
	}
	if numRead != int(storageLen) {
		return xerrors.Errorf("failed to read storage")
	}

	return nil
}

// EntrypointTag captures the possible tag values for $entrypoint.Tag
type EntrypointTag byte

// EntrypointTag values
const (
	EntrypointTagDefault        EntrypointTag = 0
	EntrypointTagRoot           EntrypointTag = 1
	EntrypointTagDo             EntrypointTag = 2
	EntrypointTagSetDelegate    EntrypointTag = 3
	EntrypointTagRemoveDelegate EntrypointTag = 4
	EntrypointTagNamed          EntrypointTag = 255
)

// Entrypoint models $entrypoint
type Entrypoint struct {
	tag  EntrypointTag
	name string
}

// Preset entrypoints (those with an implicit name)
var (
	EntrypointDefault        = Entrypoint{tag: EntrypointTagDefault}
	EntrypointRoot           = Entrypoint{tag: EntrypointTagRoot}
	EntrypointDo             = Entrypoint{tag: EntrypointTagDo}
	EntrypointSetDelegate    = Entrypoint{tag: EntrypointTagDefault}
	EntrypointRemoveDelegate = Entrypoint{tag: EntrypointTagRemoveDelegate}
)

// NewNamedEntrypoint creates a named entrypoint. This should be used when attempting to
// invoke a custom entrypoint that is not one of the reserved ones (%default, %root, %do, etcetera...).
func NewNamedEntrypoint(name string) (Entrypoint, error) {
	if len(name) > math.MaxUint8 {
		return Entrypoint{}, xerrors.Errorf("entrypoint name %s exceeds maximum length %d", math.MaxUint8)
	}
	return Entrypoint{tag: EntrypointTagNamed, name: name}, nil
}

// Tag returns the entrypoint tag
func (e Entrypoint) Tag() EntrypointTag {
	return e.tag
}

// Name returns the entrypoint name
func (e Entrypoint) Name() (string, error) {
	switch e.tag {
	case EntrypointTagDefault:
		return "default", nil
	case EntrypointTagRoot:
		return "root", nil
	case EntrypointTagDo:
		return "do", nil
	case EntrypointTagSetDelegate:
		return "set_delegate", nil
	case EntrypointTagRemoveDelegate:
		return "remove_delegate", nil
	case EntrypointTagNamed:
		if e.name == "" {
			return "", xerrors.Errorf("entrypoint is not named")
		}
		return e.name, nil
	default:
		return "", xerrors.Errorf("unrecognized entrypoint tag: %d", uint8(e.tag))
	}
}

// String implements fmt.Stringer
func (e Entrypoint) String() string {
	name, err := e.Name()
	if err != nil {
		return "<invalid entrypoint>"
	}
	return "%" + name
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (e Entrypoint) MarshalBinary() ([]byte, error) { //nolint:golint,unparam
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(e.tag))
	if e.tag == EntrypointTagNamed {
		buffer.WriteByte(uint8(len(e.name)))
		buffer.WriteString(e.name)
	}
	return buffer.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (e *Entrypoint) UnmarshalBinary(data []byte) error {
	if len(data) < 1 {
		return xerrors.Errorf("too few bytes to unmarshal Entrypoint")
	}
	e.tag = EntrypointTag(data[0])
	if e.tag == EntrypointTagNamed {
		data = data[1:]
		if len(data) < 1 {
			return xerrors.Errorf("too few bytes to unmarshal Entrypoint name length")
		}
		nameLength := data[0]
		data = data[1:]
		if len(data) < int(nameLength) {
			return xerrors.Errorf("too few bytes to unmarshal Entrypoint name")
		}
		e.name = string(data[:nameLength])
	}
	return nil
}

// TransactionParametersValue models $X_o.value
type TransactionParametersValue interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// note: want to create a rich type for this modeling Michelson instructions.
// This stopgap approach allows just using raw byte arrays in the meantime without
// sacrificing forward compatibility.

// TransactionParametersValueRawBytes is an interim way to provide the value for
// transaction parameters, until support for Michelson is added.
type TransactionParametersValueRawBytes []byte

// MarshalBinary implements encoding.BinaryMarshaler.
func (t *TransactionParametersValueRawBytes) MarshalBinary() ([]byte, error) {
	var parameters []byte
	if t != nil {
		parameters = []byte(*t)
	}
	outputBuf := new(bytes.Buffer)
	err := binary.Write(outputBuf, binary.BigEndian, uint32(len(parameters)))
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal parameters length: %w', err")
	}
	outputBuf.Write(parameters)
	return outputBuf.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *TransactionParametersValueRawBytes) UnmarshalBinary(data []byte) error {
	var length uint32
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &length)
	if err != nil {
		return xerrors.Errorf("invalid transaction parameters value: %w", err)
	}
	if len(data) != int(4+length) {
		return xerrors.Errorf("parameters should be %d bytes, but was %d", length, len(data)-4)
	}
	*t = data[4:]
	return nil
}

// TransactionParameters models $X_o.
// Reference: http://tezos.gitlab.io/babylonnet/api/p2p.html#x-0
type TransactionParameters struct {
	Entrypoint Entrypoint
	Value      TransactionParametersValue
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (t TransactionParameters) MarshalBinary() ([]byte, error) {
	buffer := new(bytes.Buffer)
	entrypointBytes, err := t.Entrypoint.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal entrypoint: %w", err)
	}
	buffer.Write(entrypointBytes)
	valueBytes, err := t.Value.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal value: %w", err)
	}
	buffer.Write(valueBytes)
	return buffer.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (t *TransactionParameters) UnmarshalBinary(data []byte) (err error) {
	// cleanly recover from out of bounds exceptions
	defer func() {
		if err == nil {
			if r := recover(); r != nil {
				err = catchOutOfRangeExceptions(r)
			}
		}
	}()
	dataPtr := data
	err = t.Entrypoint.UnmarshalBinary(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal entrypoint: %w", err)
	}
	entrypointBytes, err := t.Entrypoint.MarshalBinary()
	if err != nil {
		return err
	}
	dataPtr = dataPtr[len(entrypointBytes):]
	t.Value = &TransactionParametersValueRawBytes{}
	err = t.Value.UnmarshalBinary(dataPtr)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal value: %w", err)
	}
	return nil
}
