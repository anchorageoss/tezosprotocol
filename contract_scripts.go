package tezosprotocol

import (
	"bytes"
	"encoding/binary"

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
