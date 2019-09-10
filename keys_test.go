package tezosprotocol_test

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"testing"

	"github.com/anchorageoss/tezosprotocol"
	"github.com/btcsuite/btcd/btcec"
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
}

var keysTestCases = []keyTest{
	{
		KeyType:                 "Ed25519",
		ExpectedPrivateKey:      tezosprotocol.PrivateKey("edskRc9Pr1NKUW9x6kAZb9cFerBWMo9X9dW4fXwzzL2rvKyKPfdJaJVUcYCfR37sbBujAXJXVJZoCXsUHzfhNcWuqy9aGunQPk"),
		ExpectedPrivateKeyBytes: fromHex("01010101010101010101010101010101010101010101010101010101010101018a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c"),
		ExpectedPublicKey:       tezosprotocol.PublicKey("edpkuhEcwoLysLvodRxQLzuM3AVZvCuT6koVkUahS53mNBdE8LbuGo"),
		ExpectedPublicKeyBytes:  fromHex("008a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c"),
	}, {
		KeyType:                 "secp256k1",
		ExpectedPrivateKey:      tezosprotocol.PrivateKey("spsk1S1KpLsBEXYYw3nQEGHdNQDTjpBsJH9Y86XZVJNobHFkxezaPv"),
		ExpectedPrivateKeyBytes: fromHex("0101010101010101024798bbd525dd3cfffad755af8ea0fffbbb8dec79497fc2"),
		ExpectedPublicKey:       tezosprotocol.PublicKey("sppk7czDjVPj1o3hVLeErZTi6brjZNYGc6jFWzFVvW3oRnki3XB58Yq"),
		ExpectedPublicKeyBytes:  fromHex("0103e4f8056521e0da9cfbb85bf7023d45089588c143e7cf4f784ff319cdc9c42385"),
	}, {
		KeyType:                 "P256",
		ExpectedPrivateKey:      tezosprotocol.PrivateKey("p2sk2Mg6PgZcQ3hvj3SV6CXZvSGthGM9T91YENMMAwemHKx2AJRxU6"),
		ExpectedPrivateKeyBytes: fromHex("02020201fefefeff01445d62b55152b9866561ee015f71beb5a0b12157501662"),
		ExpectedPublicKey:       tezosprotocol.PublicKey("p2pk653txU6DqbwmfVrpRjs3kWsMfFZD2bZxuDoMbNbu3FQ4s557mHT"),
		ExpectedPublicKeyBytes:  fromHex("02023ef92fb44bb6d204854a511f775947ff762d493357c1b91205ba173171f61a2c"),
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
		}

		// private key
		privateKey, err := tezosprotocol.NewPrivateKeyFromCryptoPrivateKey(cryptoPrivateKey)
		require.NoError(err)
		require.Equal(testCase.ExpectedPrivateKey, privateKey)
		cryptoPrivateKey2, err := privateKey.CryptoPrivateKey()
		require.NoError(err)
		require.Equal(cryptoPrivateKey, cryptoPrivateKey2)
		privateKeyBytes, err := privateKey.MarshalBinary()
		require.NoError(err)
		require.Equal(testCase.ExpectedPrivateKeyBytes, privateKeyBytes, hex.EncodeToString(privateKeyBytes))

		// public key
		publicKey, err := tezosprotocol.NewPublicKeyFromCryptoPublicKey(cryptoPublicKey)
		require.NoError(err)
		require.Equal(testCase.ExpectedPublicKey, publicKey)
		publicKeyBytes, err := publicKey.MarshalBinary()
		require.NoError(err)
		require.Equal(testCase.ExpectedPublicKeyBytes, publicKeyBytes, hex.EncodeToString(publicKeyBytes))
		var publicKey2 tezosprotocol.PublicKey
		require.NoError(publicKey2.UnmarshalBinary(publicKeyBytes))
		require.Equal(publicKey, publicKey2)
	}
}
