package zarith

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/xerrors"
)

// Decode decodes a zarith encoded number from the entire input byte array.
// Assumes the input contains no extra trailing bytes.
func Decode(source []byte) (*big.Int, error) {
	if len(source) == 0 {
		return nil, xerrors.New("expected non-empty byte array")
	}

	// Split input into 8-bit bitstrings
	segments := make([]string, len(source))
	for i, curByte := range source {
		segments[i] = fmt.Sprintf("%08b", curByte)
	}

	// Trim off leading "size" bit from each segment
	for i, segment := range segments {
		segments[i] = segment[1:]
	}

	// Reverse the order of the segments.
	// Source: https://github.com/golang/go/wiki/SliceTricks#reversing
	for i := len(segments)/2 - 1; i >= 0; i-- {
		opp := len(segments) - 1 - i
		segments[i], segments[opp] = segments[opp], segments[i]
	}

	// Concat all the bits
	bitStringBuf := bytes.Buffer{}
	for _, segment := range segments {
		bitStringBuf.WriteString(segment)
	}
	bitString := bitStringBuf.String()

	// Convert from base 2 to base 10
	ret := new(big.Int)
	_, success := ret.SetString(bitString, 2)
	if !success {
		return nil, xerrors.Errorf("failed to parse bit string %s to big.Int", bitString)
	}
	return ret, nil
}

// DecodeHex decodes a zarith encoded number from the entire input hex string.
// Assumes the input contains no extra trailing bytes.
func DecodeHex(source string) (*big.Int, error) {
	bytes, err := hex.DecodeString(source)
	if err != nil {
		return nil, err
	}
	result, err := Decode(bytes)
	return result, err
}

// ReadNext reads the next variable-length zarith number from
// the given byte stream. Returns the zarith number and the count of
// bytes read. Extra bytes are ignored.
func ReadNext(byteStream []byte) (*big.Int, int, error) {
	for n := 0; n < len(byteStream); n++ {
		// if leftmost bit is zero
		if byteStream[n]&uint8(128) == 0 {
			number, err := Decode(byteStream[:n+1])
			return number, n + 1, err
		}
	}
	return nil, -1, xerrors.New("exhausted input while searching for end of next zarith number")
}

// Encode encodes a number to zarith
func Encode(value *big.Int) ([]byte, error) {
	if value == nil {
		value = big.NewInt(0)
	}
	if value.Sign() == -1 {
		return nil, xerrors.Errorf("cannot encode negative integer: %s", value)
	}

	// Convert to base 2 representation
	binaryDigits := value.Text(2)

	// Pad with leading zeros until number of bits is a multiple of 7
	numPaddingBitsRequired := (7*len(binaryDigits) - len(binaryDigits)) % 7
	paddedBinaryDigitsBuffer := bytes.Buffer{}
	for i := 0; i < numPaddingBitsRequired; i++ {
		paddedBinaryDigitsBuffer.WriteString("0")
	}
	paddedBinaryDigitsBuffer.WriteString(binaryDigits)
	paddedBinaryDigits := paddedBinaryDigitsBuffer.String()

	// Split into 7-bit segments
	numSegments := len(paddedBinaryDigits) / 7
	segments := make([]string, numSegments)
	for i := 0; i < numSegments; i++ {
		offset := 7 * i
		segments[i] = paddedBinaryDigits[offset : offset+7]
	}

	// Reverse the order of the segments
	// Source: https://github.com/golang/go/wiki/SliceTricks#reversing
	for i := len(segments)/2 - 1; i >= 0; i-- {
		opp := len(segments) - 1 - i
		segments[i], segments[opp] = segments[opp], segments[i]
	}

	// Prepend a 1 bit to each segment but the last, and a 0 bit to the last
	for i := 0; i < len(segments)-1; i++ {
		segments[i] = "1" + segments[i]
	}
	segments[len(segments)-1] = "0" + segments[len(segments)-1]

	// Concat segments to form the output bitstring
	encodedBitStringBuf := bytes.Buffer{}
	for _, segment := range segments {
		encodedBitStringBuf.WriteString(segment)
	}
	encodedBitString := encodedBitStringBuf.String()

	// Convert from bitstring to byte array
	return bitStringToBytes(encodedBitString), nil
}

// EncodeToHex encodes a number to zarith
func EncodeToHex(value *big.Int) (string, error) {
	bytes, err := Encode(value)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func bitStringToBytes(bitstring string) []byte {
	bytes := make([]byte, len(bitstring)/8)
	for i := 0; i < len(bitstring); i++ {
		bit := bitstring[i]
		if bit < '0' || bit > '1' {
			panic(xerrors.Errorf("%c is not a bit value", bit))
		}
		bytes[i>>3] |= (bit - '0') << uint(7-i&7)
	}
	return bytes
}
