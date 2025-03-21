// internal/model/wallet.go
package model

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"log"
)

type WalletMethod interface {
	generateKeys() (string, string, error)
}

// Wallet represents a user's wallet
// @Description Represents a user's wallet information
// @Model
type Wallet struct {
	Address string `json:"address"`
	PubKey  string `json:"pub_key"`
}

func NewAccount() (*Wallet, string, error) {
	wallet := &Wallet{}
	privateKey, publicKey, err := wallet.generateKeys()
	if err != nil {
		log.Println(err.Error())
		return nil, "", err
	}
	wallet.Address = wallet.Address + "::"
	wallet.PubKey = publicKey
	return wallet, privateKey, nil
}

func (w *Wallet) generateKeys() (string, string, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	publicKeyHex := hex.EncodeToString(privateKey.PublicKey.X.Bytes()) +
		hex.EncodeToString(privateKey.PublicKey.Y.Bytes())

	return privateKeyHex, publicKeyHex, nil
}
