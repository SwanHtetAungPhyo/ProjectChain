package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
)

type PublicKeyJSON struct {
	X string `json:"x"`
	Y string `json:"y"`
}

func JSONToPublicKey(publicKeyJSON []byte) (*ecdsa.PublicKey, error) {
	var pubKeyJSON PublicKeyJSON
	if err := json.Unmarshal(publicKeyJSON, &pubKeyJSON); err != nil {
		return nil, err
	}

	xBytes, err := base64.StdEncoding.DecodeString(pubKeyJSON.X)
	if err != nil {
		return nil, err
	}
	yBytes, err := base64.StdEncoding.DecodeString(pubKeyJSON.Y)
	if err != nil {
		return nil, err
	}
	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P384(),
		X:     new(big.Int).SetBytes(xBytes),
		Y:     new(big.Int).SetBytes(yBytes),
	}
	return publicKey, nil
}

func VerifySignature(publicKey *ecdsa.PublicKey, message string, signature []byte) bool {
	if len(signature) != 2*elliptic.P384().Params().BitSize/8 {
		return false
	}
	r := new(big.Int).SetBytes(signature[:len(signature)/2])
	s := new(big.Int).SetBytes(signature[len(signature)/2:])
	hash := sha256.Sum256([]byte(message))
	return ecdsa.Verify(publicKey, hash[:], r, s)
}
