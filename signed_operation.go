package tezosprotocol

import (
	"crypto"
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/btcec/v2"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/xerrors"
)

// SignedOperation represents a signed operation
type SignedOperation struct {
	Operation *Operation
	Signature Signature
}

// Watermark is the first byte of a signable payload that indicates
// the type of data represented.
type Watermark byte

// References: https://gitlab.com/tezos/tezos/blob/master/src/lib_crypto/signature.ml#L43
const (
	// BlockHeaderWatermark is the special byte prepended to serialized block headers before signing
	BlockHeaderWatermark Watermark = 1
	// EndorsementWatermark is the special byte prepended to serialized endorsements before signing
	EndorsementWatermark Watermark = 2
	// OperationWatermark is the special byte prepended to serialized operations before signing
	OperationWatermark Watermark = 3
	// CustomWatermark is for custom purposes
	CustomWatermark Watermark = 4
	// TextWatermark is the special byte prepended to plaintext messages before signing. It is not
	// yet part of the standard but has some precedent here:
	// https://tezos.stackexchange.com/questions/1177/whats-the-easiest-way-for-an-account-holder-to-verify-sign-that-they-are-the-ri/1178#1178
	TextWatermark Watermark = 5
)

// SignOperation signs the given tezos operation using the provided
// signing key. The returned bytes are the signed operation, encoded as
// (operation bytes || signature bytes).
func SignOperation(operation *Operation, privateKey PrivateKey) (SignedOperation, error) {
	// serialize operation
	operationBytes, err := operation.MarshalBinary()
	if err != nil {
		return SignedOperation{}, xerrors.Errorf("failed to marshal operation: %s: %w", operation, err)
	}

	// sign
	signature, err := signGeneric(OperationWatermark, operationBytes, privateKey)
	return SignedOperation{Operation: operation, Signature: signature}, err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (s SignedOperation) MarshalBinary() ([]byte, error) {
	opBytes, err := s.Operation.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal operation: %w", err)
	}
	sigBytes, err := s.Signature.MarshalBinary()
	if err != nil {
		return nil, xerrors.Errorf("failed to marshal signature: %w", err)
	}
	return append(opBytes, sigBytes...), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler. In cases where
// the signature type cannot be inferred, PrefixGenericSignature is used instead.
func (s *SignedOperation) UnmarshalBinary(data []byte) error {
	if len(data) < OperationSignatureLen {
		return xerrors.Errorf("signed operation too short, probably not a signed operation: %d", len(data))
	}

	// operation
	operationLen := len(data) - OperationSignatureLen
	s.Operation = &Operation{}
	err := s.Operation.UnmarshalBinary(data[:operationLen])
	if err != nil {
		return xerrors.Errorf("failed to unmarshal operation in signed operation: %w", err)
	}

	// signature
	signatureBytes := data[operationLen:]
	for _, content := range s.Operation.Contents {
		sourceableContent, ok := content.(interface{ GetSource() ContractID })
		if ok {
			sourceContract := sourceableContent.GetSource()
			var sourceContractType Base58CheckPrefix
			sourceContractType, _, err = Base58CheckDecode(string(sourceContract))
			if err != nil {
				return err
			}
			var signature string
			switch sourceContractType {
			case PrefixEd25519PublicKeyHash:
				signature, err = Base58CheckEncode(PrefixEd25519Signature, signatureBytes)
				s.Signature = Signature(signature)
				return err
			case PrefixP256PublicKeyHash:
				signature, err = Base58CheckEncode(PrefixP256Signature, signatureBytes)
				s.Signature = Signature(signature)
				return err
			case PrefixSecp256k1PublicKeyHash:
				signature, err = Base58CheckEncode(PrefixSecp256k1Signature, signatureBytes)
				s.Signature = Signature(signature)
				return err
			case PrefixContractHash:
				// manager (signer) not known -- continue searching operation contents
			}
		}
	}
	// could not determine signature type -- most likely because the source is an originated account
	signature, err := Base58CheckEncode(PrefixGenericSignature, signatureBytes)
	s.Signature = Signature(signature)
	return err
}

// GetHash returns the hash of a signed operation.
func (s SignedOperation) GetHash() (OperationHash, error) {
	signedOpBytes, err := s.MarshalBinary()
	if err != nil {
		return "", err
	}
	hashBytes := blake2b.Sum256(signedOpBytes)
	var hashEncoded OperationHash
	err = hashEncoded.UnmarshalBinary(hashBytes[:])
	return hashEncoded, err
}

// SignMessage signs the given text based message using the provided
// signing key. It returns the base58check-encoded signature which does not include the message.
// It uses the 0x04 non-standard watermark.
func SignMessage(message string, privateKey PrivateKey) (Signature, error) {
	return signGeneric(TextWatermark, []byte(message), privateKey)
}

func signGeneric(watermark Watermark, message []byte, privateKey PrivateKey) (Signature, error) {
	// prepend the tezos operation watermark
	bytesWithWatermark := append([]byte{byte(watermark)}, message...)

	// hash unsigned operation
	payloadHash := blake2b.Sum256(bytesWithWatermark)

	// sign the hash
	cryptoPrivateKey, err := privateKey.CryptoPrivateKey()
	if err != nil {
		return "", err
	}
	switch key := cryptoPrivateKey.(type) {
	case ed25519.PrivateKey:
		signatureBytes := ed25519.Sign(key, payloadHash[:])
		signature, err := Base58CheckEncode(PrefixEd25519Signature, signatureBytes)
		return Signature(signature), err
	case ecdsa.PrivateKey:
		d := &secp256k1.ModNScalar{}
		d.SetByteSlice(key.D.Bytes())
		btcecPrivKey := btcec.PrivKeyFromScalar(d)
		btcecSignature := btcecdsa.Sign(btcecPrivKey, payloadHash[:])
		signature, err := Base58CheckEncode(PrefixGenericSignature, btcecSignature.Serialize())
		return Signature(signature), err
	default:
		return "", xerrors.Errorf("unsupported private key type: %T", cryptoPrivateKey)
	}
}

// VerifyMessage verifies the signature on a human readable message
func VerifyMessage(message string, signature Signature, publicKey crypto.PublicKey) error {
	return verifyGeneric(TextWatermark, []byte(message), signature, publicKey)
}

func verifyGeneric(watermark Watermark, message []byte, signature Signature, publicKey crypto.PublicKey) error {
	// prepend the tezos operation watermark
	bytesWithWatermark := append([]byte{byte(watermark)}, message...)

	// hash
	payloadHash := blake2b.Sum256(bytesWithWatermark)

	// verify signature over hash
	sigPrefix, sigBytes, err := Base58CheckDecode(string(signature))
	if err != nil {
		return xerrors.Errorf("failed to decode signature: %s: %w", signature, err)
	}
	var ok bool
	switch key := publicKey.(type) {
	case ed25519.PublicKey:
		if sigPrefix != PrefixEd25519Signature && sigPrefix != PrefixGenericSignature {
			return xerrors.Errorf("signature type %s does not match public key type %T", sigPrefix, publicKey)
		}
		ok = ed25519.Verify(key, payloadHash[:], sigBytes)
	default:
		return xerrors.Errorf("unsupported public key type: %T", publicKey)
	}
	if !ok {
		return xerrors.Errorf("invalid signature %s for public key %s", signature, publicKey)
	}
	return nil
}
