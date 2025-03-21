package model

import (
	"crypto/sha256"
	"encoding/hex"
)

type Block struct {
	ID           string        `json:"id"`
	Timestamp    string        `json:"timestamp"`
	Hash         string        `json:"hash"`
	Parents      []string      `json:"parents"`
	Transactions []Transaction `json:"transactions"`
	Validators   Validator     `json:"validators"`
}

func GenerateBlockHash(blockID string, parents []string, transactions []Transaction) string {
	data := blockID
	for _, parent := range parents {
		data += parent
	}
	for _, tx := range transactions {
		data += tx.TransactionId
	}
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
