package zarith

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/xerrors"
)

// the rightmost 7 bits of each byte are used for encoding the value of the int. The
// leftmost bit is used to indicate whether more bytes remain
const lengthZarithBitSegment = 7

// for signed zarith integers, the leftmost bit is still the continuation bit,
// and the second-from-the-left bit is the sign flag
const lengthZarithBitSegmentWithSignFlag = lengthZarithBitSegment - 1

// Decode decodes a zarith encoded unsigned integer from the entire input byte array.
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

	// Trim off leading continuation bit from each segment
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

// DecodeHex decodes a zarith encoded unsigned integer from the entire input hex string.
// Assumes the input contains no extra trailing bytes.
func DecodeHex(source string) (*big.Int, error) {
	bytes, err := hex.DecodeString(source)
	if err != nil {
		return nil, err
	}
	result, err := Decode(bytes)
	return result, err
}

// ReadNext reads the next variable-length zarith-encoded unsigned integer from
// the given byte stream. Returns the zarith number and the count of
// bytes read. Extra bytes are ignored.
func ReadNext(byteStream []byte) (*big.Int, int, error) {
	for n := 0; n < len(byteStream); n++ {
		// if leftmost bit is zero
		if byteStream[n]&byte(128) == 0 {
			number, err := Decode(byteStream[:n+1])
			return number, n + 1, err
		}
	}
	return nil, -1, xerrors.New("exhausted input while searching for end of next zarith number")
}

// Encode encodes an unsigned integer to zarith
func Encode(value *big.Int) ([]byte, error) {
	if value == nil {
		value = big.NewInt(0)
	}
	if value.Sign() == -1 {
		return nil, xerrors.Errorf("cannot encode negative integer: %s", value)
	}

	// Convert to base 2 representation
	valueBitstring := value.Text(2)

	// Pad with leading zeros until number of bits is a multiple of 7
	numPaddingBitsRequired := (lengthZarithBitSegment*len(valueBitstring) - len(valueBitstring)) % lengthZarithBitSegment
	paddedBitstringBuffer := bytes.Buffer{}
	for i := 0; i < numPaddingBitsRequired; i++ {
		paddedBitstringBuffer.WriteString("0")
	}
	paddedBitstringBuffer.WriteString(valueBitstring)
	paddedBitString := paddedBitstringBuffer.String()

	// Split into 7-bit segments
	numSegments := len(paddedBitString) / lengthZarithBitSegment
	segments := make([]string, numSegments)
	for i := 0; i < numSegments; i++ {
		offset := lengthZarithBitSegment * i
		segments[i] = paddedBitString[offset : offset+lengthZarithBitSegment]
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
	outputBitStringBuf := bytes.Buffer{}
	for _, segment := range segments {
		outputBitStringBuf.WriteString(segment)
	}
	outputBitString := outputBitStringBuf.String()

	// Convert from bitstring to byte array
	return bitStringToBytes(outputBitString), nil
}

// EncodeToHex encodes an unsigned integer to zarith
func EncodeToHex(value *big.Int) (string, error) {
	bytes, err := Encode(value)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// EncodeSigned encodes a signed integer to zarith
func EncodeSigned(value *big.Int) ([]byte, error) {
	if value == nil || value.Sign() == 0 {
		return []byte{0}, nil
	}
	isNegative := value.Sign() == -1
	signBit := "0"
	if isNegative {
		signBit = "1"
	}

	// Convert to base 2 representation
	valueBitstring := big.NewInt(0).Abs(value).Text(2)
	numValueBits := len(valueBitstring)

	encodingFitsInOneByte := numValueBits <= lengthZarithBitSegmentWithSignFlag

	// Pad with leading zeros until number of bits is a multiple of 7
	var numPaddingBitsRequired int
	if encodingFitsInOneByte {
		numPaddingBitsRequired = lengthZarithBitSegmentWithSignFlag - numValueBits
	} else {
		numBitsAfterFirstSegment := numValueBits - lengthZarithBitSegmentWithSignFlag
		numPaddingBitsRequired = lengthZarithBitSegment - (numBitsAfterFirstSegment % lengthZarithBitSegment)
	}
	paddedBitStringBuffer := bytes.Buffer{}
	for i := 0; i < numPaddingBitsRequired; i++ {
		paddedBitStringBuffer.WriteString("0")
	}
	paddedBitStringBuffer.WriteString(valueBitstring)
	paddedBitString := paddedBitStringBuffer.String()

	// First segment is the rightmost 6 bits of the input value, prefixed with the sign bit
	segments := make([]string, 0)
	firstSegment := paddedBitString[len(paddedBitString)-lengthZarithBitSegmentWithSignFlag:]
	firstSegment = signBit + firstSegment
	segments = append(segments, firstSegment)
	paddedBitString = paddedBitString[:len(paddedBitString)-lengthZarithBitSegmentWithSignFlag] // pop 6 bits from the right

	// Remaining 7-bit segments collected from right to left
	numSevenBitSegments := len(paddedBitString) / 7
	for i := 0; i < numSevenBitSegments; i++ {
		segments = append(segments, paddedBitString[len(paddedBitString)-lengthZarithBitSegment:])
		paddedBitString = paddedBitString[:len(paddedBitString)-lengthZarithBitSegment] // pop 7 bits from the right
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

// EncodeSignedToHex encodes an unsigned integer to zarith
func EncodeSignedToHex(value *big.Int) (string, error) {
	bytes, err := EncodeSigned(value)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// DecodeSigned decodes a zarith encoded signed integer from the entire input byte array.
// Assumes the input contains no extra trailing bytes.
func DecodeSigned(source []byte) (*big.Int, error) {
	if len(source) == 0 {
		return nil, xerrors.New("expected non-empty byte array")
	}

	// Split input into 8-bit bitstrings
	segments := make([]string, len(source))
	for i, curByte := range source {
		segments[i] = fmt.Sprintf("%08b", curByte)
	}

	// Trim off leading continuation bit from each segment
	for i, segment := range segments {
		segments[i] = segment[1:]
	}

	// Trim off the sign flag from the first segment
	firstSegment := []rune(segments[0])
	isNegative := firstSegment[0] == '1'
	segments[0] = string(firstSegment[1:])

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

	// Add sign flag
	if isNegative {
		bitString = "-" + bitString
	}

	// Convert from base 2 to base 10
	ret := new(big.Int)
	_, success := ret.SetString(bitString, 2)
	if !success {
		return nil, xerrors.Errorf("failed to parse bit string %s to big.Int", bitString)
	}
	return ret, nil
}

// DecodeSignedHex decodes a zarith encoded signed integer from the entire input hex string.
// Assumes the input contains no extra trailing bytes.
func DecodeSignedHex(source string) (*big.Int, error) {
	bytes, err := hex.DecodeString(source)
	if err != nil {
		return nil, err
	}
	result, err := DecodeSigned(bytes)
	return result, err
}

// ReadNextSigned reads the next variable-length zarith-encoded signed integer from
// the given byte stream. Returns the zarith number and the count of
// bytes read. Extra bytes are ignored.
func ReadNextSigned(byteStream []byte) (*big.Int, int, error) {
	for n := 0; n < len(byteStream); n++ {
		// if leftmost bit is zero
		if byteStream[n]&byte(128) == 0 {
			number, err := DecodeSigned(byteStream[:n+1])
			return number, n + 1, err
		}
	}
	return nil, -1, xerrors.New("exhausted input while searching for end of next zarith number")
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
