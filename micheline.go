package tezosprotocol

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

// incomplete Micheline implementation based on https://gitlab.com/tezos/tezos/blob/master/src%2Flib_micheline%2Fmicheline.ml
// the "tags" come from https://gitlab.com/tezos/tezos/blob/master/src%2Flib_micheline%2Fmicheline.ml#L250

const (
	// int
	michelineTagInt byte = iota //nolint
	// string
	michelineTagString
	// sequence
	michelineTagSeq //nolint
	// Prim (no args, annot)
	michelineTagPrim0
	// Prim (no args + annot)
	michelineTagPrim0A //nolint
	// Prim (1 arg, no annot)
	michelineTagPrim1 //nolint
	// Prim (1 arg + annot)
	michelineTagPrim1A //nolint
	// Prim (2 args, no annot)
	michelineTagPrim2 //nolint
	// Prim (2 args + annot)
	michelineTagPrim2A //nolint
	// "application_encoding"
	michelineTagApplication //nolint
	// bytes
	michelineTagBytes //nolint
)

// MichelineNode represents one node in the tree of Micheline expressions
type MichelineNode interface {
	isMichelineNode()
	MarshalBinary() ([]byte, error)
	UnmarshalBinary([]byte) error
}

// MichelineInt represents an integer in a Micheline expression
type MichelineInt big.Int

func (*MichelineInt) isMichelineNode() {}

// MarshalBinary implements the MichelineNode interface
func (m MichelineInt) MarshalBinary() ([]byte, error) {
	panic("not implemented")
}

// UnmarshalBinary implements the MichelineNode interface
func (m *MichelineInt) UnmarshalBinary([]byte) error {
	panic("not implemented")
}

// MichelineString represents a string in a Micheline expression
type MichelineString string

func (*MichelineString) isMichelineNode() {}

// MarshalBinary implements the MichelineNode interface
func (m MichelineString) MarshalBinary() ([]byte, error) {
	lenBuf := new(bytes.Buffer)
	err := binary.Write(lenBuf, binary.BigEndian, uint32(len(m)))
	return append(append([]byte{michelineTagString}, lenBuf.Bytes()...), []byte(m)...), err
}

// UnmarshalBinary implements the MichelineNode interface
func (m *MichelineString) UnmarshalBinary([]byte) error {
	panic("not implemented")
}

// MichelineBytes represents a byte array in a Micheline expression
type MichelineBytes []byte

func (*MichelineBytes) isMichelineNode() {}

// MarshalBinary implements the MichelineNode interface
func (m MichelineBytes) MarshalBinary() ([]byte, error) {
	panic("not implemented")
}

// UnmarshalBinary implements the MichelineNode interface
func (m *MichelineBytes) UnmarshalBinary([]byte) error {
	panic("not implemented")
}

// MichelinePrim likely represents a Michelson primitive in a Micheline expression
type MichelinePrim struct {
	Prim   byte
	Args   []MichelineNode
	Annots []string
}

func (*MichelinePrim) isMichelineNode() {}

// MarshalBinary implements the MichelineNode interface
func (m MichelinePrim) MarshalBinary() ([]byte, error) { //nolint:unparam
	if len(m.Args) == 0 && len(m.Annots) == 0 {
		return []byte{michelineTagPrim0, m.Prim}, nil
	}
	panic("not implemented")
}

// UnmarshalBinary implements the MichelineNode interface
func (m *MichelinePrim) UnmarshalBinary([]byte) error {
	panic("not implemented")
}

// MichelineSeq represents a sequence of nodes in a Micheline expression
type MichelineSeq []MichelineNode

func (*MichelineSeq) isMichelineNode() {}

// MarshalBinary implements the MichelineNode interface
func (m MichelineSeq) MarshalBinary() ([]byte, error) {
	panic("not implemented")
}

// UnmarshalBinary implements the MichelineNode interface
func (m *MichelineSeq) UnmarshalBinary([]byte) error {
	panic("not implemented")
}
