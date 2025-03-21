package model

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
)

type Transaction struct {
	TransactionId  string      `json:"transactionId"`
	ActionTaker    string      `json:"actionTaker"`
	ActionReceiver string      `json:"actionReceiver"`
	Data           interface{} `json:"data"`
	BlockIndex     int64       `json:"blockIndex"`
	Signature      string      `json:"signature"`
}

func SignTransaction(tx *Transaction, privateKeyHex string) error {

	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key")
	}

	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	privateKey.D = new(big.Int).SetBytes(privKeyBytes)

	dataString := ""
	dataMap, ok := tx.Data.(map[string]interface{})
	if !ok {
		log.Printf("cannot convert data to map")
	}
	for key, value := range dataMap {
		var valueStr string
		switch value.(type) {
		case string:
			valueStr += value.(string)
		case float64:
			valueStr += strconv.FormatFloat(value.(float64), 'f', -1, 64)
		case int:
			valueStr += strconv.Itoa(value.(int))
		case bool:
			valueStr += strconv.FormatBool(value.(bool))
		default:
			valueStr = fmt.Sprintf("%v", value)
		}
		dataString += key + valueStr
	}
	hash := sha256.Sum256([]byte(tx.ActionTaker + tx.ActionReceiver + fmt.Sprintf("%d", tx.BlockIndex) + dataString))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return fmt.Errorf("failed to sign transaction")
	}

	tx.Signature = hex.EncodeToString(r.Bytes()) + hex.EncodeToString(s.Bytes())

	return nil
}

func VerifyTransaction(tx Transaction) bool {
	dataMap, ok := tx.Data.(map[string]interface{})
	if !ok {
		log.Printf("cannot convert data to map")
		return false
	}
	var dataString string
	for key, value := range dataMap {
		var valueStr string
		switch value.(type) {
		case string:
			valueStr += value.(string)
		case float64:
			valueStr += strconv.FormatFloat(value.(float64), 'f', -1, 64)
		case int:
			valueStr += strconv.Itoa(value.(int))
		case bool:
			valueStr += strconv.FormatBool(value.(bool))
		default:
			valueStr = fmt.Sprintf("%v", value)
		}
		dataString += key + valueStr
	}
	hash := sha256.Sum256([]byte(tx.ActionTaker + tx.ActionReceiver + fmt.Sprintf("%d", tx.BlockIndex) + dataString))

	pubKeyBytes, err := hex.DecodeString(tx.ActionTaker)
	if err != nil || len(pubKeyBytes) < 64 {
		return false
	}

	pubKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(pubKeyBytes[:32]),
		Y:     new(big.Int).SetBytes(pubKeyBytes[32:]),
	}

	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil || len(sigBytes) < 64 {
		return false
	}
	r := new(big.Int).SetBytes(sigBytes[:32])
	s := new(big.Int).SetBytes(sigBytes[32:])

	return ecdsa.Verify(pubKey, hash[:], r, s)
}
