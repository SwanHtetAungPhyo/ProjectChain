package model

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
)

func PublicKeyToJSON(publicKey *ecdsa.PublicKey) ([]byte, error) {
	xBase64 := base64.StdEncoding.EncodeToString(publicKey.X.Bytes())
	yBase64 := base64.StdEncoding.EncodeToString(publicKey.Y.Bytes())

	publicKeyJSON := PublicKeyJSON{
		X: xBase64,
		Y: yBase64,
	}
	return json.Marshal(publicKeyJSON)
}
