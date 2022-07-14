package tezosprotocol_test

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"testing"

	"github.com/anchorageoss/tezosprotocol/v2"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
)

var randSeed = bytes.Repeat([]byte{1}, 64)

func fromHex(h string) []byte {
	ret, err := hex.DecodeString(h)
	if err != nil {
		panic(err.Error())
	}
	return ret
}

type keyTest struct {
	KeyType                 string
	ExpectedPrivateKey      tezosprotocol.PrivateKey
	ExpectedPrivateKeyBytes []byte
	ExpectedPublicKey       tezosprotocol.PublicKey
	ExpectedPublicKeyBytes  []byte
	SupportedKeyType        bool
	CanDeserializePublicKey bool
}

var keysTestCases = []keyTest{
	{
		KeyType:                 "Ed25519",
		ExpectedPrivateKey:      tezosprotocol.PrivateKey("edskRc9Pr1NKUW9x6kAZb9cFerBWMo9X9dW4fXwzzL2rvKyKPfdJaJVUcYCfR37sbBujAXJXVJZoCXsUHzfhNcWuqy9aGunQPk"),
		ExpectedPrivateKeyBytes: fromHex("01010101010101010101010101010101010101010101010101010101010101018a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c"),
		ExpectedPublicKey:       tezosprotocol.PublicKey("edpkuhEcwoLysLvodRxQLzuM3AVZvCuT6koVkUahS53mNBdE8LbuGo"),
		ExpectedPublicKeyBytes:  fromHex("008a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c"),
		SupportedKeyType:        true,
		CanDeserializePublicKey: true,
	}, {
		KeyType:                 "secp256k1",
		ExpectedPrivateKey:      tezosprotocol.PrivateKey("spsk1S1KpLsBEXYYw3nQEGHdNQDTjpBsJH9Y86XZVJNobHFkxezaPv"),
		ExpectedPrivateKeyBytes: fromHex("0101010101010101024798bbd525dd3cfffad755af8ea0fffbbb8dec79497fc2"),
		ExpectedPublicKey:       tezosprotocol.PublicKey("sppk7czDjVPj1o3hVLeErZTi6brjZNYGc6jFWzFVvW3oRnki3XB58Yq"),
		ExpectedPublicKeyBytes:  fromHex("0103e4f8056521e0da9cfbb85bf7023d45089588c143e7cf4f784ff319cdc9c42385"),
		SupportedKeyType:        true,
		CanDeserializePublicKey: true,
	}, {
		KeyType:                 "P256",
		ExpectedPrivateKey:      tezosprotocol.PrivateKey("p2sk2Mg6PgZcQ3hvj3SV6CXZvSGthGM9T91YENMMAwemHKx2AJRxU6"),
		ExpectedPrivateKeyBytes: fromHex("02020201fefefeff01445d62b55152b9866561ee015f71beb5a0b12157501662"),
		ExpectedPublicKey:       tezosprotocol.PublicKey("p2pk653txU6DqbwmfVrpRjs3kWsMfFZD2bZxuDoMbNbu3FQ4s557mHT"),
		ExpectedPublicKeyBytes:  fromHex("02023ef92fb44bb6d204854a511f775947ff762d493357c1b91205ba173171f61a2c"),
		SupportedKeyType:        true,
		CanDeserializePublicKey: false,
	}, {
		KeyType:          "P224",
		SupportedKeyType: false,
	}, {
		KeyType:          "P384",
		SupportedKeyType: false,
	}, {
		KeyType:          "P521",
		SupportedKeyType: false,
	}, {
		KeyType:          "RSA4096",
		SupportedKeyType: false,
	},
}

func TestKeys(t *testing.T) {
	require := require.New(t)
	for _, testCase := range keysTestCases {
		var cryptoPrivateKey crypto.PrivateKey
		var cryptoPublicKey crypto.PublicKey
		switch testCase.KeyType {
		case "Ed25519":
			var err error
			cryptoPublicKey, cryptoPrivateKey, err = ed25519.GenerateKey(bytes.NewReader(randSeed))
			require.NoError(err)
		case "secp256k1":
			ecdsaPrivKey, err := ecdsa.GenerateKey(btcec.S256(), bytes.NewReader(randSeed))
			require.NoError(err)
			cryptoPrivateKey = ecdsaPrivKey
			cryptoPublicKey = ecdsaPrivKey.PublicKey
		case "P256":
			ecdsaPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), bytes.NewReader(randSeed))
			require.NoError(err)
			cryptoPrivateKey = ecdsaPrivKey
			cryptoPublicKey = ecdsaPrivKey.PublicKey
		case "P224":
			ecdsaPrivKey, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
			require.NoError(err)
			cryptoPrivateKey = ecdsaPrivKey
			cryptoPublicKey = ecdsaPrivKey.PublicKey
		case "P384":
			ecdsaPrivKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
			require.NoError(err)
			cryptoPrivateKey = ecdsaPrivKey
			cryptoPublicKey = ecdsaPrivKey.PublicKey
		case "P521":
			ecdsaPrivKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
			require.NoError(err)
			cryptoPrivateKey = ecdsaPrivKey
			cryptoPublicKey = ecdsaPrivKey.PublicKey
		case "RSA4096":
			rsaPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
			require.NoError(err)
			cryptoPrivateKey = rsaPrivKey
			cryptoPublicKey = rsaPrivKey.PublicKey
		}

		// private key
		privateKey, err := tezosprotocol.NewPrivateKeyFromCryptoPrivateKey(cryptoPrivateKey)
		require.Equal(err == nil, testCase.SupportedKeyType)
		// public key
		publicKey, err := tezosprotocol.NewPublicKeyFromCryptoPublicKey(cryptoPublicKey)
		require.Equal(err == nil, testCase.SupportedKeyType)

		if privateKey != "" {
			require.Equal(testCase.ExpectedPrivateKey, privateKey)
			cryptoPrivateKey2, err := privateKey.CryptoPrivateKey()
			require.NoError(err)
			require.Equal(cryptoPrivateKey, cryptoPrivateKey2)
			privateKeyBytes, err := privateKey.MarshalBinary()
			require.NoError(err)
			require.Equal(testCase.ExpectedPrivateKeyBytes, privateKeyBytes, hex.EncodeToString(privateKeyBytes))
		}
		if publicKey != "" {
			require.Equal(testCase.ExpectedPublicKey, publicKey)
			publicKeyBytes, err := publicKey.MarshalBinary()
			require.NoError(err)
			require.Equal(testCase.ExpectedPublicKeyBytes, publicKeyBytes, hex.EncodeToString(publicKeyBytes))
			var publicKey2 tezosprotocol.PublicKey
			require.NoError(publicKey2.UnmarshalBinary(publicKeyBytes))
			require.Equal(publicKey, publicKey2)
			_, err = publicKey2.CryptoPublicKey()
			if testCase.CanDeserializePublicKey {
				require.NoError(err)
			} else {
				require.Error(err)
			}
		}
	}
}
